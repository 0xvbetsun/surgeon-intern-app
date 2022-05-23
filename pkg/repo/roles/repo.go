//go:generate mockery --name IRepo
package roles

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/friendsofgo/errors"
	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/repoerrors"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type (
	RoleName string
	IRepo    interface {
		Add(ctx context.Context, role *models.Role) (*models.Role, error)
		GetByID(ctx context.Context, roleId int) (*models.Role, error)
		GetByName(ctx context.Context, name RoleName) (*models.Role, error)
		GetByNames(ctx context.Context, names []RoleName) ([]*models.Role, error)
		GetInheritedRoles(ctx context.Context, role *models.Role) ([]*models.Role, error)
	}
	Repo struct {
		db *sql.DB
	}
)

func (n RoleName) String() string {
	return string(n)
}

func (n RoleName) Valid() bool {
	switch n {
	case RESIDENT, SUPERVISOR, DIRECTOR:
		return true
	}
	return false
}

func (n RoleName) GetInheritedRoleNames() []RoleName {
	switch n {
	case DIRECTOR:
		return []RoleName{DIRECTOR, SUPERVISOR, RESIDENT}
	case SUPERVISOR:
		return []RoleName{SUPERVISOR, RESIDENT}
	case RESIDENT:
		return []RoleName{RESIDENT}
	default:
		return []RoleName{}
	}
}

const RESIDENT RoleName = "Resident"
const SUPERVISOR RoleName = "Supervisor"
const DIRECTOR RoleName = "Director"
const OgbookAdmin RoleName = "OgbookAdmin"

func NewRepo(db *sql.DB) IRepo {
	return &Repo{db: db}
}

func (r *Repo) GetInheritedRoles(ctx context.Context, role *models.Role) ([]*models.Role, error) {
	inheritedRoleNames := RoleName(role.Name).GetInheritedRoleNames()
	roles, err := r.GetByNames(ctx, inheritedRoleNames)
	if err != nil {
		return nil, err
	}
	return roles, nil

}

func (r *Repo) GetByNames(ctx context.Context, roleNames []RoleName) ([]*models.Role, error) {
	names := make([]string, 0)
	for _, name := range roleNames {
		if !name.Valid() {
			return nil, errors.New("Invalid Role")
		}
		names = append(names, name.String())
	}
	roles, err := models.Roles(models.RoleWhere.Name.IN(names)).All(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repoerrors.NewNotFoundError(names)
		}
		return nil, repoerrors.ErrorFromDbError(err)
	}
	return roles, nil
}

func (r *Repo) GetByName(ctx context.Context, name RoleName) (*models.Role, error) {
	if !name.Valid() {
		return nil, repoerrors.NewNotFoundError(fmt.Sprintf("Role with name: %s", name))
	}
	role, err := models.Roles(models.RoleWhere.Name.EQ(name.String())).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repoerrors.NewNotFoundErrorWithCause(fmt.Sprintf("Role with name: %s", name), err)
		}
		return nil, repoerrors.ErrorFromDbError(err)
	}
	return role, nil
}

func (r *Repo) Add(ctx context.Context, role *models.Role) (*models.Role, error) {
	if !RoleName(role.Name).Valid() {
		return nil, errors.New("Bad Rolename")
	}
	err := role.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return nil, repoerrors.ErrorFromDbError(err)
	}

	return role, nil
}

func (r *Repo) GetByID(ctx context.Context, roleId int) (*models.Role, error) {
	role, err := models.FindRole(ctx, r.db, roleId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repoerrors.NewNotFoundErrorWithCause(roleId, err)
		}
		return nil, repoerrors.ErrorFromDbError(err)
	}
	return role, nil
}
