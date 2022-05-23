package logbookService

import (
	"github.com/google/wire"
	"github.com/vbetsun/surgeon-intern-app/internal/logbookService/pkg/services/activities"
	"github.com/vbetsun/surgeon-intern-app/internal/logbookService/pkg/services/assessments"
	"github.com/vbetsun/surgeon-intern-app/internal/logbookService/pkg/services/clinics"
	"github.com/vbetsun/surgeon-intern-app/internal/logbookService/pkg/services/evaluationforms"
	"github.com/vbetsun/surgeon-intern-app/internal/logbookService/pkg/services/examinations"
	"github.com/vbetsun/surgeon-intern-app/internal/logbookService/pkg/services/examinationsactivities"
	"github.com/vbetsun/surgeon-intern-app/internal/logbookService/pkg/services/examinationsactivitiesreviews"
	"github.com/vbetsun/surgeon-intern-app/internal/logbookService/pkg/services/hospitals"
	"github.com/vbetsun/surgeon-intern-app/internal/logbookService/pkg/services/logbookentries"
	notifications2 "github.com/vbetsun/surgeon-intern-app/internal/logbookService/pkg/services/notifications"
	"github.com/vbetsun/surgeon-intern-app/internal/logbookService/pkg/services/organisations"
	"github.com/vbetsun/surgeon-intern-app/internal/logbookService/pkg/services/orthopedicSurgeries"
	"github.com/vbetsun/surgeon-intern-app/internal/logbookService/pkg/services/practicalActivityTypes"
	"github.com/vbetsun/surgeon-intern-app/internal/logbookService/pkg/services/procedures"
	authorization2 "github.com/vbetsun/surgeon-intern-app/internal/pkg/services/authorization"
	"github.com/vbetsun/surgeon-intern-app/internal/pkg/services/users"
	thirdparty2 "github.com/vbetsun/surgeon-intern-app/internal/pkg/thirdparty"
)

type (
	Service struct {
		ActivitiesService                   activities.IActivitiesService
		PracticalActivityTypesService       practicalActivityTypes.IService
		ExaminationsService                 examinations.IService
		ProceduresService                   procedures.IService
		ResidentExaminationsService         *examinationsactivities.Service
		SupervisorExaminationReviewsService examinationsactivitiesreviews.IService
		OrganisationsService                *organisations.Service
		UsersService                        users.IService
		HospitalService                     hospitals.IService
		ClinicsService                      clinics.IService
		AuthorizationService                authorization2.IService
		OrthopedicSurgeryService            *orthopedicSurgeries.Service
		EvaluationFormService               evaluationforms.Service
		LogbookEntriesService               logbookentries.ILogbookEntriesService
		AssessmentsService                  assessments.IAssessmentsService
		NotificationService                 notifications2.INotificationsService
	}
)

func NewService(activitiesService activities.IActivitiesService,
	practicalActivityTypesService practicalActivityTypes.IService,
	examinationsService examinations.IService,
	proceduresService procedures.IService,
	residentExaminationsService *examinationsactivities.Service,
	supervisorExaminationReviewsService examinationsactivitiesreviews.IService, organisationsService *organisations.Service,
	usersService users.IService,
	hospitalService hospitals.IService,
	clinicsService clinics.IService,
	authorizationService authorization2.IService,
	orthopedicSurgeryService *orthopedicSurgeries.Service,
	evaluationFormService evaluationforms.Service,
	logbookEntriesService logbookentries.ILogbookEntriesService,
	assessmentsSerrvice assessments.IAssessmentsService,
	notificationsService notifications2.INotificationsService) (*Service, error) {

	return &Service{
		ActivitiesService:                   activitiesService,
		PracticalActivityTypesService:       practicalActivityTypesService,
		ExaminationsService:                 examinationsService,
		ProceduresService:                   proceduresService,
		ResidentExaminationsService:         residentExaminationsService,
		SupervisorExaminationReviewsService: supervisorExaminationReviewsService,
		OrganisationsService:                organisationsService,
		UsersService:                        usersService,
		HospitalService:                     hospitalService,
		ClinicsService:                      clinicsService,
		AuthorizationService:                authorizationService,
		OrthopedicSurgeryService:            orthopedicSurgeryService,
		EvaluationFormService:               evaluationFormService,
		LogbookEntriesService:               logbookEntriesService,
		AssessmentsService:                  assessmentsSerrvice,
		NotificationService:                 notificationsService,
	}, nil
}

var ServiceSet = wire.NewSet(NewService,
	authorization2.NewService,
	activities.NewActivitiesService,
	practicalActivityTypes.NewService,
	procedures.NewService,
	examinations.NewService,
	examinationsactivities.NewService,
	examinationsactivitiesreviews.NewService,
	organisations.NewService,
	organisations.NewAdder,
	organisations.NewDeleter,
	organisations.NewGetter,
	organisations.NewUpdater,
	orthopedicSurgeries.NewActivities,
	orthopedicSurgeries.NewReviews,
	orthopedicSurgeries.NewSurgeries,
	orthopedicSurgeries.NewService,
	users.NewService,
	hospitals.NewService,
	clinics.NewService,
	evaluationforms.NewEvaluationForms,
	evaluationforms.NewDops,
	evaluationforms.NewMiniCex,
	evaluationforms.NewService,
	assessments.NewAssessmentsService,
	logbookentries.NewActivityService,
	notifications2.NewNotificationsService,
	thirdparty2.NewAuth0ManagementApi)
