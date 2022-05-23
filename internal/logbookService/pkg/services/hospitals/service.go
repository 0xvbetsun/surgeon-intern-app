package hospitals

import (
	"context"

	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/OrganizationalUnits"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/hospitals"
)

type (
	IService interface {
		GetById(ctx context.Context, hospitalID string) (*commonModel.Hospital, error)
		GetByFilter(ctx context.Context, filter HospitalFilter, userID *string, roleID *int) ([]*commonModel.Hospital, error)
		ListByOrganisationId(ctx context.Context, organisationId string) ([]*commonModel.Hospital, error)
		ListAll(ctx context.Context) ([]*commonModel.Hospital, error)
	}
	Service struct {
		repo             hospitals.IRepo
		orgUnitsRepo     *OrganizationalUnits.Repo
		orgUnitTypesRepo *OrganizationalUnits.TypesRepo
	}
	HospitalFilter struct {
		HospitalId  *string
		DisplayName *string
	}
)

func NewService(repo hospitals.IRepo,
	orgUnitsRepo *OrganizationalUnits.Repo,
	orgUnitTypesRepo *OrganizationalUnits.TypesRepo) IService {
	return &Service{repo: repo,
		orgUnitsRepo:     orgUnitsRepo,
		orgUnitTypesRepo: orgUnitTypesRepo}
}

func (s *Service) GetById(ctx context.Context, hospitalID string) (*commonModel.Hospital, error) {
	rHospital, err := s.repo.GetByID(ctx, hospitalID)
	if err != nil {
		return nil, err
	}
	hospital := s.mapGraphQlModel(rHospital)
	return hospital, nil
}

func (s *Service) GetByFilter(ctx context.Context, filter HospitalFilter, userID *string, roleID *int) ([]*commonModel.Hospital, error) {
	clinicType, err := s.orgUnitTypesRepo.GetByName(ctx, OrganizationalUnits.HOSPITAL)
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
	rHospitals, err := s.orgUnitsRepo.ListByFilter(ctx, orgUnitFilter)
	if err != nil {
		return nil, err
	}
	var hospitals = make([]*commonModel.Hospital, 0)
	for _, hospital := range rHospitals {
		hospitals = append(hospitals, commonModel.HospitalFromRepoType(hospital))
	}
	return hospitals, nil
}

func (s *Service) ListByOrganisationId(ctx context.Context, organisationId string) ([]*commonModel.Hospital, error) {
	rHospitals, err := s.repo.ListByOrganisationId(ctx, organisationId)
	if err != nil {
		return nil, err
	}
	p := make([]*commonModel.Hospital, 0)
	for _, rHospital := range rHospitals {
		hospital := s.mapGraphQlModel(rHospital)
		p = append(p, hospital)
	}
	return p, nil
}

func (s *Service) ListAll(ctx context.Context) ([]*commonModel.Hospital, error) {
	rHospitals, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	p := make([]*commonModel.Hospital, 0)
	for _, rHospital := range rHospitals {
		hospital := s.mapGraphQlModel(rHospital)
		p = append(p, hospital)
	}
	return p, nil
}

func (s Service) mapGraphQlModel(rHospital *models.OrganizationalUnit) *commonModel.Hospital {
	hospital := &commonModel.Hospital{
		HospitalID:  rHospital.ID,
		DisplayName: rHospital.DisplayName,
	}
	return hospital
}
