package examinationsactivitiesreviews

import (
	"context"
	"errors"
	"time"

	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	gin2 "github.com/vbetsun/surgeon-intern-app/internal/pkg/gin"
	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/examinationsactivitiesreviews"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/examinationsactivityreview"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/users"
	"github.com/volatiletech/null/v8"
)

type (
	IService interface {
		GetByExaminationId(ctx context.Context, examinationId string) (*commonModel.SupervisorExaminationReview, error)
		GetBySupervisorUser(ctx context.Context, supervisorUserId string, reviewed *bool) ([]*commonModel.SupervisorExaminationReview, error)
		SubmitReview(ctx context.Context, reviewedExamination commonModel.SupervisorExaminationReviewInput) (*commonModel.SupervisorExaminationReview, error)
	}
	Service struct {
		supervisorReviewsRepo examinationsactivityreview.IRepo
		connectionsRepo       examinationsactivitiesreviews.IRepo
		usersRepo             users.IRepo
	}
)

func NewService(
	supervisorReviewsRepo examinationsactivityreview.IRepo, connectionsRepo examinationsactivitiesreviews.IRepo, usersRepo users.IRepo) IService {
	return &Service{
		supervisorReviewsRepo: supervisorReviewsRepo,
		connectionsRepo:       connectionsRepo,
		usersRepo:             usersRepo,
	}
}

func (s Service) GetByExaminationId(ctx context.Context, examinationId string) (*commonModel.SupervisorExaminationReview, error) {
	rReview, err := s.connectionsRepo.GetByExaminationId(ctx, examinationId)
	if err != nil {
		return nil, err
	}
	return s.mapGraphQlModel(rReview), nil
}

func (s Service) SubmitReview(ctx context.Context, reviewedExamination commonModel.SupervisorExaminationReviewInput) (*commonModel.SupervisorExaminationReview, error) {
	ginContext, err := gin2.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}
	uid, err := gin2.FirebaseIdFromGinContext(ginContext)
	if err != nil {
		return nil, err
	}

	supervisorUser, err := s.usersRepo.GetByExternalID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if supervisorUser == nil {
		return nil, errors.New("no supervisorUser found with provided firebase ID")
	}

	connection, err := s.connectionsRepo.GetByReviewIdWithJoinedReview(ctx, reviewedExamination.SupervisorExaminationReviewID)
	if err != nil {
		return nil, err
	}

	if connection.ActivityReviewerUserID != supervisorUser.ID {
		return nil, errors.New("Reviewing user doesn't match supervisor on examination.")
	}

	supervisorExaminationReview := connection.R.ExaminationsActivityReview
	err = supervisorExaminationReview.Annotations.Marshal(reviewedExamination.Annotations) // TODO: Compare and validate that annotations match source examination.
	if reviewedExamination.Comment != nil {
		supervisorExaminationReview.Comment = null.StringFromPtr(reviewedExamination.Comment)
	}
	if err != nil {
		return nil, err
	}

	_, err = s.supervisorReviewsRepo.Update(ctx, supervisorExaminationReview)
	if err != nil {
		return nil, err
	}

	connection.IsReviewed = true
	connection.SupervisorUpdatedAt = null.TimeFrom(time.Now())
	_, err = s.connectionsRepo.Update(ctx, connection)
	if err != nil {
		return nil, err
	}

	resultReview := s.mapGraphQlModel(connection)
	return resultReview, nil
}

func (s Service) GetBySupervisorUser(ctx context.Context, supervisorUserId string, reviewed *bool) ([]*commonModel.SupervisorExaminationReview, error) {
	rReviews, err := s.connectionsRepo.GetBySupervisorUserIdWithJoinedReviews(ctx, supervisorUserId, reviewed)
	if err != nil {
		return nil, err
	}
	review := make([]*commonModel.SupervisorExaminationReview, 0)
	for _, rResidentExamination := range rReviews {
		residentExamination := s.mapGraphQlModel(rResidentExamination)
		review = append(review, residentExamination)
	}
	return review, nil
}

func (s Service) mapGraphQlModel(connection *models.ExaminationsActivitiesReview) *commonModel.SupervisorExaminationReview {
	supervisorExaminationReview := connection.R.ExaminationsActivityReview
	resultReview := &commonModel.SupervisorExaminationReview{
		SupervisorExaminationReviewID: supervisorExaminationReview.ID,
		ResidentExaminationID:         connection.ExaminationActivitiesID,
		ResidentUserID:                connection.ActivityAuthorUserID,
		SupervisorUserID:              connection.ActivityReviewerUserID,
		DisplayName:                   supervisorExaminationReview.DisplayName,
		Annotations:                   nil,
		CreatedAt:                     connection.CreatedAt,
		ResidentUpdatedAt:             &connection.ResidentUpdatedAt.Time,
		SupervisorUpdatedAt:           &connection.SupervisorUpdatedAt.Time,
		IsReviewed:                    connection.IsReviewed,
	}
	if supervisorExaminationReview.Comment.Valid {
		resultReview.Comment = &supervisorExaminationReview.Comment.String
	}
	annotations := make([]*commonModel.SupervisorExaminationReviewAnnotations, 0)
	if err := supervisorExaminationReview.Annotations.Unmarshal(&annotations); err != nil {
		// Todo handle err
	} else {
		resultReview.Annotations = annotations
	}
	return resultReview
}
