package profileApp

import (
	"log"
	"net/http"
	"retarget/configs"
	profileAppHttp "retarget/internal/profile-service/controller/http"
	profileMiddleware "retarget/internal/profile-service/controller/http/middleware"

	repoAvatar "retarget/internal/profile-service/repo/avatar"
	repoProfile "retarget/internal/profile-service/repo/profile"

	usecaseAvatar "retarget/internal/profile-service/usecase/avatar"
	usecaseProfile "retarget/internal/profile-service/usecase/profile"

	authenticate "retarget/pkg/middleware/auth"

	"go.uber.org/zap"
)

func Run(cfg *configs.Config, logger *zap.SugaredLogger) {
	authenticator, err := authenticate.NewAuthenticator(cfg.AuthRedis.EndPoint, cfg.AuthRedis.Password, cfg.AuthRedis.Database)
	if err != nil {
		log.Fatal(err.Error())
	}
	profileRepository := repoProfile.NewProfileRepository(cfg.Database.ConnectionString(), logger)
	defer func() {
		if err := profileRepository.CloseConnection(); err != nil {
			log.Println(err)
		}
	}()
	avatarRepository := repoAvatar.NewAvatarRepository(cfg.Minio.EndPoint, cfg.Minio.AccessKeyID, cfg.Minio.SecretAccesKey, cfg.Minio.Token, false, "avatar")
	profileUsecase := usecaseProfile.NewProfileUsecase(profileRepository)
	avatarUsecase := usecaseAvatar.NewAvatarUsecase(avatarRepository)

	mux := profileAppHttp.SetupRoutes(authenticator, profileUsecase, avatarUsecase, logger)

	log.Fatal(http.ListenAndServe(":8023", profileMiddleware.CORS(mux)))
}
