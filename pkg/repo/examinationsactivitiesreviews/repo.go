//go:generate mockery --name IRepo
package examinationsactivitiesreviews

import (
	"context"
	"database/sql"

	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type (
	IRepo interface {
		Add(ctx context.Context, connection *models.ExaminationsActivitiesReview) (*models.ExaminationsActivitiesReview, error)
		GetById(ctx context.Context, connectionId string) (*models.ExaminationsActivitiesReview, error)
		GetByExaminationId(ctx context.Context, examinationId string) (*models.ExaminationsActivitiesReview, error)
		GetByReviewId(ctx context.Context, reviewId int) (*models.ExaminationsActivitiesReview, error)
		GetByResidentUserId(ctx context.Context, residentUserId int) ([]*models.ExaminationsActivitiesReview, error)
		GetBySupervisorUserIdWithJoinedReviews(ctx context.Context, supervisorUserId string, reviewed *bool) ([]*models.ExaminationsActivitiesReview, error)
		Update(ctx context.Context, examination *models.ExaminationsActivitiesReview) (*models.ExaminationsActivitiesReview, error)
		GetByReviewIdWithJoinedReview(ctx context.Context, supervisorReviewId string) (*models.ExaminationsActivitiesReview, error)
	}
	Repo struct {
		db *sql.DB
	}
)

func (r Repo) GetById(ctx context.Context, residentExaminationId string) (*models.ExaminationsActivitiesReview, error) {
	residentExamination, err := models.FindExaminationsActivitiesReview(ctx, r.db, residentExaminationId)
	if err != nil {
		return nil, err
	}
	return residentExamination, nil
}

func (r Repo) GetByExaminationId(ctx context.Context, examinationId string) (*models.ExaminationsActivitiesReview, error) {
	residentExamination, err := models.ExaminationsActivitiesReviews(
		qm.Where("examination_activities_id=?", examinationId),
		qm.Load(models.ExaminationsActivitiesReviewRels.ExaminationsActivityReview)).One(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return residentExamination, nil
}

func (r Repo) GetByReviewId(ctx context.Context, reviewId int) (*models.ExaminationsActivitiesReview, error) {
	residentExamination, err := models.ExaminationsActivitiesReviews(qm.Where("activity_reviewer_user_id=?", reviewId)).One(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return residentExamination, nil
}

func (r Repo) GetByReviewIdWithJoinedReview(ctx context.Context, supervisorReviewId string) (*models.ExaminationsActivitiesReview, error) {
	residentExamination, err := models.ExaminationsActivitiesReviews(qm.Load(models.ExaminationsActivitiesReviewRels.ExaminationsActivityReview), qm.Where("examinations_activity_reviews_id=?", supervisorReviewId)).One(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return residentExamination, nil
}

func (r Repo) GetByResidentUserId(ctx context.Context, residentUserId int) ([]*models.ExaminationsActivitiesReview, error) {
	residentExaminations, err := models.ExaminationsActivitiesReviews(qm.Where("activity_author_user_id=?", residentUserId)).All(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return residentExaminations, nil
}

func (r Repo) GetBySupervisorUserIdWithJoinedReviews(ctx context.Context, supervisorUserId string, reviewed *bool) ([]*models.ExaminationsActivitiesReview, error) {
	if reviewed != nil {
		residentExaminations, err := models.ExaminationsActivitiesReviews(
			qm.Load(models.ExaminationsActivitiesReviewRels.ExaminationsActivityReview),
			qm.Where("activity_reviewer_user_id=?", supervisorUserId),
			qm.Where("is_reviewed=?", *reviewed)).All(ctx, r.db)
		if err != nil {
			return nil, err
		}
		return residentExaminations, nil
	}
	residentExaminations, err := models.ExaminationsActivitiesReviews(qm.Load(models.ExaminationsActivitiesReviewRels.ExaminationsActivityReview), qm.Where("supervisor_user_id=?", supervisorUserId)).All(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return residentExaminations, nil
}

func (r Repo) Add(ctx context.Context, connection *models.ExaminationsActivitiesReview) (*models.ExaminationsActivitiesReview, error) {
	err := connection.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return nil, err
	}

	return connection, nil
}

func (r Repo) Update(ctx context.Context, connection *models.ExaminationsActivitiesReview) (*models.ExaminationsActivitiesReview, error) {
	_, err := connection.Update(ctx, r.db, boil.Infer())
	if err != nil {
		return nil, err
	}

	return connection, nil
}

func NewRepo(db *sql.DB) IRepo {
	return &Repo{db: db}
}
