package roles

import (
	"context"

	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/roles"
)

type (
	RoleService struct {
		repo roles.IRepo
	}
)

func NewRoleService(repo roles.IRepo) *RoleService {
	return &RoleService{
		repo: repo,
	}
}

func (r *RoleService) List(ctx context.Context) ([]*commonModel.Role, error) {
	roleNames := []roles.RoleName{roles.RESIDENT, roles.SUPERVISOR, roles.DIRECTOR}
	rRoles, err := r.repo.GetByNames(ctx, roleNames)
	if err != nil {
		return nil, err
	}
	return fromRepoRoles(rRoles), nil
}
