package activities

import (
	"context"
	"database/sql"

	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/internal/pkg/graphqlhelpers"
	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/queryhelpers"
	"github.com/volatiletech/sqlboiler/v4/boil"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type (
	IActivitiesRepo interface {
		ListActivitiesByFilter(ctx context.Context, userId string, queryFilter commonModel.ActivityQueryFilter, orderBy *commonModel.QueryOrder, pagination *commonModel.QueryPaging) ([]*models.Activity, error)
		AddActivity(ctx context.Context, tx *sql.Tx, assessment *models.Activity) (*models.Activity, error)
		UpdateActivity(ctx context.Context, tx *sql.Tx, assessment *models.Activity) (*models.Activity, error)
	}

	ActivitiesRepo struct {
		db *sql.DB
	}
)

func NewActivitiesRepo(db *sql.DB) IActivitiesRepo {
	return &ActivitiesRepo{db: db}
}

func (r ActivitiesRepo) ListActivitiesByFilter(ctx context.Context, userId string, queryFilter commonModel.ActivityQueryFilter, orderBy *commonModel.QueryOrder, pagination *commonModel.QueryPaging) ([]*models.Activity, error) {
	var joinPaths []string
	var filters []QueryMod

	// Joins
	fields := graphqlhelpers.GetFields(ctx)

	logbookEntriesFilter := commonModel.LogbookEntryQueryFilter{}
	if queryFilter.LogbookEntryFilter != nil {
		logbookEntriesFilter = *queryFilter.LogbookEntryFilter
	}
	if logbookEntriesField, exists := fields["logbookEntry"]; exists && logbookEntriesFilter.LogbookEntryTypes != nil && len(logbookEntriesFilter.LogbookEntryTypes) > 0 {
		joinPaths = append(joinPaths, models.TableNames.LogbookEntries+
			" on "+
			models.ActivityTableColumns.LogbookEntryID+
			" = "+
			models.LogbookEntryTableColumns.ID)

		filters = append(filters, Load(models.ActivityRels.LogbookEntry))
		logbookEntriesFields := graphqlhelpers.GetNestedFields(ctx, logbookEntriesField)
		logbookEntriesRels, logbookEntriesJoinPaths, logbookEntriesFilters := queryhelpers.GetLogbookEntriesFilters(ctx, userId, logbookEntriesFields, logbookEntriesFilter)
		joinPaths = append(joinPaths, logbookEntriesJoinPaths...)
		if len(logbookEntriesFilters) > 0 {
			filters = append(filters, Or2(Expr(logbookEntriesFilters...)))
		}
		for _, logbookEntriesRel := range logbookEntriesRels {
			filters = append(filters, Load(Rels(models.ActivityRels.LogbookEntry, logbookEntriesRel)))
		}
	}

	assessmentFilter := commonModel.AssessmentQueryFilter{}
	if queryFilter.AssessmentFilter != nil {
		assessmentFilter = *queryFilter.AssessmentFilter
	}
	if assessmentsField, exists := fields["assessment"]; exists && assessmentFilter.AssessmentTypes != nil && len(assessmentFilter.AssessmentTypes) > 0 {
		joinPaths = append(joinPaths, models.TableNames.Assessments+
			" on "+
			models.ActivityTableColumns.AssessmentID+
			" = "+
			models.AssessmentTableColumns.ID)

		filters = append(filters, Load(models.ActivityRels.Assessment))
		assessmentFields := graphqlhelpers.GetNestedFields(ctx, assessmentsField)
		assessmentRels, assessmentsJoins, assessmentFilters := queryhelpers.GetAssessmentFilters(ctx, userId, assessmentFields, assessmentFilter)
		joinPaths = append(joinPaths, assessmentsJoins...)
		if len(assessmentFilters) > 0 {
			filters = append(filters, Or2(Expr(assessmentFilters...)))
		}
		for _, assessmentRel := range assessmentRels {
			filters = append(filters, Load(Rels(models.ActivityRels.Assessment, assessmentRel)))
		}
	}

	// If no logbook entries or assessments then don't query for any activities
	if (logbookEntriesFilter.LogbookEntryTypes == nil || len(logbookEntriesFilter.LogbookEntryTypes) < 1) &&
		(assessmentFilter.AssessmentTypes == nil || len(assessmentFilter.AssessmentTypes) < 1) {
		filters = append(filters, models.ActivityWhere.LogbookEntryID.IsNull())
		filters = append(filters, models.ActivityWhere.AssessmentID.IsNull())
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

	return activities, nil
}

func (r ActivitiesRepo) AddActivity(ctx context.Context, tx *sql.Tx, assessment *models.Activity) (*models.Activity, error) {
	err := assessment.Insert(ctx, tx, boil.Infer())
	if err != nil {
		return nil, err
	}
	return assessment, err
}

func (r ActivitiesRepo) UpdateActivity(ctx context.Context, tx *sql.Tx, assessment *models.Activity) (*models.Activity, error) {
	if _, err := assessment.Update(ctx, tx, boil.Infer()); err != nil {
		return nil, err
	}
	if err := assessment.Reload(ctx, tx); err != nil {
		return nil, err
	}
	return assessment, nil
}
