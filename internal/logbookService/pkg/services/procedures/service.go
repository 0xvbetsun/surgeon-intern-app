package procedures

import (
	"context"

	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/procedures"
)

type (
	IService interface {
		GetById(ctx context.Context, procedureID string) (*commonModel.Procedure, error)
		ListAll(ctx context.Context) ([]*commonModel.Procedure, error)
	}
	Service struct {
		repo procedures.IRepo
	}
)

func NewService(repo procedures.IRepo) IService {
	return &Service{repo: repo}
}

func (s *Service) GetById(ctx context.Context, procedureID string) (*commonModel.Procedure, error) {
	rProcedure, err := s.repo.GetByID(ctx, procedureID)
	if err != nil {
		return nil, err
	}
	procedure := s.mapGraphQlModel(rProcedure)
	return procedure, nil
}

func (s *Service) ListAll(ctx context.Context) ([]*commonModel.Procedure, error) {
	rProcedures, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	p := make([]*commonModel.Procedure, 0)
	for _, rProcedure := range rProcedures {
		procedure := s.mapGraphQlModel(rProcedure)
		p = append(p, procedure)
	}
	return p, nil
}

func (s Service) mapGraphQlModel(rProcedure *models.Procedure) *commonModel.Procedure {
	procedure := &commonModel.Procedure{
		ProcedureID: rProcedure.ID,
		DisplayName: rProcedure.DisplayName,
	}
	annotations := make([]*commonModel.ProcedureAnnotations, 0)
	if err := rProcedure.Annotations.Unmarshal(&annotations); err != nil {
		// Todo handle err
	} else {
		procedure.Annotations = annotations
	}
	return procedure
}
