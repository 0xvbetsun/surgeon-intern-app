package seeder

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/ogbook"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/users"
	"go.uber.org/zap"
)

type (
	Seeder struct {
		repo *ogbook.Repo
	}
)

func GetUuid() string {
	retUuid, _ := uuid.NewV4()
	return retUuid.String()
}

func New(repo *ogbook.Repo) *Seeder {
	return &Seeder{repo: repo}
}

func (s Seeder) SeedDevelopmentData(baseSeedJsonPath string) {

	var baseSeed = loadBaseSeedFromJson(baseSeedJsonPath)

	for _, testUnitType := range baseSeed.OrganizationalUnitTypes {
		seedTestOrganizationalUnitTypes(s, testUnitType)
	}
	for _, practicalActivityType := range baseSeed.PracticalActivityTypes {
		seedPracticalActivityTypes(s, practicalActivityType)
	}
	for _, testHospital := range TestHospitals {
		seedTestHospital(s, testHospital)
	}
	for _, testClinic := range TestClinics {
		seedTestClinic(s, testClinic)
	}
	for _, testClinicLocation := range TestOrthopedicClinicDepartments {
		seedTestClinicDepartments(s, testClinicLocation)
	}
	for _, testRole := range baseSeed.Roles {
		seedTestRole(s, testRole)
	}
	for _, testUser := range TestUsers {
		seedTestUser(s, testUser)
	}

	seedSpecialties(s, baseSeed.Specialties)
	seedSpecialtiesTypes(s, baseSeed.SpecialtiesActivityTypes)
	seedClinicSpecialties(s, GetTestClinicSpecialties(baseSeed.Specialties))
	for _, testOrthopedicSurgeryActivity := range TestOrthohopedicSurgeryActivities {
		seedOrthopedicSurgeryActivity(s, testOrthopedicSurgeryActivity)
	}

	for _, testSupervisor := range GetTestUserOrganizationalUnitRoles(baseSeed.Roles) {
		seedTestSupervisor(s, testSupervisor)
	}

	seedEvaluationForm(s, TestEvaluationForms[0], baseSeed.DopsAnnotations)
	seedEvaluationForm(s, TestEvaluationForms[1], baseSeed.MiniKexAnnotations)

	for _, focus := range baseSeed.MiniKexFocues {
		seedMiniCexFocus(s, focus)
	}
	for _, area := range baseSeed.MiniKexAreas {
		seedMiniCexArea(s, area)
	}
}

func loadBaseSeedFromJson(path string) *BaseSeed {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		zap.S().Fatalf(err.Error())
	}
	var baseSeed BaseSeed
	err = json.Unmarshal(file, &baseSeed)
	if err != nil {
		zap.S().Fatalf(err.Error())
	}
	return &baseSeed
}

func seedTestOrganizationalUnitTypes(s Seeder, ouType *models.OrganizationalUnitType) {
	_, err := s.repo.OrganizationalUnitTypesRepo.Add(context.TODO(), ouType)
	if err != nil {
		zap.S().Error(err.Error())
	}
}

func seedTestRole(s Seeder, role *models.Role) {
	_, err := s.repo.RoleRepo.GetByID(context.TODO(), role.ID)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			_, err = s.repo.RoleRepo.Add(context.TODO(), role)
			if err != nil {
				zap.S().Fatalf("Failed to insert role '%s': %s", role.DisplayName, err.Error())
			}
		} else {
			zap.S().Fatalf("Failed to get test role '%s': %s", role.DisplayName, err.Error())
		}
	}
}

func seedTestUser(s Seeder, user models.User) {
	_, err := s.repo.UsersRepo.GetByID(context.TODO(), user.ID, false)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			_, err = s.repo.UsersRepo.Add(context.TODO(), &user)
			if err != nil {
				zap.S().Fatalf("Failed to insert user '%s': %s", user.DisplayName, err.Error())
			}
		} else {
			zap.S().Fatalf("Failed to get test user '%s': %s", user.DisplayName, err.Error())
		}
	}
}

