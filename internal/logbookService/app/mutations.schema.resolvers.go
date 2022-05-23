package logbookService

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/internal/logbookService/graph/generated"
)

func (r *mutationResolver) SubmitResidentExamination(ctx context.Context, examination commonModel.ResidentExaminationInput) (*commonModel.ResidentExamination, error) {
	resultResidentExamination, err := r.service.ResidentExaminationsService.Add(ctx, examination)
	if err != nil {
		return nil, err
	}
	return resultResidentExamination, nil
}

func (r *mutationResolver) SubmitSupervisorExaminationReview(ctx context.Context, reviewedExamination commonModel.SupervisorExaminationReviewInput) (*commonModel.SupervisorExaminationReview, error) {
	resultReview, err := r.service.SupervisorExaminationReviewsService.SubmitReview(ctx, reviewedExamination)
	if err != nil {
		return nil, err
	}
	return resultReview, nil
}

func (r *mutationResolver) SubmitOrthopedicSurgeryActivity(ctx context.Context, activityInput commonModel.OrthopedicSurgeryActivityInput) (*string, error) {
	activityId, err := r.service.OrthopedicSurgeryService.ActivityService.SubmitOrthopedicSurgeryActivity(ctx, &activityInput)
	if err != nil {
		return nil, err
	}
	return activityId, err
}

func (r *mutationResolver) SubmitOrthopedicSurgeryActivityReview(ctx context.Context, reviewInput commonModel.OrthopedicSurgeryActivityReviewInput) (*string, error) {
	reviewId, err := r.service.OrthopedicSurgeryService.ReviewService.SubmitOrthopedicSurgeryReview(ctx, reviewInput)
	if err != nil {
		return nil, err
	}
	return reviewId, err
}

func (r *mutationResolver) SubmitDopsEvaluation(ctx context.Context, evaluationInput commonModel.DopsEvaluationInput) (*string, error) {
	evaluationId, err := r.service.EvaluationFormService.Dops.SubmitDopsEvaluation(ctx, evaluationInput)
	if err != nil {
		return nil, err
	}
	return evaluationId, nil
}

func (r *mutationResolver) DeleteInProgressOrthopedicSurgeryActivity(ctx context.Context, activityID string) (*bool, error) {
	err := r.service.OrthopedicSurgeryService.ActivityService.DeleteInProgressActivity(ctx, activityID)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *mutationResolver) DeleteInProgressDopsEvaluation(ctx context.Context, dopsEvaluationID string) (*bool, error) {
	err := r.service.EvaluationFormService.Dops.DeleteInProgressEvaluation(ctx, dopsEvaluationID)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *mutationResolver) DeleteInProgressMiniCexEvaluation(ctx context.Context, miniCexEvaluationID string) (*bool, error) {
	err := r.service.EvaluationFormService.MiniCex.DeleteInProgressEvaluation(ctx, miniCexEvaluationID)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *mutationResolver) ConnectActivityToDopsEvaluation(ctx context.Context, activityID string, dopsEvaluationID string) (*commonModel.OrthopedicSurgeryActivity, error) {
	activity, err := r.service.EvaluationFormService.Dops.ConnectActivityToDopsEvaluation(ctx, activityID, dopsEvaluationID)
	if err != nil {
		return nil, err
	}
	return activity, nil
}

func (r *mutationResolver) SubmitMiniCexEvaluation(ctx context.Context, evaluationInput commonModel.MiniCexEvaluationInput) (*string, error) {
	evaluationId, err := r.service.EvaluationFormService.MiniCex.SubmitMiniCexEvaluation(ctx, evaluationInput)
	if err != nil {
		return nil, err
	}
	return evaluationId, nil
}

func (r *mutationResolver) RequestMiniCexEvaluation(ctx context.Context, requestInput commonModel.MiniCexRequestInput) (*string, error) {
	evaluationId, err := r.service.EvaluationFormService.MiniCex.RequestMiniCexEvaluation(ctx, requestInput)
	if err != nil {
		return nil, err
	}
	return evaluationId, nil
}

func (r *mutationResolver) SetHasSeenNotifications(ctx context.Context, seenNotifications []*commonModel.SeenNotificationInput) ([]*commonModel.Notification, error) {
	notifications, err := r.service.NotificationService.SetHasSeenNotifications(ctx, seenNotifications)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
