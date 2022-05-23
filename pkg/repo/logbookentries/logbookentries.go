package logbookentries

import (
	"context"
	"database/sql"

	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/internal/pkg/graphqlhelpers"
	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/queryhelpers"
	"github.com/volatiletech/sqlboiler/v4/boil"
	_ "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type (
	ILogbookEntriesRepo interface {
		ListLogbookEntriesByFilter(ctx context.Context, userId string, queryFilter commonModel.LogbookEntryQueryFilter, orderBy *commonModel.QueryOrder, pagination *commonModel.QueryPaging) ([]*models.LogbookEntry, error)
		AddLogbookEntry(ctx context.Context, tx *sql.Tx, logbookEntry *models.LogbookEntry) (*models.LogbookEntry, error)
		UpdateLogbookEntry(ctx context.Context, tx *sql.Tx, logbookEntry *models.LogbookEntry) (*models.LogbookEntry, error)
	}
	LogbookEntriesRepo struct {
		db *sql.DB
	}
)

func NewLogbookEntriesRepo(db *sql.DB) ILogbookEntriesRepo {
	return &LogbookEntriesRepo{db: db}
}

func (r LogbookEntriesRepo) ListLogbookEntriesByFilter(ctx context.Context, userId string, queryFilter commonModel.LogbookEntryQueryFilter, orderBy *commonModel.QueryOrder, pagination *commonModel.QueryPaging) ([]*models.LogbookEntry, error) {
	var joinPaths []string
	var filters []QueryMod

	// Joins
	fields := graphqlhelpers.GetFields(ctx)

	joinPaths = append(joinPaths, models.TableNames.LogbookEntries+
		" on "+
		models.ActivityTableColumns.LogbookEntryID+
		" = "+
		models.LogbookEntryTableColumns.ID)

	filters = append(filters, Load(models.ActivityRels.LogbookEntry))
	logbookEntriesRels, logbookEntriesJoinPaths, logbookEntriesFilters := queryhelpers.GetLogbookEntriesFilters(ctx, userId, fields, queryFilter)
	joinPaths = append(joinPaths, logbookEntriesJoinPaths...)
	if len(logbookEntriesFilters) > 0 {
		filters = append(filters, Or2(Expr(logbookEntriesFilters...)))
	}
	for _, logbookEntriesRel := range logbookEntriesRels {
		filters = append(filters, Load(Rels(models.ActivityRels.LogbookEntry, logbookEntriesRel)))
	}

	// If no logbook entries or assessments then don't query for any activities
	if queryFilter.LogbookEntryTypes == nil || len(queryFilter.LogbookEntryTypes) < 1 {
		filters = append(filters, models.ActivityWhere.AssessmentID.IsNull())
		filters = append(filters, models.ActivityWhere.LogbookEntryID.IsNull())
	}

	// Order by
	if orderBy != nil {
		order := ""
		switch orderBy.Order {
		case commonModel.OrderAscending:
			order = "asc"
			break
		case commonModel.OrderDescending:
			order = "desc"
			break
		}
		switch orderBy.OrderBy {
		case commonModel.OrderByOcurredAt:
			filters = append(filters, OrderBy(models.ActivityTableColumns.OccurredAt+" "+order))
			break
		}
	}

	joinMap := make(map[string]string)
	for _, joinPath := range joinPaths {
		if _, exists := joinMap[joinPath]; !exists {
			joinMap[joinPath] = joinPath
			filters = append(filters, LeftOuterJoin(joinPath))
		}
	}

	activities, err := models.Activities(filters...).All(ctx, r.db)
	if err != nil {
		return nil, err
	}

	logbookEntries := make([]*models.LogbookEntry, 0)
	for _, activity := range activities {
		logbookEntries = append(logbookEntries, activity.R.LogbookEntry)
	}
	return logbookEntries, nil
}

func (r LogbookEntriesRepo) AddLogbookEntry(ctx context.Context, tx *sql.Tx, logbookEntry *models.LogbookEntry) (*models.LogbookEntry, error) {
	err := logbookEntry.Insert(ctx, tx, boil.Infer())
	if err != nil {
		return nil, err
	}
	return logbookEntry, err
}

func (r LogbookEntriesRepo) UpdateLogbookEntry(ctx context.Context, tx *sql.Tx, logbookEntry *models.LogbookEntry) (*models.LogbookEntry, error) {
	if _, err := logbookEntry.Update(ctx, tx, boil.Infer()); err != nil {
		return nil, err
	}
	if err := logbookEntry.Reload(ctx, tx); err != nil {
		return nil, err
	}
	return logbookEntry, nil
}
