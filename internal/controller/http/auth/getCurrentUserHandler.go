package auth

import (
	"encoding/json"
	// "fmt"
	"net/http"
	"retarget/internal/entity"
)

type UserResponse struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Balance  int    `json:"balance"`
	Role     int    `json:"role"`
}

func (c *AuthController) GetCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Method Not Allowed"))
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Invalid Cookie"))
		return
	}
	user, err := c.authUsecase.GetUserBySessionID(cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "User Not Found"))
		return
	}

	userResponse := UserResponse{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Avatar:   user.Avatar,
		Balance:  user.Balance,
		Role:     user.Role,
	}
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userResponse) // тут тоже надо решить как поменять структуру
}
