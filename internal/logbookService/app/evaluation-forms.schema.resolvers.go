package logbookService

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/internal/logbookService/graph/generated"
)

func (r *dopsEvaluationResolver) Description(ctx context.Context, obj *commonModel.DopsEvaluation) (*commonModel.Description, error) {
	description := &commonModel.Description{Title: "DOPS", Subrows: []*commonModel.DescriptionRow{}}

	if obj.SurgeryMetadata != nil {
		for _, surgery := range obj.SurgeryMetadata.Surgeries {
			subRow := &commonModel.DescriptionRow{}
			if surgery.Method != nil && surgery.Diagnose != nil {
				subRow.Title = surgery.Diagnose.DiagnoseName
				subRow.Subtitle = &surgery.Method.MethodName
				subRow.SubtitleHighlight = &surgery.Method.ApproachName
				description.Subrows = append(description.Subrows, subRow)
			}
		}
	}

	return description, nil
}

func (r *dopsEvaluationResolver) Supervisor(ctx context.Context, obj *commonModel.DopsEvaluation) (*commonModel.User, error) {
	user, err := r.service.UsersService.GetById(ctx, obj.SupervisorID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *dopsEvaluationResolver) Resident(ctx context.Context, obj *commonModel.DopsEvaluation) (*commonModel.User, error) {
	user, err := r.service.UsersService.GetById(ctx, obj.ResidentID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *miniCexEvaluationResolver) Description(ctx context.Context, obj *commonModel.MiniCexEvaluation) (*commonModel.Description, error) {
	description := &commonModel.Description{Title: "MINI-CEX", Subrows: []*commonModel.DescriptionRow{}}

	if obj.IsEvaluated {
		focusRow := &commonModel.DescriptionRow{}
		focuses := ""
		for index, focus := range obj.Focuses {
			if index > 0 {
				focuses = focuses + ", " + focus
			} else {
				focuses = focus
			}
		}
		focusRow.Title = "Fokus för evalueringen: " + focuses
		description.Subrows = append(description.Subrows, focusRow)

		areaRow := &commonModel.DescriptionRow{}
		areaRow.Title = "Område: " + obj.Area
		description.Subrows = append(description.Subrows, areaRow)
	} else {
		subRow := &commonModel.DescriptionRow{}
		subRow.Title = "Begärd men ej utförd"
		description.Subrows = append(description.Subrows, subRow)
	}

	return description, nil
}

func (r *miniCexEvaluationResolver) Supervisor(ctx context.Context, obj *commonModel.MiniCexEvaluation) (*commonModel.User, error) {
	user, err := r.service.UsersService.GetById(ctx, obj.SupervisorID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *miniCexEvaluationResolver) Resident(ctx context.Context, obj *commonModel.MiniCexEvaluation) (*commonModel.User, error) {
	user, err := r.service.UsersService.GetById(ctx, obj.ResidentID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// DopsEvaluation returns generated.DopsEvaluationResolver implementation.
func (r *Resolver) DopsEvaluation() generated.DopsEvaluationResolver {
	return &dopsEvaluationResolver{r}
}

// MiniCexEvaluation returns generated.MiniCexEvaluationResolver implementation.
func (r *Resolver) MiniCexEvaluation() generated.MiniCexEvaluationResolver {
	return &miniCexEvaluationResolver{r}
}

type dopsEvaluationResolver struct{ *Resolver }
type miniCexEvaluationResolver struct{ *Resolver }
