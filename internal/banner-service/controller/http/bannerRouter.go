package banner

import (
	"net/http"
	banner "retarget/internal/banner-service/usecase"

	"github.com/gorilla/mux"
)

type BannerController struct {
	BannerUsecase *banner.BannerUsecase
}

func NewBannerController(bannerUsecase *banner.BannerUsecase) *BannerController {
	return &BannerController{BannerUsecase: bannerUsecase}
}

func SetupBannerRoutes(bannerUsecase *banner.BannerUsecase) http.Handler {
	muxRouter := mux.NewRouter()
	bannerController := NewBannerController(bannerUsecase)
	// middleware.AuthMiddleware(authUsecase)()
	muxRouter.Handle("/api/v1/banner/", http.HandlerFunc(bannerController.GetUserBanners))
	muxRouter.Handle("/api/v1/banner/{banner_id}", http.HandlerFunc(bannerController.BannerHandleFunc))

	return muxRouter
}
