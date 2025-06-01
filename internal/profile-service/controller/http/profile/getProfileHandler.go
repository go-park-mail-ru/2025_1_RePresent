package profile

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	entityProfile "retarget/internal/profile-service/entity/profile"
	entity "retarget/pkg/entity"
)

func (c *ProfileController) GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(entity.СtxKeyRequestID{}).(string)
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		//nolint:errcheck
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Method Not Allowed"))
		return
	}

	user, ok := r.Context().Value(entity.UserContextKey).(entity.UserContext)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		//nolint:errcheck
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Error of authenticator"))
	}
	userID := user.UserID

	profile, err := c.profileUsecase.GetProfile(userID, requestID)
	if profile == nil {
		if errors.Is(err, entityProfile.ErrProfileNotFound) {
			w.WriteHeader(http.StatusNotFound)
			//nolint:errcheck
			json.NewEncoder(w).Encode(entity.NewResponse(true, "Profile not found"))
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		//nolint:errcheck
		json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
		return
	}
	fmt.Println("Вываливается в хендлере перед тем как собирать респонз")
	verdict := entity.NewResponse(false, "Sent")
	response := struct {
		Service *entity.ServiceResponse       `json:"service"`
		Body    entityProfile.ProfileResponse `json:"body"`
	}{
		Service: &verdict.Service,
		Body:    *profile,
	}

	w.WriteHeader(http.StatusOK)
	//nolint:errcheck
	json.NewEncoder(w).Encode(response)
}
