package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"RE/internal/usecase/auth"
)

type AuthController struct {
	usecase *auth.AuthUsecase
}

func NewAuthController(usecase *auth.AuthUsecase) *AuthController {
	return &AuthController{usecase: usecase}
}

func (c *AuthController) GetSuccessHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Получен запрос:", r.URL.Path)
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := struct {
		Success bool `json:"success"`
	}{Success: true}
	json.NewEncoder(w).Encode(response)
}

func SetupAuthRoutes(usecase *auth.AuthUsecase) *http.ServeMux {
	mux := http.NewServeMux()
	authController := NewAuthController(usecase)
	mux.HandleFunc("/", authController.GetSuccessHandler)
	return mux
}
