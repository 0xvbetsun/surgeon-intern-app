package examinations

import (
	"context"

	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/examinations"
)

type (
	IService interface {
		GetById(ctx context.Context, activityID string) (*commonModel.Examination, error)
		ListAll(ctx context.Context) ([]*commonModel.Examination, error)
	}
	Service struct {
		repo examinations.IRepo
	}
)

func NewService(repo examinations.IRepo) IService {
	return &Service{repo: repo}
}

func (s *Service) GetById(ctx context.Context, activityID string) (*commonModel.Examination, error) {
	rExamination, err := s.repo.GetByID(ctx, activityID)
	if err != nil {
		return nil, err
	}
	examination := s.mapGraphQlModel(rExamination)
	return examination, nil
}

func (s *Service) ListAll(ctx context.Context) ([]*commonModel.Examination, error) {
	rExaminations, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	examinations := make([]*commonModel.Examination, 0)
	for _, rExamination := range rExaminations {
		examination := s.mapGraphQlModel(rExamination)
		examinations = append(examinations, examination)
	}
	return examinations, nil
}

func (s Service) mapGraphQlModel(rExamination *models.Examination) *commonModel.Examination {
	examination := &commonModel.Examination{
		ExaminationID: rExamination.ID,
		ClinicID:      rExamination.DepartmentID,
		DisplayName:   rExamination.DisplayName,
	}
	annotations := make([]*commonModel.ExaminationAnnotations, 0)
	if err := rExamination.Annotations.Unmarshal(&annotations); err != nil {
		// Todo handle err
	} else {
		examination.Annotations = annotations
	}
	return examination
}
