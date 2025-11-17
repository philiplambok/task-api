package validator

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type ValidationError struct {
	Field   string
	Message string
}

type ValidationErrors struct {
	Errors []ValidationError
}

func (v ValidationErrors) Error() string {
	if len(v.Errors) == 0 {
		return ""
	}

	return v.Errors[0].Message
}

// ValidateStruct validates a struct and returns ValidationErrors if validation fails.
// The parameter s must be a struct or a pointer to a struct.
func ValidateStruct(s any) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	var validationErrs validator.ValidationErrors
	if !errors.As(err, &validationErrs) {
		return err
	}

	return buildValidationErrors(validationErrs)
}

func buildValidationErrors(validationErrs validator.ValidationErrors) ValidationErrors {
	errs := make([]ValidationError, 0, len(validationErrs))
	for _, e := range validationErrs {
		errs = append(errs, ValidationError{
			Field:   e.Field(),
			Message: getValidationMessage(e),
		})
	}
	return ValidationErrors{Errors: errs}
}

func getValidationMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", e.Field())
	default:
		return fmt.Sprintf("%s is invalid", e.Field())
	}
}
