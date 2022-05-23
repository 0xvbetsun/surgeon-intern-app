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
	IMiniCex interface {
		GetMiniCexEvaluations(ctx context.Context, queryFilter commonModel.MiniCexQueryFilter, orderBy *commonModel.QueryOrder) ([]*commonModel.MiniCexEvaluation, error)
		GetMiniCexEvaluationById(ctx context.Context, miniCexEvaluationId string) (*commonModel.MiniCexEvaluation, error)
		GetMiniCexFocuses(ctx context.Context) ([]*commonModel.MiniCexFocus, error)
		GetMiniCexAreasByClinicId(ctx context.Context, departmentID string) ([]*commonModel.MiniCexArea, error)
		SubmitMiniCexEvaluation(ctx context.Context, evaluation commonModel.MiniCexEvaluationInput) (*string, error)
		RequestMiniCexEvaluation(ctx context.Context, request commonModel.MiniCexRequestInput) (*string, error)
		DeleteInProgressEvaluation(ctx context.Context, evaluationID string) error
	}

	MiniCex struct {
		activitiesRepo          activities.IActivitiesRepo
		assessmentsRepo         assessments.IAssessmentsRepo
		evaluationFormsRepo     evaluationforms.Repo
		orthopedicSurgeriesRepo *orthopedicSurgeriesRepo.Repo
		surgeriesRepo           surgeriesRepo.IRepo
		usersRepo               users.IRepo
		dbexecutor.IDBExecutor
	}
)

func NewMiniCex(activitiesRepo activities.IActivitiesRepo, assessmentsRepo assessments.IAssessmentsRepo,
	evaluationFormsRepo evaluationforms.Repo,
	orthopedicSurgeriesRepo *orthopedicSurgeriesRepo.Repo,
	surgeriesRepo surgeriesRepo.IRepo,
	usersRepo users.IRepo, executor dbexecutor.IDBExecutor) IMiniCex {
	return &MiniCex{
		activitiesRepo:          activitiesRepo,
		assessmentsRepo:         assessmentsRepo,
		evaluationFormsRepo:     evaluationFormsRepo,
		orthopedicSurgeriesRepo: orthopedicSurgeriesRepo,
		surgeriesRepo:           surgeriesRepo,
		usersRepo:               usersRepo,
		IDBExecutor:             executor,
	}
}

func (s MiniCex) GetMiniCexEvaluations(ctx context.Context, queryFilter commonModel.MiniCexQueryFilter, orderBy *commonModel.QueryOrder) ([]*commonModel.MiniCexEvaluation, error) {
	user, err := s.usersRepo.GetByAuthenticationContext(ctx)
	if err != nil {
		return nil, err
	}

	miniCex, err := s.evaluationFormsRepo.MiniCex.ListMiniCexEvaluationsByFilter(ctx, user.ID, queryFilter, orderBy)
	if err != nil {
		return nil, err
	}

	evaluations := make([]*commonModel.MiniCexEvaluation, 0)

	for _, miniCexEvaluation := range miniCex {
		mappedMiniCexConnection := common.MapMiniCexEvaluationGraphQlModel(miniCexEvaluation)
		evaluations = append(evaluations, mappedMiniCexConnection)
	}
	return evaluations, nil
}

func (s MiniCex) GetMiniCexEvaluationById(ctx context.Context, miniCexEvaluationId string) (*commonModel.MiniCexEvaluation, error) {
	miniCexEvaluation, err := s.evaluationFormsRepo.MiniCex.GetMiniCexEvaluationById(ctx, miniCexEvaluationId)
	if err != nil {
		return nil, err
	}

	mappedMiniCexEvaluation := common.MapMiniCexEvaluationGraphQlModel(miniCexEvaluation)

	return mappedMiniCexEvaluation, nil
}

func (s MiniCex) GetMiniCexFocuses(ctx context.Context) ([]*commonModel.MiniCexFocus, error) {
	focuses, err := s.evaluationFormsRepo.MiniCex.ListAllMiniCexFocuses(ctx)
	if err != nil {
		return nil, err
	}

	mappedFocuses := mapMiniCexFocusGraphQlModel(focuses)

	return mappedFocuses, nil
}

func (s MiniCex) GetMiniCexAreasByClinicId(ctx context.Context, departmentID string) ([]*commonModel.MiniCexArea, error) {
	areas, err := s.evaluationFormsRepo.MiniCex.ListAllMiniCexAreasByClinic(ctx, departmentID)
	if err != nil {
		return nil, err
	}

	mappedAreas := mapMiniCexAreaGraphQlModel(areas)

	return mappedAreas, nil
}

