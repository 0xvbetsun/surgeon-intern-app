package hospitals

import (
	"context"

	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/OrganizationalUnits"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/roles"
)

type (
	HospitalService struct {
		unitsRepo     *OrganizationalUnits.Repo
		unitTypesRepo *OrganizationalUnits.TypesRepo
		roleRepo      roles.IRepo
	}
	Filter struct {
		HospitalName *string
	}
)

func NewHospitalService(unitsRepo *OrganizationalUnits.Repo, typesRepo *OrganizationalUnits.TypesRepo, roleRepo roles.IRepo) *HospitalService {
	return &HospitalService{
		unitsRepo:     unitsRepo,
		unitTypesRepo: typesRepo,
		roleRepo:      roleRepo,
	}
}

func (u *HospitalService) GetByFilter(ctx context.Context, filter Filter) ([]*commonModel.Hospital, error) {
	clinicUnitTypeId, err := u.unitTypesRepo.GetByName(ctx, OrganizationalUnits.HOSPITAL)
	if err != nil {
		return nil, err
	}
	hospitals, err := u.unitsRepo.ListByFilter(ctx, OrganizationalUnits.OrganizationalUnitFilter{
		TypeID:      &clinicUnitTypeId.ID,
		DisplayName: filter.HospitalName,
	})
	if err != nil {
		return nil, err
	}
	rHospitals, err := u.FromRepoHospitals(ctx, hospitals)
	if err != nil {
		return nil, err
	}
	return rHospitals, nil
}

func (u *HospitalService) AddHospital(ctx context.Context) {}

func (u *HospitalService) UpdateHospital(ctx context.Context) {}

// TODO: Add methods for adding hospitalRoles

//func (c *HospitalService) ConnectUserToClinicWithRole(ctx context.Context, input commonModel.UserClinicRoleInput) (*commonModel.ClinicRole, error) {
//	rRole, err := c.roleRepo.GetByName(ctx, roles.RoleName(input.RoleIdentifier))
//	if err != nil {
//		return nil, err
//	}
//	rOur, err := c.clinicsRepo.AddUserClinicRole(ctx, &models.UserOrganizationalUnitRole{
//		UnitID: input.ClinicID,
//		UserID: input.UserID,
//		RoleID: rRole.ID,
//	})
//	if err != nil {
//		return nil, err
//	}
//	clinicRole, err := c.FromOrganizationalUnitRole(ctx, rOur)
//	if err != nil {
//		return nil, err
//	}
//	return clinicRole, nil
//}
//
//func (c *HospitalService) DisconnectUserFromClinicWithRole(ctx context.Context, input commonModel.UserClinicRoleInput) error {
//	rRole, err := c.roleRepo.GetByName(ctx, roles.RoleName(input.RoleIdentifier))
//	if err != nil {
//		return err
//	}
//	err = c.clinicsRepo.RemoveUserClinicRole(ctx, &models.UserOrganizationalUnitRole{
//		UnitID: input.ClinicID,
//		UserID: input.UserID,
//		RoleID: rRole.ID,
//	})
//	if err != nil {
//		return err
//	}
//	return nil
//}
