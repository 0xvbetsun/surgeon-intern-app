package logbookService

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/internal/pkg/middleware"
	authorization2 "github.com/vbetsun/surgeon-intern-app/internal/pkg/services/authorization"
	"github.com/vbetsun/surgeon-intern-app/internal/pkg/services/users"
)

type (
	AuthDirectives struct {
		authorizationService authorization2.IService
		userService          users.IService
	}
)

func NewAuthDirectives(authorizationService authorization2.IService, userService users.IService) *AuthDirectives {
	return &AuthDirectives{authorizationService: authorizationService, userService: userService}
}

func (a *AuthDirectives) HasAtLeastRole(ctx context.Context, obj interface{}, next graphql.Resolver, role commonModel.RoleInput) (res interface{}, err error) {
	// Get user
	user, err := middleware.UserFromContext(ctx, a.userService)
	if err != nil {
		return nil, fmt.Errorf("Access denied")
	}

	roles, err := a.authorizationService.GetAvailableRolesForUser(ctx, user.UserID)
	for _, userRole := range roles {
		if userRole.RoleIdentifier == role.RoleIdentifier {
			return next(ctx)
		}
	}
	return nil, fmt.Errorf("Access denied")
}

func (a *AuthDirectives) HasOneOfRoles(ctx context.Context, obj interface{}, next graphql.Resolver, roles []*commonModel.RoleInput) (res interface{}, err error) {
	// Get user
	user, err := middleware.UserFromContext(ctx, a.userService)
	if err != nil {
		return nil, fmt.Errorf("Access denied")
	}

	userRoles, err := a.authorizationService.GetAvailableRolesForUser(ctx, user.UserID)
	for _, userRole := range userRoles {
		for _, acceptedRole := range roles {
			if acceptedRole.RoleIdentifier == userRole.RoleIdentifier {
				return next(ctx)
			}
		}
	}
	return nil, fmt.Errorf("Access denied")
}
