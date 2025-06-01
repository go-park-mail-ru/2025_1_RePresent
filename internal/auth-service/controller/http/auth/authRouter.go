package auth

import (
	"context"
	"net/http"
	entityAuth "retarget/internal/auth-service/entity/auth"
	usecaseAuth "retarget/internal/auth-service/usecase/auth"
	logger "retarget/pkg/middleware"
	authenticate "retarget/pkg/middleware/auth"

	"github.com/gorilla/mux"
)

type AuthUsecase interface {
	Login(ctx context.Context, email, pass string, role int, reqID string) (*entityAuth.User, error)
	AddSession(userID, role int) (*entityAuth.Session, error)
	GetUser(ctx context.Context, id int, reqID string) (*entityAuth.User, error)
	Logout(sessionID string) error
	Register(ctx context.Context, username, email, password string, role int, reqID string) (*entityAuth.User, error)
}

type AuthController struct {
	authUsecase *usecaseAuth.AuthUsecase
}

func NewAuthController(authUsecase *usecaseAuth.AuthUsecase) AuthController {
	return AuthController{authUsecase: authUsecase}
}

func SetupAuthRoutes(authenticator *authenticate.Authenticator, authUsecase *usecaseAuth.AuthUsecase) http.Handler {
	muxRouter := mux.NewRouter()
	authController := NewAuthController(authUsecase)

	muxRouter.Handle("/api/v1/auth/me", logger.LogMiddleware(authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(authController.GetCurrentUserHandler)))).Methods("GET")
	muxRouter.Handle("/api/v1/auth/logout", logger.LogMiddleware(authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(authController.LogoutHandler)))).Methods("POST")

	muxRouter.Handle("/api/v1/auth/login", logger.LogMiddleware(http.HandlerFunc(authController.LoginHandler)))
	// muxRouter.HandleFunc("/api/v1/auth/login/mail", authController.LoginConfirmHandler)

	muxRouter.Handle("/api/v1/auth/signup", logger.LogMiddleware(http.HandlerFunc(authController.RegisterHandler)))
	muxRouter.Handle("/api/v1/auth/signup/mail", logger.LogMiddleware(http.HandlerFunc(authController.RegisterConfirmHandler)))

	// muxRouter.HandleFunc("/api/v1/auth/regain", authController.RegainHandler)
	// muxRouter.HandleFunc("/api/v1/auth/regain/mail", authController.RegainConfirmHandler)

	// muxRouter.HandleFunc("/api/v1/auth/edit/password", authController.EditPasswordHandler)
	// muxRouter.HandleFunc("/api/v1/auth/edit/password/mail", authController.EditPasswordConfirmHandler)

	// muxRouter.HandleFunc("/api/v1/auth/edit/email", authController.EditMailHandler)
	// muxRouter.HandleFunc("/api/v1/auth/edit/email/mail", authController.EditMailConfirmHandler)

	return muxRouter
}
