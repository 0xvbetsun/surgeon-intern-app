package orthopedicSurgeries

import (
	"context"
	"database/sql"

	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/repoerrors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type (
	IActivitiesRepo interface {
		GetByResidentID(ctx context.Context, userID string, loadRelations bool) ([]*models.OrthopedicSurgeryActivity, error)
		GetBySupervisorID(ctx context.Context, userID string) ([]*models.OrthopedicSurgeryActivity, error)
		GetByFilter(ctx context.Context, userID string, filter *commonModel.LogbookEntryQueryFilter, loadRelations bool, orderBy *commonModel.QueryOrder) ([]*models.OrthopedicSurgeryActivity, error)
		Add(ctx context.Context, executor boil.ContextExecutor, orthopedicSurgeryActivity *models.OrthopedicSurgeryActivity) (*models.OrthopedicSurgeryActivity, error)
		Update(ctx context.Context, executor boil.ContextExecutor, orthopedicSurgeryActivity *models.OrthopedicSurgeryActivity) (*models.OrthopedicSurgeryActivity, error)
		AddSurgeriesToActivity(ctx context.Context, executor boil.ContextExecutor, orthopedicSurgeryActivityID string, surgeries []*models.Surgery) error
		UpdateActivitySurgeries(ctx context.Context, executor boil.ContextExecutor, orthopedicSurgeryActivityID string, surgeries []*models.Surgery) error
		GetByID(ctx context.Context, orthopedicSurgeryActivityID string) (*models.OrthopedicSurgeryActivity, error)
		GetByIDWithRels(ctx context.Context, orthopedicSurgeryActivityID string) (*models.OrthopedicSurgeryActivity, error)
		ListAllSurgeryActivitiesByFilter(ctx context.Context, userId string, queryFilter commonModel.SurgeryLogbookEntryQueryFilter, orderBy *commonModel.QueryOrder) ([]*models.OrthopedicSurgeryActivity, error)
		DeleteActivity(ctx context.Context, activity *models.OrthopedicSurgeryActivity) error
	}
	ActivitiesRepo struct {
		db *sql.DB
	}
)

func NewActivitiesRepo(db *sql.DB) IActivitiesRepo {
	return &ActivitiesRepo{db: db}
}

func (r *ActivitiesRepo) GetByFilter(ctx context.Context, userID string, filter *commonModel.LogbookEntryQueryFilter, loadRelations bool, orderBy *commonModel.QueryOrder) ([]*models.OrthopedicSurgeryActivity, error) {
	var queryMods []QueryMod

	if loadRelations {
		queryMods = append(queryMods,
			[]QueryMod{
				Load(models.OrthopedicSurgeryActivityRels.PracticalActivityType),
				Load(models.OrthopedicSurgeryActivityRels.Supervisor),
				Load(models.OrthopedicSurgeryActivityRels.DopsEvaluation),
				Load(models.OrthopedicSurgeryActivityRels.OrthopedicSurgeriesActivityReview),
				Load(models.OrthopedicSurgeryActivityRels.Resident),
				Load(Rels(models.OrthopedicSurgeryActivityRels.OrthopedicSurgeryActivitiesSurgeries, Rels(models.OrthopedicSurgeryActivitiesSurgeryRels.Surgery, models.SurgeryRels.Method))),
				Load(Rels(models.OrthopedicSurgeryActivityRels.OrthopedicSurgeryActivitiesSurgeries, Rels(models.OrthopedicSurgeryActivitiesSurgeryRels.Surgery, models.SurgeryRels.Diagnose)))}...,
		)
	}
	if filter == nil {
		queryMods = append(queryMods, Expr(models.OrthopedicSurgeryActivityWhere.ResidentID.EQ(userID),
			Or2(models.OrthopedicSurgeryActivityWhere.SupervisorID.EQ(null.StringFrom(userID)))))
		rOrthopedicSurgeryActivities, err := models.OrthopedicSurgeryActivities(queryMods...).All(ctx, r.db)
		if err != nil {
			return nil, err
		}
		return rOrthopedicSurgeryActivities, nil
	}
	if filter.HasDops != nil {
		queryMods = append(queryMods, models.OrthopedicSurgeryActivityWhere.DopsRequested.EQ(*filter.HasDops))
		if filter.DopsEvaluated != nil {
			queryMods = append(queryMods, InnerJoin(models.TableNames.DopsEvaluations+
				" on "+
				models.DopsEvaluationTableColumns.OrthopedicSurgeryActivityID+
				" = "+
				models.OrthopedicSurgeryActivityTableColumns.ID),
				models.DopsEvaluationWhere.IsEvaluated.EQ(*filter.DopsEvaluated))
		}
	}

	if filter.HasReview != nil && *filter.HasReview {
		queryMods = append(queryMods,
			InnerJoin(models.TableNames.OrthopedicSurgeriesActivityReview+
				" on "+
				models.OrthopedicSurgeriesActivityReviewTableColumns.OrthopedicSurgeryActivityID+
				" = "+
				models.OrthopedicSurgeryActivityTableColumns.ID))
		if filter.IsReviewed != nil {
			if *filter.IsReviewed {
				queryMods = append(queryMods, models.OrthopedicSurgeriesActivityReviewWhere.SignedAt.NEQ(null.TimeFromPtr(nil)))
			} else {
				queryMods = append(queryMods, models.OrthopedicSurgeriesActivityReviewWhere.SignedAt.EQ(null.TimeFromPtr(nil)))
			}
		}
	}
	if filter.ResidentID != nil {
		queryMods = append(queryMods, models.OrthopedicSurgeryActivityWhere.ResidentID.EQ(*filter.ResidentID))
	}
	if filter.SupervisorID != nil {
		queryMods = append(queryMods, models.OrthopedicSurgeryActivityWhere.SupervisorID.EQ(null.StringFromPtr(filter.SupervisorID)))
	}

	if filter.InProgress != nil {
		queryMods = append(queryMods, models.OrthopedicSurgeryActivityWhere.InProgress.EQ(*filter.InProgress))
	}

	if filter.SurgeryLogbookEntryFilters != nil {
		if len(filter.SurgeryLogbookEntryFilters.Surgeries) > 0 {
			queryMods = append(queryMods,
				InnerJoin(models.TableNames.OrthopedicSurgeryActivitiesSurgeries+
					" on "+
					models.OrthopedicSurgeryActivitiesSurgeryColumns.OrthopedicSurgeryActivityID+
					" = "+
					models.OrthopedicSurgeryActivityTableColumns.ID), models.OrthopedicSurgeryActivitiesSurgeryWhere.SurgeryID.IN(filter.SurgeryLogbookEntryFilters.Surgeries))
		}
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
			queryMods = append(queryMods, OrderBy(models.OrthopedicSurgeryActivityColumns.OccurredAt+" "+order))
			break
		}
	}

	rOrthopedicSurgeryActivities, err := models.OrthopedicSurgeryActivities(queryMods...).All(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return rOrthopedicSurgeryActivities, nil
}

func (r *ActivitiesRepo) GetByResidentID(ctx context.Context, userID string, loadRelations bool) ([]*models.OrthopedicSurgeryActivity, error) {
	queryMods := []QueryMod{
		Where("resident_id=?", userID),
	}
	if loadRelations {
		queryMods = append(queryMods,
			[]QueryMod{
				Load(models.OrthopedicSurgeryActivityRels.Supervisor),
				Load(models.OrthopedicSurgeryActivityRels.DopsEvaluation),
				Load(models.OrthopedicSurgeryActivityRels.OrthopedicSurgeriesActivityReview),
				Load(models.OrthopedicSurgeryActivityRels.Resident)}...,
		)
	}
	rOrthopedicSurgeryActivities, err := models.OrthopedicSurgeryActivities(queryMods...).All(ctx, r.db)

	if err != nil {
		return nil, err
	}
	return rOrthopedicSurgeryActivities, nil
}

func (r *ActivitiesRepo) GetBySupervisorID(ctx context.Context, userID string) ([]*models.OrthopedicSurgeryActivity, error) {
	rOrthopedicSurgeryActivities, err := models.OrthopedicSurgeryActivities(Where("supervisor_id=?", userID)).All(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return rOrthopedicSurgeryActivities, nil
}

func (r *ActivitiesRepo) ListAllSurgeryActivitiesByFilter(ctx context.Context, userId string, queryFilter commonModel.SurgeryLogbookEntryQueryFilter, orderBy *commonModel.QueryOrder) ([]*models.OrthopedicSurgeryActivity, error) {
	queryMods := []QueryMod{
		models.OrthopedicSurgeryActivityWhere.ResidentID.EQ(userId),
		Load(Rels(models.OrthopedicSurgeryActivityRels.OrthopedicSurgeryActivitiesSurgeries, Rels(models.OrthopedicSurgeryActivitiesSurgeryRels.Surgery, models.SurgeryRels.Method))),
		Load(Rels(models.OrthopedicSurgeryActivityRels.OrthopedicSurgeryActivitiesSurgeries, Rels(models.OrthopedicSurgeryActivitiesSurgeryRels.Surgery, models.SurgeryRels.Diagnose))),
		Load(models.OrthopedicSurgeryActivityRels.Resident),
		Load(models.OrthopedicSurgeryActivityRels.Supervisor),
	}

	if queryFilter.HasDops != nil {
		queryMods = append(queryMods, models.OrthopedicSurgeryActivityWhere.DopsRequested.EQ(*queryFilter.HasDops))
	}
	if queryFilter.InProgress != nil {
		queryMods = append(queryMods, models.OrthopedicSurgeryActivityWhere.InProgress.EQ(*queryFilter.InProgress))
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
			queryMods = append(queryMods, OrderBy(models.OrthopedicSurgeryActivityColumns.OccurredAt+" "+order))
			break
		}
	}

	rOrthopedicSurgeryActivities, err := models.OrthopedicSurgeryActivities(queryMods...).All(ctx, r.db)

	if err != nil {
		return nil, err
	}
	return rOrthopedicSurgeryActivities, nil
}

func (r *ActivitiesRepo) Add(ctx context.Context, executor boil.ContextExecutor, orthopedicSurgeryActivity *models.OrthopedicSurgeryActivity) (*models.OrthopedicSurgeryActivity, error) {
	db := executor
	if executor == nil {
		db = r.db
	}
	err := orthopedicSurgeryActivity.Insert(ctx, db, boil.Infer())

	if err != nil {
		return nil, err
	}
	return orthopedicSurgeryActivity, nil
}

func (r *ActivitiesRepo) Update(ctx context.Context, executor boil.ContextExecutor, orthopedicSurgeryActivity *models.OrthopedicSurgeryActivity) (*models.OrthopedicSurgeryActivity, error) {
	db := executor
	if executor == nil {
		db = r.db
	}
	_, err := orthopedicSurgeryActivity.Update(ctx, db, boil.Infer())
	if err != nil {
		return nil, err
	}
	err = orthopedicSurgeryActivity.Reload(ctx, db)
	if err != nil {
		return nil, err
	}
	return orthopedicSurgeryActivity, nil
}

func (r *ActivitiesRepo) GetByID(ctx context.Context, orthopedicSurgeryActivityID string) (*models.OrthopedicSurgeryActivity, error) {
	rOrthopedicSurgeryActivity, err := models.OrthopedicSurgeryActivities(
		models.OrthopedicSurgeryActivityWhere.ID.EQ(orthopedicSurgeryActivityID),
		Load(Rels(models.OrthopedicSurgeryActivityRels.OrthopedicSurgeryLogbookEntry, models.LogbookEntryRels.Activity)),
	).One(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return rOrthopedicSurgeryActivity, nil
}

func (r *ActivitiesRepo) GetByIDWithRels(ctx context.Context, orthopedicSurgeryActivityID string) (*models.OrthopedicSurgeryActivity, error) {
	rOrthopedicSurgeryActivity, err := models.OrthopedicSurgeryActivities(
		models.OrthopedicSurgeryActivityWhere.ID.EQ(orthopedicSurgeryActivityID),
		Load(Rels(models.OrthopedicSurgeryActivityRels.OrthopedicSurgeryLogbookEntry, models.LogbookEntryRels.Activity)),
		Load(models.OrthopedicSurgeryActivityRels.DopsEvaluation),
		Load(models.OrthopedicSurgeryActivityRels.OrthopedicSurgeriesActivityReview),
		Load(Rels(models.OrthopedicSurgeryActivityRels.OrthopedicSurgeryActivitiesSurgeries, Rels(models.OrthopedicSurgeryActivitiesSurgeryRels.Surgery, models.SurgeryRels.Method))),
		Load(Rels(models.OrthopedicSurgeryActivityRels.OrthopedicSurgeryActivitiesSurgeries, Rels(models.OrthopedicSurgeryActivitiesSurgeryRels.Surgery, models.SurgeryRels.Diagnose))),
		Load(models.OrthopedicSurgeryActivityRels.Resident),
		Load(models.OrthopedicSurgeryActivityRels.Supervisor),
	).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repoerrors.NewNotFoundErrorWithCause(orthopedicSurgeryActivityID, err)
		}
		return nil, repoerrors.ErrorFromDbError(err)
	}
	return rOrthopedicSurgeryActivity, nil
}

