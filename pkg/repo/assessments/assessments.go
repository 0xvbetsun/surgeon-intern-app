package assessments

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
	IAssessmentsRepo interface {
		ListAssessmentsByFilter(ctx context.Context, userId string, queryFilter commonModel.AssessmentQueryFilter, orderBy *commonModel.QueryOrder, pagination *commonModel.QueryPaging) ([]*models.Assessment, error)
		AddAssessment(ctx context.Context, tx *sql.Tx, assessment *models.Assessment) (*models.Assessment, error)
		UpdateAssessment(ctx context.Context, tx *sql.Tx, assessment *models.Assessment) (*models.Assessment, error)
	}
	AssessmentsRepo struct {
		db *sql.DB
	}
)

func NewAssessmentsRepo(db *sql.DB) IAssessmentsRepo {
	return &AssessmentsRepo{db: db}
}

func (r AssessmentsRepo) ListAssessmentsByFilter(ctx context.Context, userId string, queryFilter commonModel.AssessmentQueryFilter, orderBy *commonModel.QueryOrder, pagination *commonModel.QueryPaging) ([]*models.Assessment, error) {
	var joinPaths []string
	var filters []QueryMod

	// Joins
	fields := graphqlhelpers.GetFields(ctx)

	joinPaths = append(joinPaths, models.TableNames.Assessments+
		" on "+
		models.ActivityTableColumns.AssessmentID+
		" = "+
		models.AssessmentTableColumns.ID)

	filters = append(filters, Load(models.ActivityRels.Assessment))
	assessmentRels, assessmentsJoins, assessmentFilters := queryhelpers.GetAssessmentFilters(ctx, userId, fields, queryFilter)
	joinPaths = append(joinPaths, assessmentsJoins...)
	if len(assessmentFilters) > 0 {
		filters = append(filters, Or2(Expr(assessmentFilters...)))
	}
	for _, assessmentRel := range assessmentRels {
		filters = append(filters, Load(Rels(models.ActivityRels.Assessment, assessmentRel)))
	}

	// If no logbook entries or assessments then don't query for any activities
	if queryFilter.AssessmentTypes == nil || len(queryFilter.AssessmentTypes) < 1 {
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

	assessments := make([]*models.Assessment, 0)
	for _, activity := range activities {
		assessments = append(assessments, activity.R.Assessment)
	}
	return assessments, nil
}

func (r AssessmentsRepo) AddAssessment(ctx context.Context, tx *sql.Tx, assessment *models.Assessment) (*models.Assessment, error) {
	err := assessment.Insert(ctx, tx, boil.Infer())
	if err != nil {
		return nil, err
	}
	return assessment, err
}

func (r AssessmentsRepo) UpdateAssessment(ctx context.Context, tx *sql.Tx, assessment *models.Assessment) (*models.Assessment, error) {
	if _, err := assessment.Update(ctx, tx, boil.Infer()); err != nil {
		return nil, err
	}
	if err := assessment.Reload(ctx, tx); err != nil {
		return nil, err
	}
	return assessment, nil
}
