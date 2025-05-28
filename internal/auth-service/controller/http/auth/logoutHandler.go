package auth

import (
	"net/http"
	entity "retarget/pkg/entity"
	"time"

	"github.com/mailru/easyjson"
)

func (c *AuthController) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(entity.UserContextKey).(entity.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, "Error of authenticator"))
		resp := entity.NewResponse(true, "Error of authenticator")
		easyjson.MarshalToWriter(&resp, w)
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
		resp := entity.NewResponse(true, err.Error())
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	expires := time.Unix(0, 0)
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Expires:  expires,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   -1,
	})

	err = c.authUsecase.Logout(cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
		resp := entity.NewResponse(true, err.Error())
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	//nolint:errcheck
	// json.NewEncoder(w).Encode(entity.NewResponse(false, "Logout Successful"))
	resp := entity.NewResponse(false, "Logout Successful")
	easyjson.MarshalToWriter(&resp, w)
}
