package http

import (
	handlerAuth "retarget/internal/auth-service/controller/http/auth"
	usecaseAuth "retarget/internal/auth-service/usecase/auth"

	"github.com/gorilla/mux"
)

func SetupRoutes(authUsecase *usecaseAuth.AuthUsecase) *mux.Router {
	r := mux.NewRouter()

	authRoutes := handlerAuth.SetupAuthRoutes(authUsecase)
	r.PathPrefix("/api/v1/auth/").Handler(authRoutes)

	return r
}
