package controller

import (
	"net/http"
	banner "retarget/internal/banner-service/usecase"
	logger "retarget/pkg/middleware"
	authenticate "retarget/pkg/middleware/auth"

	"github.com/gorilla/mux"
)

type BannerController struct {
	BannerUsecase *banner.BannerUsecase
	ImageUsecase  *banner.BannerImageUsecase
	LinkBuilder   LinkBuilder
}

func NewBannerController(bannerUsecase *banner.BannerUsecase, imageUsecase *banner.BannerImageUsecase, linkBuilder LinkBuilder) *BannerController {
	return &BannerController{BannerUsecase: bannerUsecase, ImageUsecase: imageUsecase, LinkBuilder: linkBuilder}
}

func SetupBannerRoutes(authenticator *authenticate.Authenticator, bannerUsecase *banner.BannerUsecase, imageUsecase *banner.BannerImageUsecase) http.Handler {
	muxRouter := mux.NewRouter()
	linkBuilder := NewLinkBuilder(muxRouter)
	bannerController := NewBannerController(bannerUsecase, imageUsecase, linkBuilder)
	// middleware.AuthMiddleware(authUsecase)()
	muxRouter.Handle("/api/v1/banner/", logger.LogMiddleware(authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(bannerController.GetUserBanners)))).Methods("GET")
	// CRUD
	muxRouter.Handle("/api/v1/banner/create", logger.LogMiddleware(authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(bannerController.CreateBanner)))).Methods("POST")
	muxRouter.Handle("/api/v1/banner/{banner_id:[0-9]+}", logger.LogMiddleware(authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(bannerController.ReadBanner)))).Methods("GET")
	muxRouter.Handle("/api/v1/banner/{banner_id:[0-9]+}", logger.LogMiddleware(authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(bannerController.UpdateBanner)))).Methods("PUT")
	muxRouter.Handle("/api/v1/banner/{banner_id:[0-9]+}", logger.LogMiddleware(authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(bannerController.DeleteBanner)))).Methods("DELETE")
	//IFrame
	muxRouter.Handle("/api/v1/banner/iframe/{banner_id:[0-9]+}", logger.LogMiddleware(authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(bannerController.GetBannerIFrameByID)))).Methods("GET")
	// Работа с картинками
	muxRouter.Handle("/api/v1/banner/image/{image_id}", logger.LogMiddleware(http.HandlerFunc(bannerController.DownloadImage))).Methods("GET").Name("download_image")
	muxRouter.Handle("/api/v1/banner/upload", logger.LogMiddleware(authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(bannerController.UploadImageHandler)))).Methods("PUT")
	// Рандомный айфрейм юзера
	muxRouter.Handle("/api/v1/banner/uniq_link/{uniq_link}", logger.LogMiddleware(http.HandlerFunc(bannerController.RandomIFrame))).Methods("GET")

	// Маршрут для генерации описания
	muxRouter.Handle("/api/v1/banner/generate/description/{banner_id:[0-9]+}",
		logger.LogMiddleware(authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(bannerController.GenerateDescription)))).
		Methods("POST")

	return muxRouter
}
