package adminService

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/internal/pkg/commonErrors"
	"github.com/vbetsun/surgeon-intern-app/internal/pkg/middleware"
)

func HasAtLeastRole(ctx context.Context, obj interface{}, next graphql.Resolver, role commonModel.RoleInput) (res interface{}, err error) {
	if containsRole, _ := middleware.ClaimsViaContextContainsRole(ctx, role.RoleIdentifier); containsRole {
		return next(ctx)
	}
	return nil, commonErrors.NewUnauthorizedAccessError("")
}

func HasOneOfRoles(ctx context.Context, obj interface{}, next graphql.Resolver, roles []*commonModel.RoleInput) (res interface{}, err error) {
	for _, acceptedRole := range roles {
		if containsRole, _ := middleware.ClaimsViaContextContainsRole(ctx, acceptedRole.RoleIdentifier); containsRole {
			return next(ctx)
		}
	}
	return nil, commonErrors.NewUnauthorizedAccessError("")
}
