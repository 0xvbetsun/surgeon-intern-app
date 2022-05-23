//go:generate mockery --name IService
package authorization

import (
	"context"

	"github.com/casbin/casbin/v2"
	casbinpgadapter "github.com/cychiuae/casbin-pg-adapter"
	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/OrganizationalUnits"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/users"
)

type (
	IService interface {
		Authorize(requests []CasbinRequest, permissive bool) (bool, error)
		UserIsSuperAdmin(userID string) bool
		GetAvailableRolesForUser(ctx context.Context, userID string) ([]*commonModel.Role, error)
	}

	Service struct {
		Enforcer                    *casbin.Enforcer
		userRepo                    users.IRepo
		organizationalUnitsRepo     *OrganizationalUnits.Repo
		organizationalUnitTypesRepo *OrganizationalUnits.TypesRepo
	}
)

type CasbinConfigFilePath string

func NewService(adapter *casbinpgadapter.Adapter,
	userRepo users.IRepo,
	config CasbinConfigFilePath,
	orgUnitsRepo *OrganizationalUnits.Repo,
	orgUnitsTypesRepo *OrganizationalUnits.TypesRepo) (IService, error) {
	enforcer, err := casbin.NewEnforcer(string(config), adapter)
	if err != nil {
		return nil, err
	}
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, err
	}
	enforcer.AddPolicy("admin", "*", "organisations", "read")
	enforcer.AddRoleForUserInDomain("mExrbuaFhDVlzTvwVf4VIPtKJWj1", "admin", "*")
	enforcer.SavePolicy()
	//ok := enforcer.GetRolesForUserInDomain("mExrbuaFhDVlzTvwVf4VIPtKJWj1", "*")
	enforcer.EnableLog(true)
	return &Service{Enforcer: enforcer, userRepo: userRepo, organizationalUnitsRepo: orgUnitsRepo}, nil
}

func (s *Service) Authorize(requests []CasbinRequest, permissive bool) (bool, error) {
	if err := s.Enforcer.LoadPolicy(); err != nil {
		return false, err
	}

	for _, request := range requests {
		authorized, err := s.Enforcer.Enforce(request.Subject, request.Domain, request.Object, request.Action)
		if err != nil {
			return false, err
		}
		if permissive && authorized {
			return authorized, nil
		} else if authorized {
			continue
		} else {
			return false, nil
		}
	}
	return true, nil
}
func (s *Service) UserIsSuperAdmin(userID string) bool {
	roles := s.Enforcer.GetRolesForUserInDomain(userID, "*")
	return contains(roles, "admin")
}

func (s *Service) GetAvailableRolesForUser(ctx context.Context, userID string) ([]*commonModel.Role, error) {
	roles, err := s.userRepo.GetUserDistinctRoles(ctx, userID)
	if err != nil {
		return nil, err
	}
	return commonModel.RolesFromRepoTypes(roles), nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
