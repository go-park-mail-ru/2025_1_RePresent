package auth

import (
	"encoding/json"
	"net/http"
	entity "retarget/pkg/entity"
	"retarget/pkg/utils/validator"
)

type RegisterConfirmRequest struct {
	Email string `json:"email" validate:"email,required"`
}

func (c *AuthController) RegisterConfirmHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Method Not Allowed"))
		return
	}

	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
		return
	}

	validate_errors, err := validator.ValidateStruct(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewResponse(true, validate_errors))
		return
	}
	// GetUserByEmail, если нет то отправляем код
}
