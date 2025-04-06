package http

import (
	// usecaseAuth "retarget/internal/auth-service/usecase/auth"
	authenticate "pkg/middleware/auth"
	handlerBanner "retarget-bannerapp/controller/http"
	usecaseBanner "retarget-bannerapp/usecase"

	"github.com/gorilla/mux"
)

func SetupRoutes(authenticator *authenticate.Authenticator, bannerUsecase *usecaseBanner.BannerUsecase,
	imageUsecase *usecaseBanner.BannerImageUsecase) *mux.Router {
	r := mux.NewRouter()

	bannerRoutes := handlerBanner.SetupBannerRoutes(authenticator, bannerUsecase, imageUsecase)
	r.PathPrefix("/api/v1/banner/").Handler(bannerRoutes)

	return r
}