func (s MiniCex) SubmitMiniCexEvaluation(ctx context.Context, evaluation commonModel.MiniCexEvaluationInput) (*string, error) {
	// Find supervisor user
	supervisorUser, err := s.usersRepo.GetByAuthenticationContext(ctx)
	if err != nil {
		return nil, err
	}

	var miniCexEvaluation *models.MiniCexEvaluation
	var assessment *models.Assessment
	var activity *models.Activity

	if evaluation.ID != nil { // Get existing evaluation
		miniCexEvaluation, err = s.evaluationFormsRepo.MiniCex.GetMiniCexEvaluationById(ctx, *evaluation.ID)
		if miniCexEvaluation.R.MiniCexAssessment == nil || miniCexEvaluation.R.MiniCexAssessment.R.Activity == nil {
			return nil, errors.New(fmt.Sprintf("No existing miniCex evaluation found with id %s", *evaluation.ID))
		}
		assessment = miniCexEvaluation.R.MiniCexAssessment
		activity = assessment.R.Activity
		if err != nil {
			return nil, err
		}
		if miniCexEvaluation == nil {
			return nil, errors.New(fmt.Sprintf("No existing miniCex evaluation found with id %s", *evaluation.ID))
		}
	} else if evaluation.ResidentID != nil { // Get resident and make new miniCex evaluation
		// Find resident user
		residentUser, err := s.usersRepo.GetByID(ctx, *evaluation.ResidentID, false)
		if err != nil {
			return nil, err
		}
		if residentUser == nil {
			return nil, errors.New("No resident user found with provided Id.")
		}

		miniCexEvaluation = &models.MiniCexEvaluation{
			SupervisorID: supervisorUser.ID,
			ResidentID:   residentUser.ID,
			CreatedAt:    time.Now(),
		}
		assessment = &models.Assessment{}
		activity = &models.Activity{}

	} else {
		return nil, errors.New("Must supply either existing evaluation or resident user.")
	}

	miniCexEvaluation.InProgress = evaluation.InProgress
	miniCexEvaluation.ActiveStep = evaluation.ActiveStep
	miniCexEvaluation.CompletedStep = evaluation.CompletedStep
	miniCexEvaluation.IsEvaluated = !evaluation.InProgress
	if evaluation.Difficulty != nil {
		miniCexEvaluation.Difficulty = *evaluation.Difficulty
	}
	miniCexEvaluation.Focuses = evaluation.Focuses
	miniCexEvaluation.Area = evaluation.Area
	miniCexEvaluation.DepartmentID = null.StringFromPtr(evaluation.DepartmentID)
	miniCexEvaluation.OccurredAt = evaluation.OccurredAt
	activity.OccurredAt = evaluation.OccurredAt

	err = miniCexEvaluation.Annotations.Marshal(evaluation.Annotations)
	if err != nil {
		return nil, err
	}

	err = s.RunWithTX(ctx, func(tx *sql.Tx) error {
		// Update or add new
		if evaluation.ID != nil {
			miniCexEvaluation, err = s.evaluationFormsRepo.MiniCex.UpdateMiniCexEvaluation(ctx, tx, miniCexEvaluation)
			if err != nil {
				return err
			}
			activity, err = s.activitiesRepo.UpdateActivity(ctx, tx, activity)
			if err != nil {
				return err
			}
		} else {
			miniCexEvaluation, err = s.evaluationFormsRepo.MiniCex.AddMiniCexEvaluation(ctx, tx, miniCexEvaluation)
			if err != nil {
				return err
			}
			assessment.MiniCexID = null.StringFrom(miniCexEvaluation.ID)
			assessment, err = s.assessmentsRepo.AddAssessment(ctx, tx, assessment)
			if err != nil {
				return err
			}
			activity.AssessmentID = null.StringFrom(assessment.ID)
			_, err = s.activitiesRepo.AddActivity(ctx, tx, activity)
			if err != nil {
				return err
			}
		}

		if !evaluation.InProgress {
			if err = notifications.AddNotification(ctx, tx, *evaluation.ResidentID, commonModel.NotificationAnnotations{
				RelatedID:        miniCexEvaluation.ID,
				NotificationURL:  "Test",
				NotificationType: commonModel.NotificationTypeMiniCexEvaluation,
			}); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &miniCexEvaluation.ID, nil
}

func (s MiniCex) RequestMiniCexEvaluation(ctx context.Context, request commonModel.MiniCexRequestInput) (*string, error) {
	// Find resident user
	residentUser, err := s.usersRepo.GetByAuthenticationContext(ctx)
	if err != nil {
		return nil, err
	}

	var miniCexEvaluation *models.MiniCexEvaluation
	var assessment *models.Assessment
	var activity *models.Activity

	if request.SupervisorID != nil { // Get supervisor user and make new miniCex evaluation
		// Find supervisor user
		supervisorUser, err := s.usersRepo.GetByID(ctx, *request.SupervisorID, false)
		if err != nil {
			return nil, err
		}
		if supervisorUser == nil {
			return nil, errors.New("No supervisor user found with provided Id.")
		}

		miniCexEvaluation = &models.MiniCexEvaluation{
			SupervisorID: supervisorUser.ID,
			ResidentID:   residentUser.ID,
			CreatedAt:    time.Now(),
		}
		assessment = &models.Assessment{}
		activity = &models.Activity{}
	} else {
		return nil, errors.New("Must supply valid supervisor user id.")
	}

	miniCexEvaluation.Difficulty = request.Difficulty
	miniCexEvaluation.DepartmentID = null.StringFrom(request.DepartmentID)
	miniCexEvaluation.OccurredAt = request.OccurredAt
	activity.OccurredAt = request.OccurredAt
	miniCexEvaluation.ActiveStep = 1

	err = miniCexEvaluation.Annotations.Marshal(request.Annotations)
	if err != nil {
		return nil, err
	}

	err = s.RunWithTX(ctx, func(tx *sql.Tx) error {
		miniCexEvaluation, err = s.evaluationFormsRepo.MiniCex.AddMiniCexEvaluation(ctx, tx, miniCexEvaluation)
		if err != nil {
			return err
		}
		assessment.MiniCexID = null.StringFrom(miniCexEvaluation.ID)
		assessment, err = s.assessmentsRepo.AddAssessment(ctx, tx, assessment)
		if err != nil {
			return err
		}

		activity.AssessmentID = null.StringFrom(assessment.ID)
		_, err = s.activitiesRepo.AddActivity(ctx, tx, activity)
		if err != nil {
			return err
		}

		if err = notifications.AddNotification(ctx, tx, miniCexEvaluation.SupervisorID, commonModel.NotificationAnnotations{
			RelatedID:        miniCexEvaluation.ID,
			NotificationURL:  "Test",
			NotificationType: commonModel.NotificationTypeMiniCexRequest,
		}); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &miniCexEvaluation.ID, nil
}

func (s MiniCex) DeleteInProgressEvaluation(ctx context.Context, evaluationID string) error {
	miniCexEvaluation, err := s.evaluationFormsRepo.MiniCex.GetMiniCexEvaluationById(ctx, evaluationID)
	if err != nil {
		return err
	}
	user, err := s.usersRepo.GetByAuthenticationContext(ctx)
	if err != nil {
		return err
	}
	// TODO: Disable if requested by resident
	if miniCexEvaluation.SupervisorID == user.ID && miniCexEvaluation.InProgress {
		err := s.evaluationFormsRepo.MiniCex.DeleteMiniCexEvaluation(ctx, miniCexEvaluation)
		if err != nil {
			return err
		}
	}
	return nil
}

func mapMiniCexFocusGraphQlModel(focuses []*models.MiniCexFocuse) []*commonModel.MiniCexFocus {
	mappedFocuses := make([]*commonModel.MiniCexFocus, 0)

	for _, focus := range focuses {
		mappedFocus := commonModel.MiniCexFocus{
			MiniCexFocusID: focus.ID,
			Name:           focus.Name,
		}
		mappedFocuses = append(mappedFocuses, &mappedFocus)
	}
	return mappedFocuses
}

func mapMiniCexAreaGraphQlModel(areas []*models.MiniCexArea) []*commonModel.MiniCexArea {
	mappedAreas := make([]*commonModel.MiniCexArea, 0)

	for _, area := range areas {
		mappedArea := commonModel.MiniCexArea{
			MiniCexAreaID: area.ID,
			DepartmentID:  area.DepartmentID,
			Name:          area.Name,
		}
		mappedAreas = append(mappedAreas, &mappedArea)
	}
	return mappedAreas
}