func (r *ActivitiesRepo) AddSurgeriesToActivity(ctx context.Context, executor boil.ContextExecutor, orthopedicSurgeryActivityID string, surgeries []*models.Surgery) error {
	db := executor
	if executor == nil {
		db = r.db
	}
	activity, err := models.FindOrthopedicSurgeryActivity(ctx, db, orthopedicSurgeryActivityID)
	if err != nil {
		return err
	}
	for _, surgery := range surgeries {
		if err = activity.AddOrthopedicSurgeryActivitiesSurgeries(ctx, db, true, &models.OrthopedicSurgeryActivitiesSurgery{
			SurgeryID:                   surgery.ID,
			OrthopedicSurgeryActivityID: activity.ID,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (r *ActivitiesRepo) UpdateActivitySurgeries(ctx context.Context, executor boil.ContextExecutor, orthopedicSurgeryActivityID string, surgeries []*models.Surgery) error {
	db := executor
	if executor == nil {
		db = r.db
	}
	_, err := models.OrthopedicSurgeryActivitiesSurgeries(models.OrthopedicSurgeryActivitiesSurgeryWhere.OrthopedicSurgeryActivityID.EQ(orthopedicSurgeryActivityID)).DeleteAll(ctx, db)
	if err != nil {
		return err
	}
	return r.AddSurgeriesToActivity(ctx, db, orthopedicSurgeryActivityID, surgeries)
}

func (r *ActivitiesRepo) DeleteActivity(ctx context.Context, activity *models.OrthopedicSurgeryActivity) error {
	_, err := activity.Delete(ctx, r.db)
	if err != nil {
		return err
	}
	return nil
}
