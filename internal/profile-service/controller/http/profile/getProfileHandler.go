package profile

import (
	"encoding/json"
	"errors"
	"net/http"
	entityProfile "retarget/internal/profile-service/entity/profile"
	entity "retarget/pkg/entity"
)

func (c *ProfileController) GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Method Not Allowed"))
		return
	}

	user, ok := r.Context().Value(entity.UserContextKey).(entity.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Error of authenticator"))
	}
	userID := user.UserID

	profile, err := c.profileUsecase.GetProfile(userID)
	if profile == nil {
		if errors.Is(err, entityProfile.ErrProfileNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(entity.NewResponse(true, "Profile not found"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
		return
	}

	response := struct {
		Service entity.Response               `json:"service"`
		Body    entityProfile.ProfileResponse `json:"body"`
	}{
		Service: entity.NewResponse(false, "Sent"),
		Body:    *profile,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
