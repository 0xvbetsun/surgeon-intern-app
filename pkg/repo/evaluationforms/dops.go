package evaluationforms

import (
	"context"
	"database/sql"

	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/internal/pkg/graphqlhelpers"
	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/queryhelpers"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type (
	IDopsRepo interface {
		GetDopsEvaluationById(ctx context.Context, dopsId string) (*models.DopsEvaluation, error)
		GetDopsEvaluationByActivityId(ctx context.Context, activityId string) (*models.DopsEvaluation, error)
		AddDopsEvaluation(ctx context.Context, tx *sql.Tx, evaluation *models.DopsEvaluation) (*models.DopsEvaluation, error)
		UpdateDopsEvaluation(ctx context.Context, tx *sql.Tx, evaluation *models.DopsEvaluation) (*models.DopsEvaluation, error)
		UpdateDopsEvaluationSurgeries(ctx context.Context, executor boil.ContextExecutor, dopsEvaluationID string, surgeries []*models.Surgery) error
		AddSurgeriesToDopsEvaluation(ctx context.Context, executor boil.ContextExecutor, dopsEvaluationId string, surgeries []*models.Surgery) error
		ListDopsEvaluationsByFilter(ctx context.Context, userId string, queryFilter commonModel.DopsQueryFilter, orderBy *commonModel.QueryOrder) ([]*models.DopsEvaluation, error)
		DeleteDopsEvaluation(ctx context.Context, dopsEvaluation *models.DopsEvaluation) error
	}
	DopsRepo struct {
		db *sql.DB
	}
)

func NewDopsRepo(db *sql.DB) IDopsRepo {
	return &DopsRepo{db: db}
}

func (r DopsRepo) GetDopsEvaluationById(ctx context.Context, evaluationId string) (*models.DopsEvaluation, error) {
	evaluation, err := models.DopsEvaluations(
		models.DopsEvaluationWhere.ID.EQ(evaluationId),
		qm.Load(qm.Rels(models.DopsEvaluationRels.DopAssessment, models.AssessmentRels.Activity)),
		qm.Load(models.DopsEvaluationRels.OrthopedicSurgeryActivity),
		qm.Load(models.DopsEvaluationRels.Department),
		qm.Load(qm.Rels(models.DopsEvaluationRels.DopsEvaluationsSurgeries, qm.Rels(models.DopsEvaluationsSurgeryRels.Surgery, models.SurgeryRels.Method))),
		qm.Load(qm.Rels(models.DopsEvaluationRels.DopsEvaluationsSurgeries, qm.Rels(models.DopsEvaluationsSurgeryRels.Surgery, models.SurgeryRels.Diagnose))),
		qm.Load(qm.Rels(models.DopsEvaluationRels.OrthopedicSurgeryActivity, qm.Rels(models.OrthopedicSurgeryActivityRels.OrthopedicSurgeryActivitiesSurgeries, qm.Rels(models.OrthopedicSurgeryActivitiesSurgeryRels.Surgery, models.SurgeryRels.Method)))),
		qm.Load(qm.Rels(models.DopsEvaluationRels.OrthopedicSurgeryActivity, qm.Rels(models.OrthopedicSurgeryActivityRels.OrthopedicSurgeryActivitiesSurgeries, qm.Rels(models.OrthopedicSurgeryActivitiesSurgeryRels.Surgery, models.SurgeryRels.Diagnose)))),
	).One(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return evaluation, nil
}

func (r DopsRepo) GetDopsEvaluationByActivityId(ctx context.Context, activityId string) (*models.DopsEvaluation, error) {
	evaluation, err := models.DopsEvaluations(
		models.DopsEvaluationWhere.OrthopedicSurgeryActivityID.EQ(null.StringFrom(activityId)),
		qm.Load(qm.Rels(models.DopsEvaluationRels.DopAssessment, models.AssessmentRels.Activity)),
		qm.Load(models.DopsEvaluationRels.OrthopedicSurgeryActivity),
		qm.Load(models.DopsEvaluationRels.Department),
		qm.Load(qm.Rels(models.DopsEvaluationRels.DopsEvaluationsSurgeries, qm.Rels(models.DopsEvaluationsSurgeryRels.Surgery, models.SurgeryRels.Method))),
		qm.Load(qm.Rels(models.DopsEvaluationRels.DopsEvaluationsSurgeries, qm.Rels(models.DopsEvaluationsSurgeryRels.Surgery, models.SurgeryRels.Diagnose))),
		qm.Load(qm.Rels(models.DopsEvaluationRels.OrthopedicSurgeryActivity, qm.Rels(models.OrthopedicSurgeryActivityRels.OrthopedicSurgeryActivitiesSurgeries, qm.Rels(models.OrthopedicSurgeryActivitiesSurgeryRels.Surgery, models.SurgeryRels.Method)))),
		qm.Load(qm.Rels(models.DopsEvaluationRels.OrthopedicSurgeryActivity, qm.Rels(models.OrthopedicSurgeryActivityRels.OrthopedicSurgeryActivitiesSurgeries, qm.Rels(models.OrthopedicSurgeryActivitiesSurgeryRels.Surgery, models.SurgeryRels.Diagnose)))),
	).One(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return evaluation, nil
}

func (r DopsRepo) ListDopsEvaluationsByFilter(ctx context.Context, userId string, queryFilter commonModel.DopsQueryFilter, orderBy *commonModel.QueryOrder) ([]*models.DopsEvaluation, error) {
	var rels []string
	var joins []string
	var wheres []qm.QueryMod

	// Joins
	fields := graphqlhelpers.GetFields(ctx)

	dopsRels, dopsJoins, dopsFilters := queryhelpers.GetDopsFilters(ctx, userId, fields, queryFilter)

	joins = append(joins, dopsJoins...)
	wheres = append(wheres, dopsFilters...)
	for _, dopsRel := range dopsRels {
		rels = append(rels, dopsRel)
	}

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
			wheres = append(wheres, qm.OrderBy(models.DopsEvaluationColumns.OccurredAt+" "+order))
			break
		}
	}

	joinMap := make(map[string]string)
	for _, joinPath := range joins {
		if _, exists := joinMap[joinPath]; !exists {
			joinMap[joinPath] = joinPath
			wheres = append(wheres, qm.LeftOuterJoin(joinPath))
		}
	}

	dopsEvaluations, err := models.DopsEvaluations(
		wheres...).All(ctx, r.db)
	if err != nil {
		return nil, err
	}

	return dopsEvaluations, nil
}

