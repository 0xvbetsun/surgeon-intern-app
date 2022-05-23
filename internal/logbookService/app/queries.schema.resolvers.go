package logbookService

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"

	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/internal/logbookService/graph/generated"
	"github.com/vbetsun/surgeon-intern-app/internal/logbookService/pkg/services/clinics"
	"github.com/vbetsun/surgeon-intern-app/internal/logbookService/pkg/services/hospitals"
	"github.com/vbetsun/surgeon-intern-app/internal/pkg/commonErrors"
	gin2 "github.com/vbetsun/surgeon-intern-app/internal/pkg/gin"
	middleware2 "github.com/vbetsun/surgeon-intern-app/internal/pkg/middleware"
	"github.com/vbetsun/surgeon-intern-app/internal/pkg/services/users"
	"github.com/vbetsun/surgeon-intern-app/internal/util"
	"go.uber.org/zap"
)

func (r *queryResolver) Procedures(ctx context.Context) ([]*commonModel.Procedure, error) {
	procedures, err := r.service.ProceduresService.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	return procedures, nil
}

func (r *queryResolver) Examinations(ctx context.Context) ([]*commonModel.Examination, error) {
	examinations, err := r.service.ExaminationsService.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	return examinations, nil
}

func (r *queryResolver) ResidentExaminations(ctx context.Context, residentUserID *string, supervisorUserID *string) ([]*commonModel.ResidentExamination, error) {
	if residentUserID != nil {
		examinations, err := r.service.ResidentExaminationsService.GetByResidentUser(ctx, *residentUserID)
		if err != nil {
			return nil, err
		}
		return examinations, nil
	}
	if supervisorUserID != nil {
		examinations, err := r.service.ResidentExaminationsService.GetBySupervisorUser(ctx, *supervisorUserID)
		if err != nil {
			return nil, err
		}
		return examinations, nil
	}
	// Assume that the authorized user is a resident doctor.
	// TODO: think this trough...
	gc, err := gin2.GinContextFromContext(ctx)
	if err != nil {
		return nil, errors.New("Unknown server error")
	}
	if uID, ok := gc.Get(middleware2.AUTH0GinContextUserIDKey); ok != false {
		zap.S().Info(uID)
		user, err := r.service.UsersService.GetByExternalID(ctx, fmt.Sprintf("%v", uID))
		if err != nil {
			return nil, errors.New("Unknown server error")
		}
		examinations, err := r.service.ResidentExaminationsService.GetByResidentUser(ctx, user.UserID)
		if err != nil {
			return nil, err
		}
		return examinations, nil
	}

	return nil, errors.New("rEither resident user id or supervisor user id must be submitted.")
}

func (r *queryResolver) ResidentExamination(ctx context.Context, residentExaminationID string) (*commonModel.ResidentExamination, error) {
	examination, err := r.service.ResidentExaminationsService.GetById(ctx, residentExaminationID)
	if err != nil {
		return nil, err
	}
	return examination, nil
}

func (r *queryResolver) SupervisorExaminationReviews(ctx context.Context, reviewed *bool) ([]*commonModel.SupervisorExaminationReview, error) {
	// Assume that the authorized user is a supervisor.
	gc, err := gin2.GinContextFromContext(ctx)
	if err != nil {
		return nil, errors.New("Unknown server error")
	}
	if uID, ok := gc.Get(middleware2.GinContextUserIDKey); ok != false {
		user, err := r.service.UsersService.GetByExternalID(ctx, fmt.Sprintf("%v", uID))
		if err != nil {
			return nil, errors.New("Unknown server error")
		}
		reviews, err := r.service.SupervisorExaminationReviewsService.GetBySupervisorUser(ctx, user.UserID, reviewed)
		if err != nil {
			return nil, err
		}
		return reviews, nil
	}
	return nil, errors.New("Couldn't get user through firebase id.")
}

func (r *queryResolver) SupervisorExaminationReview(ctx context.Context, supervisorExaminationReviewID string) (*commonModel.SupervisorExaminationReview, error) {
	review, err := r.service.SupervisorExaminationReviewsService.GetByExaminationId(ctx, supervisorExaminationReviewID)
	if err != nil {
		return nil, err
	}
	return review, nil
}

