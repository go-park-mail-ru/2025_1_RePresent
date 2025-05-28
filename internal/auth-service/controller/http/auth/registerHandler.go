package auth

import (
	"encoding/json"
	"net/http"
	model "retarget/internal/auth-service/easyjsonModels"
	entity "retarget/pkg/entity"
	"retarget/pkg/utils/validator"

	"github.com/mailru/easyjson"
)

func (c *AuthController) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var requestID string
	if v := r.Context().Value(entity.СtxKeyRequestID{}); v != nil {
		if id, ok := v.(string); ok {
			requestID = id
		}
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, "Method Not Allowed"))
		resp := entity.NewResponse(true, "Method Not Allowed")
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	var req model.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
		resp := entity.NewResponse(true, err.Error())
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	validate_errors, err := validator.ValidateStruct(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, validate_errors))
		resp := entity.NewResponse(true, validate_errors)
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	// Login(данные пользователя), проверили данные, AddSession(user_id), поставили куки

	user, err := c.authUsecase.Register(r.Context(), req.Username, req.Email, req.Password, req.Role, requestID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
		resp := entity.NewResponse(true, err.Error())
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	session, err := c.authUsecase.AddSession(user.ID, user.Role)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
		resp := entity.NewResponse(true, err.Error())
		easyjson.MarshalToWriter(&resp, w)
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
	// json.NewEncoder(w).Encode(entity.NewResponse(false, "registration succesful"))
	resp := entity.NewResponse(false, "registration succesful")
	easyjson.MarshalToWriter(&resp, w)
}
