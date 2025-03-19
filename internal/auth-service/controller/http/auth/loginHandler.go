package auth

import (
	"encoding/json"
	"net/http"
	"time"

	entityAuth "retarget/internal/auth-service/entity/auth"
	entity "retarget/pkg/entity"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required,min=8"`
	Role     int    `json:"role" validate:"required,gte=1,lte=2"`
}

func (c *AuthController) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Method Not Allowed"))
		return
	}

	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
		return
	}

	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
		return
	}

	user, err := c.authUsecase.Login(req.Email, req.Password, req.Role)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
		return
	}

	sessionID := uuid.NewString()
	session := entityAuth.Session{
		ID:        sessionID,
		UserID:    user.ID,
		Expires:   time.Now().Add(30 * time.Minute),
		CreatedAt: time.Now(),
	}
	err = entityAuth.AddSession(session)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
		return
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  session.Expires,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}
