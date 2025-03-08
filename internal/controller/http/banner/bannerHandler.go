package banner

import (
	"RE/internal/entity"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (h *BannerController) GetBannersByUserCookie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		http.Error(w, "Cookie not found or Invalid session ID", http.StatusUnauthorized)
		return
	}

	user, err := h.AuthUsecase.GetUserBySessionID(cookie.Value)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Получаем user_id из URL
	vars := mux.Vars(r) // ЛИШНИЙ КОД, НО ПУСТЬ ПОКА БУДЕТ
	userIdStr := vars["user_id"]
	userID, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	if userID != user.ID {
		http.Error(w, "This user haven`t root on getting this content", http.StatusUnauthorized)
		return
	} // ЧТОБЫ НЕ ПЕРЕПИСЫВАТЬ FETCH

	banners, err := h.BannerUsecase.GetBannersByUserID(user.ID)
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
