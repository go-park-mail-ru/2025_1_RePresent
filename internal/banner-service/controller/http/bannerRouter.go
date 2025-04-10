package controller

import (
	"net/http"
	banner "retarget/internal/banner-service/usecase"
	authenticate "retarget/pkg/middleware/auth"

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
	muxRouter.Handle("/api/v1/banner/", authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(bannerController.GetUserBanners))).Methods("GET")
	// CRUD
	muxRouter.Handle("/api/v1/banner/create", authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(bannerController.CreateBanner))).Methods("POST")
	muxRouter.Handle("/api/v1/banner/{banner_id:[0-9]+}", authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(bannerController.ReadBanner))).Methods("GET")
	muxRouter.Handle("/api/v1/banner/{banner_id:[0-9]+}", authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(bannerController.UpdateBanner))).Methods("PUT")
	muxRouter.Handle("/api/v1/banner/{banner_id:[0-9]+}", authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(bannerController.DeleteBanner))).Methods("DELETE")
	//IFrame
	muxRouter.Handle("/api/v1/banner/iframe/{banner_id:[0-9]+}", (http.HandlerFunc(bannerController.GetBannerIFrame))).Methods("GET")
	// Работа с картинками
	muxRouter.Handle("/api/v1/banner/image/{image_id}", (http.HandlerFunc(bannerController.DownloadImage))).Methods("GET")
	muxRouter.Handle("/api/v1/banner/upload", authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(bannerController.UploadImageHandler))).Methods("POST")
	return muxRouter
}
