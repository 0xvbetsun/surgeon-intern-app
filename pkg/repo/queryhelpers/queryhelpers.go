package queryhelpers

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/internal/pkg/graphqlhelpers"
	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/volatiletech/null/v8"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func GetLogbookEntriesFilters(ctx context.Context, userId string, fields map[string]graphql.CollectedField, queryFilter commonModel.LogbookEntryQueryFilter) (relPaths []string, joinPaths []string, where []QueryMod) {
	var rels []string
	var joins []string
	var wheres []QueryMod

	for _, logbookType := range queryFilter.LogbookEntryTypes {
		switch logbookType {
		case commonModel.LogbookEntryTypeSurgery:
			if dopsField, exists := fields["surgery"]; exists {
				surgeryLogbookQueryFilter := commonModel.SurgeryLogbookEntryQueryFilter{}
				joins = append(joins, models.TableNames.OrthopedicSurgeryActivities+
					" on "+
					models.LogbookEntryTableColumns.OrthopedicSurgeryID+
					" = "+
					models.OrthopedicSurgeryActivityTableColumns.ID)
				if queryFilter.SurgeryLogbookEntryFilters != nil {
					surgeryLogbookQueryFilter = *queryFilter.SurgeryLogbookEntryFilters
				}
				rels = append(rels, models.LogbookEntryRels.OrthopedicSurgery)
				orthopedicSurgeryFields := graphqlhelpers.GetNestedFields(ctx, dopsField)
				orthopedicSurgeriesRels, orthopedicSurgeryJoinPaths, orthopedicSurgeriesFilters := GetOrthopedicSurgeryFilters(ctx, userId, orthopedicSurgeryFields, surgeryLogbookQueryFilter)
				joins = append(joins, orthopedicSurgeryJoinPaths...)
				if len(orthopedicSurgeriesFilters) > 0 {
					wheres = append(wheres, Or2(Expr(orthopedicSurgeriesFilters...)))
				} else {
					wheres = append(wheres, Or2(models.LogbookEntryWhere.OrthopedicSurgeryID.IsNotNull()))
				}
				for _, orthopedicSurgeriesRel := range orthopedicSurgeriesRels {
					rels = append(rels, Rels(models.LogbookEntryRels.OrthopedicSurgery, orthopedicSurgeriesRel))
				}
			}
			break
		}
	}

	return rels, joins, wheres
}

