package orthopedicSurgeries

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/friendsofgo/errors"
	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/activities"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/assessments"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/dbexecutor"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/orthopedicSurgeries"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/surgeries"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/users"
	"github.com/volatiletech/null/v8"
)

type (
	IReviews interface {
		GetByReviewerUserID(ctx context.Context, reviewerUserID string) ([]*commonModel.OrthopedicSurgeryActivityReview, error)
		Get(ctx context.Context, reviewId string) (*commonModel.OrthopedicSurgeryActivityReview, error)
		SubmitOrthopedicSurgeryReview(ctx context.Context, input commonModel.OrthopedicSurgeryActivityReviewInput) (*string, error)

		// TODO: GET/DELETE, Reviews can already be fetched as part of activity.
	}
	Reviews struct {
		ActivitiesRepo          activities.IActivitiesRepo
		AssessmentsRepo         assessments.IAssessmentsRepo
		SurgeriesRepo           surgeries.IRepo
		OrthopedicSurgeriesRepo *orthopedicSurgeries.Repo
		UsersRepo               users.IRepo
		dbexecutor.IDBExecutor
	}
)

// TODO: make sure dops gets joined in correctly
func NewReviews(activitiesRepo activities.IActivitiesRepo,
	assessmentsRepo assessments.IAssessmentsRepo,
	surgeriesRepo surgeries.IRepo,
	orthopedicSurgeriesRepo *orthopedicSurgeries.Repo,
	usersRepo users.IRepo,
	executor dbexecutor.IDBExecutor) IReviews {
	return &Reviews{ActivitiesRepo: activitiesRepo, AssessmentsRepo: assessmentsRepo, SurgeriesRepo: surgeriesRepo, OrthopedicSurgeriesRepo: orthopedicSurgeriesRepo, UsersRepo: usersRepo, IDBExecutor: executor}
}

func (r *Reviews) Get(ctx context.Context, reviewId string) (*commonModel.OrthopedicSurgeryActivityReview, error) {
	rReview, err := r.OrthopedicSurgeriesRepo.ReviewsRepo.GetReviewByIDWithRels(ctx, reviewId)
	if err != nil {
		return nil, err
	}
	qlReview, err := r.mapReview(rReview)
	if err != nil {
		return nil, err
	}
	return qlReview, nil
}

