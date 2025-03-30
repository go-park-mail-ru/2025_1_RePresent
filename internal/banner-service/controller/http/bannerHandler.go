package banner

import (
	"encoding/json"
	"net/http"
	sess "retarget/internal/auth-service/entity/auth" // Хардкод
	entity "retarget/internal/banner-service/entity"
	"strconv"

	// pkg "retarget/internal/pkg/entity"
	response "retarget/pkg/entity"
	// "strconv"
	"github.com/gorilla/mux"
)

type CreateUpdateBannerRequest struct {
	OwnerID     int    `json:"owner" validate:"required"`
	Title       string `json:"title" validate:"required,min=3,max=30"`
	Description string `json:"description" validate:"required"`
	Content     string `json:"content_link"`
	Link        string `json:"link" validate:"required"`
	Status      int    `json:"status" validate:"required"`
}

func (h *BannerController) GetUserBanners(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response.NewResponse(true, "Invalid Cookie"))
		return
	}
	user, err := sess.GetSession(cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response.NewResponse(true, "Invalid Cookie"))
		// http.Error(w, "Cookie not found or Invalid session ID", http.StatusUnauthorized)
		return
	}

	banners, err := h.BannerUsecase.GetBannersByUserID(user.UserID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(response.NewResponse(true, "Error fetching banners: "+err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if banners == nil {
		banners = []*entity.Banner{}
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(banners)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.NewResponse(true, "Error encoding banners: "+err.Error()))
	}
}

func (h *BannerController) ReadBanner(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response.NewResponse(true, "Invalid Cookie"))
		return
	}

	// Жду миддлваре
	user, err := sess.GetSession(cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response.NewResponse(true, "Invalid Cookie"))
		// http.Error(w, "Cookie not found or Invalid session ID", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	bannerIDstr := vars["id"]

	bannerID, err := strconv.Atoi(bannerIDstr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.NewResponse(true, "invalid banner ID"))
		return
	}

	banner, err := h.BannerUsecase.GetBannerByID(user.UserID, bannerID)
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
	var req CreateUpdateBannerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(response.NewResponse(true, err.Error()))
		// http.Error(w, "Invalid JSON", http.StatusUnprocessableEntity)
		return
	}

	// Хардкожу вытаскивание юзверя
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response.NewResponse(true, "Invalid Cookie"))
		// http.Error(w, "Cookie not found or Invalid session ID", http.StatusUnauthorized)
		return
	}
	user, err := sess.GetSession(cookie.Value)
	if err != nil || cookie.Value == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response.NewResponse(true, "Invalid Cookie"))
		// http.Error(w, "Cookie not found or Invalid session ID", http.StatusUnauthorized)
		return
	}
	UserID := user.UserID
	//Хардкод закончился

	banner := entity.Banner{
		OwnerID:     UserID,
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

	// Хардкожу вытаскивание юзверя
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response.NewResponse(true, "Invalid Cookie"))
		// http.Error(w, "Cookie not found or Invalid session ID", http.StatusUnauthorized)
		return
	}
	user, err := sess.GetSession(cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response.NewResponse(true, "Invalid Cookie"))
		// http.Error(w, "Cookie not found or Invalid session ID", http.StatusUnauthorized)
		return
	}

	UserID := user.UserID
	//Хардкод закончился

	vars := mux.Vars(r)
	bannerIDstr := vars["id"]

	bannerID, err := strconv.Atoi(bannerIDstr)
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

	err = h.BannerUsecase.UpdateBanner(UserID, banner)
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
	// Хардкожу вытаскивание юзверя
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response.NewResponse(true, "Invalid Cookie"))
		// http.Error(w, "Cookie not found or Invalid session ID", http.StatusUnauthorized)
		return
	}
	user, err := sess.GetSession(cookie.Value)
	UserID := user.UserID
	//Хардкод закончился

	vars := mux.Vars(r)
	bannerIDstr := vars["id"]
	bannerID, err := strconv.Atoi(bannerIDstr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response.NewResponse(true, "invalid banner ID"))
		return
	}

	h.BannerUsecase.BannerRepository.DeleteBannerByID(bannerID, UserID)
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
	if r.Method == http.MethodPost {
		h.CreateBanner(w, r)
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