func GetAssessmentFilters(ctx context.Context, userId string, fields map[string]graphql.CollectedField, queryFilter commonModel.AssessmentQueryFilter) (relPaths []string, joinPaths []string, where []QueryMod) {
	var rels []string
	var joins []string
	var wheres []QueryMod

	for _, assessmentType := range queryFilter.AssessmentTypes {
		switch assessmentType {
		case commonModel.AssessmentTypeDops:
			if dopsField, exists := fields["dops"]; exists {
				joins = append(joins, models.TableNames.DopsEvaluations+
					" on "+
					models.AssessmentTableColumns.DopsID+
					" = "+
					models.DopsEvaluationTableColumns.ID)
				dopsQueryFilter := commonModel.DopsQueryFilter{}
				if queryFilter.DopsFilter != nil {
					dopsQueryFilter = *queryFilter.DopsFilter
				}
				rels = append(rels, models.AssessmentRels.Dop)
				dopsFields := graphqlhelpers.GetNestedFields(ctx, dopsField)
				dopsRels, dopsJoins, dopsFilters := GetDopsFilters(ctx, userId, dopsFields, dopsQueryFilter)
				joins = append(joins, dopsJoins...)
				if len(dopsFilters) > 0 {
					wheres = append(wheres, Or2(Expr(dopsFilters...)))
				} else {
					wheres = append(wheres, Or2(models.AssessmentWhere.DopsID.IsNotNull()))
				}
				for _, dopsRel := range dopsRels {
					rels = append(rels, Rels(models.AssessmentRels.Dop, dopsRel))
				}
			}
			break
		case commonModel.AssessmentTypeMiniCex:
			if miniCexField, exists := fields["miniCex"]; exists {
				miniCexQueryFilter := commonModel.MiniCexQueryFilter{}
				joins = append(joins, models.TableNames.MiniCexEvaluations+
					" on "+
					models.AssessmentTableColumns.MiniCexID+
					" = "+
					models.MiniCexEvaluationTableColumns.ID)
				if queryFilter.MiniCexFilter != nil {
					miniCexQueryFilter = *queryFilter.MiniCexFilter
				}
				rels = append(rels, models.AssessmentRels.MiniCex)
				miniCexFields := graphqlhelpers.GetNestedFields(ctx, miniCexField)
				miniCexRels, _, miniCexFilters := GetMiniCexRels(ctx, userId, miniCexFields, miniCexQueryFilter)
				if len(miniCexFilters) > 0 {
					wheres = append(wheres, Or2(Expr(miniCexFilters...)))
				} else {
					wheres = append(wheres, Or2(models.AssessmentWhere.MiniCexID.IsNotNull()))
				}
				for _, miniCexRel := range miniCexRels {
					rels = append(rels, Rels(models.AssessmentRels.MiniCex, miniCexRel))
				}
			}
		case commonModel.AssessmentTypeSurgeryReview:
			if surgeryReviewField, exists := fields["surgeryReview"]; exists {
				surgeryReviewQueryFilter := commonModel.SurgeryReviewQueryFilter{}
				joins = append(joins, models.TableNames.OrthopedicSurgeriesActivityReview+
					" on "+
					models.AssessmentTableColumns.OrthopedicSurgeryReviewID+
					" = "+
					models.OrthopedicSurgeriesActivityReviewTableColumns.ID)
				if queryFilter.SurgeryReviewFilter != nil {
					surgeryReviewQueryFilter = *queryFilter.SurgeryReviewFilter
				}
				rels = append(rels, models.AssessmentRels.OrthopedicSurgeryReview)
				surgeryReviewFields := graphqlhelpers.GetNestedFields(ctx, surgeryReviewField)
				surgeryReviewRels, surgeryReviewJoins, surgeryReviewFilters := GetOrthopedicSurgeryReviewFilters(ctx, userId, surgeryReviewFields, surgeryReviewQueryFilter)
				joins = append(joins, surgeryReviewJoins...)
				if len(surgeryReviewFilters) > 0 {
					wheres = append(wheres, Or2(Expr(surgeryReviewFilters...)))
				} else {
					wheres = append(wheres, Or2(models.AssessmentWhere.OrthopedicSurgeryReviewID.IsNotNull()))
				}
				for _, surgeryReviewRel := range surgeryReviewRels {
					rels = append(rels, Rels(models.AssessmentRels.OrthopedicSurgeryReview, surgeryReviewRel))
				}
			}
		}
	}

	return rels, joins, wheres
}

func GetDopsFilters(ctx context.Context, userId string, fields map[string]graphql.CollectedField, queryFilter commonModel.DopsQueryFilter) (relPaths []string, joinPaths []string, where []QueryMod) {
	var rels []string
	var joins []string
	var wheres []QueryMod

	if queryFilter.ResidentID == nil && queryFilter.SupervisorID == nil {
		queryFilter.ResidentID = &userId
	}

	if queryFilter.ResidentID != nil {
		wheres = append(wheres, models.DopsEvaluationWhere.ResidentID.EQ(*queryFilter.ResidentID))
	}
	if queryFilter.SupervisorID != nil {
		wheres = append(wheres, models.DopsEvaluationWhere.SupervisorID.EQ(*queryFilter.SupervisorID))
	}
	if queryFilter.IsEvaluated != nil {
		wheres = append(wheres, models.DopsEvaluationWhere.IsEvaluated.EQ(*queryFilter.IsEvaluated))
	}
	if queryFilter.InProgress != nil {
		wheres = append(wheres, models.DopsEvaluationWhere.InProgress.EQ(*queryFilter.InProgress))
	}
	if queryFilter.IsConnected != nil {
		if *queryFilter.IsConnected {
			wheres = append(wheres, models.DopsEvaluationWhere.OrthopedicSurgeryActivityID.IsNotNull())
		} else {
			wheres = append(wheres, models.DopsEvaluationWhere.OrthopedicSurgeryActivityID.IsNull())
		}
	}
	if queryFilter.Surgeries != nil && len(queryFilter.Surgeries) > 0 {
		joins = append(joins,
			models.TableNames.DopsEvaluationsSurgeries+
				" on "+
				models.DopsEvaluationsSurgeryTableColumns.DopsEvaluationID+
				" = "+
				models.DopsEvaluationTableColumns.ID,
			models.TableNames.Surgeries+
				" as dops_surgeries "+
				" on dops_surgeries."+models.SurgeryColumns.ID+
				" = "+
				models.DopsEvaluationsSurgeryTableColumns.SurgeryID)
		surgeriesIn := make([]interface{}, 0, len(queryFilter.Surgeries))
		for _, surgeryId := range queryFilter.Surgeries {
			surgeriesIn = append(surgeriesIn, surgeryId)
		}
		wheres = append(wheres, WhereIn("dops_surgeries.id IN ?", surgeriesIn...))
	}

	if surgeryActivityField, exists := fields["surgeryActivity"]; exists {
		rels = append(rels, models.DopsEvaluationRels.OrthopedicSurgeryActivity)
		surgeryActivityFields := graphqlhelpers.GetNestedFields(ctx, surgeryActivityField)
		orthopedicSurgeryRels, _, _ := GetOrthopedicSurgeryFilters(ctx, userId, surgeryActivityFields, commonModel.SurgeryLogbookEntryQueryFilter{})
		for _, orthopedicSurgeriesRel := range orthopedicSurgeryRels {
			rels = append(rels, Rels(models.DopsEvaluationRels.OrthopedicSurgeryActivity, orthopedicSurgeriesRel))
		}
	}
	if surgeryMetadataField, exists := fields["surgeryMetadata"]; exists {
		surgeryMetadataFields := graphqlhelpers.GetNestedFields(ctx, surgeryMetadataField)
		if surgeriesField, exists := surgeryMetadataFields["surgeries"]; exists {
			rels = append(rels, Rels(models.DopsEvaluationRels.DopsEvaluationsSurgeries, models.DopsEvaluationsSurgeryRels.Surgery))
			surgeriesFields := graphqlhelpers.GetNestedFields(ctx, surgeriesField)
			surgeriesRels := GetSurgeryRels(ctx, surgeriesFields)
			for _, surgeriesRel := range surgeriesRels {
				rels = append(rels, Rels(models.DopsEvaluationRels.DopsEvaluationsSurgeries, Rels(models.DopsEvaluationsSurgeryRels.Surgery, surgeriesRel)))
			}
		}
	}
	if _, exists := fields["department"]; exists {
		rels = append(rels, models.DopsEvaluationRels.Department)
	}
	if _, exists := fields["resident"]; exists {
		rels = append(rels, models.DopsEvaluationRels.Resident)
	}
	if _, exists := fields["supervisor"]; exists {
		rels = append(rels, models.DopsEvaluationRels.Supervisor)
	}

	return rels, joins, wheres
}

