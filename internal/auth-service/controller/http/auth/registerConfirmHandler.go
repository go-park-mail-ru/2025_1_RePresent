package auth

import (
	"encoding/json"
	"net/http"
	model "retarget/internal/auth-service/easyjsonModels"
	entity "retarget/pkg/entity"
	"retarget/pkg/utils/validator"

	"github.com/mailru/easyjson"
)

func (c *AuthController) RegisterConfirmHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		//nolint:errcheck
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Method Not Allowed"))
		return
	}

	var req model.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
		resp := entity.NewResponse(true, err.Error())
		//nolint:errcheck
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	validate_errors, err := validator.ValidateStruct(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, validate_errors))
		resp := entity.NewResponse(true, validate_errors)
		//nolint:errcheck
		easyjson.MarshalToWriter(&resp, w)
		return
	}
	// GetUserByEmail, если нет то отправляем код
}
