//go:generate mockery --name IRepo
package examinationsactivityreview

import (
	"context"
	"database/sql"

	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type (
	IRepo interface {
		Add(ctx context.Context, supervisorExaminationReview *models.ExaminationsActivityReview) (*models.ExaminationsActivityReview, error)
		GetById(ctx context.Context, examinationReviewId string) (*models.ExaminationsActivityReview, error)
		Update(ctx context.Context, examination *models.ExaminationsActivityReview) (*models.ExaminationsActivityReview, error)
	}
	Repo struct {
		db *sql.DB
	}
)

func (r Repo) GetById(ctx context.Context, examinationReviewId string) (*models.ExaminationsActivityReview, error) {
	supervisorExaminationReview, err := models.FindExaminationsActivityReview(ctx, r.db, examinationReviewId)
	if err != nil {
		return nil, err
	}
	return supervisorExaminationReview, nil
}

func (r Repo) Add(ctx context.Context, supervisorExaminationReview *models.ExaminationsActivityReview) (*models.ExaminationsActivityReview, error) {
	err := supervisorExaminationReview.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return nil, err
	}

	return supervisorExaminationReview, nil
}

func (r Repo) Update(ctx context.Context, supervisorExaminationReview *models.ExaminationsActivityReview) (*models.ExaminationsActivityReview, error) {
	_, err := supervisorExaminationReview.Update(ctx, r.db, boil.Infer())
	if err != nil {
		return nil, err
	}

	return supervisorExaminationReview, nil
}

func NewRepo(db *sql.DB) IRepo {
	return &Repo{db: db}
}
