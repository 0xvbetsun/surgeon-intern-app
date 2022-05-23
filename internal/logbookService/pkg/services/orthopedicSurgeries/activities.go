package orthopedicSurgeries

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
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/logbookentries"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/notifications"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/orthopedicSurgeries"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/practicalactivitytypes"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/surgeries"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/users"
	"github.com/volatiletech/null/v8"
)

type (
	IActivities interface {
		SubmitOrthopedicSurgeryActivity(ctx context.Context, input *commonModel.OrthopedicSurgeryActivityInput) (*string, error)
		Get(ctx context.Context, orthopedicSurgeryActivityID string) (*commonModel.OrthopedicSurgeryActivity, error)
		GetActivities(ctx context.Context, queryFilter commonModel.SurgeryLogbookEntryQueryFilter, orderBy *commonModel.QueryOrder) ([]*commonModel.OrthopedicSurgeryActivity, error)
		DeleteInProgressActivity(ctx context.Context, activityID string) error
	}
	Activities struct {
		ActivitiesRepo          activities.IActivitiesRepo
		LogbookEntriesRepo      logbookentries.ILogbookEntriesRepo
		AssessmentsRepo         assessments.IAssessmentsRepo
		OrthopedicSurgeriesRepo *orthopedicSurgeries.Repo
		EvaluationFormsrepo     evaluationforms.Repo
		SurgeriesRepo           surgeries.IRepo
		UsersRepo               users.IRepo
		activityTypesRepo       practicalactivitytypes.IRepo
		dbexecutor.IDBExecutor
	}
)

func NewActivities(activitiesRepo activities.IActivitiesRepo,
	logbookentriesRepo logbookentries.ILogbookEntriesRepo,
	assessmentsRepo assessments.IAssessmentsRepo,
	orthopedicSurgeriesRepo *orthopedicSurgeries.Repo,
	evaluationFormsrepo evaluationforms.Repo,
	surgeriesRepo surgeries.IRepo,
	usersRepo users.IRepo,
	activityTypesRepo practicalactivitytypes.IRepo,
	executor dbexecutor.IDBExecutor) IActivities {
	return &Activities{ActivitiesRepo: activitiesRepo, LogbookEntriesRepo: logbookentriesRepo,
		AssessmentsRepo: assessmentsRepo, OrthopedicSurgeriesRepo: orthopedicSurgeriesRepo,
		EvaluationFormsrepo: evaluationFormsrepo,
		SurgeriesRepo:       surgeriesRepo, UsersRepo: usersRepo, activityTypesRepo: activityTypesRepo, IDBExecutor: executor}
}

