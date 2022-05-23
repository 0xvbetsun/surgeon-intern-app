package ogbook

import (
	"github.com/google/wire"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/OrganizationalUnits"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/activities"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/assessments"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/clinics"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/dbexecutor"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/evaluationforms"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/examinationactivity"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/examinations"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/examinationsactivitiesreviews"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/examinationsactivityreview"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/hospitals"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/logbookentries"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/notifications"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/organisations"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/orthopedicSurgeries"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/practicalactivitytypes"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/procedures"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/roles"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/specialties"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/surgeries"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/users"
)

type (
	Repo struct {
		ActivitiesRepo                               activities.IActivitiesRepo
		LogbookEntriesRepo                           logbookentries.ILogbookEntriesRepo
		AssessmentsRepo                              assessments.IAssessmentsRepo
		PracticalActivityTypesRepo                   practicalactivitytypes.IRepo
		ProceduresRepo                               procedures.IRepo
		ExaminationsRepo                             examinations.IRepo
		HospitalsRepo                                hospitals.IRepo
		ClinicsRepo                                  clinics.IRepo
		OrganisationsRepo                            *organisations.Repo
		UsersRepo                                    users.IRepo
		ResidentExaminationsRepo                     examinationactivity.IRepo
		SupervisorExaminationReviewsRepo             examinationsactivityreview.IRepo
		SupervisorResidentExaminationConnectionsRepo examinationsactivitiesreviews.IRepo
		OrthopedicSurgeriesRepo                      *orthopedicSurgeries.Repo
		SurgeriesRepo                                surgeries.IRepo
		DbExecutor                                   dbexecutor.IDBExecutor
		RoleRepo                                     roles.IRepo
		SpecialtiesRepo                              specialties.IRepo
		EvaluationFormsRepo                          evaluationforms.Repo
		NotificationsRepo                            notifications.INotificationsRepo
		OrganizationalUnitTypesRepo                  *OrganizationalUnits.TypesRepo
		UserOrgUnitRoleRepo                          *users.UserOrganizationalUnitRoleRepo
		OrgUnitRepo                                  *OrganizationalUnits.Repo
	}
)