func GetMiniCexRels(ctx context.Context, userId string, fields map[string]graphql.CollectedField, queryFilter commonModel.MiniCexQueryFilter) (relPaths []string, join []string, where []QueryMod) {
	var rels []string
	var joins []string
	var wheres []QueryMod

	if queryFilter.ResidentID == nil && queryFilter.SupervisorID == nil {
		queryFilter.ResidentID = &userId
	}

	if queryFilter.ResidentID != nil {
		wheres = append(wheres, models.MiniCexEvaluationWhere.ResidentID.EQ(*queryFilter.ResidentID))
	}
	if queryFilter.SupervisorID != nil {
		wheres = append(wheres, models.MiniCexEvaluationWhere.SupervisorID.EQ(*queryFilter.SupervisorID))
	}
	if queryFilter.InProgress != nil {
		wheres = append(wheres, models.MiniCexEvaluationWhere.InProgress.EQ(*queryFilter.InProgress))
	}
	if queryFilter.IsEvaluated != nil {
		wheres = append(wheres, models.MiniCexEvaluationWhere.IsEvaluated.EQ(*queryFilter.IsEvaluated))
	}

	if _, exists := fields["department"]; exists {
		rels = append(rels, models.MiniCexEvaluationRels.Department)
	}
	if _, exists := fields["resident"]; exists {
		rels = append(rels, models.MiniCexEvaluationRels.Resident)
	}
	if _, exists := fields["supervisor"]; exists {
		rels = append(rels, models.MiniCexEvaluationRels.Supervisor)
	}
	return rels, joins, wheres
}

