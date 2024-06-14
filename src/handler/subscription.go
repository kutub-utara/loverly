package handler

import (
	"loverly/src/business/usecase"
	"loverly/src/handler/verifier"
	"net/http"
)

func Subscribe(uc *usecase.Usecases) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// build and validate request body
		payload, err := verifier.BuildAndValidateSubscriptionRequest(r, Log, Verify)
		if err != nil {
			JSONError(r.Context(), w, http.StatusUnprocessableEntity, err)
			return
		}

		err = uc.Subscription.Create(r.Context(), payload)
		if err != nil {
			JSONError(r.Context(), w, http.StatusBadRequest, err)
			return
		}

		JSONSuccess(r.Context(), w, http.StatusCreated, nil)
	}
}

func GetSubscribe(uc *usecase.Usecases) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		subs, err := uc.Subscription.Get(r.Context())
		if err != nil {
			JSONError(r.Context(), w, http.StatusBadRequest, err)
			return
		}

		JSONSuccess(r.Context(), w, http.StatusOK, subs)
	}
}
