package auth

import (
	"encoding/json"
	"net/http"
	model "retarget/internal/auth-service/easyjsonModels"
	entity "retarget/pkg/entity"

	"github.com/mailru/easyjson"
)

func (c *AuthController) GetCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	var requestID string
	if v := r.Context().Value(entity.СtxKeyRequestID{}); v != nil {
		if id, ok := v.(string); ok {
			requestID = id
		}
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		//nolint:errcheck
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Method Not Allowed"))
		return
	}

	userSession, ok := r.Context().Value(entity.UserContextKey).(entity.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, "Error of authenticator"))
		resp := entity.NewResponse(true, "Error of authenticator")
		easyjson.MarshalToWriter(&resp, w)
	}
	userID := userSession.UserID

	user, err := c.authUsecase.GetUser(r.Context(), userID, requestID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, "Failed to get user"))
		resp := entity.NewResponse(true, "Failed to get user")
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	userResponse := &model.UserResponse{
		Username: user.Username,
		Email:    user.Email,
		Balance:  *user.Balance.Dec,
		Role:     user.Role,
	}

	response := model.UserResponseWithErr{
		Service: entity.NewResponse(false, "Sent"),
		Body:    *userResponse,
	}

	w.WriteHeader(http.StatusOK)
	//nolint:errcheck
	resp := response
	easyjson.MarshalToWriter(&resp, w)
}
