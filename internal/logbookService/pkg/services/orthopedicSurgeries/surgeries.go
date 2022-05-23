package orthopedicSurgeries

import (
	"context"

	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/dbexecutor"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/orthopedicSurgeries"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/surgeries"
)

type (
	ISurgeries interface {
		Get(ctx context.Context, surgeryId string) (*commonModel.Surgery, error)
		GetByClinic(ctx context.Context, clinicId string) ([]*commonModel.Surgery, error)
		List(ctx context.Context) ([]*commonModel.Surgery, error)
		GetByDiagnose(ctx context.Context, surgeryDiagnoseId string) ([]*commonModel.Surgery, error)
		GetByMethod(ctx context.Context, surgeryMethodId string) ([]*commonModel.Surgery, error)
		GetByMethodAndDiagnose(ctx context.Context, surgeryMethodId string, surgeryDiagnoseId string) (*commonModel.Surgery, error)
	}
	Surgeries struct {
		OrthopedicSurgeriesRepo *orthopedicSurgeries.Repo
		SurgeriesRepo           surgeries.IRepo
		dbexecutor.IDBExecutor
	}
)

func NewSurgeries(orthopedicSurgeriesRepo *orthopedicSurgeries.Repo,
	surgeriesRepo surgeries.IRepo,
	executor dbexecutor.IDBExecutor) ISurgeries {
	return &Surgeries{OrthopedicSurgeriesRepo: orthopedicSurgeriesRepo, SurgeriesRepo: surgeriesRepo, IDBExecutor: executor}
}

func (s *Surgeries) GetByMethodAndDiagnose(ctx context.Context, surgeryMethodId string, surgeryDiagnoseId string) (*commonModel.Surgery, error) {
	rSurgery, err := s.SurgeriesRepo.GetByDiagnoseAndMethod(ctx, surgeryDiagnoseId, surgeryMethodId, true)
	if err != nil {
		return nil, err
	}
	return s.mapDbSurgeryToDTO(rSurgery), nil
}

func (s *Surgeries) Get(ctx context.Context, surgeryId string) (*commonModel.Surgery, error) {
	rSurgery, err := s.SurgeriesRepo.GetByID(ctx, surgeryId)
	if err != nil {
		return nil, err
	}
	return s.mapDbSurgeryToDTO(rSurgery), nil
}

func (s *Surgeries) GetByClinic(ctx context.Context, clinicId string) ([]*commonModel.Surgery, error) {
	rSurgeries, err := s.SurgeriesRepo.ListByClinicWithRels(ctx, clinicId)
	if err != nil {
		return nil, err
	}
	surgeries := make([]*commonModel.Surgery, 0)
	for _, rSurgery := range rSurgeries {
		surgeries = append(surgeries, s.mapDbSurgeryToDTO(rSurgery))
	}
	return surgeries, nil
}

func (s *Surgeries) List(ctx context.Context) ([]*commonModel.Surgery, error) {
	rSurgeries, err := s.SurgeriesRepo.ListWithRels(ctx)
	if err != nil {
		return nil, err
	}
	surgeries := make([]*commonModel.Surgery, 0)
	for _, rSurgery := range rSurgeries {
		surgeries = append(surgeries, s.mapDbSurgeryToDTO(rSurgery))
	}
	return surgeries, nil
}

func (s *Surgeries) GetByDiagnose(ctx context.Context, surgeryDiagnoseId string) ([]*commonModel.Surgery, error) {
	rSurgeries, err := s.SurgeriesRepo.ListByDiagnoseWithRels(ctx, surgeryDiagnoseId)
	if err != nil {
		return nil, err
	}
	surgeries := make([]*commonModel.Surgery, 0)
	for _, rSurgery := range rSurgeries {
		surgeries = append(surgeries, s.mapDbSurgeryToDTO(rSurgery))
	}
	return surgeries, nil
}

func (s *Surgeries) GetByMethod(ctx context.Context, surgeryMethodId string) ([]*commonModel.Surgery, error) {
	rSurgeries, err := s.SurgeriesRepo.ListByMethodWithRels(ctx, surgeryMethodId)
	if err != nil {
		return nil, err
	}
	surgeries := make([]*commonModel.Surgery, 0)
	for _, rSurgery := range rSurgeries {
		surgeries = append(surgeries, s.mapDbSurgeryToDTO(rSurgery))
	}
	return surgeries, nil
}

func (s *Surgeries) mapDbSurgeryToDTO(dbSurgery *models.Surgery) *commonModel.Surgery {
	surgery := &commonModel.Surgery{
		ID: dbSurgery.ID,
		Diagnose: &commonModel.SurgeryDiagnose{
			ID:           dbSurgery.R.Diagnose.ID,
			Bodypart:     dbSurgery.R.Diagnose.Bodypart,
			DiagnoseName: dbSurgery.R.Diagnose.DiagnoseName,
			DiagnoseCode: dbSurgery.R.Diagnose.DiagnoseCode,
			ExtraCode:    dbSurgery.R.Diagnose.ExtraCode,
		},
		Method: &commonModel.SurgeryMethod{
			ID:           dbSurgery.R.Method.ID,
			MethodName:   dbSurgery.R.Method.MethodName,
			MethodCode:   dbSurgery.R.Method.MethodCode,
			ApproachName: dbSurgery.R.Method.ApproachName,
		},
	}
	return surgery
}