func (s *Activities) SubmitOrthopedicSurgeryActivity(ctx context.Context, input *commonModel.OrthopedicSurgeryActivityInput) (*string, error) {
	residentUser, err := s.UsersRepo.GetByAuthenticationContext(ctx)
	if err != nil {
		return nil, err
	}

	var surgeryActivity *models.OrthopedicSurgeryActivity
	var logbookEntry *models.LogbookEntry
	var activity *models.Activity

	if input.ID != nil { // Get existing surgeryActivity
		surgeryActivity, err = s.OrthopedicSurgeriesRepo.ActivitiesRepo.GetByIDWithRels(ctx, *input.ID)
		if err != nil {
			return nil, err
		}
		if surgeryActivity == nil {
			return nil, errors.New(fmt.Sprintf("No existing surgeryActivity found with id %s", *input.ID))
		}
		if surgeryActivity.R.OrthopedicSurgeryLogbookEntry == nil || surgeryActivity.R.OrthopedicSurgeryLogbookEntry.R.Activity == nil {
			return nil, errors.New(fmt.Sprintf("No existing activity or logbook entry found for surgeryActivity with id %s", *input.ID))
		}
		logbookEntry = surgeryActivity.R.OrthopedicSurgeryLogbookEntry
		activity = logbookEntry.R.Activity
	} else { // Create new surgeryActivity
		surgeryActivity = &models.OrthopedicSurgeryActivity{
			ResidentID: residentUser.ID,
			CreatedAt:  time.Now(),
		}
		logbookEntry = &models.LogbookEntry{}
		activity = &models.Activity{}
	}

	if input.SupervisorUserID != nil {
		// Find supervisor user
		supervisorUser, err := s.UsersRepo.GetByID(ctx, *input.SupervisorUserID, false)
		if err != nil {
			return nil, err
		}
		if supervisorUser == nil {
			return nil, errors.New("No supervisor user found with provided Id.")
		}
		surgeryActivity.SupervisorID = null.StringFrom(supervisorUser.ID)
	}

	surgeryType, err := s.activityTypesRepo.GetByType(ctx, commonModel.ActivityTypeSurgery)
	if err != nil {
		return nil, err
	}
	surgeryActivity.PracticalActivityTypeID = surgeryType.ID

	surgeryActivity.InProgress = input.InProgress

	surgeryActivity.ActiveStep = input.ActiveStep
	surgeryActivity.CompletedStep = input.CompletedStep
	if input.SurgeryMetadata != nil {
		surgeryActivity.OccurredAt = input.SurgeryMetadata.OccurredAt
		activity.OccurredAt = input.SurgeryMetadata.OccurredAt
		surgeryActivity.CaseNotes = input.SurgeryMetadata.CaseNotes
		surgeryActivity.PatientAge = input.SurgeryMetadata.PatientAge
		surgeryActivity.PatientGender = input.SurgeryMetadata.PatientGender
	}
	surgeryActivity.OperatorID = null.StringFromPtr(input.OperatorID)
	surgeryActivity.AssistantID = null.StringFromPtr(input.AssistantID)
	surgeryActivity.Comments = input.Comments
	surgeryActivity.Complications = input.Complications
	surgeryActivity.ReviewRequested = input.ReviewRequested
	surgeryActivity.DopsRequested = input.DopsRequested

	if err := surgeryActivity.Annotations.Marshal(input.Annotations); err != nil {
		return nil, err
	}

	err = s.RunWithTX(ctx, func(tx *sql.Tx) error {
		// Update or add new
		if input.ID != nil {
			surgeryActivity, err = s.OrthopedicSurgeriesRepo.ActivitiesRepo.Update(ctx, tx, surgeryActivity)
			if err != nil {
				return err
			}
			activity, err = s.ActivitiesRepo.UpdateActivity(ctx, tx, activity)
			if err != nil {
				return err
			}
		} else {
			surgeryActivity, err = s.OrthopedicSurgeriesRepo.ActivitiesRepo.Add(ctx, tx, surgeryActivity)
			if err != nil {
				return err
			}
			logbookEntry.OrthopedicSurgeryID = null.StringFrom(surgeryActivity.ID)
			logbookEntry, err = s.LogbookEntriesRepo.AddLogbookEntry(ctx, tx, logbookEntry)
			if err != nil {
				return err
			}
			activity.LogbookEntryID = null.StringFrom(logbookEntry.ID)
			_, err = s.ActivitiesRepo.AddActivity(ctx, tx, activity)
			if err != nil {
				return err
			}
		}

		surgeries := make([]*models.Surgery, 0)
		if input.SurgeryMetadata != nil {
			for _, surgeryId := range input.SurgeryMetadata.SurgeryIds {
				rSurgery, err := s.SurgeriesRepo.GetByIDWithRels(ctx, surgeryId)
				if err != nil {
					return err
				}
				surgeries = append(surgeries, rSurgery)
			}
			if err := s.OrthopedicSurgeriesRepo.ActivitiesRepo.UpdateActivitySurgeries(ctx, tx, surgeryActivity.ID, surgeries); err != nil {
				return err
			}
		}

		// TODO: Separate flag on input to determine if we create review ?
		if !input.InProgress && input.ReviewRequested {
			// TODO: Validate input. For example operator and assistant can't be nil.
			activityReview := &models.OrthopedicSurgeriesActivityReview{
				OrthopedicSurgeryActivityID: surgeryActivity.ID,
				CreatedAt:                   time.Now(),
				OccurredAt:                  surgeryActivity.OccurredAt,
				CaseNotes:                   surgeryActivity.CaseNotes,
				PatientAge:                  surgeryActivity.PatientAge,
				PatientGender:               surgeryActivity.PatientGender,
				Annotations:                 surgeryActivity.Annotations,
				ResidentID:                  surgeryActivity.ResidentID,
				SupervisorID:                surgeryActivity.SupervisorID.String,
				OperatorID:                  surgeryActivity.OperatorID.String,
				AssistantID:                 surgeryActivity.AssistantID.String,
				Comments:                    surgeryActivity.Comments,
				Complications:               surgeryActivity.Complications,
				Requested:                   true,
				ActiveStep:                  4,
				CompletedStep:               3,
			}
			activityReview, err = s.OrthopedicSurgeriesRepo.ReviewsRepo.AddOrthopedicSurgeryActivityReview(ctx, tx, activityReview)
			if err != nil {
				return err
			}
			assessment, err := s.AssessmentsRepo.AddAssessment(ctx, tx, &models.Assessment{
				OrthopedicSurgeryReviewID: null.StringFrom(activityReview.ID),
			})
			if err != nil {
				return err
			}
			_, err = s.ActivitiesRepo.AddActivity(ctx, tx, &models.Activity{
				OccurredAt:   surgeryActivity.OccurredAt,
				AssessmentID: null.StringFrom(assessment.ID),
			})
			if err != nil {
				return err
			}
			if err := s.OrthopedicSurgeriesRepo.ReviewsRepo.AddSurgeriesToReview(ctx, tx, activityReview.ID, surgeries); err != nil {
				return err
			}
			if err = notifications.AddNotification(ctx, tx, surgeryActivity.SupervisorID.String, commonModel.NotificationAnnotations{
				RelatedID:        activityReview.ID,
				NotificationURL:  "Test",
				NotificationType: commonModel.NotificationTypeReviewRequest,
			}); err != nil {
				return err
			}
		}
		if !input.InProgress && input.DopsRequested {
			dopsEvaluation := &models.DopsEvaluation{
				OrthopedicSurgeryActivityID: null.StringFrom(surgeryActivity.ID),
				ResidentID:                  surgeryActivity.ResidentID,
				SupervisorID:                surgeryActivity.SupervisorID.String,
				OccurredAt:                  surgeryActivity.OccurredAt,
				CaseNotes:                   surgeryActivity.CaseNotes,
				PatientAge:                  surgeryActivity.PatientAge,
				PatientGender:               surgeryActivity.PatientGender,
				ActiveStep:                  2,
				CompletedStep:               1,
				IsEvaluated:                 false,
				CreatedAt:                   time.Now(),
			}
			err := dopsEvaluation.Annotations.Marshal(make([]*commonModel.EvaluationFormAnnotationsInput, 0))
			if err != nil {
				return err
			}

			dopsEvaluation, err = s.EvaluationFormsrepo.Dops.AddDopsEvaluation(ctx, tx, dopsEvaluation)
			if err != nil {
				return err
			}
			if err := s.EvaluationFormsrepo.Dops.AddSurgeriesToDopsEvaluation(ctx, tx, dopsEvaluation.ID, surgeries); err != nil {
				return err
			}
			assessment, err := s.AssessmentsRepo.AddAssessment(ctx, tx, &models.Assessment{
				DopsID: null.StringFrom(dopsEvaluation.ID),
			})
			if err != nil {
				return err
			}
			_, err = s.ActivitiesRepo.AddActivity(ctx, tx, &models.Activity{
				OccurredAt:   surgeryActivity.OccurredAt,
				AssessmentID: null.StringFrom(assessment.ID),
			})
			if err != nil {
				return err
			}
			if err = notifications.AddNotification(ctx, tx, surgeryActivity.SupervisorID.String, commonModel.NotificationAnnotations{
				RelatedID:        dopsEvaluation.ID,
				NotificationURL:  "Test",
				NotificationType: commonModel.NotificationTypeDopsRequest,
			}); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &surgeryActivity.ID, nil
}

func (s *Activities) Get(ctx context.Context, orthopedicSurgeryActivityID string) (*commonModel.OrthopedicSurgeryActivity, error) {
	rOrthopedicSurgeryActivity, err := s.OrthopedicSurgeriesRepo.ActivitiesRepo.GetByIDWithRels(ctx, orthopedicSurgeryActivityID)
	if err != nil {
		return nil, err
	}
	activity := common.MapOrthopedicSurgeryActivity(rOrthopedicSurgeryActivity, true, true)
	return activity, nil
}

func (s *Activities) GetActivities(ctx context.Context, queryFilter commonModel.SurgeryLogbookEntryQueryFilter, orderBy *commonModel.QueryOrder) ([]*commonModel.OrthopedicSurgeryActivity, error) {
	user, err := s.UsersRepo.GetByAuthenticationContext(ctx)
	if err != nil {
		return nil, err
	}

	rOrthopedicSurgeryActivity, err := s.OrthopedicSurgeriesRepo.ActivitiesRepo.ListAllSurgeryActivitiesByFilter(ctx, user.ID, queryFilter, orderBy)
	if err != nil {
		return nil, err
	}
	residentOrthopedicActivities := make([]*commonModel.OrthopedicSurgeryActivity, 0)
	for _, rActivity := range rOrthopedicSurgeryActivity {
		activity := common.MapOrthopedicSurgeryActivity(rActivity, true, true)
		residentOrthopedicActivities = append(residentOrthopedicActivities, activity)
	}
	return residentOrthopedicActivities, nil
}

func (s *Activities) DeleteInProgressActivity(ctx context.Context, activityID string) error {
	activity, err := s.OrthopedicSurgeriesRepo.ActivitiesRepo.GetByIDWithRels(ctx, activityID)
	if err != nil {
		return err
	}
	user, err := s.UsersRepo.GetByAuthenticationContext(ctx)
	if err != nil {
		return err
	}
	if activity.ResidentID == user.ID && activity.InProgress {
		err := s.OrthopedicSurgeriesRepo.ActivitiesRepo.DeleteActivity(ctx, activity)
		if err != nil {
			return err
		}
	}
	return nil
}