func (r DopsRepo) AddDopsEvaluation(ctx context.Context, tx *sql.Tx, evaluation *models.DopsEvaluation) (*models.DopsEvaluation, error) {
	err := evaluation.Insert(ctx, tx, boil.Infer())
	if err != nil {
		return nil, err
	}
	return evaluation, err
}

func (r DopsRepo) UpdateDopsEvaluation(ctx context.Context, tx *sql.Tx, connection *models.DopsEvaluation) (*models.DopsEvaluation, error) {
	if _, err := connection.Update(ctx, tx, boil.Infer()); err != nil {
		return nil, err
	}
	if err := connection.Reload(ctx, tx); err != nil {
		return nil, err
	}
	return connection, nil
}

func (r DopsRepo) UpdateDopsEvaluationSurgeries(ctx context.Context, executor boil.ContextExecutor, dopsEvaluationID string, surgeries []*models.Surgery) error {
	db := executor
	if executor == nil {
		db = r.db
	}
	_, err := models.DopsEvaluationsSurgeries(models.DopsEvaluationsSurgeryWhere.DopsEvaluationID.EQ(dopsEvaluationID)).DeleteAll(ctx, db)
	if err != nil {
		return err
	}
	return r.AddSurgeriesToDopsEvaluation(ctx, db, dopsEvaluationID, surgeries)
}

func (r DopsRepo) AddSurgeriesToDopsEvaluation(ctx context.Context, executor boil.ContextExecutor, dopsEvaluationId string, surgeries []*models.Surgery) error {
	db := executor
	if executor == nil {
		db = r.db
	}
	dopsEvaluation, err := models.FindDopsEvaluation(ctx, db, dopsEvaluationId)
	if err != nil {
		return err
	}
	for _, surgery := range surgeries {
		if err = dopsEvaluation.AddDopsEvaluationsSurgeries(ctx, db, true, &models.DopsEvaluationsSurgery{
			SurgeryID:        surgery.ID,
			DopsEvaluationID: dopsEvaluation.ID,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (r DopsRepo) DeleteDopsEvaluation(ctx context.Context, dopsEvaluation *models.DopsEvaluation) error {
	_, err := dopsEvaluation.Delete(ctx, r.db)
	if err != nil {
		return err
	}
	return nil
}
