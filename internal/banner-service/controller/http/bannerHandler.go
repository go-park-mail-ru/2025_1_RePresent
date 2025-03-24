package banner

import (
	"encoding/json"
	"net/http"
	sess "retarget/internal/auth-service/entity/auth" // Хардкод
	entity "retarget/internal/banner-service/entity"
	response "retarget/pkg/entity"
	// "strconv"
	// "github.com/gorilla/mux"
)

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
