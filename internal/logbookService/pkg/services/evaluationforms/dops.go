package evaluationforms

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/friendsofgo/errors"
	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/internal/logbookService/pkg/common"
	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/activities"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/assessments"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/dbexecutor"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/evaluationforms"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/notifications"
	orthopedicSurgeriesRepo "github.com/vbetsun/surgeon-intern-app/pkg/repo/orthopedicSurgeries"
	surgeriesRepo "github.com/vbetsun/surgeon-intern-app/pkg/repo/surgeries"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/users"
	"github.com/volatiletech/null/v8"
)

type (
	IDops interface {
		SubmitDopsEvaluation(ctx context.Context, evaluation commonModel.DopsEvaluationInput) (*string, error)
		GetDopsEvaluations(ctx context.Context, queryFilter commonModel.DopsQueryFilter, orderBy *commonModel.QueryOrder) ([]*commonModel.DopsEvaluation, error)
		GetDopsEvaluationById(ctx context.Context, dopsId *string, activityId *string) (*commonModel.DopsEvaluation, error)
		ConnectActivityToDopsEvaluation(ctx context.Context, activityID string, dopsEvaluationID string) (*commonModel.OrthopedicSurgeryActivity, error)
		DeleteInProgressEvaluation(ctx context.Context, evaluationID string) error
	}

	Dops struct {
		activitiesRepo          activities.IActivitiesRepo
		assessmentsRepo         assessments.IAssessmentsRepo
		evaluationFormsRepo     evaluationforms.Repo
		orthopedicSurgeriesRepo *orthopedicSurgeriesRepo.Repo
		surgeriesRepo           surgeriesRepo.IRepo
		usersRepo               users.IRepo
		dbexecutor.IDBExecutor
	}
)

func NewDops(activitiesRepo activities.IActivitiesRepo, assessmentsRepo assessments.IAssessmentsRepo, evaluationFormsRepo evaluationforms.Repo, orthopedicSurgeriesRepo *orthopedicSurgeriesRepo.Repo,
	sugeriesRepo surgeriesRepo.IRepo, usersRepo users.IRepo, dbExecutor dbexecutor.IDBExecutor) IDops {
	return &Dops{
		activitiesRepo:          activitiesRepo,
		assessmentsRepo:         assessmentsRepo,
		evaluationFormsRepo:     evaluationFormsRepo,
		orthopedicSurgeriesRepo: orthopedicSurgeriesRepo,
		surgeriesRepo:           sugeriesRepo,
		usersRepo:               usersRepo,
		IDBExecutor:             dbExecutor,
	}
}

func (s Dops) GetDopsEvaluations(ctx context.Context, queryFilter commonModel.DopsQueryFilter, orderBy *commonModel.QueryOrder) ([]*commonModel.DopsEvaluation, error) {
	user, err := s.usersRepo.GetByAuthenticationContext(ctx)
	if err != nil {
		return nil, err
	}

	dops, err := s.evaluationFormsRepo.Dops.ListDopsEvaluationsByFilter(ctx, user.ID, queryFilter, orderBy)
	if err != nil {
		return nil, err
	}

	evaluations := make([]*commonModel.DopsEvaluation, 0)

	for _, dopsEvaluation := range dops {
		mappedDopsConnection := common.MapDopsEvaluationGraphQlModel(dopsEvaluation, true)
		evaluations = append(evaluations, mappedDopsConnection)
	}
	return evaluations, nil
}

func (s Dops) GetDopsEvaluationById(ctx context.Context, dopsId *string, activityId *string) (*commonModel.DopsEvaluation, error) {
	if dopsId != nil {
		dopsEvaluation, err := s.evaluationFormsRepo.Dops.GetDopsEvaluationById(ctx, *dopsId)
		if err != nil {
			return nil, err
		}
		mappedDopsEvaluation := common.MapDopsEvaluationGraphQlModel(dopsEvaluation, true)

		return mappedDopsEvaluation, nil
	} else if activityId != nil {
		dopsEvaluation, err := s.evaluationFormsRepo.Dops.GetDopsEvaluationByActivityId(ctx, *activityId)
		if err != nil {
			return nil, err
		}
		mappedDopsEvaluation := common.MapDopsEvaluationGraphQlModel(dopsEvaluation, true)

		return mappedDopsEvaluation, nil
	} else {
		return nil, nil
	}
}

