package banner

import (
	"net/http"
	authenticate "pkg/middleware/auth"
	banner "retarget-bannerapp/usecase"

	"github.com/gorilla/mux"
)

type BannerController struct {
	BannerUsecase *banner.BannerUsecase
}

func NewBannerController(bannerUsecase *banner.BannerUsecase) *BannerController {
	return &BannerController{BannerUsecase: bannerUsecase}
}

func SetupBannerRoutes(authenticator *authenticate.Authenticator, bannerUsecase *banner.BannerUsecase) http.Handler {
	muxRouter := mux.NewRouter()
	bannerController := NewBannerController(bannerUsecase)
	// middleware.AuthMiddleware(authUsecase)()
	muxRouter.Handle("/api/v1/banner/", authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(bannerController.GetUserBanners)))
	muxRouter.Handle("/api/v1/banner/{banner_id}", authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(bannerController.BannerHandleFunc)))

	return muxRouter
}
