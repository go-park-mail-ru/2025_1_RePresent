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

func NewAuthController(authUsecase *auth.AuthUsecase) AuthController {
	return AuthController{authUsecase: authUsecase}
}

func SetupAuthRoutes(authUsecase *auth.AuthUsecase) http.Handler {
	muxRouter := mux.NewRouter()
	authController := NewAuthController(authUsecase)

	muxRouter.Handle("/api/v1/auth/me", http.HandlerFunc(authController.GetCurrentUserHandler))
	muxRouter.Handle("/api/v1/auth/logout", middleware.AuthMiddleware(authUsecase)(http.HandlerFunc(authController.LogoutHandler)))

	muxRouter.HandleFunc("/api/v1/auth/login", authController.LoginHandler)
	muxRouter.HandleFunc("/api/v1/auth/signup", authController.RegisterHandler)

	return muxRouter
}
