package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	entity "retarget/internal/banner-service/entity"
	response "retarget/pkg/entity"
	"strconv"

	"github.com/gorilla/mux"
)

type CreateUpdateBannerRequest struct {
	// OwnerID     int    `json:"owner" validate:"required"`
	Title       string `json:"title" validate:"required,min=3,max=30"`
	Description string `json:"description" validate:"required"`
	Content     string `json:"content_link"`
	Link        string `json:"link" validate:"required"`
	Status      int    `json:"status" validate:"required"`
}

func (h *BannerController) GetUserBanners(w http.ResponseWriter, r *http.Request) {

	userSession, ok := r.Context().Value(response.UserContextKey).(response.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.NewResponse(true, "Error of authenticator"))
	}
	userID := userSession.UserID

	banners, err := h.BannerUsecase.GetBannersByUserID(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(response.NewResponse(true, "Error fetching banners: "+err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(banners)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.NewResponse(true, "Error encoding banners: "+err.Error()))
	}
}

func (h *BannerController) ReadBanner(w http.ResponseWriter, r *http.Request) {
	userSession, ok := r.Context().Value(response.UserContextKey).(response.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.NewResponse(true, "Error of authenticator"))
	}
	userID := userSession.UserID

	vars := mux.Vars(r)
	bannerIDstr := vars["banner_id"]

	bannerID, err := strconv.Atoi(bannerIDstr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.NewResponse(true, "invalid banner ID"))
		return
	}

	banner, err := h.BannerUsecase.GetBannerByID(userID, bannerID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.NewResponse(true, err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(banner)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.NewResponse(true, "error encoding banners: "+err.Error()))
	}
}

func (h *BannerController) CreateBanner(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response.NewResponse(true, "Method Not Allowed"))
	}
	var req CreateUpdateBannerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(response.NewResponse(true, err.Error()))
		// http.Error(w, "Invalid JSON", http.StatusUnprocessableEntity)
		return
	}

	userSession, ok := r.Context().Value(response.UserContextKey).(response.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.NewResponse(true, "Error of authenticator"))
	}
	userID := userSession.UserID

	//Хардкод закончился

	banner := entity.Banner{
		OwnerID:     userID,
		Title:       req.Title,
		Description: req.Description,
		Content:     req.Content,
		Link:        req.Link,
		Balance:     0,
		Status:      0,
	}

	err = h.BannerUsecase.BannerRepository.CreateNewBanner(banner)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response.NewResponse(true, err.Error()))
		// http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response.NewResponse(false, "Banner created"))
}

func (h *BannerController) UpdateBanner(w http.ResponseWriter, r *http.Request) {
	var req CreateUpdateBannerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(response.NewResponse(true, err.Error()))
		// http.Error(w, "Invalid JSON", http.StatusUnprocessableEntity)
		return
	}

	userSession, ok := r.Context().Value(response.UserContextKey).(response.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.NewResponse(true, "Error of authenticator"))
	}
	userID := userSession.UserID

	vars := mux.Vars(r)
	bannerIDstr := vars["banner_id"]
	fmt.Println("Строка баннера: ", bannerIDstr)
	bannerID, err := strconv.Atoi(bannerIDstr)
	fmt.Println("Число баннера: ", bannerID)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(response.NewResponse(true, "invalid banner ID"))
		return
	}

	banner := entity.Banner{
		ID:          bannerID,
		Title:       req.Title,
		Description: req.Description,
		Content:     req.Content,
		Status:      req.Status,
	}

	err = h.BannerUsecase.UpdateBanner(userID, banner)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(response.NewResponse(true, err.Error()))
		// http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response.NewResponse(false, "Banner updated"))
}

func (h *BannerController) DeleteBanner(w http.ResponseWriter, r *http.Request) {
	userSession, ok := r.Context().Value(response.UserContextKey).(response.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.NewResponse(true, "Error of authenticator"))
	}
	userID := userSession.UserID

	vars := mux.Vars(r)
	bannerIDstr := vars["id"]
	bannerID, err := strconv.Atoi(bannerIDstr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.NewResponse(true, "invalid banner ID"))
		return
	}

	h.BannerUsecase.BannerRepository.DeleteBannerByID(bannerID, userID)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response.NewResponse(false, "Banner deleted"))

}

func (h *BannerController) BannerHandleFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.ReadBanner(w, r)
		return
	}
	if r.Method == http.MethodDelete {
		h.DeleteBanner(w, r)
		return
	}
	if r.Method == http.MethodPut {
		h.UpdateBanner(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
	json.NewEncoder(w).Encode(response.NewResponse(true, "Method Not Allowed"))
	return
}
