package hospitals

import (
	"context"

	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/OrganizationalUnits"
)

func (u *HospitalService) FromRepoHospitals(ctx context.Context, hospitals []*models.OrganizationalUnit) ([]*commonModel.Hospital, error) {
	rHospitals := make([]*commonModel.Hospital, 0)
	for _, hospital := range hospitals {
		rHospital, err := u.FromRepoHospital(ctx, hospital)
		if err != nil {
			return nil, err
		}
		rHospitals = append(rHospitals, rHospital)
	}
	return rHospitals, nil
}

func (u *HospitalService) FromRepoHospital(ctx context.Context, hospital *models.OrganizationalUnit) (*commonModel.Hospital, error) {
	rHospital := &commonModel.Hospital{
		HospitalID:  hospital.ID,
		DisplayName: hospital.DisplayName,
	}
	clinicUnitTypeId, err := u.unitTypesRepo.GetByName(ctx, OrganizationalUnits.CLINIC)
	if err != nil {
		return nil, err
	}
	clinics, err := u.unitsRepo.ListByFilter(ctx, OrganizationalUnits.OrganizationalUnitFilter{
		TypeID:   &clinicUnitTypeId.ID,
		ParentID: &rHospital.HospitalID,
	})
	if err != nil {
		return nil, err
	}
	hospitalClinics := make([]*commonModel.Clinic, 0)
	for _, clinic := range clinics {
		// defaults to the unknown specialty
		clinicSpecialty := commonModel.SpecialtiesUnknown
		if clinic.R != nil && clinic.R.Specialties != nil && len(clinic.R.Specialties) > 0 {
			// Just grab the first one for now, does a clinic ever have more than one specialty ?
			clinicSpecialty = commonModel.Specialties(clinic.R.Specialties[0].Name)
		}
		hospitalClinics = append(hospitalClinics, &commonModel.Clinic{
			ClinicID:    clinic.ID,
			DisplayName: clinic.DisplayName,
			HospitalID:  &rHospital.HospitalID,
			Specialty:   clinicSpecialty,
		})
	}
	rHospital.Clinics = hospitalClinics

	return rHospital, nil
}
