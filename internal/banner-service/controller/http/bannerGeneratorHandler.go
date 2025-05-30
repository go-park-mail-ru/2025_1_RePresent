package controller

import (
	"net/http"
	response "retarget/pkg/entity"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
)

func (h *BannerController) GenerateDescription(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(response.СtxKeyRequestID{}).(string)

	userSession, ok := r.Context().Value(response.UserContextKey).(response.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(response.NewResponse(true, "Error of authenticator"))
		resp := response.NewResponse(true, "Error of authenticator")
		//nolint:errcheck
		easyjson.MarshalToWriter(&resp, w)
	}
	userID := userSession.UserID

	vars := mux.Vars(r)
	bannerIDstr := vars["banner_id"]
	bannerID, err := strconv.Atoi(bannerIDstr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp := response.NewResponse(true, "invalid banner ID")
		//nolint:errcheck
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	description, err := h.BannerUsecase.GenerateBannerDescription(userID, bannerID, requestID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := response.NewResponse(true, err.Error())
		//nolint:errcheck
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	resp := response.NewResponse(false, description)
	//nolint:errcheck
	easyjson.MarshalToWriter(&resp, w)
}
