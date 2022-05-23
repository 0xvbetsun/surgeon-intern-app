package clinics

import (
	"context"

	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/OrganizationalUnits"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/clinics"
)

type (
	IService interface {
		GetById(ctx context.Context, clinicID string) (*commonModel.Clinic, error)
		GetByFilter(ctx context.Context, filter ClinicFilter, userID *string, roleID *int) ([]*commonModel.Clinic, error)
		ListByHospitalId(ctx context.Context, hospitalId string) ([]*commonModel.Clinic, error)
		ListAll(ctx context.Context) ([]*commonModel.Clinic, error)
		GetClinicDepartments(ctx context.Context, clinicID string) ([]*commonModel.ClinicDepartment, error)
		GetDepartment(ctx context.Context, departmentID string) (*commonModel.ClinicDepartment, error)
	}
	Service struct {
		repo             clinics.IRepo
		orgUnitsRepo     *OrganizationalUnits.Repo
		orgUnitTypesRepo *OrganizationalUnits.TypesRepo
	}
	ClinicFilter struct {
		HospitalId  *string
		DisplayName *string
	}
)

func NewService(repo clinics.IRepo,
	orgUnitRepo *OrganizationalUnits.Repo,
	orgUnitTypesRepo *OrganizationalUnits.TypesRepo) IService {
	return &Service{repo: repo, orgUnitsRepo: orgUnitRepo, orgUnitTypesRepo: orgUnitTypesRepo}
}

func (s *Service) GetByFilter(ctx context.Context, filter ClinicFilter, userID *string, roleID *int) ([]*commonModel.Clinic, error) {
	clinicType, err := s.orgUnitTypesRepo.GetByName(ctx, OrganizationalUnits.CLINIC)
	if err != nil {
		return nil, err
	}
	orgUnitFilter := OrganizationalUnits.OrganizationalUnitFilter{
		TypeID:      &clinicType.ID,
		ParentID:    filter.HospitalId,
		DisplayName: filter.DisplayName,
	}
	if userID != nil {
		orgUnitFilter.UserFilter = &OrganizationalUnits.UserFilter{
			UserID: *userID,
			RoleID: roleID,
		}
	}
	rClinics, err := s.orgUnitsRepo.ListByFilter(ctx, orgUnitFilter)
	if err != nil {
		return nil, err
	}
	var clinics = make([]*commonModel.Clinic, 0)
	for _, clinic := range rClinics {
		clinics = append(clinics, commonModel.ClinicFromRepoType(clinic))
	}
	return clinics, nil
}

func (s *Service) GetClinicsByUserAndAtLeastRole(ctx context.Context, userID int, roleId int) ([]*commonModel.Clinic, error) {
	panic("implement me")
}

func (s *Service) GetById(ctx context.Context, clinicID string) (*commonModel.Clinic, error) {
	rClinic, err := s.repo.GetByID(ctx, clinicID)
	if err != nil {
		return nil, err
	}
	clinic := s.mapClinicGraphQlModel(rClinic)
	return clinic, nil
}

func (s *Service) ListByHospitalId(ctx context.Context, hospitalId string) ([]*commonModel.Clinic, error) {
	rClinics, err := s.repo.ListByHospitalId(ctx, hospitalId)
	if err != nil {
		return nil, err
	}
	p := make([]*commonModel.Clinic, 0)
	for _, rClinic := range rClinics {
		clinic := s.mapClinicGraphQlModel(rClinic)
		p = append(p, clinic)
	}
	return p, nil
}

func (s *Service) ListAll(ctx context.Context) ([]*commonModel.Clinic, error) {
	rClinics, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	p := make([]*commonModel.Clinic, 0)
	for _, rClinic := range rClinics {
		clinic := s.mapClinicGraphQlModel(rClinic)
		p = append(p, clinic)
	}
	return p, nil
}

func (s *Service) GetClinicDepartments(ctx context.Context, clinicID string) ([]*commonModel.ClinicDepartment, error) {
	rLocations, err := s.repo.GetClinicDepartments(ctx, clinicID)
	if err != nil {
		return nil, err
	}

	locations := make([]*commonModel.ClinicDepartment, 0)
	for _, rLocation := range rLocations {
		location := s.mapClinicDepartmentGraphQlModel(rLocation)
		locations = append(locations, location)
	}
	return locations, nil
}

func (s *Service) GetDepartment(ctx context.Context, departmentID string) (*commonModel.ClinicDepartment, error) {
	rLocation, err := s.repo.GetClinicDepartmentByID(ctx, departmentID)
	if err != nil {
		return nil, err
	}

	location := s.mapClinicDepartmentGraphQlModel(rLocation)
	return location, nil
}

func (s Service) mapClinicGraphQlModel(rClinic *models.OrganizationalUnit) *commonModel.Clinic {
	clinic := &commonModel.Clinic{
		ClinicID:    rClinic.ID,
		DisplayName: rClinic.DisplayName,
		HospitalID:  rClinic.ParentID.Ptr(),
	}
	return clinic
}

func (s Service) mapClinicDepartmentGraphQlModel(rDepartment *models.OrganizationalUnit) *commonModel.ClinicDepartment {
	clinic := &commonModel.ClinicDepartment{
		DepartmentID:   rDepartment.ID,
		DepartmentName: rDepartment.DisplayName,
		ClinicID:       rDepartment.ParentID.String,
	}
	return clinic
}
