package orthopedicSurgeries

import (
	"context"
	"database/sql"
	"errors"

	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type (
	IReviewsRepo interface {
		GetByFilter(ctx context.Context, userID string, filter *commonModel.SurgeryReviewQueryFilter, loadRelations bool, orderBy *commonModel.QueryOrder) ([]*models.OrthopedicSurgeriesActivityReview, error)
		GetReviewsBySupervisorUserIDWithRels(ctx context.Context, reviewUserID string) ([]*models.OrthopedicSurgeriesActivityReview, error)
		GetReviewByActivity(ctx context.Context, orthopedicSurgeryActivityID string) (*models.OrthopedicSurgeriesActivityReview, error)
		GetReviewByIDWithRels(ctx context.Context, orthopedicSurgeryActivityReviewID string) (*models.OrthopedicSurgeriesActivityReview, error)
		AddOrthopedicSurgeryActivityReview(ctx context.Context, executor boil.ContextExecutor, review *models.OrthopedicSurgeriesActivityReview) (*models.OrthopedicSurgeriesActivityReview, error)
		UpdateOrthopedicSurgeryActivityReview(ctx context.Context, executor boil.ContextExecutor, review *models.OrthopedicSurgeriesActivityReview) (*models.OrthopedicSurgeriesActivityReview, error)
		UpdateSurgeriesForReview(ctx context.Context, executor boil.ContextExecutor, reviewId string, surgeries []*models.Surgery) error
		AddSurgeriesToReview(ctx context.Context, executor boil.ContextExecutor, reviewId string, surgeries []*models.Surgery) error
	}
	ReviewsRepo struct {
		db *sql.DB
	}
)

func NewReviewsRepo(db *sql.DB) IReviewsRepo {
	return &ReviewsRepo{db: db}
}

func (r *ReviewsRepo) GetByFilter(ctx context.Context, userID string, filter *commonModel.SurgeryReviewQueryFilter, loadRelations bool, orderBy *commonModel.QueryOrder) ([]*models.OrthopedicSurgeriesActivityReview, error) {
	var queryMods []QueryMod

	if loadRelations {
		queryMods = append(queryMods,
			[]QueryMod{
				Load(models.OrthopedicSurgeriesActivityReviewRels.OrthopedicSurgeryActivity),
				Load(Rels(models.OrthopedicSurgeriesActivityReviewRels.OrthopedicSurgeryActivity, models.OrthopedicSurgeryActivityRels.DopsEvaluation)),
				Load(Rels(models.OrthopedicSurgeriesActivityReviewRels.OrthopedicSurgeryActivity, models.OrthopedicSurgeryActivityRels.PracticalActivityType)),
				Load(models.OrthopedicSurgeriesActivityReviewRels.Supervisor),
				Load(models.OrthopedicSurgeriesActivityReviewRels.Resident),
				Load(Rels(models.OrthopedicSurgeriesActivityReviewRels.OrthopedicSurgeriesActivityReviewSurgeries, Rels(models.OrthopedicSurgeriesActivityReviewSurgeryRels.Surgery, models.SurgeryRels.Method))),
				Load(Rels(models.OrthopedicSurgeriesActivityReviewRels.OrthopedicSurgeriesActivityReviewSurgeries, Rels(models.OrthopedicSurgeriesActivityReviewSurgeryRels.Surgery, models.SurgeryRels.Diagnose)))}...,
		)
	}
	if filter == nil {
		queryMods = append(queryMods, models.OrthopedicSurgeriesActivityReviewWhere.SupervisorID.EQ(userID))
		reviews, err := models.OrthopedicSurgeriesActivityReviews(queryMods...).All(ctx, r.db)
		if err != nil {
			return nil, err
		}
		return reviews, nil
	}

	if filter.IsEvaluated != nil {
		if *filter.IsEvaluated {
			queryMods = append(queryMods, models.OrthopedicSurgeriesActivityReviewWhere.SignedAt.NEQ(null.TimeFromPtr(nil)))
		} else {
			queryMods = append(queryMods, models.OrthopedicSurgeriesActivityReviewWhere.SignedAt.EQ(null.TimeFromPtr(nil)))
		}
	}

	if filter.ResidentID != nil {
		queryMods = append(queryMods, models.OrthopedicSurgeriesActivityReviewWhere.ResidentID.EQ(*filter.ResidentID))
	}
	if filter.SupervisorID != nil {
		queryMods = append(queryMods, models.OrthopedicSurgeriesActivityReviewWhere.SupervisorID.EQ(*filter.SupervisorID))
	}

	if filter.InProgress != nil {
		queryMods = append(queryMods, models.OrthopedicSurgeriesActivityReviewWhere.InProgress.EQ(*filter.InProgress))
	}

	if orderBy != nil {
		order := ""
		switch orderBy.Order {
		case commonModel.OrderAscending:
			order = "asc"
		case commonModel.OrderDescending:
			order = "desc"
		}
		switch orderBy.OrderBy {
		case commonModel.OrderByOcurredAt:
			queryMods = append(queryMods, OrderBy(models.OrthopedicSurgeriesActivityReviewColumns.OccurredAt+" "+order))
		}
	}

	reviews, err := models.OrthopedicSurgeriesActivityReviews(queryMods...).All(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return reviews, nil
}

func (r *ReviewsRepo) GetReviewsBySupervisorUserIDWithRels(ctx context.Context, supervisorUserID string) ([]*models.OrthopedicSurgeriesActivityReview, error) {
	reviews, err := models.OrthopedicSurgeriesActivityReviews(
		models.OrthopedicSurgeriesActivityReviewWhere.SupervisorID.EQ(supervisorUserID),
		Load(models.OrthopedicSurgeriesActivityReviewRels.OrthopedicSurgeryActivity),
		Load(Rels(models.OrthopedicSurgeriesActivityReviewRels.OrthopedicSurgeryActivity, models.OrthopedicSurgeryActivityRels.DopsEvaluation)),
		Load(Rels(models.OrthopedicSurgeriesActivityReviewRels.OrthopedicSurgeryActivity, models.OrthopedicSurgeryActivityRels.PracticalActivityType)),
		Load(models.OrthopedicSurgeriesActivityReviewRels.Supervisor),
		Load(models.OrthopedicSurgeriesActivityReviewRels.Resident),
		Load(Rels(models.OrthopedicSurgeriesActivityReviewRels.OrthopedicSurgeriesActivityReviewSurgeries, Rels(models.OrthopedicSurgeriesActivityReviewSurgeryRels.Surgery, models.SurgeryRels.Method))),
		Load(Rels(models.OrthopedicSurgeriesActivityReviewRels.OrthopedicSurgeriesActivityReviewSurgeries, Rels(models.OrthopedicSurgeriesActivityReviewSurgeryRels.Surgery, models.SurgeryRels.Diagnose))),
	).All(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return reviews, nil
}

func (r *ReviewsRepo) GetReviewByIDWithRels(ctx context.Context, orthopedicSurgeryActivityReviewID string) (*models.OrthopedicSurgeriesActivityReview, error) {
	review, err := models.OrthopedicSurgeriesActivityReviews(
		models.OrthopedicSurgeriesActivityReviewWhere.ID.EQ(orthopedicSurgeryActivityReviewID),
		Load(Rels(models.OrthopedicSurgeriesActivityReviewRels.OrthopedicSurgeryReviewAssessment, models.AssessmentRels.Activity)),
		Load(models.OrthopedicSurgeriesActivityReviewRels.OrthopedicSurgeryActivity),
		Load(Rels(models.OrthopedicSurgeriesActivityReviewRels.OrthopedicSurgeryActivity, models.OrthopedicSurgeryActivityRels.DopsEvaluation)),
		Load(Rels(models.OrthopedicSurgeriesActivityReviewRels.OrthopedicSurgeryActivity, models.OrthopedicSurgeryActivityRels.PracticalActivityType)),
		Load(models.OrthopedicSurgeriesActivityReviewRels.Supervisor),
		Load(models.OrthopedicSurgeriesActivityReviewRels.Resident),
		Load(Rels(models.OrthopedicSurgeriesActivityReviewRels.OrthopedicSurgeriesActivityReviewSurgeries, Rels(models.OrthopedicSurgeriesActivityReviewSurgeryRels.Surgery, models.SurgeryRels.Method))),
		Load(Rels(models.OrthopedicSurgeriesActivityReviewRels.OrthopedicSurgeriesActivityReviewSurgeries, Rels(models.OrthopedicSurgeriesActivityReviewSurgeryRels.Surgery, models.SurgeryRels.Diagnose))),
	).One(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return review, nil
}

func (r *ReviewsRepo) AddOrthopedicSurgeryActivityReview(ctx context.Context, executor boil.ContextExecutor, review *models.OrthopedicSurgeriesActivityReview) (*models.OrthopedicSurgeriesActivityReview, error) {
	db := executor
	if executor == nil {
		db = r.db
	}
	if err := review.Insert(ctx, db, boil.Infer()); err != nil {
		return nil, err
	}
	return review, nil
}
func (r *ReviewsRepo) UpdateOrthopedicSurgeryActivityReview(ctx context.Context, executor boil.ContextExecutor, review *models.OrthopedicSurgeriesActivityReview) (*models.OrthopedicSurgeriesActivityReview, error) {
	db := executor
	if executor == nil {
		db = r.db
	}
	if _, err := review.Update(ctx, db, boil.Infer()); err != nil {
		return nil, err
	}
	if err := review.Reload(ctx, db); err != nil {
		return nil, err
	}
	return review, nil
}

func (r *ReviewsRepo) GetReviewByActivity(ctx context.Context, orthopedicSurgeryActivityID string) (*models.OrthopedicSurgeriesActivityReview, error) {
	// We only look for one here since we expect there to be one review per activity, even if db permits 1 - many from activity to assessments.
	// TODO: evaluate above.
	reviews, err := models.OrthopedicSurgeriesActivityReviews(models.OrthopedicSurgeriesActivityReviewWhere.OrthopedicSurgeryActivityID.EQ(orthopedicSurgeryActivityID),
		Load(Rels(models.OrthopedicSurgeriesActivityReviewRels.OrthopedicSurgeriesActivityReviewSurgeries, Rels(models.OrthopedicSurgeriesActivityReviewSurgeryRels.Surgery, models.SurgeryRels.Method))),
		Load(Rels(models.OrthopedicSurgeriesActivityReviewRels.OrthopedicSurgeriesActivityReviewSurgeries, Rels(models.OrthopedicSurgeriesActivityReviewSurgeryRels.Surgery, models.SurgeryRels.Diagnose))),
		Load(models.OrthopedicSurgeriesActivityReviewRels.Resident),
		Load(models.OrthopedicSurgeriesActivityReviewRels.Supervisor),
	).All(ctx, r.db)
	if err != nil {
		return nil, err
	}
	if len(reviews) > 0 {
		return reviews[0], nil
	} else {
		return nil, errors.New("No Reviews Found")
	}
}
func (r *ReviewsRepo) UpdateSurgeriesForReview(ctx context.Context, executor boil.ContextExecutor, reviewId string, surgeries []*models.Surgery) error {
	db := executor
	if executor == nil {
		db = r.db
	}
	_, err := models.OrthopedicSurgeriesActivityReviewSurgeries(models.OrthopedicSurgeriesActivityReviewSurgeryWhere.OrthopedicSurgeriesActivityReviewID.EQ(reviewId)).DeleteAll(ctx, db)
	if err != nil {
		return err
	}
	return r.AddSurgeriesToReview(ctx, db, reviewId, surgeries)
}

func (r *ReviewsRepo) AddSurgeriesToReview(ctx context.Context, executor boil.ContextExecutor, reviewId string, surgeries []*models.Surgery) error {
	db := executor
	if executor == nil {
		db = r.db
	}
	review, err := models.FindOrthopedicSurgeriesActivityReview(ctx, db, reviewId)
	if err != nil {
		return err
	}
	for _, surgery := range surgeries {
		if err = review.AddOrthopedicSurgeriesActivityReviewSurgeries(ctx, db, true, &models.OrthopedicSurgeriesActivityReviewSurgery{
			SurgeryID:                           surgery.ID,
			OrthopedicSurgeriesActivityReviewID: review.ID,
		}); err != nil {
			return err
		}
	}
	return nil
}
