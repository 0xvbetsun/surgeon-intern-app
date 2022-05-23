package examinationsactivities

import (
	"context"

	"github.com/friendsofgo/errors"
	"github.com/gin-gonic/gin"
	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/examinationactivity"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/examinations"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/examinationsactivitiesreviews"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/examinationsactivityreview"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/users"
	"github.com/volatiletech/null/v8"
)

const GinContextUserIDKey = "UID"
const GinContextKey = "GinContextKey"

type (
	IService interface {
		Add(ctx context.Context, input commonModel.ResidentExaminationInput) (*commonModel.ResidentExamination, error)
		GetByUser(ctx context.Context, residentUserId string) ([]*commonModel.ResidentExamination, error)
	}
	Service struct {
		examinationsRepo         examinations.IRepo
		usersRepo                users.IRepo
		residentExaminationsRepo examinationactivity.IRepo
		supervisorReviewsRepo    examinationsactivityreview.IRepo
		connectionsRepo          examinationsactivitiesreviews.IRepo
	}
)

func NewService(
	examinationsRepo examinations.IRepo,
	usersRepo users.IRepo,
	residentExaminationsRepo examinationactivity.IRepo,
	supervisorReviewsRepo examinationsactivityreview.IRepo,
	connectionsRepo examinationsactivitiesreviews.IRepo) *Service {
	return &Service{
		residentExaminationsRepo: residentExaminationsRepo,
		examinationsRepo:         examinationsRepo,
		usersRepo:                usersRepo,
		supervisorReviewsRepo:    supervisorReviewsRepo,
		connectionsRepo:          connectionsRepo,
	}
}

func (s Service) Add(ctx context.Context, examination commonModel.ResidentExaminationInput) (*commonModel.ResidentExamination, error) {
	residentUser, err := s.usersRepo.GetByAuthenticationContext(ctx)
	if err != nil {
		return nil, err
	}

	originalExamination, err := s.examinationsRepo.GetByID(ctx, examination.ExaminationID)
	if err != nil {
		return nil, err
	}

	var supervisorId null.String
	if examination.SupervisorUserID != nil {
		supervisorId = null.NewString(*examination.SupervisorUserID, true)
	}
	residentExamination := models.ExaminationActivity{
		ResidentUserID:   residentUser.ID,
		SupervisorUserID: supervisorId, // TODO: Remove when new queries and mutations are done.
		ExaminationID:    originalExamination.ID,
		DisplayName:      originalExamination.DisplayName,
	}
	err = residentExamination.Annotations.Marshal(examination.Annotations) // TODO: Compare and validate that annotations looks like root examination.
	if err != nil {
		return nil, err
	}

	_, err = s.residentExaminationsRepo.Add(ctx, &residentExamination)
	if err != nil {
		return nil, err
	}

	// Supervisor provided. We add a resident - supervisor relation for this examination.
	if examination.SupervisorUserID != nil {
		supervisorExaminationReview := models.ExaminationsActivityReview{
			SupervisorUserID: *examination.SupervisorUserID,
			DisplayName:      originalExamination.DisplayName,
		}
		err = supervisorExaminationReview.Annotations.Marshal(examination.Annotations)
		if err != nil {
			return nil, err
		}

		_, err = s.supervisorReviewsRepo.Add(ctx, &supervisorExaminationReview)
		if err != nil {
			return nil, err
		}

		connection := models.ExaminationsActivitiesReview{
			ActivityAuthorUserID:          residentUser.ID,
			ActivityReviewerUserID:        *examination.SupervisorUserID,
			ExaminationActivitiesID:       residentExamination.ID,
			ExaminationsActivityReviewsID: null.StringFrom(supervisorExaminationReview.ID),
			ResidentUpdatedAt:             null.Time{},
			SupervisorUpdatedAt:           null.Time{},
			IsReviewed:                    false,
		}
		_, err = s.connectionsRepo.Add(ctx, &connection)
		if err != nil {
			return nil, err
		}
	}

	resultUserActivity := s.mapGraphQlModel(&residentExamination)
	return resultUserActivity, nil
}