func (s Dops) SubmitDopsEvaluation(ctx context.Context, evaluation commonModel.DopsEvaluationInput) (*string, error) {
	// Find supervisor supervisorUser
	supervisorUser, err := s.usersRepo.GetByAuthenticationContext(ctx)
	if err != nil {
		return nil, err
	}

	var dopsEvaluation *models.DopsEvaluation
	var activity *models.Activity
	var assessment *models.Assessment

	if evaluation.DopsEvaluationID != nil { // Get existing evaluation
		dopsEvaluation, err = s.evaluationFormsRepo.Dops.GetDopsEvaluationById(ctx, *evaluation.DopsEvaluationID)
		if err != nil {
			return nil, err
		}
		if dopsEvaluation == nil {
			return nil, errors.New(fmt.Sprintf("No existing dops evaluation found with id %s", *evaluation.DopsEvaluationID))
		}
		if dopsEvaluation.R.DopAssessment == nil || dopsEvaluation.R.DopAssessment.R.Activity == nil {
			return nil, errors.New(fmt.Sprintf("No existing activity or assessment found for dops with id %s", *evaluation.DopsEvaluationID))
		}
		assessment = dopsEvaluation.R.DopAssessment
		activity = assessment.R.Activity

	} else if evaluation.ResidentID != nil { // Get resident and make new dops evaluation
		// Find resident user
		residentUser, err := s.usersRepo.GetByID(ctx, *evaluation.ResidentID, false)
		if err != nil {
			return nil, err
		}
		if residentUser == nil {
			return nil, errors.New("No resident user found with provided Id.")
		}

		dopsEvaluation = &models.DopsEvaluation{
			SupervisorID:                supervisorUser.ID,
			ResidentID:                  residentUser.ID,
			OrthopedicSurgeryActivityID: null.StringFromPtr(evaluation.SurgeryActivityID),
			CreatedAt:                   time.Now(),
		}
		assessment = &models.Assessment{}
		activity = &models.Activity{}
	} else {
		return nil, errors.New("Must supply either connection or resident user.")
	}

	dopsEvaluation.InProgress = evaluation.InProgress

	dopsEvaluation.ActiveStep = evaluation.ActiveStep
	dopsEvaluation.CompletedStep = evaluation.CompletedStep
	if evaluation.SurgeryMetadata != nil {
		dopsEvaluation.OccurredAt = evaluation.SurgeryMetadata.OccurredAt
		activity.OccurredAt = evaluation.SurgeryMetadata.OccurredAt
		dopsEvaluation.CaseNotes = evaluation.SurgeryMetadata.CaseNotes
		dopsEvaluation.PatientAge = evaluation.SurgeryMetadata.PatientAge
		dopsEvaluation.PatientGender = evaluation.SurgeryMetadata.PatientGender
	}
	dopsEvaluation.IsEvaluated = !evaluation.InProgress
	if evaluation.Difficulty != nil {
		dopsEvaluation.Difficulty = *evaluation.Difficulty
	}
	dopsEvaluation.DepartmentID = null.StringFromPtr(evaluation.DepartmentID)

	err = dopsEvaluation.Annotations.Marshal(evaluation.Annotations)
	if err != nil {
		return nil, err
	}

	err = s.RunWithTX(ctx, func(tx *sql.Tx) error {
		// Update or add new
		if evaluation.DopsEvaluationID != nil {
			dopsEvaluation, err = s.evaluationFormsRepo.Dops.UpdateDopsEvaluation(ctx, tx, dopsEvaluation)
			activity, err = s.activitiesRepo.UpdateActivity(ctx, tx, activity)
			if err != nil {
				return err
			}
		} else {
			dopsEvaluation, err = s.evaluationFormsRepo.Dops.AddDopsEvaluation(ctx, tx, dopsEvaluation)
			if err != nil {
				return err
			}
			assessment.DopsID = null.StringFrom(dopsEvaluation.ID)
			assessment, err = s.assessmentsRepo.AddAssessment(ctx, tx, assessment)
			if err != nil {
				return err
			}
			activity.AssessmentID = null.StringFrom(assessment.ID)
			activity, err = s.activitiesRepo.AddActivity(ctx, tx, activity)
			if err != nil {
				return err
			}
		}

		if evaluation.SurgeryMetadata != nil {
			// Update surgeries if changed
			surgeries := make([]*models.Surgery, 0)
			for _, surgeryId := range evaluation.SurgeryMetadata.SurgeryIds {
				rSurgery, err := s.surgeriesRepo.GetByID(ctx, surgeryId)
				if err != nil {
					return err
				}
				surgeries = append(surgeries, rSurgery)
			}
			if err := s.evaluationFormsRepo.Dops.UpdateDopsEvaluationSurgeries(ctx, tx, dopsEvaluation.ID, surgeries); err != nil {
				return err
			}
		}

		if !evaluation.InProgress {
			if err = notifications.AddNotification(ctx, tx, *evaluation.ResidentID, commonModel.NotificationAnnotations{
				RelatedID:        dopsEvaluation.ID,
				NotificationURL:  "Test",
				NotificationType: commonModel.NotificationTypeDopsEvaluation,
			}); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return &dopsEvaluation.ID, nil
}

func (s Dops) ConnectActivityToDopsEvaluation(ctx context.Context, activityID string, dopsEvaluationID string) (*commonModel.OrthopedicSurgeryActivity, error) {
	activity, err := s.orthopedicSurgeriesRepo.ActivitiesRepo.GetByID(ctx, activityID)
	if err != nil {
		return nil, err
	}
	dopsEvaluation, err := s.evaluationFormsRepo.Dops.GetDopsEvaluationById(ctx, dopsEvaluationID)
	if err != nil {
		return nil, err
	}
	if dopsEvaluation.OrthopedicSurgeryActivityID.Valid {
		return nil, errors.New("Dops evaluation is already connected to an activity")
	}
	err = s.RunWithTX(ctx, func(tx *sql.Tx) error {
		activity.DopsRequested = true
		if !activity.SupervisorID.Valid {
			activity.SupervisorID = null.StringFrom(dopsEvaluation.SupervisorID)
		}
		activity, err := s.orthopedicSurgeriesRepo.ActivitiesRepo.Update(ctx, tx, activity)
		if err != nil {
			return err
		}

		dopsEvaluation.OrthopedicSurgeryActivityID = null.StringFrom(activity.ID)
		_, err = s.evaluationFormsRepo.Dops.UpdateDopsEvaluation(ctx, tx, dopsEvaluation)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	var mappedActivity = common.MapOrthopedicSurgeryActivity(activity, true, true)
	return mappedActivity, nil
}

func (s Dops) DeleteInProgressEvaluation(ctx context.Context, evaluationID string) error {
	dopsEvaluation, err := s.evaluationFormsRepo.Dops.GetDopsEvaluationById(ctx, evaluationID)
	if err != nil {
		return err
	}
	user, err := s.usersRepo.GetByAuthenticationContext(ctx)
	if err != nil {
		return err
	}
	if dopsEvaluation.SupervisorID == user.ID && dopsEvaluation.InProgress {
		err := s.evaluationFormsRepo.Dops.DeleteDopsEvaluation(ctx, dopsEvaluation)
		if err != nil {
			return err
		}
	}
	return nil
}
