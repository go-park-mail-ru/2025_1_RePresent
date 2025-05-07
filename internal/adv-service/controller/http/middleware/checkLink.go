package middleware

import (
	"encoding/json"
	"net/http"
	usecaseSlot "retarget/internal/adv-service/usecase/slot"
	entity "retarget/pkg/entity"

	"github.com/gorilla/mux"
)

func LinkMiddleware(slotUsecase *usecaseSlot.SlotUsecase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			link := vars["link"]
			if err := slotUsecase.CheckLink(link); err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
