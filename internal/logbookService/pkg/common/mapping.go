package common

import (
	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/pkg/models"
)

func MapActivityGraphQlModel(activity *models.Activity) *commonModel.Activity {
	mappedActivity := &commonModel.Activity{
		OccurredAt: activity.OccurredAt,
	}

	if activity.R != nil {
		if activity.R.Assessment != nil {
			mappedActivity.Assessment = MapAssessmentGraphQlModel(activity.R.Assessment)
		}
		if activity.R.LogbookEntry != nil {
			mappedActivity.LogbookEntry = MapLogbookEntryGraphQlModel(activity.R.LogbookEntry)
		}
	}

	return mappedActivity
}

func MapAssessmentGraphQlModel(assessment *models.Assessment) *commonModel.Assessment {
	mappedAssessment := &commonModel.Assessment{}

	if assessment.R != nil {
		if assessment.R.Dop != nil {
			mappedAssessment.Dops = MapDopsEvaluationGraphQlModel(assessment.R.Dop, true)
		}
		if assessment.R.MiniCex != nil {
			mappedAssessment.MiniCex = MapMiniCexEvaluationGraphQlModel(assessment.R.MiniCex)
		}
		if assessment.R.OrthopedicSurgeryReview != nil {
			mappedAssessment.SurgeryReview = MapOrthopedicSurgeryReviewGraphQlModel(assessment.R.OrthopedicSurgeryReview, true)
		}
	}

	return mappedAssessment
}

func MapLogbookEntryGraphQlModel(logbookEntry *models.LogbookEntry) *commonModel.LogbookEntry {
	mappedLogbookEntry := &commonModel.LogbookEntry{}

	if logbookEntry.R != nil {
		if logbookEntry.R.OrthopedicSurgery != nil {
			mappedLogbookEntry.Surgery = MapOrthopedicSurgeryActivity(logbookEntry.R.OrthopedicSurgery, true, true)
		}
	}

	return mappedLogbookEntry
}

func MapMiniCexEvaluationGraphQlModel(evaluation *models.MiniCexEvaluation) *commonModel.MiniCexEvaluation {
	mappedEvaluation := &commonModel.MiniCexEvaluation{
		ID:            evaluation.ID,
		SupervisorID:  evaluation.SupervisorID,
		ResidentID:    evaluation.ResidentID,
		Difficulty:    &evaluation.Difficulty,
		Area:          evaluation.Area,
		Focuses:       evaluation.Focuses,
		DepartmentID:  &evaluation.DepartmentID.String,
		IsEvaluated:   evaluation.IsEvaluated,
		InProgress:    evaluation.InProgress,
		ActiveStep:    evaluation.ActiveStep,
		CompletedStep: evaluation.CompletedStep,
		OccurredAt:    evaluation.OccurredAt,
		CreatedAt:     evaluation.CreatedAt,
	}

	annotations := make([]*commonModel.EvaluationFormAnnotations, 0)
	if err := evaluation.Annotations.Unmarshal(&annotations); err != nil {
		// Todo handle err
	} else {
		mappedEvaluation.Annotations = annotations
	}

	mappedEvaluation.Description = &commonModel.Description{Title: "MINI-CEX", Subrows: []*commonModel.DescriptionRow{}}

	if evaluation.IsEvaluated {
		focusRow := &commonModel.DescriptionRow{}
		focuses := ""
		for index, focus := range evaluation.Focuses {
			if index > 0 {
				focuses = focuses + ", " + focus
			} else {
				focuses = focus
			}
		}
		focusRow.Title = "Fokus för evalueringen: " + focuses
		mappedEvaluation.Description.Subrows = append(mappedEvaluation.Description.Subrows, focusRow)

		areaRow := &commonModel.DescriptionRow{}
		areaRow.Title = "Område: " + evaluation.Area
		mappedEvaluation.Description.Subrows = append(mappedEvaluation.Description.Subrows, areaRow)
	} else {
		subRow := &commonModel.DescriptionRow{}
		subRow.Title = "Begärd men ej utförd"
		mappedEvaluation.Description.Subrows = append(mappedEvaluation.Description.Subrows, subRow)
	}

	if evaluation.R != nil {
		if evaluation.R.Department != nil {
			mappedEvaluation.Department = &commonModel.ClinicDepartment{
				DepartmentID:   evaluation.R.Department.ID,
				DepartmentName: evaluation.R.Department.DisplayName,
				ClinicID:       evaluation.R.Department.ParentID.String,
			}
		}

		if evaluation.R.Supervisor != nil && evaluation.R.Resident != nil {
			mappedEvaluation.Resident = &commonModel.User{
				UserID:      evaluation.R.Resident.ID,
				DisplayName: evaluation.R.Resident.DisplayName,
			}
			mappedEvaluation.Supervisor = &commonModel.User{
				UserID:      evaluation.R.Supervisor.ID,
				DisplayName: evaluation.R.Supervisor.DisplayName,
			}
		}
	}

	return mappedEvaluation
}

