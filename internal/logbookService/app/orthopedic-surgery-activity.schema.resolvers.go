package logbookService

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/internal/logbookService/graph/generated"
)

func (r *orthopedicSurgeryActivityResolver) Description(ctx context.Context, obj *commonModel.OrthopedicSurgeryActivity) (*commonModel.Description, error) {
	description := &commonModel.Description{Title: "Ortopedisk operation", Subrows: []*commonModel.DescriptionRow{}}

	if obj.SurgeryMetadata != nil {
		for _, surgery := range obj.SurgeryMetadata.Surgeries {
			subRow := &commonModel.DescriptionRow{}
			if surgery.Diagnose != nil {
				subRow.Title = surgery.Diagnose.DiagnoseName
			}
			if surgery.Method != nil {
				subRow.Subtitle = &surgery.Method.MethodName
				subRow.SubtitleHighlight = &surgery.Method.ApproachName
			}
			if surgery.Method != nil && surgery.Diagnose != nil {
				description.Subrows = append(description.Subrows, subRow)
			}
		}
	}

	return description, nil
}

func (r *orthopedicSurgeryActivityResolver) Resident(ctx context.Context, obj *commonModel.OrthopedicSurgeryActivity) (*commonModel.User, error) {
	user, err := r.service.UsersService.GetById(ctx, obj.ResidentUserID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *orthopedicSurgeryActivityResolver) Supervisor(ctx context.Context, obj *commonModel.OrthopedicSurgeryActivity) (*commonModel.User, error) {
	if obj.SupervisorUserID == nil {
		return nil, nil
	}
	user, err := r.service.UsersService.GetById(ctx, *obj.SupervisorUserID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// OrthopedicSurgeryActivity returns generated.OrthopedicSurgeryActivityResolver implementation.
func (r *Resolver) OrthopedicSurgeryActivity() generated.OrthopedicSurgeryActivityResolver {
	return &orthopedicSurgeryActivityResolver{r}
}

type orthopedicSurgeryActivityResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *orthopedicSurgeryActivityResolver) ResidentUserID(ctx context.Context, obj *commonModel.OrthopedicSurgeryActivity) (string, error) {
	panic(fmt.Errorf("not implemented"))
}
func (r *orthopedicSurgeryActivityResolver) SupervisorUserID(ctx context.Context, obj *commonModel.OrthopedicSurgeryActivity) (*string, error) {
	panic(fmt.Errorf("not implemented"))
}
func (r *orthopedicSurgeryActivityResolver) OperatorID(ctx context.Context, obj *commonModel.OrthopedicSurgeryActivity) (*string, error) {
	panic(fmt.Errorf("not implemented"))
}
func (r *orthopedicSurgeryActivityResolver) AssistantID(ctx context.Context, obj *commonModel.OrthopedicSurgeryActivity) (*string, error) {
	panic(fmt.Errorf("not implemented"))
}
func (r *orthopedicSurgeryActivityReviewResolver) OperatorID(ctx context.Context, obj *commonModel.OrthopedicSurgeryActivityReview) (string, error) {
	panic(fmt.Errorf("not implemented"))
}
func (r *orthopedicSurgeryActivityReviewResolver) AssistantID(ctx context.Context, obj *commonModel.OrthopedicSurgeryActivityReview) (string, error) {
	panic(fmt.Errorf("not implemented"))
}
func (r *orthopedicSurgeryActivityReviewResolver) ResidentUserID(ctx context.Context, obj *commonModel.OrthopedicSurgeryActivityReview) (string, error) {
	panic(fmt.Errorf("not implemented"))
}
func (r *orthopedicSurgeryActivityReviewResolver) SupervisorUserID(ctx context.Context, obj *commonModel.OrthopedicSurgeryActivityReview) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

type orthopedicSurgeryActivityReviewResolver struct{ *Resolver }