func (r *queryResolver) Examination(ctx context.Context, id string) (*commonModel.Examination, error) {
	examination, err := r.service.ExaminationsService.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return examination, nil
}

func (r *queryResolver) Activities(ctx context.Context, queryFilter commonModel.ActivityQueryFilter, orderBy *commonModel.QueryOrder, pagination *commonModel.QueryPaging) ([]*commonModel.Activity, error) {
	activities, err := r.service.ActivitiesService.GetActivities(ctx, queryFilter, orderBy, pagination)
	if err != nil {
		return nil, err
	}

	return activities, nil
}

func (r *queryResolver) LogbookEntries(ctx context.Context, queryFilters commonModel.LogbookEntryQueryFilter, orderBy *commonModel.QueryOrder, pagination *commonModel.QueryPaging) ([]*commonModel.LogbookEntry, error) {
	logbookEntries, err := r.service.LogbookEntriesService.GetLogbookEntries(ctx, queryFilters, orderBy, pagination)
	if err != nil {
		return nil, err
	}
	return logbookEntries, nil
}

func (r *queryResolver) Assessments(ctx context.Context, queryFilters commonModel.AssessmentQueryFilter, orderBy *commonModel.QueryOrder, pagination *commonModel.QueryPaging) ([]*commonModel.Assessment, error) {
	reviews, err := r.service.AssessmentsService.GetAssessments(ctx, queryFilters, orderBy, pagination)
	if err != nil {
		return nil, err
	}
	return reviews, nil
}

func (r *queryResolver) Hospitals(ctx context.Context, organisationID *string) ([]*commonModel.Hospital, error) {
	user, err := middleware2.UserFromContext(ctx, r.service.UsersService)
	if err != nil {
		return nil, errors.New("Unknown")
	}
	hospitals, err := r.service.HospitalService.GetByFilter(ctx, hospitals.HospitalFilter{
		HospitalId: organisationID,
	}, &user.UserID, nil)
	if err != nil {
		return nil, err
	}
	return hospitals, err
}

func (r *queryResolver) Clinics(ctx context.Context, hospitalID *string) ([]*commonModel.Clinic, error) {
	user, err := middleware2.UserFromContext(ctx, r.service.UsersService)
	if err != nil {
		return nil, errors.New("Unknown")
	}
	clinics, err := r.service.ClinicsService.GetByFilter(ctx, clinics.ClinicFilter{
		HospitalId: hospitalID,
	}, &user.UserID, nil)
	if err != nil {
		return nil, err
	}
	return clinics, err
}

func (r *queryResolver) ClinicDepartments(ctx context.Context, clinicID string) ([]*commonModel.ClinicDepartment, error) {
	departments, err := r.service.ClinicsService.GetClinicDepartments(ctx, clinicID)
	if err != nil {
		return nil, err
	}
	return departments, nil
}

func (r *queryResolver) ClinicDepartment(ctx context.Context, departmentID string) (*commonModel.ClinicDepartment, error) {
	department, err := r.service.ClinicsService.GetDepartment(ctx, departmentID)
	if err != nil {
		return nil, err
	}
	return department, nil
}

func (r *queryResolver) Supervisors(ctx context.Context, clinicID *string) ([]*commonModel.User, error) {
	user, err := middleware2.UserFromContext(ctx, r.service.UsersService)
	if err != nil {
		return nil, errors.New("Unknown")
	}
	isSupervisor := true
	supervisors, err := r.service.UsersService.ListByUserAndFilter(ctx, user.UserID, users.UserServiceFilter{
		OrganizationalUnitID: clinicID,
		IsSupervisor:         &isSupervisor,
		IsResident:           nil,
	})
	if err != nil {
		return nil, err
	}
	return supervisors, nil
}

func (r *queryResolver) Residents(ctx context.Context, clinicID *string) ([]*commonModel.User, error) {
	user, err := middleware2.UserFromContext(ctx, r.service.UsersService)
	if err != nil {
		return nil, errors.New("Unknown")
	}
	isResident := true
	residents, err := r.service.UsersService.ListByUserAndFilter(ctx, user.UserID, users.UserServiceFilter{
		OrganizationalUnitID: clinicID,
		IsSupervisor:         nil,
		IsResident:           &isResident,
	})
	if err != nil {
		return nil, err
	}
	return residents, nil
}

