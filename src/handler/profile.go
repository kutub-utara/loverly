package handler

import (
	"loverly/src/business/usecase"
	"net/http"
)

func GetProfile(uc *usecase.Usecases) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		profile, err := uc.Profile.Get(r.Context())
		if err != nil {
			JSONError(r.Context(), w, http.StatusBadRequest, err)
			return
		}

		JSONSuccess(r.Context(), w, http.StatusOK, profile)
	}
}
