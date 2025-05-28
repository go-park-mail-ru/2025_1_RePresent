package auth

import (
	"io"
	"net/http"
	model "retarget/internal/auth-service/easyjsonModels"
	entity "retarget/pkg/entity"
	"retarget/pkg/utils/validator"

	"github.com/mailru/easyjson"
)

func (c *AuthController) LoginHandler(w http.ResponseWriter, r *http.Request) {
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
		//nolint:errcheck
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	var req model.LoginRequest
	data, _ := io.ReadAll(r.Body)
	err := req.UnmarshalJSON(data)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
		resp := entity.NewResponse(true, err.Error())
		//nolint:errcheck
		easyjson.MarshalToWriter(&resp, w)
		return
	}
	errors, err := validator.ValidateStruct(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, errors))
		resp := entity.NewResponse(true, errors)
		//nolint:errcheck
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	user, err := c.authUsecase.Login(r.Context(), req.Email, req.Password, req.Role, requestID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
		resp := entity.NewResponse(true, err.Error())
		//nolint:errcheck
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	session, err := c.authUsecase.AddSession(user.ID, user.Role)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
		resp := entity.NewResponse(true, err.Error())
		//nolint:errcheck
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

	w.WriteHeader(http.StatusOK)
	//nolint:errcheck
	// json.NewEncoder(w).Encode(entity.NewResponse(false, "Login Successful"))
	resp := entity.NewResponse(false, "Login Succesful")
	//nolint:errcheck
	easyjson.MarshalToWriter(&resp, w)
}
