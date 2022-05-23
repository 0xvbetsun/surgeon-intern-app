package seeder

import (
	"time"

	"github.com/thoas/go-funk"
	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	qlmodel "github.com/vbetsun/surgeon-intern-app/internal/logbookService/graph/model"
	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/volatiletech/null/v8"
)

type BaseSeed struct {
	OrganizationalUnitTypes  []*models.OrganizationalUnitType       `json:"organizational_unit_types"`
	PracticalActivityTypes   []*models.PracticalActivityType        `json:"practical_activity_types"`
	Specialties              []*models.Specialty                    `json:"specialties"`
	SpecialtiesActivityTypes []*SpecialtiesActivityTypes            `json:"specialties_activity_types"`
	Roles                    []*models.Role                         `json:"roles"`
	ExaminationAnnotations   []commonModel.ExaminationAnnotations   `json:"examination_annotations"`
	ProcedureAnnotations     []qlmodel.ProcedureAnnotations         `json:"procedure_annotations"`
	DopsAnnotations          []commonModel.EvaluationFormAnnotation `json:"dops_annotations"`
	MiniKexAnnotations       []commonModel.EvaluationFormAnnotation `json:"mini_kex_annotations"`
	MiniKexAreas             []*models.MiniCexArea                  `json:"mini_kex_areas"`
	MiniKexFocues            []*models.MiniCexFocuse                `json:"mini_kex_focues"`
}

type SpecialtiesActivityTypes struct {
	Specialty             *models.Specialty
	PracticalActivityType *models.PracticalActivityType
}

var TestClinics = []*models.OrganizationalUnit{
	{
		ID:          "ede88e55-b7c9-4b37-8a6e-b90480e1b4e6",
		DisplayName: "Akutsjukvård",
		ParentID:    null.StringFrom(TestHospitals[0].ID),
	},
	{
		ID:          "991418c5-0096-4442-8ac8-c39ffdbcebba",
		DisplayName: "Ortopedi",
		ParentID:    null.StringFrom(TestHospitals[0].ID),
	},
	{
		ID:          "48347e6e-fa52-44a7-81ad-95fd34bd56a6",
		DisplayName: "Akutsjukvård",
		ParentID:    null.StringFrom(TestHospitals[1].ID),
	},
	{
		ID:          "895efb93-54ce-403e-8d9f-c2ba492e98fc",
		DisplayName: "Ortopedi",
		ParentID:    null.StringFrom(TestHospitals[1].ID),
	},
}

func GetTestClinicSpecialties(specialties []*models.Specialty) map[*models.OrganizationalUnit]*models.Specialty {
	return map[*models.OrganizationalUnit]*models.Specialty{
		funk.Find(TestClinics, func(clinic *models.OrganizationalUnit) bool {
			return clinic.ID == "895efb93-54ce-403e-8d9f-c2ba492e98fc"
		}).(*models.OrganizationalUnit): funk.Find(specialties, func(specialty *models.Specialty) bool {
			return specialty.Name == commonModel.SpecialtiesOrthopedics.String()
		}).(*models.Specialty),
		funk.Find(TestClinics, func(clinic *models.OrganizationalUnit) bool {
			return clinic.ID == "991418c5-0096-4442-8ac8-c39ffdbcebba"
		}).(*models.OrganizationalUnit): funk.Find(specialties, func(specialty *models.Specialty) bool {
			return specialty.Name == commonModel.SpecialtiesOrthopedics.String()
		}).(*models.Specialty),
	}
}

var TestHospitals = []models.OrganizationalUnit{
	{
		ID:          "14dc6697-a14e-4c42-b012-092cd1902e3a",
		DisplayName: "NÄL",
	},
	{
		ID:          "7d8ecbd5-68d9-4a43-8ff2-965ef5528639",
		DisplayName: "Sahlgrenska",
	},
}

