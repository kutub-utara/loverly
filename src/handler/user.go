package handler

import (
	"loverly/src/business/usecase"
	"loverly/src/handler/verifier"
	"net/http"
)

func SignIn(uc *usecase.Usecases) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// build and validate request body
		payload, err := verifier.BuildAndValidateLoginRequest(r, Log, Verify)
		if err != nil {
			JSONError(r.Context(), w, http.StatusUnprocessableEntity, err)
			return
		}

		// service to authenticate user
		res, err := uc.User.SignIn(r.Context(), payload)
		if err != nil {
			JSONError(r.Context(), w, http.StatusUnauthorized, err)
			return
		}

		JSONSuccess(r.Context(), w, http.StatusOK, res)
	}
}

func SignUp(uc *usecase.Usecases) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// build and validate request body
		payload, err := verifier.BuildAndValidateRegisterRequest(r, Log, Verify)
		if err != nil {
			JSONError(r.Context(), w, http.StatusUnprocessableEntity, err)
			return
		}

		res, err := uc.User.SignUp(r.Context(), payload)
		if err != nil {
			JSONError(r.Context(), w, http.StatusBadRequest, err)
			return
		}

		JSONSuccess(r.Context(), w, http.StatusCreated, res)
	}
}
