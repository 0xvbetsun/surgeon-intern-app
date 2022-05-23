package logbookentries

import (
	"context"

	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/internal/logbookService/pkg/common"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/logbookentries"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/users"
)

type (
	ILogbookEntriesService interface {
		GetLogbookEntries(ctx context.Context, queryFilter commonModel.LogbookEntryQueryFilter, orderBy *commonModel.QueryOrder, pagination *commonModel.QueryPaging) ([]*commonModel.LogbookEntry, error)
	}
	LogbookEntriesService struct {
		logbookEntriesRepo logbookentries.ILogbookEntriesRepo
		usersRepo          users.IRepo
	}
)

func NewActivityService(logbookEntriesRepo logbookentries.ILogbookEntriesRepo, usersRepo users.IRepo) ILogbookEntriesService {
	return &LogbookEntriesService{
		logbookEntriesRepo: logbookEntriesRepo,
		usersRepo:          usersRepo,
	}
}

func (a *LogbookEntriesService) GetLogbookEntries(ctx context.Context, queryFilter commonModel.LogbookEntryQueryFilter, orderBy *commonModel.QueryOrder, pagination *commonModel.QueryPaging) ([]*commonModel.LogbookEntry, error) {
	user, err := a.usersRepo.GetByAuthenticationContext(ctx)
	if err != nil {
		return nil, err
	}

	logbookEntries, err := a.logbookEntriesRepo.ListLogbookEntriesByFilter(ctx, user.ID, queryFilter, orderBy, pagination)
	if err != nil {
		return nil, err
	}
	qlLogbookEntries := make([]*commonModel.LogbookEntry, 0)

	for _, logbookEntry := range logbookEntries {
		mappedLogbookEntry := common.MapLogbookEntryGraphQlModel(logbookEntry)
		qlLogbookEntries = append(qlLogbookEntries, mappedLogbookEntry)
	}
	return qlLogbookEntries, nil
}
