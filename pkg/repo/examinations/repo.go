//go:generate mockery --name IRepo
package examinations

import (
	"context"
	"database/sql"

	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type (
	IRepo interface {
		Add(ctx context.Context, examination *models.Examination) (*models.Examination, error)
		GetByID(ctx context.Context, examinationID string) (*models.Examination, error)
		ListAll(ctx context.Context) ([]*models.Examination, error)
	}
	Repo struct {
		db *sql.DB
	}
)

func NewRepo(db *sql.DB) IRepo {
	return &Repo{db: db}
}

func (r *Repo) Add(ctx context.Context, examination *models.Examination) (*models.Examination, error) {
	err := examination.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return nil, err
	}

	return examination, nil
}

func (r *Repo) GetByID(ctx context.Context, examinationID string) (*models.Examination, error) {
	examination, err := models.FindExamination(ctx, r.db, examinationID)
	if err != nil {
		return nil, err
	}
	return examination, nil
}

func (r *Repo) ListAll(ctx context.Context) ([]*models.Examination, error) {
	activities, err := models.Examinations().All(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return activities, nil
}
