package clinics

import (
	"context"

	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/OrganizationalUnits"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/clinics"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/roles"
)

type (
	ClinicService struct {
		unitsRepo     *OrganizationalUnits.Repo
		unitTypesRepo *OrganizationalUnits.TypesRepo
		clinicsRepo   clinics.IRepo
		roleRepo      roles.IRepo
	}
	Filter struct {
		HospitalId *string
		ClinicName *string
	}
)

func NewClinicService(unitsRepo *OrganizationalUnits.Repo, typesRepo *OrganizationalUnits.TypesRepo, clinicsRepo clinics.IRepo, roleRepo roles.IRepo) *ClinicService {
	return &ClinicService{
		unitsRepo:     unitsRepo,
		unitTypesRepo: typesRepo,
		clinicsRepo:   clinicsRepo,
		roleRepo:      roleRepo,
	}
}

func (c *ClinicService) GetByFilter(ctx context.Context, filter Filter) ([]*commonModel.Clinic, error) {
	clinicUnitTypeId, err := c.unitTypesRepo.GetByName(ctx, OrganizationalUnits.CLINIC)
	if err != nil {
		return nil, err
	}
	clinics, err := c.unitsRepo.ListByFilter(ctx, OrganizationalUnits.OrganizationalUnitFilter{
		TypeID:      &clinicUnitTypeId.ID,
		ParentID:    filter.HospitalId,
		DisplayName: filter.ClinicName,
	})
	if err != nil {
		return nil, err
	}
	return FromRepoClinics(clinics), nil
}

func (c *ClinicService) ConnectUserToClinicWithRole(ctx context.Context, input commonModel.UserClinicRoleInput) (*commonModel.ClinicRole, error) {
	rRole, err := c.roleRepo.GetByName(ctx, roles.RoleName(input.RoleIdentifier))
	if err != nil {
		return nil, err
	}
	rOur, err := c.clinicsRepo.AddUserClinicRole(ctx, &models.UserOrganizationalUnitRole{
		UnitID: input.ClinicID,
		UserID: input.UserID,
		RoleID: rRole.ID,
	})
	if err != nil {
		return nil, err
	}
	clinicRole, err := c.FromOrganizationalUnitRole(ctx, rOur)
	if err != nil {
		return nil, err
	}
	return clinicRole, nil
}

func (c *ClinicService) DisconnectUserFromClinicWithRole(ctx context.Context, input commonModel.UserClinicRoleInput) error {
	rRole, err := c.roleRepo.GetByName(ctx, roles.RoleName(input.RoleIdentifier))
	if err != nil {
		return err
	}
	err = c.clinicsRepo.RemoveUserClinicRole(ctx, &models.UserOrganizationalUnitRole{
		UnitID: input.ClinicID,
		UserID: input.UserID,
		RoleID: rRole.ID,
	})
	if err != nil {
		return err
	}
	return nil
}
