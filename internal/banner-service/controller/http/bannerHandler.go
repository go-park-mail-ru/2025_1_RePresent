package controller

import (
	"encoding/json"
	"io"
	"net/http"
	model "retarget/internal/banner-service/easyjsonModels"
	response "retarget/pkg/entity"
	validator "retarget/pkg/utils/validator"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
)

func (h *BannerController) GetUserBanners(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(response.СtxKeyRequestID{}).(string)
	userSession, ok := r.Context().Value(response.UserContextKey).(response.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(response.NewResponse(true, "Error of authenticator"))
		resp := response.NewResponse(true, "Error of authenticator")
		easyjson.MarshalToWriter(&resp, w)
	}
	userID := userSession.UserID

	banners, err := h.BannerUsecase.GetBannersByUserID(userID, requestID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(response.NewResponse(true, "Error fetching banners: "+err.Error()))
		resp := response.NewResponse(true, "Error fetching banners: "+err.Error())
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = easyjson.MarshalToWriter(&banners, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(response.NewResponse(true, "Error encoding banners: "+err.Error()))
		resp := response.NewResponse(true, "Error encoding banners: "+err.Error())
		easyjson.MarshalToWriter(&resp, w)
	}
}

func (h *BannerController) ReadBanner(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(response.СtxKeyRequestID{}).(string)
	userSession, ok := r.Context().Value(response.UserContextKey).(response.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(response.NewResponse(true, "Error of authenticator"))
		resp := response.NewResponse(true, "Error of authenticator")
		easyjson.MarshalToWriter(&resp, w)
	}
	userID := userSession.UserID

	vars := mux.Vars(r)
	bannerIDstr := vars["banner_id"]

	bannerID, err := strconv.Atoi(bannerIDstr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(response.NewResponse(true, "invalid banner ID"))
		resp := response.NewResponse(true, "invalid banner ID")
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	banner, err := h.BannerUsecase.GetBannerByID(userID, bannerID, requestID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(response.NewResponse(true, err.Error()))
		resp := response.NewResponse(true, err.Error())
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(banner); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(response.NewResponse(true, "error encoding banners: "+err.Error()))
		resp := response.NewResponse(true, err.Error())
		easyjson.MarshalToWriter(&resp, w)
	}
}

func (h *BannerController) CreateBanner(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(response.СtxKeyRequestID{}).(string)
	var req model.CreateUpdateBannerRequest
	data, _ := io.ReadAll(r.Body)
	err := req.UnmarshalJSON(data)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(response.NewResponse(true, err.Error()))
		resp := response.NewResponse(true, err.Error())
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	validate_errors, err := validator.ValidateStruct(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(response.NewResponse(true, validate_errors))
		resp := response.NewResponse(true, validate_errors)
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	userSession, ok := r.Context().Value(response.UserContextKey).(response.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(response.NewResponse(true, "Error of authenticator"))
		resp := response.NewResponse(true, "Error of authenticator")
		easyjson.MarshalToWriter(&resp, w)
	}
	userID := userSession.UserID

	banner := model.Banner{
		OwnerID:     userID,
		Title:       req.Title,
		Description: req.Description,
		Content:     req.Content,
		Link:        req.Link,
		Balance:     0,
		Status:      req.Status,
		MaxPrice:    req.MaxPrice,
	}

	if err := h.BannerUsecase.BannerRepository.CreateNewBanner(banner, requestID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(response.NewResponse(true, err.Error()))
		resp := response.NewResponse(true, err.Error())
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	w.WriteHeader(http.StatusCreated)
	//nolint:errcheck
	// json.NewEncoder(w).Encode(response.NewResponse(false, "Banner created"))
	resp := response.NewResponse(false, "Banner created")
	easyjson.MarshalToWriter(&resp, w)
}

func (h *BannerController) UpdateBanner(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(response.СtxKeyRequestID{}).(string)
	var req model.CreateUpdateBannerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(response.NewResponse(true, err.Error()))
		resp := response.NewResponse(true, err.Error())
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	if validate_errors, err := validator.ValidateStruct(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(response.NewResponse(true, validate_errors))
		resp := response.NewResponse(true, validate_errors)
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	userSession, ok := r.Context().Value(response.UserContextKey).(response.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(response.NewResponse(true, "Error of authenticator"))
		resp := response.NewResponse(true, "Error of authenticator")
		easyjson.MarshalToWriter(&resp, w)
	}
	userID := userSession.UserID

	vars := mux.Vars(r)
	bannerIDstr := vars["banner_id"]
	bannerID, err := strconv.Atoi(bannerIDstr)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(response.NewResponse(true, "invalid banner ID"))
		resp := response.NewResponse(true, "Invalid banner ID")
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	banner := model.Banner{
		ID:          bannerID,
		Title:       req.Title,
		Description: req.Description,
		Link:        req.Link,
		Content:     req.Content,
		Status:      req.Status,
		MaxPrice:    req.MaxPrice,
	}

	if err := h.BannerUsecase.UpdateBanner(userID, banner, requestID); err != nil {
		w.WriteHeader(http.StatusForbidden)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(response.NewResponse(true, err.Error()))
		resp := response.NewResponse(true, err.Error())
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	//nolint:errcheck
	// json.NewEncoder(w).Encode(response.NewResponse(false, "Banner updated"))
	resp := response.NewResponse(false, "Banner updated")
	easyjson.MarshalToWriter(&resp, w)
}

func (h *BannerController) DeleteBanner(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(response.СtxKeyRequestID{}).(string)
	userSession, ok := r.Context().Value(response.UserContextKey).(response.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(response.NewResponse(true, "Error of authenticator"))
		resp := response.NewResponse(true, "Error of authenticator")
		easyjson.MarshalToWriter(&resp, w)
		return
	}
	userID := userSession.UserID

	vars := mux.Vars(r)
	bannerIDstr := vars["banner_id"]
	bannerID, err := strconv.Atoi(bannerIDstr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		// json.NewEncoder(w).Encode(response.NewResponse(true, "invalid banner ID"))
		resp := response.NewResponse(true, "invalid banner ID")
		easyjson.MarshalToWriter(&resp, w)
		return
	}

	if err := h.BannerUsecase.BannerRepository.DeleteBannerByID(userID, bannerID, requestID); err != nil {
		// response := response.NewResponse(true, "failed to delete banner")
		w.WriteHeader(http.StatusBadRequest)
		// _ = json.NewEncoder(w).Encode(response)
		resp := response.NewResponse(true, "failed to delete banner")
		easyjson.MarshalToWriter(&resp, w)
		return
	}
	w.WriteHeader(http.StatusOK)
	//nolint:errcheck
	// json.NewEncoder(w).Encode(response.NewResponse(false, "Banner deleted"))
	resp := response.NewResponse(false, "banner deleted")
	easyjson.MarshalToWriter(&resp, w)

}
