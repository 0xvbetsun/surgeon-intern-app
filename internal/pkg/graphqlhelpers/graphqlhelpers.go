package graphqlhelpers

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
)

func GetFields(ctx context.Context) map[string]graphql.CollectedField {
	fields := graphql.CollectFieldsCtx(ctx, nil)
	fieldsMap := make(map[string]graphql.CollectedField)
	for _, field := range fields {
		fieldsMap[field.Name] = field
	}
	return fieldsMap
}
func GetNestedFields(ctx context.Context, field graphql.CollectedField) map[string]graphql.CollectedField {
	graphqlCtx := graphql.GetOperationContext(ctx)
	nestedFields := graphql.CollectFields(graphqlCtx, field.Selections, nil)
	fieldsMap := make(map[string]graphql.CollectedField)
	for _, nestedField := range nestedFields {
		fieldsMap[nestedField.Name] = nestedField
	}
	return fieldsMap
}
