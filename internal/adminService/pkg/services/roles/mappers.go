package roles

import (
	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/pkg/models"
)

func fromRepoRoles(roles []*models.Role) []*commonModel.Role {
	rRoles := make([]*commonModel.Role, 0)
	for _, role := range roles {
		rRoles = append(rRoles, fromRepoRole(role))
	}
	return rRoles
}

func fromRepoRole(role *models.Role) *commonModel.Role {
	return &commonModel.Role{
		RoleIdentifier: role.Name,
		DisplayName:    role.DisplayName,
	}
}