func (r *Reviews) GetByReviewerUserID(ctx context.Context, reviewerUserID string) ([]*commonModel.OrthopedicSurgeryActivityReview, error) {
	rReviews, err := r.OrthopedicSurgeriesRepo.ReviewsRepo.GetReviewsBySupervisorUserIDWithRels(ctx, reviewerUserID)
	if err != nil {
		return nil, err
	}
	var reviews = make([]*commonModel.OrthopedicSurgeryActivityReview, 0)
	for _, dbReview := range rReviews {
		review, err := r.mapReview(dbReview)
		if err != nil {
			continue
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}

func (r *Reviews) SubmitOrthopedicSurgeryReview(ctx context.Context, input commonModel.OrthopedicSurgeryActivityReviewInput) (*string, error) {
	// check if surgeryActivity exists
	surgeryActivity, err := r.OrthopedicSurgeriesRepo.ActivitiesRepo.GetByID(ctx, input.ActivityID)
	if err != nil {
		return nil, err
	}

	// Find supervisor user
	supervisorUser, err := r.UsersRepo.GetByAuthenticationContext(ctx)
	if err != nil {
		return nil, err
	}

	var review *models.OrthopedicSurgeriesActivityReview
	var activity *models.Activity
	var assessment *models.Assessment

	if input.ReviewID != nil { // Get existing review
		review, err = r.OrthopedicSurgeriesRepo.ReviewsRepo.GetReviewByIDWithRels(ctx, *input.ReviewID)
		if err != nil {
			return nil, err
		}
		if review == nil {
			return nil, errors.New(fmt.Sprintf("No existing review found with id %s", *input.ReviewID))
		}
		if review.R.OrthopedicSurgeryReviewAssessment == nil || review.R.OrthopedicSurgeryReviewAssessment.R.Activity == nil {
			return nil, errors.New(fmt.Sprintf("No existing activity or assessment found for review with id %s", *input.ReviewID))
		}
		assessment = review.R.OrthopedicSurgeryReviewAssessment
		activity = assessment.R.Activity
	} else {
		review = &models.OrthopedicSurgeriesActivityReview{
			OrthopedicSurgeryActivityID: input.ActivityID,
			CreatedAt:                   time.Now(),
			UpdatedAt:                   null.TimeFrom(time.Now()),
			SupervisorID:                supervisorUser.ID,
		}
		assessment = &models.Assessment{}
		activity = &models.Activity{}
	}

	review.ResidentID = surgeryActivity.ResidentID

	review.InProgress = input.InProgress

	review.ActiveStep = input.ActiveStep
	review.CompletedStep = input.CompletedStep
	review.OrthopedicSurgeryActivityID = surgeryActivity.ID

	if input.ShouldSign {
		review.SignedAt = null.TimeFrom(time.Now())
	}

	// Map surgery metadata
	review.CaseNotes = input.SurgeryMetadata.CaseNotes
	review.PatientAge = input.SurgeryMetadata.PatientAge
	review.PatientGender = input.SurgeryMetadata.PatientGender
	review.OccurredAt = input.SurgeryMetadata.OccurredAt
	activity.OccurredAt = input.SurgeryMetadata.OccurredAt

	// TODO: Sanity check that users exist
	review.OperatorID = input.OperatorID
	review.AssistantID = input.AssistantID

	review.Complications = input.Complications
	review.Comments = input.Comments
	review.ReviewComment = input.ReviewComment

	review.Requested = surgeryActivity.ReviewRequested

	if err := review.Annotations.Marshal(&input.Annotations); err != nil {
		return nil, err
	}

	err = r.RunWithTX(ctx, func(tx *sql.Tx) error {
		// Update or add new
		if input.ReviewID != nil {
			review, err = r.OrthopedicSurgeriesRepo.ReviewsRepo.UpdateOrthopedicSurgeryActivityReview(ctx, tx, review)
			if err != nil {
				return err
			}
			activity, err = r.ActivitiesRepo.UpdateActivity(ctx, tx, activity)
			if err != nil {
				return err
			}
		} else {
			review, err = r.OrthopedicSurgeriesRepo.ReviewsRepo.AddOrthopedicSurgeryActivityReview(ctx, tx, review)
			if err != nil {
				return err
			}
			assessment.OrthopedicSurgeryReviewID = null.StringFrom(review.ID)
			assessment, err = r.AssessmentsRepo.AddAssessment(ctx, tx, assessment)
			if err != nil {
				return err
			}
			activity.AssessmentID = null.StringFrom(assessment.ID)
			_, err = r.ActivitiesRepo.AddActivity(ctx, tx, activity)
			if err != nil {
				return err
			}
		}

		surgeries := make([]*models.Surgery, 0)
		for _, surgeryId := range input.SurgeryMetadata.SurgeryIds {
			rSurgery, err := r.SurgeriesRepo.GetByIDWithRels(ctx, surgeryId)
			if err != nil {
				return err
			}
			surgeries = append(surgeries, rSurgery)
		}

		if err := r.OrthopedicSurgeriesRepo.ReviewsRepo.UpdateSurgeriesForReview(ctx, tx, review.ID, surgeries); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &review.ID, nil
}

func (r *Reviews) mapReview(dbReview *models.OrthopedicSurgeriesActivityReview) (*commonModel.OrthopedicSurgeryActivityReview, error) {

	review := &commonModel.OrthopedicSurgeryActivityReview{
		ReviewID:   dbReview.ID,
		ActivityID: dbReview.OrthopedicSurgeryActivityID,
		CreatedAt:  dbReview.CreatedAt,
		UpdatedAt:  dbReview.UpdatedAt.Time,
		SignedAt:   dbReview.SignedAt.Ptr(),
		SurgeryMetadata: &commonModel.OrthopedicSurgeryMetadata{
			OccurredAt:    dbReview.OccurredAt,
			CaseNotes:     dbReview.CaseNotes,
			PatientAge:    dbReview.PatientAge,
			PatientGender: dbReview.PatientGender,
		},
		OperatorID:       dbReview.OperatorID,
		AssistantID:      dbReview.AssistantID,
		Comments:         dbReview.Comments,
		Complications:    dbReview.Complications,
		Annotations:      &commonModel.OrthopedicSurgeryActivityAnnotations{},
		ResidentUserID:   dbReview.ResidentID,
		SupervisorUserID: dbReview.SupervisorID,
		ReviewComment:    dbReview.ReviewComment,
		InProgress:       dbReview.InProgress,
		ActiveStep:       dbReview.ActiveStep,
		CompletedStep:    dbReview.CompletedStep,
	}

	err := dbReview.Annotations.Unmarshal(review.Annotations)
	if err != nil {
		return nil, err
	}

	if dbReview.R != nil && dbReview.R.OrthopedicSurgeriesActivityReviewSurgeries != nil {
		// TODO: Map surgeries
		surgeries := make([]*commonModel.Surgery, 0)
		for _, surgery := range dbReview.R.OrthopedicSurgeriesActivityReviewSurgeries {
			mappedSurgery := &commonModel.Surgery{
				ID: surgery.SurgeryID,
				Diagnose: &commonModel.SurgeryDiagnose{
					ID:           surgery.R.Surgery.R.Diagnose.ID,
					Bodypart:     surgery.R.Surgery.R.Diagnose.Bodypart,
					DiagnoseName: surgery.R.Surgery.R.Diagnose.DiagnoseName,
					DiagnoseCode: surgery.R.Surgery.R.Diagnose.DiagnoseCode,
					ExtraCode:    surgery.R.Surgery.R.Diagnose.ExtraCode,
				},
				Method: &commonModel.SurgeryMethod{
					ID:           surgery.R.Surgery.R.Method.ID,
					MethodName:   surgery.R.Surgery.R.Method.MethodName,
					MethodCode:   surgery.R.Surgery.R.Method.MethodCode,
					ApproachName: surgery.R.Surgery.R.Method.ApproachName,
				},
			}
			surgeries = append(surgeries, mappedSurgery)
		}
		review.SurgeryMetadata.Surgeries = surgeries
	}

	if dbReview.R != nil && dbReview.R.Supervisor != nil && dbReview.R.Resident != nil {
		review.Resident = &commonModel.User{
			UserID:      dbReview.R.Resident.ID,
			DisplayName: dbReview.R.Resident.DisplayName,
		}
		review.Supervisor = &commonModel.User{
			UserID:      dbReview.R.Supervisor.ID,
			DisplayName: dbReview.R.Supervisor.DisplayName,
		}
	}

	return review, nil
}

func (r *Reviews) updateDbReview(existingReview *models.OrthopedicSurgeriesActivityReview,
	input *commonModel.OrthopedicSurgeryActivityReviewInput) (*models.OrthopedicSurgeriesActivityReview, error) {
	// Map surgery metadata
	existingReview.CaseNotes = input.SurgeryMetadata.CaseNotes
	existingReview.PatientAge = input.SurgeryMetadata.PatientAge
	existingReview.PatientGender = input.SurgeryMetadata.PatientGender
	existingReview.OccurredAt = input.SurgeryMetadata.OccurredAt

	// TODO: Sanity check that users exist
	existingReview.OperatorID = input.OperatorID
	existingReview.AssistantID = input.AssistantID

	existingReview.Complications = input.Complications
	existingReview.Comments = input.Comments
	existingReview.ReviewComment = input.ReviewComment

	existingReview.UpdatedAt = null.TimeFrom(time.Now())

	if input.ShouldSign {
		existingReview.SignedAt = null.TimeFrom(time.Now())
	}

	if err := existingReview.Annotations.Marshal(&input.Annotations); err != nil {
		return nil, err
	}

	return existingReview, nil
}
