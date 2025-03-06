package http

import (
	"fmt"
	"net/http"

	handlerAuth "RE/internal/controller/http/auth"
	usecaseAuth "RE/internal/usecase/auth"
)

func SetupRoutes(usecase *usecaseAuth.AuthUsecase) *http.ServeMux {
	mux := http.NewServeMux()
	authRoutes := handlerAuth.SetupAuthRoutes(usecase)
	fmt.Println("authRouter успешно подключен")
	mux.Handle("/auth", authRoutes)
	return mux
}
