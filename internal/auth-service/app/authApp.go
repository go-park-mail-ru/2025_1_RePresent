package authApp

import (
	"log"
	"net/http"
	"retarget/configs"
	authAppHttp "retarget/internal/auth-service/controller/http"
	authMiddleware "retarget/internal/auth-service/controller/http/middleware"
	repoAuth "retarget/internal/auth-service/repo/auth"
	usecaseAuth "retarget/internal/auth-service/usecase/auth"
)

func Run(cfg *configs.Config) {
	userRepository := repoAuth.NewUserRepository(cfg.Database.Username, cfg.Database.Password, cfg.Database.Dbname, cfg.Database.Host, cfg.Database.Port, cfg.Database.Sslmode)
	defer func() {
		if err := userRepository.CloseConnection(); err != nil {
			log.Println(err)
		}
	}()

	authUsecase := usecaseAuth.NewAuthUsecase(userRepository)

	mux := authAppHttp.SetupRoutes(authUsecase)

	log.Fatal(http.ListenAndServe(":8080", authMiddleware.CORS(mux)))
}
