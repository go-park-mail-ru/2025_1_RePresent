package auth

import (
	"net/http"
	authMiddleware "retarget/internal/auth-service/controller/http/middleware"
	usecaseAuth "retarget/internal/auth-service/usecase/auth"

	"github.com/gorilla/mux"
)

type AuthController struct {
	authUsecase *usecaseAuth.AuthUsecase
}

func NewAuthController(authUsecase *usecaseAuth.AuthUsecase) AuthController {
	return AuthController{authUsecase: authUsecase}
}

func SetupAuthRoutes(authUsecase *usecaseAuth.AuthUsecase) http.Handler {
	muxRouter := mux.NewRouter()
	authController := NewAuthController(authUsecase)

	muxRouter.Handle("/api/v1/auth/me", http.HandlerFunc(authController.GetCurrentUserHandler))
	muxRouter.Handle("/api/v1/auth/logout", authMiddleware.AuthMiddleware(authUsecase)(http.HandlerFunc(authController.LogoutHandler)))

	muxRouter.HandleFunc("/api/v1/auth/login", authController.LoginHandler)
	muxRouter.HandleFunc("/api/v1/auth/signup", authController.RegisterHandler)

	return muxRouter
}
