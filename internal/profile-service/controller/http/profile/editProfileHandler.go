package profile

import (
	"encoding/json"
	"net/http"
	entityProfile "retarget/internal/profile-service/entity/profile"
	entity "retarget/pkg/entity"
	"retarget/pkg/utils/validator"

	"github.com/mailru/easyjson"
)

func (c *ProfileController) EditProfileHandler(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(entity.СtxKeyRequestID{}).(string)
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, "Method Not Allowed"))
		resp := entity.NewResponse(true, "Method Not Allowed")
		//nolint:errcheck
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	user, ok := r.Context().Value(entity.UserContextKey).(entity.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		resp := entity.NewResponse(true, "Error of authenticator")
		//nolint:errcheck
		_, err := easyjson.MarshalToWriter(&resp, w)
		if err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}
	}
	userID := user.UserID

	var profileRequest entityProfile.ProfileRequest
	err := json.NewDecoder(r.Body).Decode(&profileRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, "Invalid request body"))
		resp := entity.NewResponse(true, "Invalid request body")
		//nolint:errcheck
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	errorMessages, err := validator.ValidateStruct(profileRequest)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, errorMessages))
		resp := entity.NewResponse(true, errorMessages)
		//nolint:errcheck
		easyjson.MarshalToWriter(&resp, w)
		return
	}
	err = c.profileUsecase.PutProfile(userID, profileRequest.Username, profileRequest.Description, requestID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
		resp := entity.NewResponse(true, err.Error())
		//nolint:errcheck
		easyjson.MarshalToWriter(&resp, w)
		return
	}
	w.WriteHeader(http.StatusOK)
	//nolint:errcheck
	// json.NewEncoder(w).Encode(entity.NewResponse(false, "Got and Saved"))
	resp := entity.NewResponse(false, "Got and Saved")
	//nolint:errcheck
	easyjson.MarshalToWriter(&resp, w)
}