func (r *queryResolver) UserRoles(ctx context.Context, userID *string) ([]*commonModel.Role, error) {
	user, err := middleware2.UserFromContext(ctx, r.service.UsersService)
	if err != nil {
		return nil, errors.New("couldn't get gin context from context")
	}
	roles, err := r.service.AuthorizationService.GetAvailableRolesForUser(ctx, user.UserID)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *queryResolver) User(ctx context.Context, userID *string) (*commonModel.User, error) {
	gc, err := gin2.GinContextFromContext(ctx)
	if err != nil {
		return nil, errors.New("couldn't get gin context from context")
	}
	if userID != nil {
		user, err := r.service.UsersService.GetById(ctx, *userID)
		if err != nil {
			return nil, err
		}
		return user, nil
	} else if uID, ok := gc.Get(middleware2.AUTH0GinContextUserIDKey); ok != false {
		user, err := r.service.UsersService.GetByExternalID(ctx, uID.(string))
		if err != nil {
			return nil, err
		}
		return user, nil
	} else {
		return nil, errors.New("has to supply user id or authorization")
	}
}

func (r *queryResolver) PracticalActivityTypes(ctx context.Context) ([]*commonModel.PracticalActivityType, error) {
	user, err := middleware2.UserFromContext(ctx, r.service.UsersService)
	if err != nil {
		return nil, err
	}
	practicalActivityTypes, err := r.service.PracticalActivityTypesService.ListByUser(ctx, user.UserID)
	if err != nil {
		return nil, err
	}
	return practicalActivityTypes, nil
}

func (r *queryResolver) OrthopedicSurgeryActivityByID(ctx context.Context, activityID string) (*commonModel.OrthopedicSurgeryActivity, error) {
	if !util.IsValidUUID(activityID) {
		return nil, commonErrors.NewInvalidInputError(fmt.Sprintf("ActivityID: %s is not a valid UUID", activityID))
	}
	orthopedicSurgeryActivity, err := r.service.OrthopedicSurgeryService.ActivityService.Get(ctx, activityID)
	if err != nil {
		zap.S().Errorf(err.Error())
		return nil, err
	}
	return orthopedicSurgeryActivity, nil
}

func (r *queryResolver) OrthopedicSurgeryActivityReviewByID(ctx context.Context, reviewID string) (*commonModel.OrthopedicSurgeryActivityReview, error) {
	orthopedicSurgeryActivityReview, err := r.service.OrthopedicSurgeryService.ReviewService.Get(ctx, reviewID)
	if err != nil {
		zap.S().Errorf(err.Error())
		return nil, err
	}
	return orthopedicSurgeryActivityReview, nil
}

func (r *queryResolver) OrthopedicSurgeryActivities(ctx context.Context, queryFilter commonModel.SurgeryLogbookEntryQueryFilter, orderBy *commonModel.QueryOrder) ([]*commonModel.OrthopedicSurgeryActivity, error) {
	surgeries, err := r.service.OrthopedicSurgeryService.ActivityService.GetActivities(ctx, queryFilter, orderBy)
	if err != nil {
		zap.S().Errorf(err.Error())
		return nil, err
	}
	return surgeries, nil
}

func (r *queryResolver) Surgeries(ctx context.Context, clinicID *string, orderBy *commonModel.SurgeryOrder, paging *commonModel.QueryPaging) ([]*commonModel.Surgery, error) {
	if clinicID != nil {
		surgeries, err := r.service.OrthopedicSurgeryService.SurgeryService.GetByClinic(ctx, *clinicID)
		if err != nil {
			return nil, err
		}
		return surgeries, nil
	}
	surgeries, err := r.service.OrthopedicSurgeryService.SurgeryService.List(ctx)
	if err != nil {
		return nil, err
	}
	return surgeries, nil
}

