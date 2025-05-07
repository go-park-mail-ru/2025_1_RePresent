package avatar

import (
	"net/http"
	usecaseAvatar "retarget/internal/profile-service/usecase/avatar"
	logger "retarget/pkg/middleware"
	authenticate "retarget/pkg/middleware/auth"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type AvatarController struct {
	avatarUsecase *usecaseAvatar.AvatarUsecase
	logger        *zap.SugaredLogger
}

func NewAvatarController(avatarUsecase *usecaseAvatar.AvatarUsecase, sugar *zap.SugaredLogger) *AvatarController {
	return &AvatarController{avatarUsecase: avatarUsecase, logger: sugar}
}

func SetupAvatarRoutes(authenticator *authenticate.Authenticator, avatarUsecase *usecaseAvatar.AvatarUsecase, sugar *zap.SugaredLogger) http.Handler {
	muxRouter := mux.NewRouter()
	avatarController := NewAvatarController(avatarUsecase, sugar)

	muxRouter.Handle("/api/v1/avatar/download", logger.LogMiddleware(authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(avatarController.DownloadAvatarHandler))))
	muxRouter.Handle("/api/v1/avatar/upload", logger.LogMiddleware(authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(avatarController.UploadAvatarHandler))))

	return muxRouter
}
