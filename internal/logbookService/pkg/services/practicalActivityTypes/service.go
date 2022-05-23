package practicalActivityTypes

import (
	"context"

	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/practicalactivitytypes"
)

type IService interface {
	List(ctx context.Context) ([]*commonModel.PracticalActivityType, error)
	ListByUser(ctx context.Context, userId string) ([]*commonModel.PracticalActivityType, error)
	GetByType(ctx context.Context, activityType commonModel.ActivityType) (*commonModel.PracticalActivityType, error)
}

type Service struct {
	Repo practicalactivitytypes.IRepo
}

type ActivityTypeFilter struct {
	activityType *commonModel.ActivityType
}

func NewService(repo practicalactivitytypes.IRepo) IService {
	return &Service{
		Repo: repo,
	}
}

func (s Service) ListByUser(ctx context.Context, userId string) ([]*commonModel.PracticalActivityType, error) {
	userActivityTypes, err := s.Repo.GetByAuthUserAndFilter(ctx, userId, practicalactivitytypes.Filter{TypeNames: nil})
	if err != nil {
		return nil, err
	}
	return commonModel.PracticalActivityTypesFromRepoType(userActivityTypes), nil
}

func (s Service) List(ctx context.Context) ([]*commonModel.PracticalActivityType, error) {
	allActivityTypes, err := s.Repo.All(ctx)
	if err != nil {
		return nil, err
	}

	return commonModel.PracticalActivityTypesFromRepoType(allActivityTypes), nil
}

func (s Service) GetByType(ctx context.Context, activityType commonModel.ActivityType) (*commonModel.PracticalActivityType, error) {
	practicalActivityType, err := s.Repo.GetByType(ctx, activityType)
	if err != nil {
		return nil, err
	}
	practicalActivity := commonModel.PracticalActivityTypeFromRepoType(practicalActivityType)
	return practicalActivity, nil
}
