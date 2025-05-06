package profile

import (
	"encoding/json"
	"errors"
	"net/http"
	entityProfile "retarget/internal/profile-service/entity/profile"
	entity "retarget/pkg/entity"
)

func (c *ProfileController) GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(entity.СtxKeyRequestID{}).(string)
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Method Not Allowed"))
		return
	}

	user, ok := r.Context().Value(entity.UserContextKey).(entity.UserContext)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Error of authenticator"))
	}
	userID := user.UserID

	profile, err := c.profileUsecase.GetProfile(userID, requestID)
	if profile == nil {
		if errors.Is(err, entityProfile.ErrProfileNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(entity.NewResponse(true, "Profile not found"))
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
		return
	}

	verdict := entity.NewResponse(false, "Sent")
	response := struct {
		Service *entity.ServiceResponse       `json:"service"`
		Body    entityProfile.ProfileResponse `json:"body"`
	}{
		Service: &verdict.Service,
		Body:    *profile,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
