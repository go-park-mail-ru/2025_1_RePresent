package auth

import (
	"net/http"
	"retarget/internal/controller/http/middleware"
	"retarget/internal/usecase/auth"

	"github.com/gorilla/mux"
)

type AuthController struct {
	authUsecase *auth.AuthUsecase
}

func NewAuthController(authUsecase *auth.AuthUsecase) *AuthController {
	return &AuthController{authUsecase: authUsecase}
}

func SetupAuthRoutes(authUsecase *auth.AuthUsecase) http.Handler {
	muxRouter := mux.NewRouter()
	authController := NewAuthController(authUsecase)

	muxRouter.Handle("/auth/me", http.HandlerFunc(authController.getCurrentUserHandler))
	muxRouter.Handle("/auth/logout", middleware.AuthMiddleware(authUsecase)(http.HandlerFunc(authController.LogoutHandler)))

	muxRouter.HandleFunc("/auth/login", authController.LoginHandler)
	muxRouter.HandleFunc("/auth/signup", authController.RegisterHandler)

	return muxRouter
}
