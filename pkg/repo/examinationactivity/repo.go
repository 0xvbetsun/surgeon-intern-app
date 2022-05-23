//go:generate mockery --name IRepo
package examinationactivity

import (
	"context"
	"database/sql"

	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type (
	IRepo interface {
		Add(ctx context.Context, residentExamination *models.ExaminationActivity) (*models.ExaminationActivity, error)
		GetByResidentUserId(ctx context.Context, residentUserId string) ([]*models.ExaminationActivity, error)
		GetBySupervisorUserId(ctx context.Context, supervisorUserId string) ([]*models.ExaminationActivity, error)
		GetById(ctx context.Context, residentExaminationId string) (*models.ExaminationActivity, error)
		Update(ctx context.Context, examination *models.ExaminationActivity) (*models.ExaminationActivity, error)
	}
	Repo struct {
		db *sql.DB
	}
)

func (r Repo) GetByResidentUserId(ctx context.Context, residentUserId string) ([]*models.ExaminationActivity, error) {
	residentExaminations, err := models.ExaminationActivities(qm.Where("resident_user_id=?", residentUserId)).All(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return residentExaminations, nil
}
func (r Repo) GetById(ctx context.Context, residentExaminationId string) (*models.ExaminationActivity, error) {
	residentExamination, err := models.FindExaminationActivity(ctx, r.db, residentExaminationId)
	if err != nil {
		return nil, err
	}
	return residentExamination, nil
}

func (r Repo) GetBySupervisorUserId(ctx context.Context, supervisorUserId string) ([]*models.ExaminationActivity, error) {
	residentExaminations, err := models.ExaminationActivities(qm.Where("supervisor_user_id=?", supervisorUserId)).All(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return residentExaminations, nil
}

func (r Repo) Add(ctx context.Context, residentExamination *models.ExaminationActivity) (*models.ExaminationActivity, error) {
	err := residentExamination.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return nil, err
	}

	return residentExamination, nil
}

func (r Repo) Update(ctx context.Context, residentExamination *models.ExaminationActivity) (*models.ExaminationActivity, error) {
	_, err := residentExamination.Update(ctx, r.db, boil.Infer())
	if err != nil {
		return nil, err
	}

	return residentExamination, nil
}

func NewRepo(db *sql.DB) IRepo {
	return &Repo{db: db}
}
