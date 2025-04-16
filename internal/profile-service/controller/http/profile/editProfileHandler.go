package profile

import (
	"encoding/json"
	"net/http"
	entityProfile "retarget/internal/profile-service/entity/profile"
	entity "retarget/pkg/entity"
	"retarget/pkg/utils/validator"
)

func (c *ProfileController) EditProfileHandler(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(entity.СtxKeyRequestID{}).(string)
	if r.Method != http.MethodPut {
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

	var profileRequest entityProfile.ProfileRequest
	err := json.NewDecoder(r.Body).Decode(&profileRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Invalid request body"))
		return
	}

	errorMessages, err := validator.ValidateStruct(profileRequest)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(entity.NewResponse(true, errorMessages))
		return
	}
	err = c.profileUsecase.PutProfile(userID, profileRequest.Username, profileRequest.Description, requestID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(entity.NewResponse(false, "Got and Saved"))
}
