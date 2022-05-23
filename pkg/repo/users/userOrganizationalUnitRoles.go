package users

import (
	"context"
	"database/sql"

	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/repoerrors"
	_ "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type (
	UserOrganizationalUnitRoleRepo struct {
		db *sql.DB
	}
	UserOrganizationalUnitRoleFilter struct {
		UserId *string
		RoleId *int
		UnitId *string
	}
)

func NewUserOrganizationalUnitRoleRepo(db *sql.DB) *UserOrganizationalUnitRoleRepo {
	return &UserOrganizationalUnitRoleRepo{db: db}
}

func (r *UserOrganizationalUnitRoleRepo) GetByFilter(ctx context.Context, filter UserOrganizationalUnitRoleFilter, loadRelations bool) ([]*models.UserOrganizationalUnitRole, error) {
	querymods := make([]QueryMod, 0)
	if loadRelations {
		querymods = append(querymods,
			Load(models.UserOrganizationalUnitRoleRels.User),
			Load(models.UserOrganizationalUnitRoleRels.Unit),
			Load(models.UserOrganizationalUnitRoleRels.Role),
		)
	}
	if filter.UnitId != nil {
		querymods = append(querymods, models.UserOrganizationalUnitRoleWhere.UnitID.EQ(*filter.UnitId))
	}
	if filter.UserId != nil {
		querymods = append(querymods, models.UserOrganizationalUnitRoleWhere.UserID.EQ(*filter.UserId))
	}
	if filter.RoleId != nil {
		querymods = append(querymods, models.UserOrganizationalUnitRoleWhere.RoleID.EQ(*filter.RoleId))
	}
	UserOrgRoles, err := models.UserOrganizationalUnitRoles(querymods...).All(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return make([]*models.UserOrganizationalUnitRole, 0), nil
		}
		return nil, repoerrors.ErrorFromDbError(err)
	}
	return UserOrgRoles, nil
}
