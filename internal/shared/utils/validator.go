package validator

import (
	"strings"

	apperrors "music/internal/common/errors"

	govalidator "github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *govalidator.Validate
}

func New() *Validator {
	v := govalidator.New()

	return &Validator{
		validate: v,
	}
}

func (v *Validator) Struct(data interface{}) error {
	if err := v.validate.Struct(data); err != nil {
		return apperrors.Validation("validation failed", formatValidationErrors(err))
	}

	return nil
}

func formatValidationErrors(err error) map[string]string {
	result := map[string]string{}

	validationErrors, ok := err.(govalidator.ValidationErrors)
	if !ok {
		result["_error"] = err.Error()
		return result
	}

	for _, fieldError := range validationErrors {
		fieldName := toJSONFieldName(fieldError.Field())

		switch fieldError.Tag() {
		case "required":
			result[fieldName] = "is required"
		case "email":
			result[fieldName] = "must be a valid email"
		case "min":
			result[fieldName] = "must be at least " + fieldError.Param() + " characters"
		case "max":
			result[fieldName] = "must be at most " + fieldError.Param() + " characters"
		case "oneof":
			result[fieldName] = "must be one of: " + fieldError.Param()
		default:
			result[fieldName] = "is invalid"
		}
	}

	return result
}

func toJSONFieldName(field string) string {
	if field == "" {
		return field
	}

	return strings.ToLower(field[:1]) + field[1:]
}
func (v *Validator) Validate(data interface{}) map[string]string {
	if err := v.validate.Struct(data); err != nil {
		return formatValidationErrors(err)
	}

	return nil
}