func GetOrthopedicSurgeryFilters(ctx context.Context, userId string, fields map[string]graphql.CollectedField, queryFilter commonModel.SurgeryLogbookEntryQueryFilter) (relPaths []string, joinPaths []string, where []QueryMod) {
	var rels []string
	var joins []string
	var wheres []QueryMod

	if queryFilter.ResidentID == nil && queryFilter.SupervisorID == nil {
		queryFilter.ResidentID = &userId
	}

	if queryFilter.ResidentID != nil {
		wheres = append(wheres, models.OrthopedicSurgeryActivityWhere.ResidentID.EQ(*queryFilter.ResidentID))
	}
	if queryFilter.SupervisorID != nil {
		wheres = append(wheres, models.OrthopedicSurgeryActivityWhere.SupervisorID.EQ(null.StringFromPtr(queryFilter.SupervisorID)))
	}
	if queryFilter.InProgress != nil {
		wheres = append(wheres, models.OrthopedicSurgeryActivityWhere.InProgress.EQ(*queryFilter.InProgress))
	}
	if queryFilter.HasDops != nil {
		joins = append(joins,
			models.TableNames.DopsEvaluations+
				" as surgery_logbook_entry_dops "+
				" on surgery_logbook_entry_dops."+models.DopsEvaluationColumns.OrthopedicSurgeryActivityID+
				" = "+
				models.OrthopedicSurgeryActivityTableColumns.ID)
		if *queryFilter.HasDops {
			wheres = append(wheres, Where("surgery_logbook_entry_dops."+models.DopsEvaluationColumns.ID+" IS NOT NULL"))
		} else {
			wheres = append(wheres, Where("surgery_logbook_entry_dops."+models.DopsEvaluationColumns.ID+" IS NULL"))
		}
	}
	if queryFilter.HasReview != nil {
		joins = append(joins,
			models.TableNames.OrthopedicSurgeriesActivityReview+
				" as surgery_logbook_entry_reviews "+
				" on surgery_logbook_entry_reviews."+models.OrthopedicSurgeriesActivityReviewColumns.OrthopedicSurgeryActivityID+
				" = "+
				models.OrthopedicSurgeryActivityTableColumns.ID)
		if *queryFilter.HasReview {
			wheres = append(wheres, Where("surgery_logbook_entry_reviews."+models.OrthopedicSurgeriesActivityReviewColumns.OrthopedicSurgeryActivityID+" IS NOT NULL"))
		} else {
			wheres = append(wheres, Where("surgery_logbook_entry_reviews."+models.OrthopedicSurgeriesActivityReviewColumns.OrthopedicSurgeryActivityID+" IS NULL"))
		}
	}
	if queryFilter.Surgeries != nil && len(queryFilter.Surgeries) > 0 {
		joins = append(joins,
			models.TableNames.OrthopedicSurgeryActivitiesSurgeries+
				" on "+
				models.OrthopedicSurgeryActivitiesSurgeryTableColumns.OrthopedicSurgeryActivityID+
				" = "+
				models.OrthopedicSurgeryActivityTableColumns.ID,
			models.TableNames.Surgeries+
				" as surgery_logbook_entry_surgeries "+
				" on surgery_logbook_entry_surgeries."+models.SurgeryColumns.ID+
				" = "+
				models.OrthopedicSurgeryActivitiesSurgeryTableColumns.SurgeryID)
		surgeriesIn := make([]interface{}, 0, len(queryFilter.Surgeries))
		for _, surgeryId := range queryFilter.Surgeries {
			surgeriesIn = append(surgeriesIn, surgeryId)
		}
		wheres = append(wheres, WhereIn("surgery_logbook_entry_surgeries."+models.SurgeryColumns.ID+" IN ?", surgeriesIn...))
	}

	if surgeryMetadataField, exists := fields["surgeryMetadata"]; exists {
		surgeryMetadataFields := graphqlhelpers.GetNestedFields(ctx, surgeryMetadataField)
		if surgeriesField, exists := surgeryMetadataFields["surgeries"]; exists {
			rels = append(rels, Rels(models.OrthopedicSurgeryActivityRels.OrthopedicSurgeryActivitiesSurgeries, models.OrthopedicSurgeryActivitiesSurgeryRels.Surgery))
			surgeriesFields := graphqlhelpers.GetNestedFields(ctx, surgeriesField)
			surgeriesRels := GetSurgeryRels(ctx, surgeriesFields)
			for _, surgeriesRel := range surgeriesRels {
				rels = append(rels, Rels(models.OrthopedicSurgeryActivityRels.OrthopedicSurgeryActivitiesSurgeries, Rels(models.OrthopedicSurgeryActivitiesSurgeryRels.Surgery, surgeriesRel)))
			}
		}
	}
	if _, exists := fields["dops"]; exists {
		rels = append(rels, models.OrthopedicSurgeryActivityRels.DopsEvaluation)
	}
	if _, exists := fields["review"]; exists {
		rels = append(rels, models.OrthopedicSurgeryActivityRels.OrthopedicSurgeriesActivityReview)
	}
	if _, exists := fields["resident"]; exists {
		rels = append(rels, models.OrthopedicSurgeryActivityRels.Resident)
	}
	if _, exists := fields["supervisor"]; exists {
		rels = append(rels, models.OrthopedicSurgeryActivityRels.Supervisor)
	}

	return rels, joins, wheres
}

