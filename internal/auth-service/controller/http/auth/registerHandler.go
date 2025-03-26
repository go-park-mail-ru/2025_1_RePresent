package auth

import (
	"encoding/json"
	"net/http"
	entity "retarget/pkg/entity"
	"retarget/pkg/utils/validator"
)

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Email    string `json:"email" validate:"email,required"`
	Code     string `json:"code" validate:"required,len=6"`
	Password string `json:"password" validate:"required,min=8"`
	Role     int    `json:"role" validate:"required,gte=1,lte=2"`
}

func (c *AuthController) RegisterHandler(w http.ResponseWriter, r *http.Request) {
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
	// Login(данные пользователя), проверили данные, AddSession(user_id), поставили куки
}
