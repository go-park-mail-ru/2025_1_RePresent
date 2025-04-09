package authApp

import (
	"log"
	"net/http"
	configs "retarget/configs"
	controller "retarget/internal/banner-service/controller"
	middleware "retarget/internal/banner-service/controller/http/middleware"
	"retarget/internal/banner-service/repo"
	authenticate "retarget/pkg/middleware/auth"

	// repoSession "retarget/internal/banner-service/repo"
	usecase "retarget/internal/banner-service/usecase"
)

func Run(cfg *configs.Config) {
	authenticator, err := authenticate.NewAuthenticator(cfg.AuthRedis.EndPoint, cfg.AuthRedis.Password, cfg.AuthRedis.Database)
	if err != nil {
		log.Fatal(err.Error())
	}

	imageRepository := repo.NewBannerImageRepository(cfg.Minio.EndPoint, cfg.Minio.AccessKeyID, cfg.Minio.SecretAccesKey, cfg.Minio.Token, false, "image")
	bannerRepository := repo.NewBannerRepository(cfg.Database.ConnectionString())
	defer func() {
		if err := bannerRepository.CloseConnection(); err != nil {
			log.Println(err)
		}
	}()

	image := usecase.NewBannerImageUsecase(imageRepository)
	banner := usecase.NewBannerUsecase(bannerRepository)

	mux := controller.SetupRoutes(authenticator, banner, image)

	log.Fatal(http.ListenAndServe(":8024", middleware.CORS(mux)))
}
