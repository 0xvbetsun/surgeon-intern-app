package users

import (
	"context"

	"github.com/friendsofgo/errors"
	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	thirdparty2 "github.com/vbetsun/surgeon-intern-app/internal/pkg/thirdparty"
	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/OrganizationalUnits"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/roles"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/users"
)

const residentRoleId = 1
const supervisorRoleId = 2

type (
	IService interface {
		GetById(ctx context.Context, userId string) (*commonModel.User, error)
		GetByExternalID(ctx context.Context, externalID string) (*commonModel.User, error)
		ListByUserAndFilter(ctx context.Context, authorizedUserId string, filter UserServiceFilter) ([]*commonModel.User, error)
		ListAll(ctx context.Context) ([]*commonModel.User, error)
		ListSupervisorsByClinicId(ctx context.Context, clinicId string) ([]*commonModel.User, error)
		ListAllSupervisors(ctx context.Context) ([]*commonModel.User, error)
		ListResidentsByClinicId(ctx context.Context, clinicId string) ([]*commonModel.User, error)
		ListAllResidents(ctx context.Context) ([]*commonModel.User, error)
	}
	Service struct {
		repo                        users.IRepo
		organizationalUnitsRepo     *OrganizationalUnits.Repo
		organizationalUnitTypesRepo *OrganizationalUnits.TypesRepo
		userOrgUnitRoleRepo         *users.UserOrganizationalUnitRoleRepo
		roleRepo                    roles.IRepo
		auth0                       *thirdparty2.Auth0
	}
	UserServiceFilter struct {
		OrganizationalUnitID *string // i.e Clinic, Hospital
		IsSupervisor         *bool
		IsResident           *bool
	}
)

func NewService(repo users.IRepo, auth0 *thirdparty2.Auth0, orgUnitsrepo *OrganizationalUnits.Repo, orgUnitTypesRepo *OrganizationalUnits.TypesRepo, roleRepo roles.IRepo) IService {
	return &Service{repo: repo, auth0: auth0, organizationalUnitsRepo: orgUnitsrepo, organizationalUnitTypesRepo: orgUnitTypesRepo, roleRepo: roleRepo}
}

func (s *Service) GetByExternalID(ctx context.Context, fbID string) (*commonModel.User, error) {
	rUser, err := s.repo.GetByExternalID(ctx, fbID)
	if err != nil {
		return nil, err
	}
	user := s.mapUserGraphQlModel(rUser)
	return user, nil
}

func (s *Service) GetById(ctx context.Context, userID string) (*commonModel.User, error) {
	rUser, err := s.repo.GetByID(ctx, userID, false)
	if err != nil {
		return nil, err
	}
	user := s.mapUserGraphQlModel(rUser)
	return user, nil
}

func (s *Service) ListByUserAndFilter(ctx context.Context, authorizedUserId string, filter UserServiceFilter) ([]*commonModel.User, error) {
	var unitIds []string
	var roleNames []roles.RoleName
	if filter.OrganizationalUnitID != nil {
		unitIds = []string{*filter.OrganizationalUnitID}
	}
	if filter.IsResident != nil && *filter.IsResident {
		roleNames = append(roleNames, roles.RESIDENT)
	}
	if filter.IsSupervisor != nil && *filter.IsSupervisor {
		roleNames = append(roleNames, roles.SUPERVISOR)
	}
	if roleNames != nil && len(roleNames) == 0 {
		return nil, errors.New("bad filter")
	}
	if unitIds != nil && len(unitIds) == 0 {
		return nil, errors.New("bad filter")
	}
	userUnitRoles, err := s.repo.ListByAuthorizedUserAndFilter(ctx, authorizedUserId, users.UserRepoFilter{
		UnitIds:   unitIds,
		RoleNames: roleNames,
	})
	if err != nil {
		return nil, err
	}
	qlUsers := commonModel.UsersFromUserUnitRoles(userUnitRoles)
	return qlUsers, nil
}

