package http

import (
	// usecaseAuth "retarget/internal/auth-service/usecase/auth"
	handlerBanner "retarget/internal/banner-service/controller/http"
	usecaseBanner "retarget/internal/banner-service/usecase"
	logs "retarget/pkg/middleware"

	"github.com/gorilla/mux"
)

func SetupRoutes(bannerUsecase *usecaseBanner.BannerUsecase) *mux.Router {
	r := mux.NewRouter()

	r.Use(logs.LogMiddleware)

	bannerRoutes := handlerBanner.SetupBannerRoutes(bannerUsecase)
	r.PathPrefix("/api/v1/banner/").Handler(bannerRoutes)

	return r
}
