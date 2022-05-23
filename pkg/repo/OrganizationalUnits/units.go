package OrganizationalUnits

import (
	"context"
	"database/sql"

	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/repoerrors"
	"github.com/volatiletech/null/v8"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type (
	Repo struct {
		db *sql.DB
	}
	OrganizationalUnitFilter struct {
		TypeID      *string
		ParentID    *string
		DisplayName *string
		UserFilter  *UserFilter
	}
	UserFilter struct {
		UserID string
		RoleID *int
	}
)

func NewRepo(db *sql.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (o *Repo) ListByFilter(ctx context.Context, filter OrganizationalUnitFilter) ([]*models.OrganizationalUnit, error) {
	var queryMods = make([]QueryMod, 0)
	if filter.UserFilter != nil {
		queryMods = append(queryMods,
			InnerJoin("user_organizational_unit_roles on user_organizational_unit_roles.unit_id = organizational_units.id"))
		queryMods = append(queryMods, models.UserOrganizationalUnitRoleWhere.UserID.EQ(filter.UserFilter.UserID))
		if filter.UserFilter.RoleID != nil {
			queryMods = append(queryMods, models.UserOrganizationalUnitRoleWhere.RoleID.EQ(*filter.UserFilter.RoleID))
		}
	}
	if filter.TypeID != nil {
		queryMods = append(queryMods, models.OrganizationalUnitWhere.TypeID.EQ(*filter.TypeID))
	}
	if filter.ParentID != nil {
		queryMods = append(queryMods, models.OrganizationalUnitWhere.ParentID.EQ(null.StringFromPtr(filter.ParentID)))
	}
	if filter.DisplayName != nil {
		queryMods = append(queryMods, models.OrganizationalUnitWhere.DisplayName.EQ(*filter.DisplayName))
	}
	queryMods = append(queryMods, Load(models.OrganizationalUnitRels.Parent))
	queryMods = append(queryMods, Load(models.OrganizationalUnitRels.Specialties))
	units, err := models.OrganizationalUnits(queryMods...).All(ctx, o.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return make([]*models.OrganizationalUnit, 0), nil
		}
		return nil, repoerrors.ErrorFromDbError(err)
	}
	return units, nil
}

func (o *Repo) ConnectUnitToSpecialty(ctx context.Context, unit *models.OrganizationalUnit, specialty *models.Specialty) error {
	if err := unit.AddSpecialties(ctx, o.db, false, specialty); err != nil {
		return err
	}
	return nil
}
