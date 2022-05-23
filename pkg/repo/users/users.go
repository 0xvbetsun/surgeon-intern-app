//go:generate mockery --name IRepo
package users

import (
	"context"
	"database/sql"

	"github.com/friendsofgo/errors"
	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/repoerrors"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/roles"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	_ "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

const GinContextUserIDKey = "UID"
const GinContextKey = "GinContextKey"

type (
	IRepo interface {
		Add(ctx context.Context, user *models.User) (*models.User, error)
		Delete(ctx context.Context, id string) error
		Update(ctx context.Context, user *models.User) (*models.User, error)
		GetByID(ctx context.Context, userId string, loadOrganizationalUnitRoles bool) (*models.User, error)
		GetInactiveByActivationCode(ctx context.Context, activationCode string, loadOrganizationalUnitRoles bool) (*models.User, error)
		GetByAuthenticationContext(ctx context.Context) (*models.User, error)
		GetByExternalID(ctx context.Context, fbID string) (*models.User, error)
		ListByAuthorizedUserAndFilter(ctx context.Context, authorizedUserId string, filter UserRepoFilter) ([]*UserUnitRole, error)
		ListByFilter(ctx context.Context, filter UserRepoFilter) ([]*UserWithUnitRoles, error)
		ListAll(ctx context.Context) ([]*models.User, error)
		GetUserDistinctRoles(ctx context.Context, userId string) ([]*models.Role, error)
	}
	Repo struct {
		db          *sql.DB
		userOrgRepo *UserOrganizationalUnitRoleRepo
		roleRepo    roles.IRepo
	}
	UserUnitRole struct {
		models.User               `boil:"users,bind"`
		models.Role               `boil:"roles,bind"`
		models.OrganizationalUnit `boil:"organizational_units,bind"`
	}

	UserWithUnitRoles struct {
		User      *models.User
		UnitRoles []*UnitRole
	}

	UnitRole struct {
		Role *models.Role
		Unit *models.OrganizationalUnit
	}

	UserUnitRoles []*UserUnitRole

	UserRepoFilter struct {
		UnitIds     []string
		RoleNames   []roles.RoleName
		UnitTypeIds []string
		Page        *int
		PerPage     *int
	}
)

func NewRepo(db *sql.DB, userOrgRepo *UserOrganizationalUnitRoleRepo, roleRepo roles.IRepo) IRepo {
	return &Repo{db: db, userOrgRepo: userOrgRepo, roleRepo: roleRepo}
}

func (r *Repo) GetInactiveByActivationCode(ctx context.Context, activationCode string, loadOrganizationalUnitRoles bool) (*models.User, error) {
	querymods := make([]QueryMod, 0)
	if loadOrganizationalUnitRoles {
		querymods = append(querymods, Load(Rels(models.UserRels.UserOrganizationalUnitRoles, models.UserOrganizationalUnitRoleRels.Role)), Load(Rels(models.UserRels.UserOrganizationalUnitRoles, models.UserOrganizationalUnitRoleRels.Unit, models.OrganizationalUnitRels.Parent)))
	}
	querymods = append(querymods, models.UserWhere.Activationcode.EQ(null.StringFrom(activationCode)), models.UserWhere.Activated.EQ(false))
	user, err := models.Users(querymods...).One(ctx, r.db)
	if err != nil {
		return nil, repoerrors.ErrorFromDbError(err)
	}
	return user, nil
}

func (r *Repo) Delete(ctx context.Context, id string) error {
	user, err := models.FindUser(ctx, r.db, id)
	if err != nil {
		return repoerrors.ErrorFromDbError(err)
	}
	_, err = user.Delete(ctx, r.db)
	if err != nil {
		return repoerrors.ErrorFromDbError(err)
	}
	return nil
}

func (r *Repo) Update(ctx context.Context, user *models.User) (*models.User, error) {
	_, err := user.Update(ctx, r.db, boil.Infer())
	if err != nil {
		return nil, repoerrors.ErrorFromDbError(err)
	}
	return user, nil
}

func (r *Repo) ListByAuthorizedUserAndFilter(ctx context.Context, authorizedUserId string, filter UserRepoFilter) ([]*UserUnitRole, error) {
	authorizedUserRoles, err := r.GetUserDistinctRoles(ctx, authorizedUserId)
	if err != nil {
		return nil, err
	}
	if len(authorizedUserRoles) == 0 {
		return make([]*UserUnitRole, 0), nil
	}
	userUnitRoles := make(UserUnitRoles, 0)
	// Loop the roles available to authorized user
	for _, authorizedRole := range authorizedUserRoles {
		// Find units (Clinics, Hospitals etc ) that the user has current role on.
		unitsWithRole, err := r.userOrgRepo.GetByFilter(ctx, UserOrganizationalUnitRoleFilter{
			UserId: &authorizedUserId,
			RoleId: &authorizedRole.ID,
			UnitId: nil,
		}, true)
		if err != nil {
			return nil, err
		}
		unitIds := funk.Map(unitsWithRole, func(unitWithRole *models.UserOrganizationalUnitRole) string {
			return unitWithRole.UnitID
		}).([]string)
		// intersect authorized user units with filterUnits to get a combined set
		if filter.UnitIds != nil {
			unitIds = funk.IntersectString(unitIds, filter.UnitIds)
		}
		// Get roles that is inherited via current role
		rolesToCheckFor, err := r.roleRepo.GetInheritedRoles(ctx, authorizedRole)
		if err != nil {
			return nil, err
		}
		roleIdsToCheckFor := funk.Map(rolesToCheckFor, func(role *models.Role) int {
			return role.ID
		}).([]int)
		if filter.RoleNames != nil {
			filterRoles, err := r.roleRepo.GetByNames(ctx, filter.RoleNames)
			if err != nil {
				return nil, err
			}
			filterRoleIds := funk.Map(filterRoles, func(role *models.Role) int {
				return role.ID
			}).([]int)
			// Get combined set of inherited roles and filter roles
			roleIdsToCheckFor = funk.InnerJoinInt(filterRoleIds, roleIdsToCheckFor)
		}
		var userUnitRole []*UserUnitRole
		err = models.NewQuery(
			Select(models.UserTableColumns.ID, models.UserTableColumns.UserExternalID, models.UserTableColumns.DisplayName,
				models.RoleTableColumns.ID, models.RoleTableColumns.Name, models.RoleTableColumns.DisplayName,
				models.OrganizationalUnitTableColumns.ID, models.OrganizationalUnitTableColumns.DisplayName,
				models.OrganizationalUnitTableColumns.ParentID, models.OrganizationalUnitTableColumns.TypeID),
			From(models.TableNames.UserOrganizationalUnitRoles),
			InnerJoin(models.TableNames.Users+" ON "+models.UserTableColumns.ID+" = "+models.UserOrganizationalUnitRoleTableColumns.UserID),
			InnerJoin(models.TableNames.Roles+" ON "+models.RoleTableColumns.ID+" = "+models.UserOrganizationalUnitRoleTableColumns.RoleID),
			InnerJoin(models.TableNames.OrganizationalUnits+" ON "+models.OrganizationalUnitTableColumns.ID+" = "+models.UserOrganizationalUnitRoleTableColumns.UnitID),
			models.UserOrganizationalUnitRoleWhere.UnitID.IN(unitIds),
			models.UserOrganizationalUnitRoleWhere.RoleID.IN(roleIdsToCheckFor),
		).Bind(ctx, r.db, &userUnitRole)
		if err != nil {
			return nil, err
		}
		userUnitRoles = append(userUnitRoles, userUnitRole...)
	}
	return userUnitRoles, nil
}

func (r *Repo) ListByFilter(ctx context.Context, filter UserRepoFilter) ([]*UserWithUnitRoles, error) {
	querymods := make([]QueryMod, 0)

	if filter.UnitIds != nil && len(filter.UnitIds) > 0 {
		querymods = append(querymods,
			Load(Rels(models.UserRels.UserOrganizationalUnitRoles, models.UserOrganizationalUnitRoleRels.Unit, models.OrganizationalUnitRels.Parent),
				models.UserOrganizationalUnitRoleWhere.UnitID.IN(filter.UnitIds)),
		)
	} else {
		querymods = append(querymods,
			Load(Rels(models.UserRels.UserOrganizationalUnitRoles, models.UserOrganizationalUnitRoleRels.Unit, models.OrganizationalUnitRels.Parent)),
		)
	}
	if filter.RoleNames != nil && len(filter.RoleNames) > 0 {
		filterRoles, err := r.roleRepo.GetByNames(ctx, filter.RoleNames)
		if err != nil {
			return nil, err
		}
		filterRoleIds := funk.Map(filterRoles, func(role *models.Role) int {
			return role.ID
		}).([]int)
		querymods = append(querymods,
			Load(Rels(models.UserRels.UserOrganizationalUnitRoles, models.UserOrganizationalUnitRoleRels.Role),
				models.UserOrganizationalUnitRoleWhere.RoleID.IN(filterRoleIds)),
		)
	} else {
		querymods = append(querymods,
			Load(Rels(models.UserRels.UserOrganizationalUnitRoles, models.UserOrganizationalUnitRoleRels.Role)),
		)
	}

	if filter.Page != nil && filter.PerPage != nil {
		querymods = append(querymods, Limit(*filter.PerPage), Offset((*filter.Page-1)**filter.PerPage))
	} else {
		querymods = append(querymods, Limit(20))
	}

	userUnitRoles := make([]*UserWithUnitRoles, 0)
	users, err := models.Users(querymods...).All(ctx, r.db)
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		userWithUnitRoles := &UserWithUnitRoles{
			User: user,
		}
		unitRoles := make([]*UnitRole, 0)
		if user.R.UserOrganizationalUnitRoles != nil {
			for _, our := range user.R.UserOrganizationalUnitRoles {
				if our.R.Unit != nil && our.R.Role != nil {
					unitRoles = append(unitRoles, &UnitRole{
						Role: our.R.Role,
						Unit: our.R.Unit,
					})
				}
			}
		}
		userWithUnitRoles.UnitRoles = unitRoles
		userUnitRoles = append(userUnitRoles, userWithUnitRoles)
	}
	return userUnitRoles, nil
}

