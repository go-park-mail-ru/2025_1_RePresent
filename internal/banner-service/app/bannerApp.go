package authApp

import (
	"log"
	"net/http"
	configs "retarget/configs"
	controller "retarget/internal/banner-service/controller"
	middleware "retarget/internal/banner-service/controller/http/middleware"
	"retarget/internal/banner-service/repo"
	authenticate "retarget/pkg/middleware/auth"

	server "retarget/internal/banner-service/grpc"
	usecase "retarget/internal/banner-service/usecase"

	"go.uber.org/zap"
)

func Run(cfg *configs.Config, logger *zap.SugaredLogger) {
	authenticator, err := authenticate.NewAuthenticator(cfg.AuthRedis.EndPoint, cfg.AuthRedis.Password, cfg.AuthRedis.Database)
	if err != nil {
		log.Fatal(err.Error())
	}

	imageRepository := repo.NewBannerImageRepository(cfg.Minio.EndPoint, cfg.Minio.AccessKeyID, cfg.Minio.SecretAccesKey, cfg.Minio.Token, false, "image")
	bannerRepository := repo.NewBannerRepository(cfg.Database.ConnectionString(), logger)
	defer func() {
		if err := bannerRepository.CloseConnection(); err != nil {
			log.Println(err)
		}
	}()

	image := usecase.NewBannerImageUsecase(imageRepository)
	banner := usecase.NewBannerUsecase(bannerRepository)

	mux := controller.SetupRoutes(authenticator, banner, image)

	log.Fatal(http.ListenAndServe(":8024", middleware.CORS(mux)))
	server.RunGRPCServer(*banner)
}
