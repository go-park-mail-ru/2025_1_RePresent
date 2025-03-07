package http

import (
	"net/http"

	handlerAuth "RE/internal/controller/http/auth"
	usecaseAuth "RE/internal/usecase/auth"
)

func SetupRoutes(usecase *usecaseAuth.AuthUsecase) *http.ServeMux {
	mux := http.NewServeMux()

	authRoutes := handlerAuth.SetupAuthRoutes(usecase)
	mux.Handle("/auth/", authRoutes)
	// mux.Handle("/auth/", http.StripPrefix("/auth", authRoutes)) ТАК НЕ ДЕЛАТЬ !!!
	return mux
}