func (r *Repo) GetUserDistinctRoles(ctx context.Context, userId string) ([]*models.Role, error) {
	userRoles, err := models.Roles(Distinct(models.TableNames.Roles+".*"),
		InnerJoin("user_organizational_unit_roles on user_organizational_unit_roles.role_id = roles.id"),
		models.UserOrganizationalUnitRoleWhere.UserID.EQ(userId),
	).All(ctx, r.db)

	if err != nil {
		if err == sql.ErrNoRows {
			return make([]*models.Role, 0), nil
		}
		return nil, repoerrors.ErrorFromDbError(err)
	}
	return userRoles, nil
}

func (r *Repo) GetByAuthenticationContext(ctx context.Context) (*models.User, error) {
	ginContext := ctx.Value(GinContextKey).(*gin.Context)
	if ginContext == nil {
		return nil, errors.New("Gin context is missing.")
	}
	uid := ginContext.Value(GinContextUserIDKey)
	if uid == nil {
		return nil, repoerrors.NewNotFoundError(uid.(string))
	}
	user, err := r.GetByExternalID(ctx, uid.(string))
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, repoerrors.NewNotFoundError(uid.(string))
	}
	return user, nil
}

func (r *Repo) GetByExternalID(ctx context.Context, fbID string) (*models.User, error) {
	user, err := models.Users(Where("user_external_id = ?", fbID)).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repoerrors.NewNotFoundErrorWithCause(fbID, err)
		}
		return nil, repoerrors.ErrorFromDbError(err)
	}
	return user, nil
}

