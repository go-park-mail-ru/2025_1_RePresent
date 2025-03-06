package banner

import (
	"net/http"
	"retarget/internal/controller/http/middleware"
	"retarget/internal/usecase/auth"
	"retarget/internal/usecase/banner"

	"github.com/gorilla/mux"
)

type BannerController struct {
	BannerUsecase *banner.BannerUsecase
	AuthUsecase   *auth.AuthUsecase
}

func NewBannerController(authUsecase *auth.AuthUsecase, bannerUsecase *banner.BannerUsecase) *BannerController {
	return &BannerController{AuthUsecase: authUsecase, BannerUsecase: bannerUsecase}
}

func SetupBannerRoutes(authUsecase *auth.AuthUsecase, bannerUsecase *banner.BannerUsecase) http.Handler {
	muxRouter := mux.NewRouter()
	bannerController := NewBannerController(authUsecase, bannerUsecase)

	muxRouter.Handle("/api/v1/banner/user/{user_id}/all", middleware.AuthMiddleware(authUsecase)(http.HandlerFunc(bannerController.GetBannersByUserCookie)))

	return muxRouter
}
