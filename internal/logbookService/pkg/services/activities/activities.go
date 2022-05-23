package activities

import (
	"context"

	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/internal/logbookService/pkg/common"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/activities"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/dbexecutor"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/users"
)

type (
	IActivitiesService interface {
		GetActivities(ctx context.Context, queryFilter commonModel.ActivityQueryFilter, orderBy *commonModel.QueryOrder, pagination *commonModel.QueryPaging) ([]*commonModel.Activity, error)
	}

	ActivitiesService struct {
		activitiesRepo activities.IActivitiesRepo
		usersRepo      users.IRepo
		dbexecutor.IDBExecutor
	}
)

func NewActivitiesService(activitiesRepo activities.IActivitiesRepo, usersRepo users.IRepo, dbExecutor dbexecutor.IDBExecutor) IActivitiesService {
	return &ActivitiesService{
		activitiesRepo: activitiesRepo,
		usersRepo:      usersRepo,
		IDBExecutor:    dbExecutor,
	}
}

func (s ActivitiesService) GetActivities(ctx context.Context, queryFilter commonModel.ActivityQueryFilter, orderBy *commonModel.QueryOrder, pagination *commonModel.QueryPaging) ([]*commonModel.Activity, error) {
	user, err := s.usersRepo.GetByAuthenticationContext(ctx)
	if err != nil {
		return nil, err
	}

	activities, err := s.activitiesRepo.ListActivitiesByFilter(ctx, user.ID, queryFilter, orderBy, pagination)
	if err != nil {
		return nil, err
	}

	qlActivities := make([]*commonModel.Activity, 0)

	for _, activity := range activities {
		mappedActivity := common.MapActivityGraphQlModel(activity)
		qlActivities = append(qlActivities, mappedActivity)
	}
	return qlActivities, nil
}
