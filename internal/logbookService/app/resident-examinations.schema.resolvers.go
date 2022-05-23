package logbookService

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/internal/logbookService/graph/generated"
)

func (r *residentExaminationResolver) Resident(ctx context.Context, obj *commonModel.ResidentExamination) (*commonModel.User, error) {
	user, err := r.service.UsersService.GetById(ctx, obj.ResidentUserID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *residentExaminationResolver) Supervisor(ctx context.Context, obj *commonModel.ResidentExamination) (*commonModel.User, error) {
	user, err := r.service.UsersService.GetById(ctx, obj.SupervisorUserID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// ResidentExamination returns generated.ResidentExaminationResolver implementation.
func (r *Resolver) ResidentExamination() generated.ResidentExaminationResolver {
	return &residentExaminationResolver{r}
}

type residentExaminationResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *residentExaminationResolver) ResidentUserID(ctx context.Context, obj *commonModel.ResidentExamination) (string, error) {
	panic(fmt.Errorf("not implemented"))
}
func (r *residentExaminationResolver) SupervisorUserID(ctx context.Context, obj *commonModel.ResidentExamination) (string, error) {
	panic(fmt.Errorf("not implemented"))
}
