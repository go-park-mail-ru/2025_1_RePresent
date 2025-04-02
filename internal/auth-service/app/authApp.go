package authApp

import (
	"log"
	"net/http"
	"retarget/configs"
	authAppHttp "retarget/internal/auth-service/controller/http"
	authMiddleware "retarget/internal/auth-service/controller/http/middleware"
	repoAuth "retarget/internal/auth-service/repo/auth"
	repoSession "retarget/internal/auth-service/repo/auth"
	usecaseAuth "retarget/internal/auth-service/usecase/auth"
	authenticate "retarget/pkg/middleware/auth"
	"time"
)

func Run(cfg *configs.Config) {
	authenticator, err := authenticate.NewAuthenticator(cfg.AuthRedis.EndPoint, cfg.AuthRedis.Password, cfg.AuthRedis.Database)
	if err != nil {
		log.Fatal(err.Error())
	}

	sessionRepository := repoSession.NewSessionRepository(
		cfg.AuthRedis.EndPoint,
		cfg.AuthRedis.Password,
		cfg.AuthRedis.Database,
		30*time.Minute,
	)
	defer func() {
		if err := sessionRepository.CloseConnection(); err != nil {
			log.Printf("error closing session repository: %v", err)
		}
	}()

	// userRepository := repoAuth.NewAuthRepository(cfg.Database.Username, cfg.Database.Password, cfg.Database.Dbname, cfg.Database.Host, cfg.Database.Port, cfg.Database.Sslmode)
	userRepository := repoAuth.NewAuthRepository(cfg.Database.ConnectionString())
	defer func() {
		if err := userRepository.CloseConnection(); err != nil {
			log.Println(err)
		}
	}()

	authUsecase := usecaseAuth.NewAuthUsecase(userRepository, sessionRepository)

	mux := authAppHttp.SetupRoutes(authenticator, authUsecase)

	log.Fatal(http.ListenAndServe(":8020", authMiddleware.CORS(mux)))
}
