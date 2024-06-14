package verifier

import (
	"encoding/json"
	"fmt"
	"io"
	"loverly/lib/log"
	"loverly/src/business/entity"
	"net/http"

	appErr "loverly/src/errors"

	"github.com/go-playground/validator/v10"
)

func BuildAndValidateLoginRequest(r *http.Request, log log.Interface, validate *validator.Validate) (entity.SignInParam, error) {
	var signIn entity.SignInParam

	bodyByte, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(r.Context(), fmt.Sprintf("read request body err: %v", err))
		return signIn, err
	}

	if err := json.Unmarshal(bodyByte, &signIn); err != nil {
		log.Error(r.Context(), fmt.Sprintf("unmarshal request body err: %v", err))
		return signIn, err
	}

	if err := validate.Struct(signIn); err != nil {
		log.Error(r.Context(), fmt.Sprintf("validate request body err: %v", err))

		// Validation failed, handle the errors here
		if errors, ok := err.(validator.ValidationErrors); ok {
			// Check for a specific error on the "email" field with type "required"
			if hasSpecificFieldError(errors, "Password", "min") {
				return signIn, appErr.ErrInvalidEmailOrPassword
			}
		}
		// Handle other types of errors (e.g., type conversion failures)

		return signIn, err
	}

	return signIn, nil
}

func BuildAndValidateRegisterRequest(r *http.Request, log log.Interface, validate *validator.Validate) (entity.SignUpParam, error) {
	var signUp entity.SignUpParam

	bodyByte, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(r.Context(), fmt.Sprintf("read request body err: %v", err))
		return signUp, err
	}

	if err := json.Unmarshal(bodyByte, &signUp); err != nil {
		log.Error(r.Context(), fmt.Sprintf("unmarshal request body err: %v", err))
		return signUp, err
	}

	if err := validate.Struct(signUp); err != nil {
		log.Error(r.Context(), fmt.Sprintf("validate request body err: %v", err))

		// Validation failed, handle the errors here
		if errors, ok := err.(validator.ValidationErrors); ok {
			// Check for a specific error on the "email" field with type "required"
			if hasSpecificFieldError(errors, "Email", "email") {
				return signUp, appErr.ErrInvalidEmailFormat
			} else if hasSpecificFieldError(errors, "Password", "min") {
				return signUp, appErr.ErrInvalidEmailOrPassword
			} else if hasSpecificFieldError(errors, "ConfirmPassword", "eqfield") {
				return signUp, appErr.ErrPasswordNotMatch
			}
		}
		// Handle other types of errors (e.g., type conversion failures)

		return signUp, err
	}

	return signUp, nil
}
