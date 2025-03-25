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

type CreateBannerRequest struct {
	OwnerID     int    `json:"owner"`
	Title       string `json:"title" validate:"required,min=3,max=30"`
	Description string `json:"description" validate:"required"`
	Content     string `json:"content_link"`
	Link        string `json:"link" validate:"required"`
}

func (h *BannerController) GetBannersByUserCookie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response.NewResponse(true, "Method Not Allowed"))
		// http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response.NewResponse(true, "Invalid Cookie"))
		// http.Error(w, "Cookie not found or Invalid session ID", http.StatusUnauthorized)
		return
	}

	// Хардкожу вытаскивание баннеров юзверя
	// if err != nil {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	json.NewEncoder(w).Encode(response.NewResponse(true, "User Not Found"))
	// 	// http.Error(w, "User not found", http.StatusUnauthorized)
	// 	return
	// }

	// // Получаем user_id из URL
	// vars := mux.Vars(r) // ЛИШНИЙ КОД, НО ПУСТЬ ПОКА БУДЕТ
	// userIdStr := vars["user_id"]
	// userID, err := strconv.Atoi(userIdStr)
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	json.NewEncoder(w).Encode(response.NewResponse(true, "Invalid User ID"))
	// 	// http.Error(w, "Invalid user_id", http.StatusBadRequest)
	// 	return
	// }
	user, err := sess.GetSession(cookie.Value)

	// userID := user.UserID

	// if userID != user.UserID {
	// 	w.WriteHeader(http.StatusUnauthorized)

	// 	json.NewEncoder(w).Encode(response.NewResponse(true, "This user haven`t root on getting this content"))
	// 	// http.Error(w, "This user haven`t root on getting this content", http.StatusUnauthorized)
	// 	return
	// } // ЧТОБЫ НЕ ПЕРЕПИСЫВАТЬ FETCH

	banners, err := h.BannerUsecase.GetBannersByUserID(user.UserID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(response.NewResponse(true, "Error fetching banners: "+err.Error()))
		// http.Error(w, "Error fetching banners: "+err.Error(), http.StatusInternalServerError)
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
		// http.Error(w, "Error encoding banners: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *BannerController) ReadBanner(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response.NewResponse(true, "Method Not Allowed"))
		// http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response.NewResponse(true, "Invalid Cookie"))
		return
	}

	// Жду миддлваре
	user, err := sess.GetSession(cookie.Value)

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
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response.NewResponse(true, "Method Not Allowed"))
		// http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateBannerRequest
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
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response.NewResponse(true, "Method Not Allowed"))
		// http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateBannerRequest
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

func (h *BannerController) DeleteBanner(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response.NewResponse(true, "Method Not Allowed"))
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
