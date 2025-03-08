package banner

import (
	"RE/internal/usecase/banner"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type BannerHandler struct {
	BannerUsecase *banner.BannerUsecase
}

// NewBannerHandler создает новый экземпляр BannerHandler
func NewBannerHandler(bannerUsecase *banner.BannerUsecase) *BannerHandler {
	return &BannerHandler{BannerUsecase: bannerUsecase}
}

// GetBannersByUserID обрабатывает запрос для получения всех баннеров пользователя по user_id
func (h *BannerHandler) GetBannersByUserID(w http.ResponseWriter, r *http.Request) {
	// Получаем user_id из URL
	vars := mux.Vars(r)
	userID := vars["user_id"]

	// Получаем баннеры из usecase
	banners, err := h.BannerUsecase.GetBannersByUserID(userID)
	if err != nil {
		http.Error(w, "Error fetching banners: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем ответ в формате JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(banners); err != nil {
		http.Error(w, "Error encoding banners: "+err.Error(), http.StatusInternalServerError)
	}
}
