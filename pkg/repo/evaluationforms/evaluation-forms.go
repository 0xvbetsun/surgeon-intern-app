package evaluationforms

import (
	"context"
	"database/sql"

	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type (
	IEvaluationFormsRepo interface {
		ListAllEvaluationForms(ctx context.Context) ([]*models.EvaluationForm, error)
		GetEvaluationFormById(ctx context.Context, evaluationFormId int) (*models.EvaluationForm, error)
		AddEvaluationForm(ctx context.Context, form *models.EvaluationForm) (*models.EvaluationForm, error)
	}
	EvaluationFormsRepo struct {
		db *sql.DB
	}
)

func NewEvaluationFormsRepo(db *sql.DB) IEvaluationFormsRepo {
	return &EvaluationFormsRepo{db: db}
}

func (r EvaluationFormsRepo) GetEvaluationFormById(ctx context.Context, evaluationFormId int) (*models.EvaluationForm, error) {
	form, err := models.EvaluationForms(models.EvaluationFormWhere.ID.EQ(evaluationFormId)).One(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return form, nil
}

func (r EvaluationFormsRepo) ListAllEvaluationForms(ctx context.Context) ([]*models.EvaluationForm, error) {
	forms, err := models.EvaluationForms().All(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return forms, nil
}

func (r EvaluationFormsRepo) AddEvaluationForm(ctx context.Context, form *models.EvaluationForm) (*models.EvaluationForm, error) {
	err := form.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return nil, err
	}
	return form, err
}
