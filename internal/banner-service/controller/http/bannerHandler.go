package controller

import (
	"encoding/json"
	"net/http"
	entity "retarget/internal/banner-service/entity"
	response "retarget/pkg/entity"
	validator "retarget/pkg/utils/validator"
	"strconv"

	"github.com/gorilla/mux"
)

type CreateUpdateBannerRequest struct {
	Title       string         `json:"title" validate:"required,min=3,max=30"`
	Description string         `json:"description" validate:"max=100"`
	Content     string         `json:"content" validate:"required,len=32"`
	Link        string         `json:"link" validate:"required,max=100"`
	Status      int            `json:"status"`
	MaxPrice    entity.Decimal `json:"max_price"`
}

func (h *BannerController) GetUserBanners(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(response.СtxKeyRequestID{}).(string)
	userSession, ok := r.Context().Value(response.UserContextKey).(response.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		json.NewEncoder(w).Encode(response.NewResponse(true, "Error of authenticator"))
	}
	userID := userSession.UserID

	banners, err := h.BannerUsecase.GetBannersByUserID(userID, requestID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		json.NewEncoder(w).Encode(response.NewResponse(true, "Error fetching banners: "+err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(banners); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		json.NewEncoder(w).Encode(response.NewResponse(true, "Error encoding banners: "+err.Error()))
	}
}

func (h *BannerController) ReadBanner(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(response.СtxKeyRequestID{}).(string)
	userSession, ok := r.Context().Value(response.UserContextKey).(response.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		json.NewEncoder(w).Encode(response.NewResponse(true, "Error of authenticator"))
	}
	userID := userSession.UserID

	vars := mux.Vars(r)
	bannerIDstr := vars["banner_id"]

	bannerID, err := strconv.Atoi(bannerIDstr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		json.NewEncoder(w).Encode(response.NewResponse(true, "invalid banner ID"))
		return
	}

	banner, err := h.BannerUsecase.GetBannerByID(userID, bannerID, requestID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		json.NewEncoder(w).Encode(response.NewResponse(true, err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(banner); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		json.NewEncoder(w).Encode(response.NewResponse(true, "error encoding banners: "+err.Error()))
	}
}

func (h *BannerController) CreateBanner(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(response.СtxKeyRequestID{}).(string)
	var req CreateUpdateBannerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		//nolint:errcheck
		json.NewEncoder(w).Encode(response.NewResponse(true, err.Error()))
		return
	}

	validate_errors, err := validator.ValidateStruct(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		//nolint:errcheck
		json.NewEncoder(w).Encode(response.NewResponse(true, validate_errors))
		return
	}

	userSession, ok := r.Context().Value(response.UserContextKey).(response.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		json.NewEncoder(w).Encode(response.NewResponse(true, "Error of authenticator"))
	}
	userID := userSession.UserID

	banner := entity.Banner{
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
		json.NewEncoder(w).Encode(response.NewResponse(true, err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	//nolint:errcheck
	json.NewEncoder(w).Encode(response.NewResponse(false, "Banner created"))
}

func (h *BannerController) UpdateBanner(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(response.СtxKeyRequestID{}).(string)
	var req CreateUpdateBannerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		//nolint:errcheck
		json.NewEncoder(w).Encode(response.NewResponse(true, err.Error()))
		return
	}

	if validate_errors, err := validator.ValidateStruct(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		//nolint:errcheck
		json.NewEncoder(w).Encode(response.NewResponse(true, validate_errors))
		return
	}

	userSession, ok := r.Context().Value(response.UserContextKey).(response.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		json.NewEncoder(w).Encode(response.NewResponse(true, "Error of authenticator"))
	}
	userID := userSession.UserID

	vars := mux.Vars(r)
	bannerIDstr := vars["banner_id"]
	bannerID, err := strconv.Atoi(bannerIDstr)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		//nolint:errcheck
		json.NewEncoder(w).Encode(response.NewResponse(true, "invalid banner ID"))
		return
	}

	banner := entity.Banner{
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
		json.NewEncoder(w).Encode(response.NewResponse(true, err.Error()))
		return
	}

	w.WriteHeader(http.StatusAccepted)
	//nolint:errcheck
	json.NewEncoder(w).Encode(response.NewResponse(false, "Banner updated"))
}

func (h *BannerController) DeleteBanner(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value(response.СtxKeyRequestID{}).(string)
	userSession, ok := r.Context().Value(response.UserContextKey).(response.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		json.NewEncoder(w).Encode(response.NewResponse(true, "Error of authenticator"))
	}
	userID := userSession.UserID

	vars := mux.Vars(r)
	bannerIDstr := vars["banner_id"]
	bannerID, err := strconv.Atoi(bannerIDstr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//nolint:errcheck
		json.NewEncoder(w).Encode(response.NewResponse(true, "invalid banner ID"))
		return
	}

	if err := h.BannerUsecase.BannerRepository.DeleteBannerByID(userID, bannerID, requestID); err != nil {
		response := response.NewResponse(true, "failed to delete banner")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusOK)
	//nolint:errcheck
	json.NewEncoder(w).Encode(response.NewResponse(false, "Banner deleted"))

}
