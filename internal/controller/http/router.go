package http

import (
	handlerAuth "retarget/internal/controller/http/auth"

	handlerBanner "retarget/internal/controller/http/banner"
	logs "retarget/internal/controller/http/middleware"
	usecaseAuth "retarget/internal/usecase/auth"
	usecaseBanner "retarget/internal/usecase/banner"

	"github.com/gorilla/mux"
)

func SetupRoutes(authUsecase *usecaseAuth.AuthUsecase, bannerUsecase *usecaseBanner.BannerUsecase) *mux.Router {
	r := mux.NewRouter()

	r.Use(logs.ErrorMiddleware)

	authRoutes := handlerAuth.SetupAuthRoutes(authUsecase)
	r.PathPrefix("/api/v1/auth/").Handler(authRoutes)

	bannerRoutes := handlerBanner.SetupBannerRoutes(authUsecase, bannerUsecase)
	r.PathPrefix("/api/v1/banner/").Handler(bannerRoutes)

	return r
}