func GetOrthopedicSurgeryReviewFilters(ctx context.Context, userId string, fields map[string]graphql.CollectedField, queryFilter commonModel.SurgeryReviewQueryFilter) (relPaths []string, joinPaths []string, where []QueryMod) {
	var rels []string
	var joins []string
	var wheres []QueryMod

	if queryFilter.ResidentID == nil && queryFilter.SupervisorID == nil {
		queryFilter.ResidentID = &userId
	}

	if queryFilter.ResidentID != nil {
		wheres = append(wheres, models.OrthopedicSurgeriesActivityReviewWhere.ResidentID.EQ(*queryFilter.ResidentID))
	}
	if queryFilter.SupervisorID != nil {
		wheres = append(wheres, models.OrthopedicSurgeriesActivityReviewWhere.SupervisorID.EQ(*queryFilter.SupervisorID))
	}
	if queryFilter.InProgress != nil {
		wheres = append(wheres, models.OrthopedicSurgeriesActivityReviewWhere.InProgress.EQ(*queryFilter.InProgress))
	}
	if queryFilter.IsEvaluated != nil {
		if *queryFilter.IsEvaluated {
			wheres = append(wheres, models.OrthopedicSurgeriesActivityReviewWhere.SignedAt.IsNotNull())
		} else {
			wheres = append(wheres, models.OrthopedicSurgeriesActivityReviewWhere.SignedAt.IsNull())
		}
	}

	if queryFilter.Surgeries != nil && len(queryFilter.Surgeries) > 0 {
		joins = append(joins,
			models.TableNames.OrthopedicSurgeriesActivityReviewSurgeries+
				" on "+
				models.OrthopedicSurgeriesActivityReviewSurgeryTableColumns.OrthopedicSurgeriesActivityReviewID+
				" = "+
				models.OrthopedicSurgeriesActivityReviewTableColumns.ID,
			models.TableNames.Surgeries+
				" as surgery_review_surgeries "+
				" on surgery_review_surgeries."+models.SurgeryColumns.ID+
				" = "+
				models.OrthopedicSurgeriesActivityReviewSurgeryTableColumns.SurgeryID)
		surgeriesIn := make([]interface{}, 0, len(queryFilter.Surgeries))
		for _, surgeryId := range queryFilter.Surgeries {
			surgeriesIn = append(surgeriesIn, surgeryId)
		}
		wheres = append(wheres, WhereIn("surgery_review_surgeries."+models.SurgeryColumns.ID+" IN ?", surgeriesIn...))
	}

	if surgeryMetadataField, exists := fields["surgeryMetadata"]; exists {
		surgeryMetadataFields := graphqlhelpers.GetNestedFields(ctx, surgeryMetadataField)
		if surgeriesField, exists := surgeryMetadataFields["surgeries"]; exists {
			rels = append(rels, Rels(models.OrthopedicSurgeriesActivityReviewRels.OrthopedicSurgeriesActivityReviewSurgeries, models.OrthopedicSurgeriesActivityReviewSurgeryRels.Surgery))
			surgeriesFields := graphqlhelpers.GetNestedFields(ctx, surgeriesField)
			surgeriesRels := GetSurgeryRels(ctx, surgeriesFields)
			for _, surgeriesRel := range surgeriesRels {
				rels = append(rels, Rels(models.OrthopedicSurgeriesActivityReviewRels.OrthopedicSurgeriesActivityReviewSurgeries, Rels(models.OrthopedicSurgeriesActivityReviewSurgeryRels.Surgery, surgeriesRel)))
			}
		}
	}
	if _, exists := fields["activity"]; exists {
		rels = append(rels, models.OrthopedicSurgeriesActivityReviewRels.OrthopedicSurgeryActivity)
	}
	if _, exists := fields["resident"]; exists {
		rels = append(rels, models.OrthopedicSurgeriesActivityReviewRels.Resident)
	}
	if _, exists := fields["supervisor"]; exists {
		rels = append(rels, models.OrthopedicSurgeriesActivityReviewRels.Supervisor)
	}

	return rels, joins, wheres
}

func GetSurgeryRels(ctx context.Context, fields map[string]graphql.CollectedField) []string {
	var joinRels []string

	if _, exists := fields["diagnose"]; exists {
		joinRels = append(joinRels, models.SurgeryRels.Diagnose)
	}
	if _, exists := fields["method"]; exists {
		joinRels = append(joinRels, models.SurgeryRels.Method)
	}
	return joinRels
}
