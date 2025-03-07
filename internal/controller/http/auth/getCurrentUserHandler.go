package auth

import (
	"encoding/json"
	"net/http"
)

type UserResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Balance  int    `json:"balance"`
	Role     int    `json:"role"`
}

func (c *AuthController) getCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Cookie not found", http.StatusUnauthorized)
		return
	}

	sessionID := cookie.Value

	user, err := c.usecase.GetUserBySessionID(sessionID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	userResponse := UserResponse{
		Username: user.Username,
		Email:    user.Email,
		Avatar:   user.Avatar,
		Balance:  user.Balance,
		Role:     user.Role,
	}
	w.Header().Set("Cache-Control", "private, max-age=60, must-revalidate")
	json.NewEncoder(w).Encode(userResponse)
}
