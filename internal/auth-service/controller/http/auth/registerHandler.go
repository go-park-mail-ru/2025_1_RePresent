package auth

import (
	"encoding/json"
	"net/http"
	model "retarget/internal/auth-service/easyjsonModels"
	entity "retarget/pkg/entity"
	"retarget/pkg/utils/validator"
)

func (c *AuthController) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(entity.СtxKeyRequestID{}).(string)
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
		json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
		return
	}

	validate_errors, err := validator.ValidateStruct(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		//nolint:errcheck
		json.NewEncoder(w).Encode(entity.NewResponse(true, validate_errors))
		return
	}

	// Login(данные пользователя), проверили данные, AddSession(user_id), поставили куки

	user, err := c.authUsecase.Register(r.Context(), req.Username, req.Email, req.Password, req.Role, requestID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		//nolint:errcheck
		json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
		return
	}

	session, err := c.authUsecase.AddSession(user.ID, user.Role)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
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
	//nolint:errcheck
	json.NewEncoder(w).Encode(entity.NewResponse(false, "registration succesful"))
}
