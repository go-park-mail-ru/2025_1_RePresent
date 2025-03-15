package auth

import (
	"encoding/json"
	"net/http"
	entity "retarget/internal/auth-service/entity"
	entityAuth "retarget/internal/auth-service/entity/auth"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required,min=8"`
	Role     int    `json:"role" validate:"required,gte=1,lte=2"`
}

func (c *AuthController) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Method Not Allowed"))
		// http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
		// http.Error(w, "Invalid JSON", http.StatusUnprocessableEntity)
		return
	}

	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
		// http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := c.authUsecase.Register(req.Username, req.Email, req.Password, req.Role)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
		// http.Error(w, err.Error(), http.StatusBadRequest)
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

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("{response: Registration successful}")) // пока костылём
}