func NewRepo(activitiesRepo activities.IActivitiesRepo,
	logbookEntriesRepo logbookentries.ILogbookEntriesRepo,
	assessmentsRepo assessments.IAssessmentsRepo,
	practicalStepsRepo practicalactivitytypes.IRepo,
	proceduresRepo procedures.IRepo,
	examinationsRepo examinations.IRepo,
	hospitalsRepo hospitals.IRepo,
	clinicsRepo clinics.IRepo,
	organisationsRepo *organisations.Repo,
	usersRepo users.IRepo,
	residentExaminationsRepo examinationactivity.IRepo,
	supervisorExaminationReviewsRepo examinationsactivityreview.IRepo,
	supervisorResidentExaminationConnectionsRepo examinationsactivitiesreviews.IRepo,
	orthopedicSurgeriesRepo *orthopedicSurgeries.Repo,
	surgeriesRepo surgeries.IRepo,
	dbExecutor dbexecutor.IDBExecutor,
	roleRepo roles.IRepo,
	specialtiesRepo specialties.IRepo,
	evaluationFormsRepo evaluationforms.Repo,
	notificationsRepo notifications.INotificationsRepo,
	orgTypesRepo *OrganizationalUnits.TypesRepo,
	userOrgUnitRoleRepo *users.UserOrganizationalUnitRoleRepo,
	orgUnitRepo *OrganizationalUnits.Repo) *Repo {
	//// examinations
	//examinationsRepo := examinations.NewRepo(db)
	//
	//// procedures
	//proceduresRepo := procedures.NewRepo(db)
	//
	//// practical steps
	//practicalStepsRepo := practicalactivitytypes.NewRepo(db)
	//
	//// user examinations
	//residentExaminationsRepo := examinationactivity.NewRepo(db)
	//supervisorExaminationReviewsRepo := examinationsactivityreview.NewRepo(db)
	//supervisorResidentExaminationConnectionsRepo := examinationsactivitiesreviews.NewRepo(db)
	//
	//// orthopedicSurgeriesRepo
	//orthopedicSurgeriesActivitiesRepo := orthopedicSurgeries.NewActivitiesRepo(db)
	//orthopedicSurgeriesReviewsRepo := orthopedicSurgeries.NewReviewsRepo(db)
	//orthopedicSurgeriesRepo := orthopedicSurgeries.NewRepo(orthopedicSurgeriesActivitiesRepo, orthopedicSurgeriesReviewsRepo)
	//
	//// Specialty
	//specialtiesRepo := specialties.NewRepo(db)
	//
	//// Surgerys
	//// Surgeries
	//surgeriesRepo := surgeries.NewRepo(db)
	//
	//// organisations
	//organisationsGetter := organisations.NewGetter(db, ctx)
	//organisationsAdder := organisations.NewAdder(db, ctx)
	//organisationsRepo := organisations.NewRepo(organisationsGetter, nil, organisationsAdder, nil)
	//
	//// users
	//usersRepo := users.NewRepo(db)
	//
	//// roles
	//roleRepo := roles.NewRepo(db)
	//
	//// hospitals
	//hospitalsRepo := hospitals.NewRepo(db)
	//clinicsRepo := clinics.NewRepo(db)
	//dbExecutor := dbexecutor.NewDBExecutor(db)
	//
	//// Evaluation forms
	//evaluationForms := evaluationforms.NewEvaluationFormsRepo(db)
	//dops := evaluationforms.NewDopsRepo(db)
	//miniCex := evaluationforms.NewMiniCexRepo(db)
	//evaluationFormsRepo := evaluationforms.NewRepo(evaluationForms, dops, miniCex)

	// Notifications
	//notificationsrepo := notifications.NewNotificationsRepo(db)

	return &Repo{
		ActivitiesRepo:                               activitiesRepo,
		LogbookEntriesRepo:                           logbookEntriesRepo,
		AssessmentsRepo:                              assessmentsRepo,
		PracticalActivityTypesRepo:                   practicalStepsRepo,
		ProceduresRepo:                               proceduresRepo,
		ExaminationsRepo:                             examinationsRepo,
		ResidentExaminationsRepo:                     residentExaminationsRepo,
		SupervisorExaminationReviewsRepo:             supervisorExaminationReviewsRepo,
		SupervisorResidentExaminationConnectionsRepo: supervisorResidentExaminationConnectionsRepo,
		OrganisationsRepo:                            organisationsRepo,
		HospitalsRepo:                                hospitalsRepo,
		ClinicsRepo:                                  clinicsRepo,
		UsersRepo:                                    usersRepo,
		OrthopedicSurgeriesRepo:                      orthopedicSurgeriesRepo,
		SurgeriesRepo:                                surgeriesRepo,
		DbExecutor:                                   dbExecutor,
		RoleRepo:                                     roleRepo,
		SpecialtiesRepo:                              specialtiesRepo,
		EvaluationFormsRepo:                          evaluationFormsRepo,
		NotificationsRepo:                            notificationsRepo,
		OrganizationalUnitTypesRepo:                  orgTypesRepo,
		UserOrgUnitRoleRepo:                          userOrgUnitRoleRepo,
		OrgUnitRepo:                                  orgUnitRepo,
	}

}

var RepoSet = wire.NewSet(NewRepo,
	activities.NewActivitiesRepo,
	logbookentries.NewLogbookEntriesRepo,
	assessments.NewAssessmentsRepo,
	examinations.NewRepo,
	procedures.NewRepo,
	users.NewUserOrganizationalUnitRoleRepo,
	practicalactivitytypes.NewRepo,
	OrganizationalUnits.NewTypesRepo,
	OrganizationalUnits.NewRepo,
	examinationactivity.NewRepo,
	examinationsactivityreview.NewRepo,
	examinationsactivitiesreviews.NewRepo,
	orthopedicSurgeries.NewActivitiesRepo,
	orthopedicSurgeries.NewReviewsRepo,
	orthopedicSurgeries.NewRepo,
	specialties.NewRepo,
	surgeries.NewRepo,
	organisations.NewGetter,
	organisations.NewAdder,
	organisations.NewRepo,
	organisations.NewUpdater,
	organisations.NewDeleter,
	users.NewRepo,
	roles.NewRepo,
	hospitals.NewRepo,
	clinics.NewRepo,
	dbexecutor.NewDBExecutor,
	evaluationforms.NewEvaluationFormsRepo,
	evaluationforms.NewDopsRepo,
	evaluationforms.NewMiniCexRepo,
	evaluationforms.NewRepo,
	notifications.NewNotificationsRepo)
