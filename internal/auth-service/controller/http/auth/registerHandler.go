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
	// Code     string `json:"code" validate:"required,len=6"`
	Password string `json:"password" validate:"required,min=8"`
	Role     int    `json:"role" validate:"required,gte=1,lte=2"`
}

func (c *AuthController) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(entity.СtxKeyRequestID{}).(string)
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

	user, err := c.authUsecase.Register(req.Username, req.Email, req.Password, req.Role, requestID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
		return
	}

	session, err := c.authUsecase.AddSession(user.ID, user.Role)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
		return
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Expires:  session.Expires,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(entity.NewResponse(false, "registration succesful"))
}
