package http

import (
	"net/http"

	handlerAuth "RE/internal/controller/http/auth"

	handlerBanner "RE/internal/controller/http/banner"
	usecaseAuth "RE/internal/usecase/auth"
	usecaseBanner "RE/internal/usecase/banner"
)

func SetupRoutes(usecase *usecaseAuth.AuthUsecase, bannerUsecase *usecaseBanner.BannerUsecase) *http.ServeMux {
	mux := http.NewServeMux()

	authRoutes := handlerAuth.SetupAuthRoutes(usecase)
	mux.Handle("/auth/", authRoutes)
	// mux.Handle("/auth/", http.StripPrefix("/auth", authRoutes)) ТАК НЕ ДЕЛАТЬ !!!

	bannerRoutes := handlerBanner.SetupBannerRoutes(usecase, bannerUsecase)
	mux.Handle("/banner/", bannerRoutes)

	return mux
}
