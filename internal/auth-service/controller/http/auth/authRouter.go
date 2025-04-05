package auth

import (
	"net/http"
	authenticate "pkg/middleware/auth"
	usecaseAuth "retarget-authapp/usecase/auth"

	"github.com/gorilla/mux"
)

type AuthController struct {
	authUsecase *usecaseAuth.AuthUsecase
}

func NewAuthController(authUsecase *usecaseAuth.AuthUsecase) AuthController {
	return AuthController{authUsecase: authUsecase}
}

func SetupAuthRoutes(authenticator *authenticate.Authenticator, authUsecase *usecaseAuth.AuthUsecase) http.Handler {
	muxRouter := mux.NewRouter()
	authController := NewAuthController(authUsecase)

	muxRouter.Handle("/api/v1/auth/me", authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(authController.GetCurrentUserHandler)))
	muxRouter.Handle("/api/v1/auth/logout", authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(authController.LogoutHandler)))

	muxRouter.HandleFunc("/api/v1/auth/login", authController.LoginHandler)
	// muxRouter.HandleFunc("/api/v1/auth/login/mail", authController.LoginConfirmHandler)

	muxRouter.HandleFunc("/api/v1/auth/signup", authController.RegisterHandler)
	muxRouter.HandleFunc("/api/v1/auth/signup/mail", authController.RegisterConfirmHandler)

	// muxRouter.HandleFunc("/api/v1/auth/regain", authController.RegainHandler)
	// muxRouter.HandleFunc("/api/v1/auth/regain/mail", authController.RegainConfirmHandler)

	// muxRouter.HandleFunc("/api/v1/auth/edit/password", authController.EditPasswordHandler)
	// muxRouter.HandleFunc("/api/v1/auth/edit/password/mail", authController.EditPasswordConfirmHandler)

	// muxRouter.HandleFunc("/api/v1/auth/edit/email", authController.EditMailHandler)
	// muxRouter.HandleFunc("/api/v1/auth/edit/email/mail", authController.EditMailConfirmHandler)

	return muxRouter
}
