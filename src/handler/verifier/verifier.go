package verifier

import (
	"github.com/go-playground/validator/v10"
)

func hasSpecificFieldError(errs validator.ValidationErrors, fieldName string, errType string) bool {
	for _, fieldError := range errs {
		if fieldError.Field() == fieldName && fieldError.Tag() == errType {
			return true
		}
	}
	return false
}
