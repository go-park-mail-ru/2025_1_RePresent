package adv

import (
	"encoding/json"
	"fmt"
	"net/http"
	entity "retarget/pkg/entity"
)

func (c *AdvController) GetLinkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Method Not Allowed"))
		return
	} // Получение списка ссылок
	/*
		cookie, err := r.Cookie("session_id")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(entity.NewResponse(true, "Unauthorized"))
			return
		}

		// TODO: по gRPC получить user_id по cookie из auth-service
		resp, err := http.Get("http://re-target.ru:/api/v1/auth/user-by-cookie?cookie=" + cookie.Value)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
			return
		}
		defer resp.Body.Close()

		var userData map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&userData)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
			return
		}

		userID, ok := userData["user_id"].(string)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(entity.NewResponse(true, "Unauthorized"))
			return
		}*/
	userID := 1
	secretLink, err := c.advUsecase.GetLink(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
		return
	}
	fullSecretLink := fmt.Sprintf("http://re-target.ru/api/v1/adv/iframe/%s", secretLink)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"service": entity.NewResponse(false, "Sent").Service,
		"body": map[string]interface{}{
			"secret-link": fullSecretLink,
		},
	})
}
