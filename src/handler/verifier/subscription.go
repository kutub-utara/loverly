package verifier

import (
	"encoding/json"
	"fmt"
	"io"
	"loverly/lib/log"
	"loverly/src/business/entity"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func BuildAndValidateSubscriptionRequest(r *http.Request, log log.Interface, validate *validator.Validate) (entity.SubscriptionParam, error) {
	var subs entity.SubscriptionParam

	bodyByte, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(r.Context(), fmt.Sprintf("read request body err: %v", err))
		return subs, err
	}

	if err := json.Unmarshal(bodyByte, &subs); err != nil {
		log.Error(r.Context(), fmt.Sprintf("unmarshal request body err: %v", err))
		return subs, err
	}

	if err := validate.Struct(subs); err != nil {
		log.Error(r.Context(), fmt.Sprintf("validate request body err: %v", err))
		return subs, err
	}

	return subs, nil
}