func (r *queryResolver) SurgeryByID(ctx context.Context, surgeryID string) (*commonModel.Surgery, error) {
	surgery, err := r.service.OrthopedicSurgeryService.SurgeryService.Get(ctx, surgeryID)
	if err != nil {
		return nil, err
	}
	return surgery, nil
}

func (r *queryResolver) SurgeriesByDiagnose(ctx context.Context, diagnoseID string) ([]*commonModel.Surgery, error) {
	surgeries, err := r.service.OrthopedicSurgeryService.SurgeryService.GetByDiagnose(ctx, diagnoseID)
	if err != nil {
		return nil, err
	}
	return surgeries, nil
}

func (r *queryResolver) SurgeriesByMethod(ctx context.Context, methodID string) ([]*commonModel.Surgery, error) {
	surgeries, err := r.service.OrthopedicSurgeryService.SurgeryService.GetByMethod(ctx, methodID)
	if err != nil {
		return nil, err
	}
	return surgeries, nil
}

func (r *queryResolver) SurgeryByDiagnoseAndMethod(ctx context.Context, methodID string, diagnoseID string) (*commonModel.Surgery, error) {
	surgery, err := r.service.OrthopedicSurgeryService.SurgeryService.GetByMethodAndDiagnose(ctx, methodID, diagnoseID)
	if err != nil {
		return nil, err
	}
	return surgery, nil
}

func (r *queryResolver) EvaluationForms(ctx context.Context) ([]*commonModel.EvaluationForm, error) {
	forms, err := r.service.EvaluationFormService.EvaluationForms.ListAllEvaluationForms(ctx)
	if err != nil {
		return nil, err
	}
	return forms, nil
}

func (r *queryResolver) DopsEvaluations(ctx context.Context, queryFilter commonModel.DopsQueryFilter, orderBy *commonModel.QueryOrder) ([]*commonModel.DopsEvaluation, error) {
	dops, err := r.service.EvaluationFormService.Dops.GetDopsEvaluations(ctx, queryFilter, orderBy)
	if err != nil {
		return nil, err
	}
	return dops, nil
}

func (r *queryResolver) DopsEvaluation(ctx context.Context, dopsID *string, activityID *string) (*commonModel.DopsEvaluation, error) {
	dopsEvaluation, err := r.service.EvaluationFormService.Dops.GetDopsEvaluationById(ctx, dopsID, activityID)
	if err != nil {
		return nil, err
	}
	return dopsEvaluation, nil
}

func (r *queryResolver) MiniCexEvaluations(ctx context.Context, queryFilter commonModel.MiniCexQueryFilter, orderBy *commonModel.QueryOrder) ([]*commonModel.MiniCexEvaluation, error) {
	dops, err := r.service.EvaluationFormService.MiniCex.GetMiniCexEvaluations(ctx, queryFilter, orderBy)
	if err != nil {
		return nil, err
	}
	return dops, nil
}

func (r *queryResolver) MiniCexEvaluation(ctx context.Context, miniCexEvaluationID string) (*commonModel.MiniCexEvaluation, error) {
	dopsEvaluation, err := r.service.EvaluationFormService.MiniCex.GetMiniCexEvaluationById(ctx, miniCexEvaluationID)
	if err != nil {
		return nil, err
	}
	return dopsEvaluation, nil
}

func (r *queryResolver) MiniCexFocuses(ctx context.Context) ([]*commonModel.MiniCexFocus, error) {
	focuses, err := r.service.EvaluationFormService.MiniCex.GetMiniCexFocuses(ctx)
	if err != nil {
		return nil, err
	}
	return focuses, nil
}

func (r *queryResolver) MiniCexAreas(ctx context.Context, departmentID string) ([]*commonModel.MiniCexArea, error) {
	areas, err := r.service.EvaluationFormService.MiniCex.GetMiniCexAreasByClinicId(ctx, departmentID)
	if err != nil {
		return nil, err
	}
	return areas, nil
}

func (r *queryResolver) Notifications(ctx context.Context, notificationType *commonModel.NotificationType, hasSeen bool) ([]*commonModel.Notification, error) {
	notifications, err := r.service.NotificationService.GetNotifications(ctx, notificationType, hasSeen)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
