package authApp

import (
	"log"
	"net/http"
	authenticate "pkg/middleware/auth"
	configs "retarget-bannerapp/configs"
	controller "retarget-bannerapp/controller"
	middleware "retarget-bannerapp/controller/http/middleware"
	"retarget-bannerapp/repo"

	// repoSession "retarget-bannerapp/repo"
	usecase "retarget-bannerapp/usecase"
)

func Run(cfg *configs.Config) {
	authenticator, err := authenticate.NewAuthenticator(cfg.AuthRedis.EndPoint, cfg.AuthRedis.Password, cfg.AuthRedis.Database)
	if err != nil {
		log.Fatal(err.Error())
	}

	imageRepository := repo.NewBannerImageRepository("ReTargetMiniO:9000", "minioadmin", "minioadmin", "", false, "image")
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
