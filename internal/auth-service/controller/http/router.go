package http

import (
	authenticate "pkg/middleware/auth"
	handlerAuth "retarget-authapp/controller/http/auth"
	usecaseAuth "retarget-authapp/usecase/auth"

	"github.com/gorilla/mux"
)

func SetupRoutes(authenticator *authenticate.Authenticator, authUsecase *usecaseAuth.AuthUsecase) *mux.Router {
	r := mux.NewRouter()

	authRoutes := handlerAuth.SetupAuthRoutes(authenticator, authUsecase)
	r.PathPrefix("/api/v1/auth/").Handler(authRoutes)

	return r
}
