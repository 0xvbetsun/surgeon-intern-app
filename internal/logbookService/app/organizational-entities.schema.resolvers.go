package logbookService

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/internal/logbookService/graph/generated"
	"github.com/vbetsun/surgeon-intern-app/internal/logbookService/pkg/services/clinics"
)

func (r *clinicResolver) Hospital(ctx context.Context, obj *commonModel.Clinic) (*commonModel.Hospital, error) {
	hospital, err := r.service.HospitalService.GetById(ctx, *obj.HospitalID)
	if err != nil {
		return nil, err
	}

	return hospital, nil
}

func (r *hospitalResolver) Clinics(ctx context.Context, obj *commonModel.Hospital) ([]*commonModel.Clinic, error) {
	clinics, err := r.service.ClinicsService.GetByFilter(ctx, clinics.ClinicFilter{
		HospitalId:  &obj.HospitalID,
		DisplayName: nil,
	}, nil, nil)
	if err != nil {
		return nil, err
	}
	return clinics, nil
}

// Clinic returns generated.ClinicResolver implementation.
func (r *Resolver) Clinic() generated.ClinicResolver { return &clinicResolver{r} }

// Hospital returns generated.HospitalResolver implementation.
func (r *Resolver) Hospital() generated.HospitalResolver { return &hospitalResolver{r} }

type clinicResolver struct{ *Resolver }
type hospitalResolver struct{ *Resolver }
