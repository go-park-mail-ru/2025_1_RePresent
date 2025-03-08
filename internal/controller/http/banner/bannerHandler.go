package banner

import (
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
	println(vars)
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
