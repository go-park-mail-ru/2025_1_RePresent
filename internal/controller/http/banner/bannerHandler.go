package banner

import (
	"RE/internal/entity"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// GetBannersByUserID обрабатывает запрос для получения всех баннеров пользователя по user_id
func (h *BannerController) GetBannersByUserId(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем user_id из URL
	vars := mux.Vars(r)
	userID := vars["user_id"]

	// Получаем баннеры из usecase
	banners, err := h.BannerUsecase.GetBannersByUserID(userID)
	if err != nil {
		http.Error(w, "Error fetching banners: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if banners == nil {
		banners = []*entity.Banner{}
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(banners)
	if err != nil {
		http.Error(w, "Error encoding banners: "+err.Error(), http.StatusInternalServerError)
	}
}
