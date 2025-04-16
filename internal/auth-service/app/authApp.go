package authApp

import (
	"log"
	"net/http"
	configs "retarget/configs"
	authAppHttp "retarget/internal/auth-service/controller/http"
	authMiddleware "retarget/internal/auth-service/controller/http/middleware"
	repoAuth "retarget/internal/auth-service/repo/auth"
	repoSession "retarget/internal/auth-service/repo/auth"
	usecaseAuth "retarget/internal/auth-service/usecase/auth"
	authenticate "retarget/pkg/middleware/auth"
	"time"

	"go.uber.org/zap"
)

func Run(cfg *configs.Config, logger *zap.SugaredLogger) {
	authenticator, err := authenticate.NewAuthenticator(cfg.AuthRedis.EndPoint, cfg.AuthRedis.Password, cfg.AuthRedis.Database)
	if err != nil {
		log.Fatal(err.Error())
	}

	sessionRepository := repoSession.NewSessionRepository(
		cfg.AuthRedis.EndPoint,
		cfg.AuthRedis.Password,
		cfg.AuthRedis.Database,
		24*time.Hour,
	)
	defer func() {
		if err := sessionRepository.CloseConnection(); err != nil {
			log.Printf("error closing session repository: %v", err)
		}
	}()

	userRepository := repoAuth.NewAuthRepository(cfg.Database.ConnectionString(), logger)
	defer func() {
		if err := userRepository.CloseConnection(); err != nil {
			log.Println(err)
		}
	}()

	authUsecase := usecaseAuth.NewAuthUsecase(userRepository, sessionRepository)

	mux := authAppHttp.SetupRoutes(authenticator, authUsecase)

	log.Fatal(http.ListenAndServe(":8025", authMiddleware.CORS(mux)))
}
