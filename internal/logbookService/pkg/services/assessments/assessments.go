package assessments

import (
	"context"

	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/internal/logbookService/pkg/common"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/assessments"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/users"
)

type (
	IAssessmentsService interface {
		GetAssessments(ctx context.Context, queryFilter commonModel.AssessmentQueryFilter, orderBy *commonModel.QueryOrder, pagination *commonModel.QueryPaging) ([]*commonModel.Assessment, error)
	}
	AssessmentsService struct {
		assessmentsRepo assessments.IAssessmentsRepo
		usersRepo       users.IRepo
	}
)

func NewAssessmentsService(assessmentsRepo assessments.IAssessmentsRepo, usersRepo users.IRepo) IAssessmentsService {
	return &AssessmentsService{
		assessmentsRepo: assessmentsRepo,
		usersRepo:       usersRepo,
	}
}

func (a *AssessmentsService) GetAssessments(ctx context.Context, queryFilter commonModel.AssessmentQueryFilter, orderBy *commonModel.QueryOrder, pagination *commonModel.QueryPaging) ([]*commonModel.Assessment, error) {
	user, err := a.usersRepo.GetByAuthenticationContext(ctx)
	if err != nil {
		return nil, err
	}

	assessments, err := a.assessmentsRepo.ListAssessmentsByFilter(ctx, user.ID, queryFilter, orderBy, pagination)
	if err != nil {
		return nil, err
	}
	qlAssessments := make([]*commonModel.Assessment, 0)

	for _, assessment := range assessments {
		mappedAssessment := common.MapAssessmentGraphQlModel(assessment)
		qlAssessments = append(qlAssessments, mappedAssessment)
	}
	return qlAssessments, nil
}