func seedTestSupervisor(s Seeder, clinicRole models.UserOrganizationalUnitRole) {
	userRoles, _ := s.repo.UserOrgUnitRoleRepo.GetByFilter(context.TODO(), users.UserOrganizationalUnitRoleFilter{
		UserId: &clinicRole.UserID,
		RoleId: nil,
		UnitId: nil,
	}, false)
	for _, v := range userRoles {
		if v.UnitID == clinicRole.UnitID && v.RoleID == clinicRole.RoleID {
			return
		}
	}
	_, err := s.repo.ClinicsRepo.AddUserClinicRole(context.TODO(), &clinicRole)
	if err != nil {
		zap.S().Fatalf(err.Error())
	}
}

func seedTestHospital(s Seeder, hospital models.OrganizationalUnit) {
	_, err := s.repo.HospitalsRepo.GetByID(context.TODO(), hospital.ID)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			_, err = s.repo.HospitalsRepo.Add(context.TODO(), &hospital)
			if err != nil {
				zap.S().Fatalf("Failed to insert hospital '%s': %s", hospital.DisplayName, err.Error())
			}
		} else {
			zap.S().Fatalf("Failed to get test hospital '%s': %s", hospital.DisplayName, err.Error())
		}
	}
}

func seedTestClinic(s Seeder, clinic *models.OrganizationalUnit) {
	_, err := s.repo.ClinicsRepo.GetByID(context.TODO(), clinic.ID)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			_, err = s.repo.ClinicsRepo.Add(context.TODO(), clinic)
			if err != nil {
				zap.S().Fatalf("Failed to insert clinic '%s': %s", clinic.DisplayName, err.Error())
			}
		} else {
			zap.S().Fatalf("Failed to get test clinic '%s': %s", clinic.DisplayName, err.Error())
		}
	}
}

func seedTestClinicDepartments(s Seeder, location models.OrganizationalUnit) {
	_, err := s.repo.ClinicsRepo.GetClinicDepartmentByID(context.TODO(), location.ID)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			_, err = s.repo.ClinicsRepo.AddClinicDepartment(context.TODO(), &location)
			if err != nil {
				zap.S().Fatalf("Failed to insert clinic location '%s': %s", location.DisplayName, err.Error())
			}
		} else {
			zap.S().Fatalf("Failed to get test clinic location '%s': %s", location.DisplayName, err.Error())
		}
	}
}

func seedMiniCexFocus(s Seeder, focus *models.MiniCexFocuse) {
	_, err := s.repo.EvaluationFormsRepo.MiniCex.GetMiniCexFocusById(context.TODO(), focus.ID)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			_, err = s.repo.EvaluationFormsRepo.MiniCex.AddMiniCexFocus(context.TODO(), focus)
			if err != nil {
				zap.S().Fatalf("Failed to insert mini-cex focus '%s': %s", focus.Name, err.Error())
			}
		} else {
			zap.S().Fatalf("Failed to get mini-cex focus '%s': %s", focus.Name, err.Error())
		}
	}
}

func seedMiniCexArea(s Seeder, area *models.MiniCexArea) {
	_, err := s.repo.EvaluationFormsRepo.MiniCex.GetMiniCexAreaById(context.TODO(), area.ID)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			_, err = s.repo.EvaluationFormsRepo.MiniCex.AddMiniCexArea(context.TODO(), area)
			if err != nil {
				zap.S().Fatalf("Failed to insert mini-cex area '%s': %s", area.Name, err.Error())
			}
		} else {
			zap.S().Fatalf("Failed to get mini-cex area '%s': %s", area.Name, err.Error())
		}
	}
}

func seedPracticalActivityTypes(s Seeder, practicalActivityType *models.PracticalActivityType) {
	_, err := s.repo.PracticalActivityTypesRepo.Get(context.TODO(), practicalActivityType.ID)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			_, err = s.repo.PracticalActivityTypesRepo.Add(context.TODO(), practicalActivityType)
			if err != nil {
				zap.S().Fatalf("Failed to insert user '%s': %s", practicalActivityType.DisplayName, err.Error())
			}
		} else {
			zap.S().Fatalf("Failed to get test user '%s': %s", practicalActivityType.DisplayName, err.Error())
		}
	}
}

