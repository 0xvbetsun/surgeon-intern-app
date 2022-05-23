package evaluationforms

import (
	"context"
	"database/sql"

	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/internal/pkg/graphqlhelpers"
	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/queryhelpers"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type (
	IMiniCexRepo interface {
		GetMiniCexEvaluationById(ctx context.Context, connectionId string) (*models.MiniCexEvaluation, error)
		GetMiniCexFocusById(ctx context.Context, focusId int) (*models.MiniCexFocuse, error)
		GetMiniCexAreaById(ctx context.Context, areaId int) (*models.MiniCexArea, error)
		AddMiniCexEvaluation(ctx context.Context, tx *sql.Tx, evaluation *models.MiniCexEvaluation) (*models.MiniCexEvaluation, error)
		AddMiniCexFocus(ctx context.Context, focus *models.MiniCexFocuse) (*models.MiniCexFocuse, error)
		AddMiniCexArea(ctx context.Context, area *models.MiniCexArea) (*models.MiniCexArea, error)
		UpdateMiniCexEvaluation(ctx context.Context, tx *sql.Tx, evaluation *models.MiniCexEvaluation) (*models.MiniCexEvaluation, error)
		ListMiniCexEvaluationsByFilter(ctx context.Context, userId string, queryFilter commonModel.MiniCexQueryFilter, orderBy *commonModel.QueryOrder) ([]*models.MiniCexEvaluation, error)
		ListAllMiniCexFocuses(ctx context.Context) ([]*models.MiniCexFocuse, error)
		ListAllMiniCexAreasByClinic(ctx context.Context, departmentID string) ([]*models.MiniCexArea, error)
		DeleteMiniCexEvaluation(ctx context.Context, miniCexEvaluation *models.MiniCexEvaluation) error
	}
	MiniCexRepo struct {
		db *sql.DB
	}
)

func NewMiniCexRepo(db *sql.DB) IMiniCexRepo {
	return &MiniCexRepo{db: db}
}

func (r MiniCexRepo) GetMiniCexEvaluationById(ctx context.Context, evaluationId string) (*models.MiniCexEvaluation, error) {
	evaluation, err := models.MiniCexEvaluations(
		models.MiniCexEvaluationWhere.ID.EQ(evaluationId),
		qm.Load(qm.Rels(models.MiniCexEvaluationRels.MiniCexAssessment, models.AssessmentRels.Activity)),
		qm.Load(models.MiniCexEvaluationRels.Department),
	).One(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return evaluation, nil
}

func (r MiniCexRepo) GetMiniCexFocusById(ctx context.Context, focusId int) (*models.MiniCexFocuse, error) {
	focus, err := models.MiniCexFocuses(
		models.MiniCexFocuseWhere.ID.EQ(focusId),
	).One(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return focus, nil
}

func (r MiniCexRepo) GetMiniCexAreaById(ctx context.Context, areaId int) (*models.MiniCexArea, error) {
	area, err := models.MiniCexAreas(
		models.MiniCexAreaWhere.ID.EQ(areaId),
	).One(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return area, nil
}

func (r MiniCexRepo) ListMiniCexEvaluationsByFilter(ctx context.Context, userId string, queryFilter commonModel.MiniCexQueryFilter, orderBy *commonModel.QueryOrder) ([]*models.MiniCexEvaluation, error) {
	var rels []string
	var joins []string
	var wheres []qm.QueryMod

	// Joins
	fields := graphqlhelpers.GetFields(ctx)

	miniCexRels, miniCexJoins, miniCexFilters := queryhelpers.GetMiniCexRels(ctx, userId, fields, queryFilter)

	joins = append(joins, miniCexJoins...)
	wheres = append(wheres, miniCexFilters...)
	for _, miniCexRel := range miniCexRels {
		rels = append(rels, miniCexRel)
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

	miniCexEvaluations, err := models.MiniCexEvaluations(
		wheres...).All(ctx, r.db)
	if err != nil {
		return nil, err
	}

	return miniCexEvaluations, nil
}

func (r MiniCexRepo) ListAllMiniCexFocuses(ctx context.Context) ([]*models.MiniCexFocuse, error) {
	focuses, err := models.MiniCexFocuses().All(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return focuses, nil
}

func (r MiniCexRepo) ListAllMiniCexAreasByClinic(ctx context.Context, departmentID string) ([]*models.MiniCexArea, error) {
	areas, err := models.MiniCexAreas(models.MiniCexAreaWhere.DepartmentID.EQ(departmentID)).All(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return areas, nil
}

func (r MiniCexRepo) AddMiniCexEvaluation(ctx context.Context, tx *sql.Tx, evaluation *models.MiniCexEvaluation) (*models.MiniCexEvaluation, error) {
	err := evaluation.Insert(ctx, tx, boil.Infer())
	if err != nil {
		return nil, err
	}
	return evaluation, err
}

func (r MiniCexRepo) AddMiniCexFocus(ctx context.Context, focus *models.MiniCexFocuse) (*models.MiniCexFocuse, error) {
	err := focus.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return nil, err
	}
	return focus, err
}

func (r MiniCexRepo) AddMiniCexArea(ctx context.Context, area *models.MiniCexArea) (*models.MiniCexArea, error) {
	err := area.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return nil, err
	}
	return area, err
}

func (r MiniCexRepo) UpdateMiniCexEvaluation(ctx context.Context, tx *sql.Tx, connection *models.MiniCexEvaluation) (*models.MiniCexEvaluation, error) {
	if _, err := connection.Update(ctx, tx, boil.Infer()); err != nil {
		return nil, err
	}
	if err := connection.Reload(ctx, tx); err != nil {
		return nil, err
	}
	return connection, nil
}

func (r MiniCexRepo) DeleteMiniCexEvaluation(ctx context.Context, miniCexEvaluation *models.MiniCexEvaluation) error {
	_, err := miniCexEvaluation.Delete(ctx, r.db)
	if err != nil {
		return err
	}
	return nil
}
