package auth

import (
	"RE/internal/controller/http/middleware"
	"RE/internal/usecase/auth"
	"net/http"
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

	mux.Handle("/auth/me", middleware.AuthMiddleware(usecase)(http.HandlerFunc(authController.getCurrentUserHandler)))
	mux.Handle("/auth/logout", middleware.AuthMiddleware(usecase)(http.HandlerFunc(authController.LogoutHandler)))

	mux.HandleFunc("/auth/login", authController.LoginHandler)
	mux.HandleFunc("/auth/signup", authController.RegisterHandler)

	return mux
}
