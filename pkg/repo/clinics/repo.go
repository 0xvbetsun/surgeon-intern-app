//go:generate mockery --name IRepo
package clinics

import (
	"context"
	"database/sql"

	"github.com/vbetsun/surgeon-intern-app/pkg/repo/OrganizationalUnits"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/repoerrors"
	"github.com/volatiletech/null/v8"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type (
	IRepo interface {
		Add(ctx context.Context, clinic *models.OrganizationalUnit) (*models.OrganizationalUnit, error)
		AddUserClinicRole(ctx context.Context, userClinicRole *models.UserOrganizationalUnitRole) (*models.UserOrganizationalUnitRole, error)
		RemoveUserClinicRole(ctx context.Context, userClinicRole *models.UserOrganizationalUnitRole) error
		GetByID(ctx context.Context, clinicId string) (*models.OrganizationalUnit, error)
		ListByHospitalId(ctx context.Context, hospitalId string) ([]*models.OrganizationalUnit, error)
		ListAll(ctx context.Context) ([]*models.OrganizationalUnit, error)
		GetClinicDepartments(ctx context.Context, clinicID string) ([]*models.OrganizationalUnit, error)
		GetClinicDepartmentByID(ctx context.Context, clinicDepartmentId string) (*models.OrganizationalUnit, error)
		AddClinicDepartment(ctx context.Context, location *models.OrganizationalUnit) (*models.OrganizationalUnit, error)
	}
	Repo struct {
		db       *sql.DB
		typeRepo *OrganizationalUnits.TypesRepo
	}
)

func NewRepo(db *sql.DB, typeRepo *OrganizationalUnits.TypesRepo) IRepo {
	return &Repo{db: db, typeRepo: typeRepo}
}

func (r *Repo) Add(ctx context.Context, clinic *models.OrganizationalUnit) (*models.OrganizationalUnit, error) {
	clinicType, err := r.typeRepo.GetByName(ctx, OrganizationalUnits.CLINIC)
	if err != nil {
		return nil, err
	}
	clinic.TypeID = clinicType.ID
	err = clinic.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return nil, repoerrors.ErrorFromDbError(err)
	}
	return clinic, nil
}

func (r *Repo) AddUserClinicRole(ctx context.Context, userClinicRole *models.UserOrganizationalUnitRole) (*models.UserOrganizationalUnitRole, error) {
	_, err := r.GetByID(ctx, userClinicRole.UnitID)
	if err != nil {
		return nil, repoerrors.NewBadArgumentError("userClinicRole.UnitId")
	}
	err = userClinicRole.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return nil, repoerrors.ErrorFromDbError(err)
	}

	return userClinicRole, nil
}

func (r *Repo) RemoveUserClinicRole(ctx context.Context, userClinicRole *models.UserOrganizationalUnitRole) error {
	uor, err := models.UserOrganizationalUnitRoles(
		models.UserOrganizationalUnitRoleWhere.UserID.EQ(userClinicRole.UserID),
		models.UserOrganizationalUnitRoleWhere.UnitID.EQ(userClinicRole.UnitID),
		models.UserOrganizationalUnitRoleWhere.RoleID.EQ(userClinicRole.RoleID),
	).One(ctx, r.db)
	if err != nil {
		return repoerrors.ErrorFromDbError(err)
	}
	_, err = uor.Delete(ctx, r.db)
	if err != nil {
		return repoerrors.ErrorFromDbError(err)
	}
	return nil
}

func (r *Repo) GetByID(ctx context.Context, clinicId string) (*models.OrganizationalUnit, error) {
	clinic, err := models.OrganizationalUnits(
		models.OrganizationalUnitWhere.ID.EQ(clinicId),
		Load(Rels(models.OrganizationalUnitRels.Type),
			Where("organizational_unit_types.name = ?", OrganizationalUnits.CLINIC)),
	).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repoerrors.NewNotFoundErrorWithCause(clinicId, err)
		}
		return nil, repoerrors.ErrorFromDbError(err)
	}
	return clinic, nil
}

// We assume that a hospital is the parent of a clinic
func (r *Repo) ListByHospitalId(ctx context.Context, hospitalId string) ([]*models.OrganizationalUnit, error) {
	clinics, err := models.OrganizationalUnits(models.OrganizationalUnitWhere.ParentID.EQ(null.StringFrom(hospitalId)), Load(Rels(models.OrganizationalUnitRels.Type),
		Where("organizational_unit_types.name = ?", OrganizationalUnits.CLINIC))).All(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return make([]*models.OrganizationalUnit, 0), nil
		}
		return nil, repoerrors.ErrorFromDbError(err)
	}
	return clinics, nil
}

func (r *Repo) ListAll(ctx context.Context) ([]*models.OrganizationalUnit, error) {
	clinics, err := models.OrganizationalUnits(Load(Rels(models.OrganizationalUnitRels.Type),
		Where("organizational_unit_types.name = ?", OrganizationalUnits.CLINIC))).All(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return empty list if no results.
			return make([]*models.OrganizationalUnit, 0), nil
		}
		return nil, repoerrors.ErrorFromDbError(err)
	}
	return clinics, nil
}

func (r *Repo) AddClinicDepartment(ctx context.Context, department *models.OrganizationalUnit) (*models.OrganizationalUnit, error) {
	departmentType, err := r.typeRepo.GetByName(ctx, OrganizationalUnits.CLINIC_DEPARTMENT)
	if err != nil {
		return nil, err
	}
	department.TypeID = departmentType.ID
	err = department.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return nil, repoerrors.ErrorFromDbError(err)
	}
	return department, nil
}

func (r *Repo) GetClinicDepartmentByID(ctx context.Context, clinicDepartmentId string) (*models.OrganizationalUnit, error) {
	location, err := models.FindOrganizationalUnit(ctx, r.db, clinicDepartmentId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repoerrors.NewNotFoundErrorWithCause(clinicDepartmentId, err)
		}
		return nil, repoerrors.ErrorFromDbError(err)
	}
	return location, nil
}

func (r *Repo) GetClinicDepartments(ctx context.Context, clinicID string) ([]*models.OrganizationalUnit, error) {
	locations, err := models.OrganizationalUnits(models.OrganizationalUnitWhere.ParentID.EQ(null.StringFrom(clinicID)),
		Load(Rels(models.OrganizationalUnitRels.Type),
			Where("organizational_unit_types.name = ?", OrganizationalUnits.CLINIC_DEPARTMENT))).All(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return empty list if no results.
			return make([]*models.OrganizationalUnit, 0), nil
		}
		return nil, repoerrors.ErrorFromDbError(err)
	}
	return locations, nil
}