func (s *Service) ListAll(ctx context.Context) ([]*commonModel.User, error) {
	rUsers, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	p := make([]*commonModel.User, 0)
	for _, rUser := range rUsers {
		user := s.mapUserGraphQlModel(rUser)
		p = append(p, user)
	}
	return p, nil
}

func (s *Service) ListSupervisorsByClinicId(ctx context.Context, clinicId string) ([]*commonModel.User, error) {
	roleId := supervisorRoleId
	rUserClinicRoles, err := s.userOrgUnitRoleRepo.GetByFilter(ctx, users.UserOrganizationalUnitRoleFilter{
		UserId: nil,
		RoleId: &roleId,
		UnitId: &clinicId,
	}, true)
	if err != nil {
		return nil, err
	}
	p := make([]*commonModel.User, 0)
	for _, userClinicRole := range rUserClinicRoles {
		supervisor := s.mapSupervisorGraphQlModel(userClinicRole)
		p = append(p, supervisor)
	}
	return p, nil
}

func (s *Service) ListAllSupervisors(ctx context.Context) ([]*commonModel.User, error) {
	roleId := supervisorRoleId
	rUserClinicRoles, err := s.userOrgUnitRoleRepo.GetByFilter(ctx, users.UserOrganizationalUnitRoleFilter{
		UserId: nil,
		RoleId: &roleId,
		UnitId: nil,
	}, true)
	if err != nil {
		return nil, err
	}
	p := make([]*commonModel.User, 0)
	for _, userClinicRole := range rUserClinicRoles {
		supervisor := s.mapSupervisorGraphQlModel(userClinicRole)
		p = append(p, supervisor)
	}
	return p, nil
}

func (s *Service) ListResidentsByClinicId(ctx context.Context, clinicId string) ([]*commonModel.User, error) {
	roleId := residentRoleId
	rUserClinicRoles, err := s.userOrgUnitRoleRepo.GetByFilter(ctx, users.UserOrganizationalUnitRoleFilter{
		UserId: nil,
		RoleId: &roleId,
		UnitId: &clinicId,
	}, true)
	if err != nil {
		return nil, err
	}
	p := make([]*commonModel.User, 0)
	for _, userClinicRole := range rUserClinicRoles {
		supervisor := s.mapResidentGraphQlModel(userClinicRole)
		p = append(p, supervisor)
	}
	return p, nil
}

func (s *Service) ListAllResidents(ctx context.Context) ([]*commonModel.User, error) {
	roleId := residentRoleId
	rUserClinicRoles, err := s.userOrgUnitRoleRepo.GetByFilter(ctx, users.UserOrganizationalUnitRoleFilter{
		UserId: nil,
		RoleId: &roleId,
		UnitId: nil,
	}, true)
	if err != nil {
		return nil, err
	}
	p := make([]*commonModel.User, 0)
	for _, userClinicRole := range rUserClinicRoles {
		supervisor := s.mapResidentGraphQlModel(userClinicRole)
		p = append(p, supervisor)
	}
	return p, nil
}

func (s Service) mapResidentGraphQlModel(userClinicRole *models.UserOrganizationalUnitRole) *commonModel.User {
	supervisor := &commonModel.User{
		DisplayName: userClinicRole.R.User.DisplayName,
		UserID:      userClinicRole.UserID,
		ClinicIds:   nil,
	}
	return supervisor
}

func (s Service) mapSupervisorGraphQlModel(userClinicRole *models.UserOrganizationalUnitRole) *commonModel.User {
	supervisor := &commonModel.User{
		DisplayName: userClinicRole.R.User.DisplayName,
		UserID:      userClinicRole.UserID,
		ClinicIds:   nil,
	}
	return supervisor
}

func (s Service) mapUserGraphQlModel(rUser *models.User) *commonModel.User {
	user := &commonModel.User{
		UserID:      rUser.ID,
		DisplayName: rUser.DisplayName,
	}
	return user
}