func (r *Repo) Add(ctx context.Context, user *models.User) (*models.User, error) {
	err := user.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return nil, repoerrors.ErrorFromDbError(err)
	}

	return user, nil
}

func (r *Repo) GetByID(ctx context.Context, userId string, loadOrganizationalUnitRoles bool) (*models.User, error) {
	querymods := make([]QueryMod, 0)
	querymods = append(querymods, models.UserWhere.ID.EQ(userId))

	if loadOrganizationalUnitRoles {
		querymods = append(querymods, Load(Rels(models.UserRels.UserOrganizationalUnitRoles, models.UserOrganizationalUnitRoleRels.Role)), Load(Rels(models.UserRels.UserOrganizationalUnitRoles, models.UserOrganizationalUnitRoleRels.Unit, models.OrganizationalUnitRels.Parent)))
	}

	user, err := models.Users(querymods...).One(ctx, r.db)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repoerrors.NewNotFoundErrorWithCause(userId, err)
		}
		return nil, repoerrors.ErrorFromDbError(err)
	}
	return user, nil
}

func (r *Repo) ListAll(ctx context.Context) ([]*models.User, error) {
	users, err := models.Users().All(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return make([]*models.User, 0), nil
		}
		return nil, repoerrors.ErrorFromDbError(err)
	}
	return users, nil
}
