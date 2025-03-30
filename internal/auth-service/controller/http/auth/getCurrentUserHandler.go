package auth

import (
	"encoding/json"

	"net/http"
	entity "retarget/pkg/entity"
)

type UserResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Balance  int    `json:"balance"`
	Role     int    `json:"role"`
}

func (c *AuthController) GetCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Method Not Allowed"))
		return
	}

	user, ok := r.Context().Value(entity.UserContextKey).(entity.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Error of authenticator"))
	}
	userID := user.UserID
	// Получили пользователя, в usecase GetUser(user_id), получили пользователя, вернули
}