func MapDopsEvaluationGraphQlModel(evaluation *models.DopsEvaluation, mapActivity bool) *commonModel.DopsEvaluation {
	mappedEvaluation := &commonModel.DopsEvaluation{
		DopsEvaluationID: evaluation.ID,
		SupervisorID:     evaluation.SupervisorID,
		ResidentID:       evaluation.ResidentID,
		SurgeryMetadata: &commonModel.OrthopedicSurgeryMetadata{
			OccurredAt:    evaluation.OccurredAt,
			CaseNotes:     evaluation.CaseNotes,
			PatientAge:    evaluation.PatientAge,
			PatientGender: evaluation.PatientGender,
		},
		Difficulty:    &evaluation.Difficulty,
		DepartmentID:  evaluation.DepartmentID.Ptr(),
		IsEvaluated:   evaluation.IsEvaluated,
		InProgress:    evaluation.InProgress,
		ActiveStep:    evaluation.ActiveStep,
		CompletedStep: evaluation.CompletedStep,
		CreatedAt:     evaluation.CreatedAt,
	}

	annotations := make([]*commonModel.EvaluationFormAnnotations, 0)
	if err := evaluation.Annotations.Unmarshal(&annotations); err != nil {
		// Todo handle err
	} else {
		mappedEvaluation.Annotations = annotations
	}

	mappedEvaluation.Description = &commonModel.Description{Title: "DOPS", Subrows: []*commonModel.DescriptionRow{}}

	if evaluation.R != nil {
		if evaluation.R.Department != nil {
			mappedEvaluation.Department = &commonModel.ClinicDepartment{
				DepartmentID:   evaluation.DepartmentID.String,
				DepartmentName: evaluation.R.Department.DisplayName,
				ClinicID:       evaluation.R.Department.ParentID.String,
			}
		}

		if evaluation.R.DopsEvaluationsSurgeries != nil {
			mappedEvaluation.SurgeryMetadata.Surgeries = []*commonModel.Surgery{}
			for _, surgery := range evaluation.R.DopsEvaluationsSurgeries {
				mappedEvaluation.SurgeryMetadata.Surgeries = append(mappedEvaluation.SurgeryMetadata.Surgeries, MapDopsSurgery(surgery))
				subRow := &commonModel.DescriptionRow{}
				if surgery.R != nil && surgery.R.Surgery != nil && surgery.R.Surgery.R.Method != nil && surgery.R.Surgery.R.Diagnose != nil {
					subRow.Title = surgery.R.Surgery.R.Diagnose.DiagnoseName
					subRow.Subtitle = &surgery.R.Surgery.R.Method.MethodName
					subRow.SubtitleHighlight = &surgery.R.Surgery.R.Method.ApproachName
					mappedEvaluation.Description.Subrows = append(mappedEvaluation.Description.Subrows, subRow)
				}
			}
		}

		if evaluation.R.OrthopedicSurgeryActivity != nil && mapActivity {
			mappedEvaluation.SurgeryActivity = MapOrthopedicSurgeryActivity(evaluation.R.OrthopedicSurgeryActivity, false, true)
		}

		if evaluation.R.Supervisor != nil && evaluation.R.Resident != nil {
			mappedEvaluation.Resident = &commonModel.User{
				UserID:      evaluation.R.Resident.ID,
				DisplayName: evaluation.R.Resident.DisplayName,
			}
			mappedEvaluation.Supervisor = &commonModel.User{
				UserID:      evaluation.R.Supervisor.ID,
				DisplayName: evaluation.R.Supervisor.DisplayName,
			}
		}
	}

	return mappedEvaluation
}

