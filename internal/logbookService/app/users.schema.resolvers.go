package logbookService

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

//import (
//	"context"
//	"fmt"
//
//	"gitlab.com/st-appen/logbookService/internal/logbookService/graph/generated"
//	qlmodel "gitlab.com/st-appen/logbookService/internal/logbookService/graph/model"
//)
//
//func (r *userResolver) ClinicRoles(ctx context.Context, obj *qlmodel.User) ([]*qlmodel.ClinicRole, error) {
//	panic(fmt.Errorf("not implemented"))
//}
//
//// User returns generated.UserResolver implementation.
//func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }
//
//type userResolver struct{ *Resolver }
//
//// !!! WARNING !!!
//// The code below was going to be deleted when updating resolvers. It has been copied here so you have
//// one last chance to move it out of harms way if you want. There are two reasons this happens:
////  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
////    it when you're done.
////  - You have helper methods in this file. Move them out to keep these resolver files clean.
//func (r *userResolver) Clinics(ctx context.Context, obj *qlmodel.User) ([]*qlmodel.Clinic, error) {
//	// TODO: this is broken for now.
//	return nil, nil
//}
