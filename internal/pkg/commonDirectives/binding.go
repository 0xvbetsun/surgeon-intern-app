package commonDirectives

import (
	"context"
	"fmt"

	"github.com/vbetsun/surgeon-intern-app/internal/pkg/commonErrors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	validate *validator.Validate
	trans    ut.Translator
)

func init() {
	validate = validator.New()
	translator := en.New()
	uni := ut.New(translator, translator)
	trans, _ = uni.GetTranslator("translator")
	en_translations.RegisterDefaultTranslations(validate, trans)
}

func Binding(ctx context.Context, obj interface{}, next graphql.Resolver, constraint string) (interface{}, error) {
	val, err := next(ctx)
	if err != nil {
		panic(err)
	}
	fieldName := *graphql.GetPathContext(ctx).Field

	err = validate.Var(val, constraint)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		transErr := fmt.Sprintf("%s%+v", fieldName, validationErrors[0].Translate(trans))
		return val, commonErrors.NewInvalidInputError(transErr)
	}

	return val, nil
}

func ValidateAddTranslation(tag string, message string) {
	validate.RegisterTranslation(tag, trans, func(ut ut.Translator) error {
		return ut.Add(tag, message, true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(tag, fe.Field())

		return t
	})
}
