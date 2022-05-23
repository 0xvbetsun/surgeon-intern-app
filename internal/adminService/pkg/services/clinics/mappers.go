package clinics

import (
	"context"

	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/pkg/models"
)

func FromRepoClinics(clinics []*models.OrganizationalUnit) []*commonModel.Clinic {
	rClinics := make([]*commonModel.Clinic, 0)
	for _, clinic := range clinics {
		rClinics = append(rClinics, FromRepoClinic(clinic))
	}
	return rClinics
}

func FromRepoClinic(clinic *models.OrganizationalUnit) *commonModel.Clinic {
	// defaults to the unknown specialty
	clinicSpecialty := commonModel.SpecialtiesUnknown
	if clinic.R != nil && clinic.R.Specialties != nil && len(clinic.R.Specialties) > 0 {
		// Just grab the first one for now, does a clinic ever have more than one specialty ?
		clinicSpecialty = commonModel.Specialties(clinic.R.Specialties[0].Name)
	}
	rClinic := &commonModel.Clinic{
		ClinicID:    clinic.ID,
		DisplayName: clinic.DisplayName,
		HospitalID:  clinic.ParentID.Ptr(),
		Specialty:   clinicSpecialty,
	}
	if clinic.R != nil && clinic.R.Parent != nil {
		rClinic.Hospital = &commonModel.Hospital{
			HospitalID:  clinic.R.Parent.ID,
			DisplayName: clinic.R.Parent.DisplayName,
		}
	}
	return rClinic
}

func (c *ClinicService) FromOrganizationalUnitRole(ctx context.Context, our *models.UserOrganizationalUnitRole) (*commonModel.ClinicRole, error) {
	role, err := c.roleRepo.GetByID(ctx, our.RoleID)
	if err != nil {
		return nil, err
	}
	clinic, err := c.clinicsRepo.GetByID(ctx, our.UnitID)
	if err != nil {
		return nil, err
	}
	clinicSpecialty := commonModel.SpecialtiesUnknown
	if clinic.R != nil && clinic.R.Specialties != nil && len(clinic.R.Specialties) > 0 {
		// Just grab the first one for now, does a clinic ever have more than one specialty ?
		clinicSpecialty = commonModel.Specialties(clinic.R.Specialties[0].Name)
	}
	return &commonModel.ClinicRole{
		Clinic: &commonModel.Clinic{
			ClinicID:    clinic.ID,
			DisplayName: clinic.DisplayName,
			Specialty:   clinicSpecialty,
		},
		Role: &commonModel.Role{
			RoleIdentifier: role.Name,
			DisplayName:    role.DisplayName,
		},
	}, nil
}
