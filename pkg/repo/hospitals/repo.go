//go:generate mockery --name IRepo
package hospitals

import (
	"context"
	"database/sql"

	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/OrganizationalUnits"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/repoerrors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type (
	IRepo interface {
		Add(ctx context.Context, hospital *models.OrganizationalUnit) (*models.OrganizationalUnit, error)
		GetByID(ctx context.Context, hospitalId string) (*models.OrganizationalUnit, error)
		ListByOrganisationId(ctx context.Context, organisationId string) ([]*models.OrganizationalUnit, error)
		ListAll(ctx context.Context) ([]*models.OrganizationalUnit, error)
	}
	Repo struct {
		db       *sql.DB
		typeRepo *OrganizationalUnits.TypesRepo
	}
)

func NewRepo(db *sql.DB, typeRepo *OrganizationalUnits.TypesRepo) IRepo {
	return &Repo{db: db, typeRepo: typeRepo}
}

func (r *Repo) Add(ctx context.Context, hospital *models.OrganizationalUnit) (*models.OrganizationalUnit, error) {
	hospitalType, err := r.typeRepo.GetByName(ctx, OrganizationalUnits.HOSPITAL)
	if err != nil {
		return nil, err
	}
	hospital.TypeID = hospitalType.ID
	err = hospital.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return nil, repoerrors.ErrorFromDbError(err)
	}
	return hospital, nil
}

func (r *Repo) GetByID(ctx context.Context, hospitalId string) (*models.OrganizationalUnit, error) {
	clinic, err := models.OrganizationalUnits(
		models.OrganizationalUnitWhere.ID.EQ(hospitalId),
		Load(Rels(models.OrganizationalUnitRels.Type),
			Where("organizational_unit_types.name = ?", OrganizationalUnits.HOSPITAL)),
	).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repoerrors.NewNotFoundErrorWithCause(hospitalId, err)
		}
		return nil, repoerrors.ErrorFromDbError(err)
	}
	return clinic, nil
}

func (r *Repo) ListByOrganisationId(ctx context.Context, organisationId string) ([]*models.OrganizationalUnit, error) {
	hospitals, err := models.OrganizationalUnits(
		models.OrganizationalUnitWhere.ParentID.EQ(null.StringFrom(organisationId)),
		Load(Rels(models.OrganizationalUnitRels.Type),
			Where("organizational_unit_types.name = ?", OrganizationalUnits.HOSPITAL))).All(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return make([]*models.OrganizationalUnit, 0), nil
		}
		return nil, repoerrors.ErrorFromDbError(err)
	}
	return hospitals, nil
}

func (r *Repo) ListAll(ctx context.Context) ([]*models.OrganizationalUnit, error) {
	clinics, err := models.OrganizationalUnits(Load(Rels(models.OrganizationalUnitRels.Type),
		Where("organizational_unit_types.name = ?", OrganizationalUnits.HOSPITAL))).All(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return empty list if no results.
			return make([]*models.OrganizationalUnit, 0), nil
		}
		return nil, repoerrors.ErrorFromDbError(err)
	}
	return clinics, nil
}
