package http

import (
	handlerAuth "RE/internal/controller/http/auth"

	handlerBanner "RE/internal/controller/http/banner"
	logs "RE/internal/controller/http/middleware"
	usecaseAuth "RE/internal/usecase/auth"
	usecaseBanner "RE/internal/usecase/banner"

	"github.com/gorilla/mux"
)

func SetupRoutes(authUsecase *usecaseAuth.AuthUsecase, bannerUsecase *usecaseBanner.BannerUsecase) *mux.Router {
	r := mux.NewRouter()

	r.Use(logs.ErrorMiddleware)

	authRoutes := handlerAuth.SetupAuthRoutes(authUsecase)
	r.PathPrefix("/auth/").Handler(authRoutes)

	bannerRoutes := handlerBanner.SetupBannerRoutes(authUsecase, bannerUsecase)
	r.PathPrefix("/banner/").Handler(bannerRoutes)

	return r
}
