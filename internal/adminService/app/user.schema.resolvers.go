package adminService

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/vbetsun/surgeon-intern-app/internal/adminService/graph/generated"
	qlmodel "github.com/vbetsun/surgeon-intern-app/internal/adminService/graph/model"
	"github.com/vbetsun/surgeon-intern-app/internal/adminService/pkg/services/clinics"
	"github.com/vbetsun/surgeon-intern-app/internal/adminService/pkg/services/hospitals"
	"github.com/vbetsun/surgeon-intern-app/internal/adminService/pkg/services/users"
	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/internal/pkg/commonErrors"
	"github.com/vbetsun/surgeon-intern-app/internal/pkg/middleware"
	users2 "github.com/vbetsun/surgeon-intern-app/internal/pkg/services/users"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/roles"
	"go.uber.org/zap"
)

func (r *mutationResolver) UpdateUser(ctx context.Context, updateUser qlmodel.UpdateUserInput) (*qlmodel.DetailedUser, error) {
	if containsRole, _ := middleware.ClaimsViaContextContainsRole(ctx, roles.OgbookAdmin.String()); containsRole {
		retUser, err := r.service.adminUserService.UpdateUser(ctx, &updateUser)
		if err != nil {
			return nil, err
		}
		return retUser, nil
	}

	return nil, commonErrors.NewUnauthorizedAccessError("")
}

func (r *mutationResolver) InviteUser(ctx context.Context, inviteUser qlmodel.InviteUserInput) (*qlmodel.DetailedUser, error) {
	if containsRole, _ := middleware.ClaimsViaContextContainsRole(ctx, roles.OgbookAdmin.String()); containsRole {
		retUser, err := r.service.adminUserService.InviteUser(ctx, inviteUser)
		if err != nil {
			return nil, err
		}
		return retUser, nil
	}

	return nil, commonErrors.NewUnauthorizedAccessError("")
}

func (r *mutationResolver) DeleteUser(ctx context.Context, userID string) (*string, error) {
	if containsRole, _ := middleware.ClaimsViaContextContainsRole(ctx, roles.OgbookAdmin.String()); containsRole {
		err := r.service.adminUserService.DeleteUser(ctx, userID)
		if err != nil {
			return nil, err
		}
		return &userID, nil
	}

	return nil, commonErrors.NewUnauthorizedAccessError("")
}

func (r *mutationResolver) ConnectUserToClinic(ctx context.Context, clinicRole commonModel.UserClinicRoleInput) (*commonModel.ClinicRole, error) {
	if containsRole, _ := middleware.ClaimsViaContextContainsRole(ctx, roles.OgbookAdmin.String()); containsRole {
		retClinicRole, err := r.service.clinicService.ConnectUserToClinicWithRole(ctx, clinicRole)
		if err != nil {
			return nil, err
		}

		x, err := r.service.UserService.ListByUserAndFilter(ctx, clinicRole.UserID, users2.UserServiceFilter{})
		if err != nil {
			return nil, err
		}
		zap.S().Info(x)

		return retClinicRole, nil
	}

	return nil, commonErrors.NewUnauthorizedAccessError("")
}

func (r *mutationResolver) DisconnectUserFromClinic(ctx context.Context, clinicRole commonModel.UserClinicRoleInput) (*int, error) {
	if containsRole, _ := middleware.ClaimsViaContextContainsRole(ctx, roles.OgbookAdmin.String()); containsRole {
		err := r.service.clinicService.DisconnectUserFromClinicWithRole(ctx, clinicRole)
		ret := 0
		if err != nil {
			return nil, err
		}
		return &ret, nil
	}

	return nil, commonErrors.NewUnauthorizedAccessError("")
}

func (r *queryResolver) Users(ctx context.Context) ([]*qlmodel.DetailedUser, error) {
	if containsRole, _ := middleware.ClaimsViaContextContainsRole(ctx, roles.OgbookAdmin.String()); containsRole {
		retUsers, err := r.service.adminUserService.GetByFilter(ctx, users.Filter{})
		if err != nil {
			return nil, err
		}
		return retUsers, nil
	}

	return nil, commonErrors.NewUnauthorizedAccessError("")
}

func (r *queryResolver) User(ctx context.Context, userID *string) (*qlmodel.DetailedUser, error) {
	if containsRole, _ := middleware.ClaimsViaContextContainsRole(ctx, roles.OgbookAdmin.String()); containsRole {
		if userID == nil {
			externalId, err := middleware.UserExternalIdFromContext(ctx)
			if err != nil {
				return nil, err
			}
			user, err := r.service.UserService.GetByExternalID(ctx, externalId)
			userID = &user.UserID
		}
		retUser, err := r.service.adminUserService.GetByID(ctx, *userID)
		if err != nil {
			return nil, err
		}
		return retUser, nil
	}

	return nil, commonErrors.NewUnauthorizedAccessError("")
}

func (r *queryResolver) Clinics(ctx context.Context) ([]*commonModel.Clinic, error) {
	if containsRole, _ := middleware.ClaimsViaContextContainsRole(ctx, roles.OgbookAdmin.String()); containsRole {
		retClinics, err := r.service.clinicService.GetByFilter(ctx, clinics.Filter{
			HospitalId: nil,
			ClinicName: nil,
		})
		if err != nil {
			return nil, err
		}
		return retClinics, nil
	}

	return nil, commonErrors.NewUnauthorizedAccessError("")
}

func (r *queryResolver) Hospitals(ctx context.Context) ([]*commonModel.Hospital, error) {
	if containsRole, _ := middleware.ClaimsViaContextContainsRole(ctx, roles.OgbookAdmin.String()); containsRole {
		retHospitals, err := r.service.hospitalService.GetByFilter(ctx, hospitals.Filter{
			HospitalName: nil,
		})
		if err != nil {
			return nil, err
		}
		return retHospitals, nil
	}

	return nil, commonErrors.NewUnauthorizedAccessError("")
}

func (r *queryResolver) Roles(ctx context.Context) ([]*commonModel.Role, error) {
	if containsRole, _ := middleware.ClaimsViaContextContainsRole(ctx, roles.OgbookAdmin.String()); containsRole {
		retRoles, err := r.service.roleService.List(ctx)
		if err != nil {
			return nil, err
		}
		return retRoles, nil
	}

	return nil, commonErrors.NewUnauthorizedAccessError("")
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *mutationResolver) ActivateUser(ctx context.Context, activateUserInput qlmodel.ActivateUserInput) (*qlmodel.DetailedUser, error) {
	resUser, err := r.service.adminUserService.ActivateUser(ctx, activateUserInput)
	if err != nil {
		return nil, err
	}
	return resUser, nil
}
func (r *queryResolver) UserByActivationCode(ctx context.Context, activationCode *string) (*qlmodel.DetailedUser, error) {
	panic(fmt.Errorf("not implemented"))
}
func (r *mutationResolver) CreateUser(ctx context.Context, createUser qlmodel.InviteUserInput) (*qlmodel.DetailedUser, error) {
	if containsRole, _ := middleware.ClaimsViaContextContainsRole(ctx, roles.OgbookAdmin.String()); containsRole {
		retUser, err := r.service.adminUserService.InviteUser(ctx, createUser)
		if err != nil {
			return nil, err
		}
		return retUser, nil
	}

	return nil, commonErrors.NewUnauthorizedAccessError("")
}
func (r *queryResolver) Test(ctx context.Context) ([]*commonModel.User, error) {
	panic(fmt.Errorf("not implemented"))
}
