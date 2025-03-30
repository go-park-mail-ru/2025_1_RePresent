package adv

import (
	"encoding/json"
	"fmt"
	"net/http"
	entity "retarget/pkg/entity"
)

func (c *AdvController) GenerateLinkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Method Not Allowed"))
		return
	}

	/*cookie, err := r.Cookie("session_id")
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

	secretLink, isCreated, err := c.advUsecase.PutLink(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewResponse(true, err.Error()))
		return
	}
	fullSecretLink := fmt.Sprintf("http://re-target.ru/api/v1/adv/iframe/%s", secretLink)

	if isCreated {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"service": entity.NewResponse(false, "Created").Service,
			"body": map[string]interface{}{
				"secret-link": fullSecretLink,
			},
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"service": entity.NewResponse(false, "Sent").Service,
		"body": map[string]interface{}{
			"secret-link": fullSecretLink,
		},
	})
}
