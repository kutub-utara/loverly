package handler

import (
	"loverly/src/business/usecase"
	"loverly/src/handler/verifier"
	"net/http"
)

func Discovery(uc *usecase.Usecases) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := uc.Dating.Discovery(r.Context())
		if err != nil {
			JSONError(r.Context(), w, http.StatusBadRequest, err)
			return
		}

		JSONSuccess(r.Context(), w, http.StatusOK, res)
	}
}

func Swipe(uc *usecase.Usecases) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := verifier.BuildAndValidateSwipeRequest(r, Log, Verify)
		if err != nil {
			JSONError(r.Context(), w, http.StatusUnprocessableEntity, err)
			return
		}

		res, err := uc.Dating.Swipe(r.Context(), payload)
		if err != nil {
			JSONError(r.Context(), w, http.StatusBadRequest, err)
			return
		}

		JSONSuccess(r.Context(), w, http.StatusOK, res)
	}
}