func seedClinicSpecialties(s Seeder, clinicSpecialties map[*models.OrganizationalUnit]*models.Specialty) {
	for clinic, specialty := range clinicSpecialties {
		err := s.repo.OrgUnitRepo.ConnectUnitToSpecialty(context.TODO(), clinic, specialty)
		if err != nil {
			zap.S().Fatal(err.Error())
		}
	}
}

func seedSpecialties(s Seeder, specialties []*models.Specialty) {
	for _, specialty := range specialties {
		err := s.repo.SpecialtiesRepo.Add(context.TODO(), s.repo.DbExecutor.GetDB(), specialty)
		if err != nil {
			zap.S().Fatalf(err.Error())
		}
	}
}

func seedSpecialtiesTypes(s Seeder, specialtiesTypes []*SpecialtiesActivityTypes) {
	for _, specialtyActivityType := range specialtiesTypes {
		err := s.repo.SpecialtiesRepo.ConnectSpecialtyToActivityType(context.TODO(), specialtyActivityType.Specialty, specialtyActivityType.PracticalActivityType)
		if err != nil {
			zap.S().Fatalf("Failed to connect activityType to specialty: %s", err.Error())
		}
	}
}

func seedOrthopedicSurgeryActivity(s Seeder, activity *models.OrthopedicSurgeryActivity) {
	_, err := s.repo.OrthopedicSurgeriesRepo.ActivitiesRepo.GetByID(context.TODO(), activity.ID)
	if err != nil && errors.Cause(err) == sql.ErrNoRows {
		err = activity.Annotations.Marshal(commonModel.OrthopedicSurgeryActivityAnnotations{
			OperationPlanning:  nil,
			PatientPositioning: nil,
			Incision:           nil,
			Opening:            nil,
			Inspection:         nil,
			Repositioning:      nil,
			SawingAndDrillning: nil,
			Osteosyntes:        nil,
			Closing:            nil,
			Plastering:         nil,
			PostOp:             nil,
		})
		if err != nil {
			zap.S().Fatalf("Failed to marshal examination annotations for: %s", err.Error())
		}
		err = s.repo.DbExecutor.RunWithTX(context.TODO(), func(tx *sql.Tx) error {
			_, err = s.repo.OrthopedicSurgeriesRepo.ActivitiesRepo.Add(context.TODO(), tx, activity)
			if err != nil {
				zap.S().Info("Failed to insert orthopedicSurgeryActivity ")
				zap.S().Fatalf(err.Error())
			}
			var surgeries = make([]*models.Surgery, 0)
			surgeriesForSeeding := getSurgeriesForSeeding()
			for _, surgeryForSeeding := range surgeriesForSeeding {
				surgery, err := s.repo.SurgeriesRepo.Add(context.TODO(), surgeryForSeeding.method, surgeryForSeeding.diagnose, "89c91507-a2c9-4d58-b5de-d47054876192")

				if err != nil {
					zap.S().Fatalf(err.Error())
				}
				surgeries = append(surgeries, surgery)
			}
			err = s.repo.OrthopedicSurgeriesRepo.ActivitiesRepo.AddSurgeriesToActivity(context.TODO(), tx, activity.ID, surgeries)
			if err != nil {
				zap.S().Fatalf(err.Error())
			}

			return nil
		})
		if err != nil {
			zap.S().Fatalf(err.Error())
		}
	}

	return
}

func seedEvaluationForm(s Seeder, form *models.EvaluationForm, annotations []commonModel.EvaluationFormAnnotation) {
	_, err := s.repo.EvaluationFormsRepo.EvaluationForms.GetEvaluationFormById(context.TODO(), form.ID)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			err = form.Annotations.Marshal(annotations)
			if err != nil {
				zap.S().Fatalf("Failed to marshal activitiy annotations for: %s", form.Name, err.Error())
			}
			_, err = s.repo.EvaluationFormsRepo.EvaluationForms.AddEvaluationForm(context.TODO(), form)
			if err != nil {
				zap.S().Fatalf("Failed to insert DOPS form '%s': %s", form.Name, err.Error())
			}
		} else {
			zap.S().Fatalf("Failed to get test DOPS form '%s': %s", form.Name, err.Error())
		}
	}
}