func MapOrthopedicSurgeryActivity(activity *models.OrthopedicSurgeryActivity, mapDops bool, mapReview bool) *commonModel.OrthopedicSurgeryActivity {
	var mappedActivity = &commonModel.OrthopedicSurgeryActivity{
		ID:        activity.ID,
		CreatedAt: activity.CreatedAt,
		SurgeryMetadata: &commonModel.OrthopedicSurgeryMetadata{
			OccurredAt:    activity.OccurredAt,
			CaseNotes:     activity.CaseNotes,
			PatientAge:    activity.PatientAge,
			PatientGender: activity.PatientGender,
		},
		ResidentUserID:   activity.ResidentID,
		SupervisorUserID: activity.SupervisorID.Ptr(),
		OperatorID:       activity.OperatorID.Ptr(),
		AssistantID:      activity.AssistantID.Ptr(),
		Comments:         activity.Comments,
		Complications:    activity.Complications,
		InProgress:       activity.InProgress,
		ActiveStep:       activity.ActiveStep,
		CompletedStep:    activity.CompletedStep,
		DopsRequested:    activity.DopsRequested,
		ReviewRequested:  activity.ReviewRequested,
	}
	annotations := &commonModel.OrthopedicSurgeryActivityAnnotations{}
	if err := activity.Annotations.Unmarshal(&annotations); err != nil {
		// Todo handle err
	} else {
		mappedActivity.Annotations = annotations
	}

	mappedActivity.Description = &commonModel.Description{Title: "Ortopedisk operation", Subrows: []*commonModel.DescriptionRow{}}

	if activity.R != nil {
		if activity.R.OrthopedicSurgeryActivitiesSurgeries != nil {
			mappedActivity.SurgeryMetadata.Surgeries = []*commonModel.Surgery{}
			for _, surgery := range activity.R.OrthopedicSurgeryActivitiesSurgeries {
				mappedActivity.SurgeryMetadata.Surgeries = append(mappedActivity.SurgeryMetadata.Surgeries, MapOrthopedicSurgerySurgery(surgery))
				subRow := &commonModel.DescriptionRow{}
				if surgery.R != nil && surgery.R.Surgery != nil && surgery.R.Surgery.R.Method != nil && surgery.R.Surgery.R.Diagnose != nil {
					subRow.Title = surgery.R.Surgery.R.Diagnose.DiagnoseName
					subRow.Subtitle = &surgery.R.Surgery.R.Method.MethodName
					subRow.SubtitleHighlight = &surgery.R.Surgery.R.Method.ApproachName
					mappedActivity.Description.Subrows = append(mappedActivity.Description.Subrows, subRow)
				}
			}
		}
		if activity.R.DopsEvaluation != nil && mapDops {
			mappedActivity.Dops = MapDopsEvaluationGraphQlModel(activity.R.DopsEvaluation, false)
		}
		if activity.R.OrthopedicSurgeriesActivityReview != nil && mapReview {
			mappedActivity.Review = MapOrthopedicSurgeryReviewGraphQlModel(activity.R.OrthopedicSurgeriesActivityReview, false)
		}

		if activity.R.Supervisor != nil && activity.R.Resident != nil {
			mappedActivity.Resident = &commonModel.User{
				UserID:      activity.R.Resident.ID,
				DisplayName: activity.R.Resident.DisplayName,
			}
			mappedActivity.Supervisor = &commonModel.User{
				UserID:      activity.R.Supervisor.ID,
				DisplayName: activity.R.Supervisor.DisplayName,
			}
		}
	}

	return mappedActivity
}

