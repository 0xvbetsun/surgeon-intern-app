package specialties

import (
	"context"
	"database/sql"

	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/repoerrors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	_ "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type (
	IRepo interface {
		Add(ctx context.Context, executor boil.ContextExecutor, specialty *models.Specialty) error
		Get(ctx context.Context, specialtyId string) (*models.Specialty, error)
		GetUserSpecialties(ctx context.Context, userId string) ([]*models.Specialty, error)
		ConnectSpecialtyToActivityType(ctx context.Context, specialty *models.Specialty, activityType *models.PracticalActivityType) error
	}
	Repo struct {
		db *sql.DB
	}
)

func NewRepo(db *sql.DB) IRepo {
	return &Repo{db: db}
}

func (r *Repo) GetUserSpecialties(ctx context.Context, userId string) ([]*models.Specialty, error) {
	specialties, err := models.Specialties(
		InnerJoin(models.TableNames.OrganizationalUnitSpecialties+" ON organizational_unit_specialties.specialty_id = "+models.SpecialtyTableColumns.ID),
		InnerJoin(models.TableNames.UserOrganizationalUnitRoles+" ON "+models.UserOrganizationalUnitRoleTableColumns.UnitID+" = organizational_unit_specialties.unit_id"),
		models.UserOrganizationalUnitRoleWhere.UserID.EQ(userId)).All(ctx, r.db)
	if err != nil {
		return nil, repoerrors.ErrorFromDbError(err)
	}
	if specialties == nil {
		return make([]*models.Specialty, 0), nil
	}
	return specialties, nil
}

func (r *Repo) Add(ctx context.Context, executor boil.ContextExecutor, specialty *models.Specialty) error {
	db := executor
	if executor == nil {
		db = r.db
	}
	err := specialty.Insert(ctx, db, boil.Infer())
	if err != nil {
		return repoerrors.ErrorFromDbError(err)
	}
	return nil
}

func (r *Repo) Get(ctx context.Context, specialtyId string) (*models.Specialty, error) {
	specialty, err := models.FindSpecialty(ctx, r.db, specialtyId)
	if err != nil {
		return nil, err
	}
	return specialty, err
}

func (o *Repo) ConnectSpecialtyToActivityType(ctx context.Context, specialty *models.Specialty, activityType *models.PracticalActivityType) error {
	if err := specialty.AddActivityTypePracticalActivityTypes(ctx, o.db, false, activityType); err != nil {
		return err
	}
	return nil
}
