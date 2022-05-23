package evaluationforms

import (
	"context"

	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/dbexecutor"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/evaluationforms"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/orthopedicSurgeries"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/users"
)

type (
	IEvaluationForms interface {
		ListAllEvaluationForms(ctx context.Context) ([]*commonModel.EvaluationForm, error)
	}
	EvaluationForms struct {
		evaluationFormsRepo evaluationforms.Repo
		usersRepo           users.IRepo
		activitiesRepo      orthopedicSurgeries.IActivitiesRepo
		dbexecutor.IDBExecutor
	}
)

func NewEvaluationForms(evaluationFormsRepo evaluationforms.Repo, usersRepo users.IRepo, activitiesRepo orthopedicSurgeries.IActivitiesRepo, executor dbexecutor.IDBExecutor) IEvaluationForms {
	return &EvaluationForms{
		evaluationFormsRepo: evaluationFormsRepo,
		activitiesRepo:      activitiesRepo,
		usersRepo:           usersRepo,
		IDBExecutor:         executor,
	}
}

func (s EvaluationForms) ListAllEvaluationForms(ctx context.Context) ([]*commonModel.EvaluationForm, error) {
	rForms, err := s.evaluationFormsRepo.EvaluationForms.ListAllEvaluationForms(ctx)
	if err != nil {
		return nil, err
	}

	forms := make([]*commonModel.EvaluationForm, 0)
	for _, rForm := range rForms {
		form := s.mapEvaluationFormGraphQlModel(rForm)
		forms = append(forms, form)
	}
	return forms, nil
}

func (s EvaluationForms) mapEvaluationFormGraphQlModel(rForm *models.EvaluationForm) *commonModel.EvaluationForm {
	form := &commonModel.EvaluationForm{
		EvaluationFormID: rForm.ID,
		DisplayName:      rForm.Name,
		Difficulty:       rForm.Difficulty,
		Citations:        rForm.Citations,
	}
	annotations := make([]*commonModel.EvaluationFormAnnotation, 0)
	if err := rForm.Annotations.Unmarshal(&annotations); err != nil {
		// Todo handle err
	} else {
		form.Annotations = annotations
	}
	return form
}
