package practicalactivitytypes

import (
	"context"
	"database/sql"

	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/repoerrors"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/specialties"
	"github.com/volatiletech/sqlboiler/v4/boil"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type (
	IRepo interface {
		Get(ctx context.Context, practicalActivityId int) (*models.PracticalActivityType, error)
		GetByAuthUserAndFilter(ctx context.Context, userId string, filter Filter) ([]*models.PracticalActivityType, error)
		All(ctx context.Context) ([]*models.PracticalActivityType, error)
		Add(ctx context.Context, practicalActivity *models.PracticalActivityType) (*models.PracticalActivityType, error)
		GetByType(ctx context.Context, activityType commonModel.ActivityType) (*models.PracticalActivityType, error)
	}
	Repo struct {
		db              *sql.DB
		specialtiesRepo specialties.IRepo
	}
	Filter struct {
		TypeNames []string
	}
)

func NewRepo(db *sql.DB, specialtiesRepo specialties.IRepo) IRepo {
	return &Repo{
		db:              db,
		specialtiesRepo: specialtiesRepo,
	}
}

func (r *Repo) GetByAuthUserAndFilter(ctx context.Context, userId string, filter Filter) ([]*models.PracticalActivityType, error) {
	querymods := make([]QueryMod, 0)
	userSpecialties, err := r.specialtiesRepo.GetUserSpecialties(ctx, userId)
	if err != nil {
		return nil, err
	}
	if len(userSpecialties) == 0 {
		return make([]*models.PracticalActivityType, 0), nil
	}
	specialtiesIds := make([]interface{}, 0)
	for _, userSpecialty := range userSpecialties {
		specialtiesIds = append(specialtiesIds, userSpecialty.ID)
	}
	querymods = append(querymods, Distinct(models.TableNames.PracticalActivityTypes+".*"),
		InnerJoin(models.TableNames.SpecialtiesActivityTypes+" on specialties_activity_types.activity_type_id = "+models.PracticalActivityTypeTableColumns.ID),
		WhereIn("specialties_activity_types.specialty_id in ?", specialtiesIds...))
	if filter.TypeNames != nil {
		querymods = append(querymods, models.PracticalActivityTypeWhere.Name.IN(filter.TypeNames))
	}
	types, err := models.PracticalActivityTypes(querymods...).
		All(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return make([]*models.PracticalActivityType, 0), nil
		}
		return nil, repoerrors.ErrorFromDbError(err)
	}
	return types, nil
}

func (r *Repo) Get(ctx context.Context, practicalActivityId int) (*models.PracticalActivityType, error) {
	practicalActivity, err := models.FindPracticalActivityType(ctx, r.db, practicalActivityId)
	if err != nil {
		return nil, err
	}
	return practicalActivity, nil
}

func (r *Repo) GetByType(ctx context.Context, activityType commonModel.ActivityType) (*models.PracticalActivityType, error) {
	practicalActivity, err := models.PracticalActivityTypes(models.PracticalActivityTypeWhere.Name.EQ(activityType.String())).One(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return practicalActivity, nil
}

func (r *Repo) All(ctx context.Context) ([]*models.PracticalActivityType, error) {
	rPracticalActivities, err := models.PracticalActivityTypes().All(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return rPracticalActivities, nil
}

func (r *Repo) Add(ctx context.Context, practicalActivity *models.PracticalActivityType) (*models.PracticalActivityType, error) {
	err := practicalActivity.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return nil, err
	}

	return practicalActivity, nil
}
