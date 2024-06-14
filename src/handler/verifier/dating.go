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

func BuildAndValidateSwipeRequest(r *http.Request, log log.Interface, validate *validator.Validate) (entity.SwipeParam, error) {
	var swipe entity.SwipeParam

	bodyByte, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(r.Context(), fmt.Sprintf("read request body err: %v", err))
		return swipe, err
	}

	if err := json.Unmarshal(bodyByte, &swipe); err != nil {
		log.Error(r.Context(), fmt.Sprintf("unmarshal request body err: %v", err))
		return swipe, err
	}

	if err := validate.Struct(swipe); err != nil {
		log.Error(r.Context(), fmt.Sprintf("validate request body err: %v", err))
		return swipe, err
	}

	return swipe, nil
}
