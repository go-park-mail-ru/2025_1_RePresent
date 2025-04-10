package http

import (
	// usecaseAuth "retarget/internal/auth-service/usecase/auth"
	handlerBanner "retarget/internal/banner-service/controller/http"
	usecaseBanner "retarget/internal/banner-service/usecase"
	authenticate "retarget/pkg/middleware/auth"

	"github.com/gorilla/mux"
)

func SetupRoutes(authenticator *authenticate.Authenticator, bannerUsecase *usecaseBanner.BannerUsecase,
	imageUsecase *usecaseBanner.BannerImageUsecase) *mux.Router {
	r := mux.NewRouter()

	bannerRoutes := handlerBanner.SetupBannerRoutes(authenticator, bannerUsecase, imageUsecase)
	r.PathPrefix("/api/v1/banner/").Handler(bannerRoutes)

	return r
}
