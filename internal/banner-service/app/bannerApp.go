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

	// sessionRepository := repoSession.NewSessionRepository(
	// 	cfg.AuthRedis.EndPoint,
	// 	cfg.AuthRedis.Password,
	// 	cfg.AuthRedis.Database,
	// 	30*time.Minute,
	// )
	// defer func() {
	// 	if err := sessionRepository.CloseConnection(); err != nil {
	// 		log.Printf("error closing session repository: %v", err)
	// 	}
	// }()

	// userRepository := repoAuth.NewAuthRepository(cfg.Database.Username, cfg.Database.Password, cfg.Database.Dbname, cfg.Database.Host, cfg.Database.Port, cfg.Database.Sslmode)
	bannerRepository := repo.NewBannerRepository(cfg.Database.ConnectionString())
	defer func() {
		if err := bannerRepository.CloseConnection(); err != nil {
			log.Println(err)
		}
	}()

	banner := usecase.NewBannerUsecase(bannerRepository)

	mux := controller.SetupRoutes(authenticator, banner)

	log.Fatal(http.ListenAndServe(":8024", middleware.CORS(mux)))
}
