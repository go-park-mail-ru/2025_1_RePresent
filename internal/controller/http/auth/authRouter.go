package auth

import (
	"RE/internal/controller/http/middleware"
	"RE/internal/usecase/auth"
	"net/http"

	"github.com/gorilla/mux"
)

type AuthController struct {
	usecase *auth.AuthUsecase
}

func NewAuthController(usecase *auth.AuthUsecase) *AuthController {
	return &AuthController{usecase: usecase}
}

func SetupAuthRoutes(usecase *auth.AuthUsecase) http.Handler {
	muxRouter := mux.NewRouter()
	authController := NewAuthController(usecase)

	muxRouter.Handle("/auth/me", http.HandlerFunc(authController.getCurrentUserHandler))
	muxRouter.Handle("/auth/logout", middleware.AuthMiddleware(usecase)(http.HandlerFunc(authController.LogoutHandler)))

	muxRouter.HandleFunc("/auth/login", authController.LoginHandler)
	muxRouter.HandleFunc("/auth/signup", authController.RegisterHandler)

	return muxRouter
}