func MapOrthopedicSurgeryReviewGraphQlModel(review *models.OrthopedicSurgeriesActivityReview, mapActivity bool) *commonModel.OrthopedicSurgeryActivityReview {

	mappedReview := &commonModel.OrthopedicSurgeryActivityReview{
		ReviewID:   review.ID,
		ActivityID: review.OrthopedicSurgeryActivityID,
		CreatedAt:  review.CreatedAt,
		UpdatedAt:  review.UpdatedAt.Time,
		SignedAt:   review.SignedAt.Ptr(),
		SurgeryMetadata: &commonModel.OrthopedicSurgeryMetadata{
			OccurredAt:    review.OccurredAt,
			CaseNotes:     review.CaseNotes,
			PatientAge:    review.PatientAge,
			PatientGender: review.PatientGender,
		},
		OperatorID:       review.OperatorID,
		AssistantID:      review.AssistantID,
		Comments:         review.Comments,
		Complications:    review.Complications,
		Annotations:      &commonModel.OrthopedicSurgeryActivityAnnotations{},
		ResidentUserID:   review.ResidentID,
		SupervisorUserID: review.SupervisorID,
		ReviewComment:    review.ReviewComment,
		InProgress:       review.InProgress,
		ActiveStep:       review.ActiveStep,
		CompletedStep:    review.CompletedStep,
	}

	err := review.Annotations.Unmarshal(mappedReview.Annotations)
	if err != nil {
		return nil
	}

	mappedReview.Description = &commonModel.Description{Title: "Ortopedisk operation", Subrows: []*commonModel.DescriptionRow{}}

	if review.R != nil {
		if review.R.OrthopedicSurgeryActivity != nil && mapActivity {
			mappedReview.Activity = MapOrthopedicSurgeryActivity(review.R.OrthopedicSurgeryActivity, true, false)
		}
		if review.R.OrthopedicSurgeriesActivityReviewSurgeries != nil {
			surgeries := make([]*commonModel.Surgery, 0)
			for _, surgery := range review.R.OrthopedicSurgeriesActivityReviewSurgeries {
				surgeries = append(surgeries, MapOrthopedicSurgeryReviewSurgery(surgery))

				subRow := &commonModel.DescriptionRow{}
				if surgery.R != nil && surgery.R.Surgery != nil && surgery.R.Surgery.R.Method != nil && surgery.R.Surgery.R.Diagnose != nil {
					subRow.Title = surgery.R.Surgery.R.Diagnose.DiagnoseName
					subRow.Subtitle = &surgery.R.Surgery.R.Method.MethodName
					subRow.SubtitleHighlight = &surgery.R.Surgery.R.Method.ApproachName
					mappedReview.Description.Subrows = append(mappedReview.Description.Subrows, subRow)
				}
			}
			mappedReview.SurgeryMetadata.Surgeries = surgeries
		}

		if review.R.Supervisor != nil && review.R.Resident != nil {
			mappedReview.Resident = &commonModel.User{
				UserID:      review.R.Resident.ID,
				DisplayName: review.R.Resident.DisplayName,
			}
			mappedReview.Supervisor = &commonModel.User{
				UserID:      review.R.Supervisor.ID,
				DisplayName: review.R.Supervisor.DisplayName,
			}
		}
	}

	return mappedReview
}