var TestOrthopedicClinicDepartments = []models.OrganizationalUnit{
	{
		ID:          GetUuid(),
		DisplayName: "Mottagning",
		ParentID:    null.StringFrom(TestClinics[1].ID),
	},
	{
		ID:          GetUuid(),
		DisplayName: "Akuten",
		ParentID:    null.StringFrom(TestClinics[1].ID),
	},
	{
		ID:          GetUuid(),
		DisplayName: "Operation",
		ParentID:    null.StringFrom(TestClinics[1].ID),
	},
	{
		ID:          GetUuid(),
		DisplayName: "Traumarum",
		ParentID:    null.StringFrom(TestClinics[1].ID),
	},
	{
		ID:          GetUuid(),
		DisplayName: "Mottagning",
		ParentID:    null.StringFrom(TestClinics[3].ID),
	},
	{
		ID:          GetUuid(),
		DisplayName: "Akuten",
		ParentID:    null.StringFrom(TestClinics[3].ID),
	},
	{
		ID:          GetUuid(),
		DisplayName: "Operation",
		ParentID:    null.StringFrom(TestClinics[3].ID),
	},
	{
		ID:          GetUuid(),
		DisplayName: "Traumarum",
		ParentID:    null.StringFrom(TestClinics[3].ID),
	},
}

var TestUsers = []models.User{
	{
		ID:             "8c1f3668-7e08-451a-9638-cb2f1d56f93a",
		DisplayName:    "Linus Swahn",
		UserExternalID: null.StringFrom("auth0|60cf4185ade114007146edd1"),
		Email:          "linus@swahn.io",
		Activated:      true,
	},
	{
		ID:             "6d357302-5b72-4369-ba72-0da19aef301e",
		DisplayName:    "Pontus Pohl",
		UserExternalID: null.StringFrom("auth0|60cf41b663682e0070471e07"),
		Email:          "pontus.pohl@gmail.com",
		Activated:      true,
	},
	{
		ID:             "76171d5d-0255-4591-ac34-c900b23e32a2",
		DisplayName:    "Sadesh",
		UserExternalID: null.StringFrom("auth0|60cf4957b72a9f006a28e1e7"),
		Email:          "sadeshb@gmail.com",
		Activated:      true,
	},
	{
		ID:             "44708f63-6aa7-4578-a0b8-11f9b6a01624",
		DisplayName:    "Divyen",
		UserExternalID: null.StringFrom("auth0|60cf497163682e0070471fec"),
		Email:          "divyen.sisodia@gmail.com",
		Activated:      true,
	},
	{
		ID:             "89147aab-6bed-4dc1-83dc-8d3aec040c38",
		DisplayName:    "David Häggå",
		UserExternalID: null.StringFrom("auth0|60cf4989ef2a5e0068d63673"),
		Email:          "davidhaegga@gmail.com",
		Activated:      true,
	},
	{
		ID:             "c155aec7-941c-4aae-b4aa-48d7b966a7e8",
		DisplayName:    "Mr Cypress",
		UserExternalID: null.StringFrom("auth0|60cf4d970eb20b00699833bf"),
		Email:          "cypress@cypress.io",
		Activated:      true,
	},
}

