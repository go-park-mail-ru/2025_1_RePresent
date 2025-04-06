package controller

import (
	"net/http"
	authenticate "pkg/middleware/auth"
	banner "retarget-bannerapp/usecase"

	"github.com/gorilla/mux"
)

type BannerController struct {
	BannerUsecase *banner.BannerUsecase
	ImageUsecase  *banner.BannerImageUsecase
}

func NewBannerController(bannerUsecase *banner.BannerUsecase, imageUsecase *banner.BannerImageUsecase) *BannerController {
	return &BannerController{BannerUsecase: bannerUsecase, ImageUsecase: imageUsecase}
}

func SetupBannerRoutes(authenticator *authenticate.Authenticator, bannerUsecase *banner.BannerUsecase, imageUsecase *banner.BannerImageUsecase) http.Handler {
	muxRouter := mux.NewRouter()
	bannerController := NewBannerController(bannerUsecase, imageUsecase)
	// middleware.AuthMiddleware(authUsecase)()
	muxRouter.Handle("/api/v1/banner/", authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(bannerController.GetUserBanners)))
	muxRouter.Handle("/api/v1/banner/create", authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(bannerController.CreateBanner)))
	muxRouter.Handle("/api/v1/banner/upload", (http.HandlerFunc(bannerController.UploadImageHandler)))
	muxRouter.Handle("/api/v1/banner/image/{image_id}", (http.HandlerFunc(bannerController.DownloadImage)))
	muxRouter.Handle("/api/v1/banner/{banner_id}", authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(bannerController.BannerHandleFunc)))

	return muxRouter
}
