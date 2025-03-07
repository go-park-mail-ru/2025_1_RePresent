package auth

import (
	"net/http"

	"RE/internal/usecase/auth"
)

type AuthController struct {
	usecase *auth.AuthUsecase
}

func NewAuthController(usecase *auth.AuthUsecase) *AuthController {
	return &AuthController{usecase: usecase}
}

func SetupAuthRoutes(usecase *auth.AuthUsecase) http.Handler {
	mux := http.NewServeMux()
	authController := NewAuthController(usecase)

	mux.HandleFunc("/auth/login", authController.LoginHandler)
	mux.HandleFunc("/auth/signup", authController.RegisterHandler)
	mux.HandleFunc("/auth/logout", authController.LogoutHandler)
	mux.HandleFunc("/auth/me", authController.getCurrentUserHandler)

	return mux
}