func (s Service) Review(ctx context.Context, reviewedExamination *commonModel.SupervisorExaminationReviewInput) (*commonModel.ResidentExamination, error) {
	ginContext := ctx.Value(GinContextKey).(*gin.Context)
	if ginContext == nil {
		return nil, errors.New("Gin context is missing.")
	}
	uid := ginContext.Value(GinContextUserIDKey)
	if uid == nil {
		return nil, errors.New("No supervisorUser UID found.")
	}

	supervisorUser, err := s.usersRepo.GetByExternalID(ctx, uid.(string))
	if err != nil {
		return nil, err
	}
	if supervisorUser == nil {
		return nil, errors.New("No supervisorUser found with provided external ID.")
	}

	residentExamination, err := s.residentExaminationsRepo.GetById(ctx, reviewedExamination.SupervisorExaminationReviewID)
	if err != nil {
		return nil, err
	}
	if residentExamination == nil {
		return nil, errors.New("Couldn't find original resident examination.")
	}

	if residentExamination.SupervisorUserID.Valid && residentExamination.SupervisorUserID.String != supervisorUser.ID {
		return nil, errors.New("Reviewing user doesn't match supervisor on examination.")
	}

	err = residentExamination.Annotations.Marshal(reviewedExamination.Annotations) // TODO: Compare and validate that annotations match source examination.
	if err != nil {
		return nil, err
	}

	_, err = s.residentExaminationsRepo.Update(ctx, residentExamination)
	if err != nil {
		return nil, err
	}

	resultUserActivity := s.mapGraphQlModel(residentExamination)
	return resultUserActivity, nil
}

func (s Service) GetById(ctx context.Context, residentExaminationID string) (*commonModel.ResidentExamination, error) {
	rResidentExamination, err := s.residentExaminationsRepo.GetById(ctx, residentExaminationID)
	if err != nil {
		return nil, err
	}
	residentExamination := s.mapGraphQlModel(rResidentExamination)
	return residentExamination, nil
}

func (s Service) GetByResidentUser(ctx context.Context, residentUserId string) ([]*commonModel.ResidentExamination, error) {
	rResidentExaminations, err := s.residentExaminationsRepo.GetByResidentUserId(ctx, residentUserId)
	if err != nil {
		return nil, err
	}
	residentExaminations := make([]*commonModel.ResidentExamination, 0)
	for _, rResidentExamination := range rResidentExaminations {
		residentExamination := s.mapGraphQlModel(rResidentExamination)
		residentExaminations = append(residentExaminations, residentExamination)
	}
	return residentExaminations, nil
}

func (s Service) GetBySupervisorUser(ctx context.Context, supervisorUserId string) ([]*commonModel.ResidentExamination, error) {
	rResidentExaminations, err := s.residentExaminationsRepo.GetBySupervisorUserId(ctx, supervisorUserId)
	if err != nil {
		return nil, err
	}
	residentExaminations := make([]*commonModel.ResidentExamination, 0)
	for _, rResidentExamination := range rResidentExaminations {
		residentExamination := s.mapGraphQlModel(rResidentExamination)
		residentExaminations = append(residentExaminations, residentExamination)
	}
	return residentExaminations, nil
}

func (s Service) mapGraphQlModel(residentExamination *models.ExaminationActivity) *commonModel.ResidentExamination {
	resultUserActivity := &commonModel.ResidentExamination{
		ResidentExaminationID: residentExamination.ID,
		DisplayName:           residentExamination.DisplayName,
		CreatedAt:             residentExamination.CreatedAt,
		ResidentUserID:        residentExamination.ResidentUserID,
		SupervisorUserID:      residentExamination.SupervisorUserID.String,
		ExaminationID:         residentExamination.ExaminationID,
	}
	annotations := make([]*commonModel.ResidentExaminationAnnotations, 0)
	if err := residentExamination.Annotations.Unmarshal(&annotations); err != nil {
		// Todo handle err
	} else {
		resultUserActivity.Annotations = annotations
	}
	return resultUserActivity
}