func MapDopsSurgery(surgery *models.DopsEvaluationsSurgery) *commonModel.Surgery {
	mappedSurgery := &commonModel.Surgery{
		ID: surgery.SurgeryID,
	}
	if surgery.R.Surgery != nil {
		if surgery.R.Surgery.R.Diagnose != nil {
			mappedSurgery.Diagnose = &commonModel.SurgeryDiagnose{
				ID:           surgery.R.Surgery.R.Diagnose.ID,
				Bodypart:     surgery.R.Surgery.R.Diagnose.Bodypart,
				DiagnoseName: surgery.R.Surgery.R.Diagnose.DiagnoseName,
				DiagnoseCode: surgery.R.Surgery.R.Diagnose.DiagnoseCode,
				ExtraCode:    surgery.R.Surgery.R.Diagnose.ExtraCode,
			}
		}
		if surgery.R.Surgery.R.Method != nil {
			mappedSurgery.Method = &commonModel.SurgeryMethod{
				ID:           surgery.R.Surgery.R.Method.ID,
				MethodName:   surgery.R.Surgery.R.Method.MethodName,
				MethodCode:   surgery.R.Surgery.R.Method.MethodCode,
				ApproachName: surgery.R.Surgery.R.Method.ApproachName,
			}
		}
	}
	return mappedSurgery
}

func MapOrthopedicSurgerySurgery(surgery *models.OrthopedicSurgeryActivitiesSurgery) *commonModel.Surgery {
	mappedSurgery := &commonModel.Surgery{
		ID: surgery.SurgeryID,
	}
	if surgery.R.Surgery != nil {
		if surgery.R.Surgery.R.Diagnose != nil {
			mappedSurgery.Diagnose = &commonModel.SurgeryDiagnose{
				ID:           surgery.R.Surgery.R.Diagnose.ID,
				Bodypart:     surgery.R.Surgery.R.Diagnose.Bodypart,
				DiagnoseName: surgery.R.Surgery.R.Diagnose.DiagnoseName,
				DiagnoseCode: surgery.R.Surgery.R.Diagnose.DiagnoseCode,
				ExtraCode:    surgery.R.Surgery.R.Diagnose.ExtraCode,
			}
		}
		if surgery.R.Surgery.R.Method != nil {
			mappedSurgery.Method = &commonModel.SurgeryMethod{
				ID:           surgery.R.Surgery.R.Method.ID,
				MethodName:   surgery.R.Surgery.R.Method.MethodName,
				MethodCode:   surgery.R.Surgery.R.Method.MethodCode,
				ApproachName: surgery.R.Surgery.R.Method.ApproachName,
			}
		}
	}
	return mappedSurgery
}

func MapOrthopedicSurgeryReviewSurgery(surgery *models.OrthopedicSurgeriesActivityReviewSurgery) *commonModel.Surgery {
	mappedSurgery := &commonModel.Surgery{
		ID: surgery.SurgeryID,
	}
	if surgery.R.Surgery != nil {
		if surgery.R.Surgery.R.Diagnose != nil {
			mappedSurgery.Diagnose = &commonModel.SurgeryDiagnose{
				ID:           surgery.R.Surgery.R.Diagnose.ID,
				Bodypart:     surgery.R.Surgery.R.Diagnose.Bodypart,
				DiagnoseName: surgery.R.Surgery.R.Diagnose.DiagnoseName,
				DiagnoseCode: surgery.R.Surgery.R.Diagnose.DiagnoseCode,
				ExtraCode:    surgery.R.Surgery.R.Diagnose.ExtraCode,
			}
		}
		if surgery.R.Surgery.R.Method != nil {
			mappedSurgery.Method = &commonModel.SurgeryMethod{
				ID:           surgery.R.Surgery.R.Method.ID,
				MethodName:   surgery.R.Surgery.R.Method.MethodName,
				MethodCode:   surgery.R.Surgery.R.Method.MethodCode,
				ApproachName: surgery.R.Surgery.R.Method.ApproachName,
			}
		}
	}
	return mappedSurgery
}