func GetTestUserOrganizationalUnitRoles(roles []*models.Role) []models.UserOrganizationalUnitRole {
	return []models.UserOrganizationalUnitRole{
		{
			UserID: TestUsers[0].ID,
			UnitID: TestClinics[0].ID,
			RoleID: funk.Find(roles, func(role *models.Role) bool {
				return role.Name == "Resident"
			}).(*models.Role).ID,
		},
		{
			UserID: TestUsers[0].ID,
			UnitID: TestClinics[1].ID,
			RoleID: funk.Find(roles, func(role *models.Role) bool {
				return role.Name == "Resident"
			}).(*models.Role).ID,
		},
		{
			UserID: TestUsers[1].ID,
			UnitID: TestClinics[0].ID,
			RoleID: funk.Find(roles, func(role *models.Role) bool {
				return role.Name == "Resident"
			}).(*models.Role).ID,
		},
		{
			UserID: TestUsers[1].ID,
			UnitID: TestClinics[1].ID,
			RoleID: funk.Find(roles, func(role *models.Role) bool {
				return role.Name == "Resident"
			}).(*models.Role).ID,
		},
		{
			UserID: TestUsers[2].ID,
			UnitID: TestClinics[0].ID,
			RoleID: funk.Find(roles, func(role *models.Role) bool {
				return role.Name == "Resident"
			}).(*models.Role).ID,
		},
		{
			UserID: TestUsers[2].ID,
			UnitID: TestClinics[1].ID,
			RoleID: funk.Find(roles, func(role *models.Role) bool {
				return role.Name == "Resident"
			}).(*models.Role).ID,
		},
		{
			UserID: TestUsers[3].ID,
			UnitID: TestClinics[0].ID,
			RoleID: funk.Find(roles, func(role *models.Role) bool {
				return role.Name == "Resident"
			}).(*models.Role).ID,
		},
		{
			UserID: TestUsers[3].ID,
			UnitID: TestClinics[1].ID,
			RoleID: funk.Find(roles, func(role *models.Role) bool {
				return role.Name == "Resident"
			}).(*models.Role).ID,
		},
		{
			UserID: TestUsers[4].ID,
			UnitID: TestClinics[0].ID,
			RoleID: funk.Find(roles, func(role *models.Role) bool {
				return role.Name == "Resident"
			}).(*models.Role).ID,
		},
		{
			UserID: TestUsers[4].ID,
			UnitID: TestClinics[1].ID,
			RoleID: funk.Find(roles, func(role *models.Role) bool {
				return role.Name == "Resident"
			}).(*models.Role).ID,
		},
		{
			UserID: TestUsers[0].ID,
			UnitID: TestClinics[0].ID,
			RoleID: funk.Find(roles, func(role *models.Role) bool {
				return role.Name == "Supervisor"
			}).(*models.Role).ID,
		},
		{
			UserID: TestUsers[0].ID,
			UnitID: TestClinics[1].ID,
			RoleID: funk.Find(roles, func(role *models.Role) bool {
				return role.Name == "Supervisor"
			}).(*models.Role).ID,
		},
		{
			UserID: TestUsers[1].ID,
			UnitID: TestClinics[0].ID,
			RoleID: funk.Find(roles, func(role *models.Role) bool {
				return role.Name == "Supervisor"
			}).(*models.Role).ID,
		},
		{
			UserID: TestUsers[1].ID,
			UnitID: TestClinics[1].ID,
			RoleID: funk.Find(roles, func(role *models.Role) bool {
				return role.Name == "Supervisor"
			}).(*models.Role).ID,
		},
		{
			UserID: TestUsers[2].ID,
			UnitID: TestClinics[0].ID,
			RoleID: funk.Find(roles, func(role *models.Role) bool {
				return role.Name == "Supervisor"
			}).(*models.Role).ID,
		},
		{
			UserID: TestUsers[2].ID,
			UnitID: TestClinics[1].ID,
			RoleID: funk.Find(roles, func(role *models.Role) bool {
				return role.Name == "Supervisor"
			}).(*models.Role).ID,
		},
		{
			UserID: TestUsers[3].ID,
			UnitID: TestClinics[0].ID,
			RoleID: funk.Find(roles, func(role *models.Role) bool {
				return role.Name == "Supervisor"
			}).(*models.Role).ID,
		},
		{
			UserID: TestUsers[3].ID,
			UnitID: TestClinics[1].ID,
			RoleID: funk.Find(roles, func(role *models.Role) bool {
				return role.Name == "Supervisor"
			}).(*models.Role).ID,
		},
		{
			UserID: TestUsers[4].ID,
			UnitID: TestClinics[0].ID,
			RoleID: funk.Find(roles, func(role *models.Role) bool {
				return role.Name == "Supervisor"
			}).(*models.Role).ID,
		},
		{
			UserID: TestUsers[4].ID,
			UnitID: TestClinics[1].ID,
			RoleID: funk.Find(roles, func(role *models.Role) bool {
				return role.Name == "Supervisor"
			}).(*models.Role).ID,
		},
		{
			UserID: TestUsers[5].ID,
			UnitID: TestClinics[0].ID,
			RoleID: funk.Find(roles, func(role *models.Role) bool {
				return role.Name == "Supervisor"
			}).(*models.Role).ID,
		},
		{
			UserID: TestUsers[5].ID,
			UnitID: TestClinics[1].ID,
			RoleID: funk.Find(roles, func(role *models.Role) bool {
				return role.Name == "Supervisor"
			}).(*models.Role).ID,
		},
		{
			UserID: TestUsers[5].ID,
			UnitID: TestClinics[0].ID,
			RoleID: funk.Find(roles, func(role *models.Role) bool {
				return role.Name == "Resident"
			}).(*models.Role).ID,
		},
		{
			UserID: TestUsers[5].ID,
			UnitID: TestClinics[1].ID,
			RoleID: funk.Find(roles, func(role *models.Role) bool {
				return role.Name == "Resident"
			}).(*models.Role).ID,
		},
	}
}

