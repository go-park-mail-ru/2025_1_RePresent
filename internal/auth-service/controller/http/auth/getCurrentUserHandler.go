package auth

import (
	"encoding/json"

	"net/http"
	entity "pkg/entity"
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

	userSession, ok := r.Context().Value(entity.UserContextKey).(entity.UserContext)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewResponse(true, "Error of authenticator"))
	}
	userID := userSession.UserID

	user, err := c.authUsecase.GetUser(userID)
	if err != nil {
		panic("UNIMPLIMENTED")
	}

	userResponse := &UserResponse{
		Username: user.Username,
		Email:    user.Email,
		Balance:  user.Balance,
		Role:     user.Role,
	}

	response := struct {
		Service entity.Response `json:"service"`
		Body    UserResponse    `json:"body"`
	}{
		Service: entity.NewResponse(false, "Sent"),
		Body:    *userResponse,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
