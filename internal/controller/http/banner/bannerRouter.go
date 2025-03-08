package banner

import (
	"RE/internal/controller/http/middleware"
	"RE/internal/usecase/auth"
	"RE/internal/usecase/banner"
	"net/http"

	"github.com/gorilla/mux"
)

type BannerController struct {
	BannerUsecase *banner.BannerUsecase
}

func NewBannerController(bannerUsecase *banner.BannerUsecase) *BannerController {
	return &BannerController{BannerUsecase: bannerUsecase}
}

func SetupBannerRoutes(usecase *auth.AuthUsecase, bannerUsecase *banner.BannerUsecase) http.Handler {
	muxRouter := mux.NewRouter()
	bannerController := NewBannerController(bannerUsecase)

	muxRouter.Handle("/banner/user/{user_id}/all", middleware.AuthMiddleware(usecase)(http.HandlerFunc(bannerController.GetBannersByUserId)))

	return muxRouter
}