var TestOrthohopedicSurgeryActivities = []*models.OrthopedicSurgeryActivity{{
	ID:                      "5ae1d208-a92e-4683-b249-5fe89d090307",
	OccurredAt:              time.Now(),
	CaseNotes:               "123",
	PatientAge:              60,
	PatientGender:           "male",
	ResidentID:              TestUsers[5].ID,
	SupervisorID:            null.StringFrom(TestUsers[5].ID),
	OperatorID:              null.StringFrom(TestUsers[1].ID),
	AssistantID:             null.StringFrom(TestUsers[0].ID),
	Comments:                "Bla Bla",
	Complications:           "Aj Aj AJ",
	PracticalActivityTypeID: 3,
}}

var TestEvaluationForms = []*models.EvaluationForm{
	{
		ID:           1,
		DepartmentID: TestClinics[1].ID,
		Name:         "DOPS",
		Difficulty: []string{
			"Lätt",
			"Medel",
			"Svår",
			"Specialistnivå",
			"Subspecialistnivå",
		},
		Citations: []string{
			"Lowry J, Cripps J. Results of the online EWTD trainee survey. Ann R Coll Surg Engl Suppl 2005; 87:86-87.",
			"Delamothe T. Modernising Medical Careers: final report. BMJ 2008; 336:54-5.",
			"Beard J, Strachan A, Davies H, Patterson F, Stark P, Ball S, et al. Developing an education and assessment framework for the Foundation Programme. Med Educ 2005; 39:841-51.",
			"Siau K, Crossley J, Dunckley P, Johnson G, Feeney M, Hawkes ND, Beales ILP; Joint Advisory Group on Gastrointestinal Endoscopy (JAG). Direct observation of procedural skills (DOPS) assessment in diagnostic gastroscopy: nationwide evidence of validity and competency development during training. Surg Endosc. 2020 Jan;34(1):105-114. doi: 10.1007/s00464-019-06737-7. Epub 2019 Mar 25. Erratum in: Surg Endosc. 2019 Apr 1;: PMID: 30911922; PMCID: PMC6946748.",
			"Siau K, Crossley J, Dunckley P, Johnson G, Feeney M, Hawkes ND, Beales ILP; Joint Advisory Group on Gastrointestinal Endoscopy (JAG). Correction to: Direct observation of procedural skills (DOPS) assessment in diagnostic gastroscopy: nationwide evidence of validity and competency development during training. Surg Endosc. 2020 Jan;34(1):115. doi: 10.1007/s00464-019-06776-0. Erratum for: Surg Endosc. 2020 Jan;34(1):105-114. PMID: 30937617; PMCID: PMC7050556.",
			"Siau K, Crossley J, Dunckley P, et alPWE-109 Direct observation of procedural skills (dops) assessment in colonoscopy: validity and competency development during trainingGut 2019;68:A253-A254.",
			"Morris A, Hewitt J, Roberts CM. Practical experience of using directly observed procedures, mini clinical evaluation examinations, and peer observation in pre-registration house officer (FY1) trainees. Postgrad Med J 2006; 82:285-8.",
			"Naeem N. Validity, reliability, feasibility, acceptability and educational impact of direct observation of procedural skills (DOPS). J Coll Physicians Surg Pak. 2013 Jan;23(1):77-82. PMID: 23286629.",
			"Erfani Khanghahi, Masoumeh, and Farbod Ebadi Fard Azar. “Direct observation of procedural skills (DOPS) evaluation method: Systematic review of evidence.” Medical journal of the Islamic Republic of Iran vol. 32 45. 3 Jun. 2018, doi:10.14196/mjiri.32.45",
		},
	},
	{
		ID:           2,
		DepartmentID: TestClinics[1].ID,
		Name:         "MINI-CEX",
		Difficulty: []string{
			"Lätt",
			"Medel",
			"Svår",
			"Specialistnivå",
			"Subspecialistnivå",
		},
		Citations: []string{
			"Norcini JJ, Blank LL, Arnold GK, Kimball HR: The mini-CEX (clinical evaluation exercise): a preliminary investigation. Ann Intern Med. 1995, 123 (10): 795-799.",
			"Frank JR, Danoff D: The CanMEDS initiative: implementing an outcomes-based framework of physician competencies. Med Teach. 2007, 29 (7): 642-647. 10.1080/01421590701746983.",
			"Liao, KC., Pu, SJ., Liu, MS. et al. Development and implementation of a mini-Clinical Evaluation Exercise (mini-CEX) program to assess the clinical competencies of internal medicine residents: from faculty development to curriculum evaluation. BMC Med Educ 13, 31 (2013). https://doi.org/10.1186/1472-6920-13-31",
			"Swing SR: The ACGME outcome project: retrospective and prospective. Med Teach. 2007, 29 (7): 648-654. 10.1080/01421590701392903.",
			"Kogan JR, Holmboe ES, Hauer KE: Tools for direct observation and assessment of clinical skills of medical trainees: a systematic review. JAMA. 2009, 302 (12): 1316-1326. 10.1001/jama.2009.1365.",
			"Pelgrim EA, Kramer AW, Mokkink HG, van den Elsen L, Grol RP, van der Vleuten CP: In-training assessment using direct observation of single-patient encounters: a literature review. Adv Health Sci Educ Theory Pract. 2011, 16 (1): 131-142. 10.1007/s10459-010-9235-6.",
			"Gozu A, Windish DM, Knight AM, Thomas PA, Kolodner K, Bass EB, Sisson SD, Kern DE: Long-term follow-up of a 10-month programme in curriculum development for medical educators: a cohort study. Med Educ. 2008, 42 (7): 684-692. 10.1111/j.1365-2923.2008.03090.x.",
			"Kogan JR, Bellini LM, Shea JA: Implementation of the mini-CEX to evaluate medical students’ clinical skills. Acad Med. 2002, 77 (11): 1156-1157.",
			"Weller JM1, Misur M2, Nicolson S2, Morris J3, Ure S4, Crossley J5, Jolly B6. Can I leave the theatre? A key to more reliable workplace-based assessment. Br J Anaesth. 2014 Jun;112(6):1083-91. doi: 10.1093/bja/aeu052. Epub 2014 Mar 17.",
			"Crossley J1, Johnson G, Booth J, Wade W. Good questions, good answers: construct alignment improves the performance of workplace- based assessment scales. Med Educ. 2011 Jun;45(6):560-9. doi: 10.1111/j.1365-2923.2010.03913.x. Epub 2011 Apr 18.",
			"F. Walentin, Modifierad MINI-Cex, hjarntanken.se",
		},
	},
}
