package http

import (
	usecaseAuth "retarget/internal/auth-service/usecase/auth"
	handlerBanner "retarget/internal/controller/http/banner"
	usecaseBanner "retarget/internal/usecase/banner"
	logs "retarget/pkg/middleware"

	"github.com/gorilla/mux"
)

func SetupRoutes(authUsecase *usecaseAuth.AuthUsecase, bannerUsecase *usecaseBanner.BannerUsecase) *mux.Router {
	r := mux.NewRouter()

	r.Use(logs.LogMiddleware)

	bannerRoutes := handlerBanner.SetupBannerRoutes(authUsecase, bannerUsecase)
	r.PathPrefix("/api/v1/banner/").Handler(bannerRoutes)

	return r
}
