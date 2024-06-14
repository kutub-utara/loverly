package handler

import (
	"loverly/src/business/usecase"
	"net/http"
)

func Match(uc *usecase.Usecases) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		matchs, err := uc.Match.GetList(r.Context())
		if err != nil {
			JSONError(r.Context(), w, http.StatusBadRequest, err)
			return
		}

		JSONSuccess(r.Context(), w, http.StatusOK, matchs)
	}
}
