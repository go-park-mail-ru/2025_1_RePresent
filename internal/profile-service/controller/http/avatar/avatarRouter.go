package avatar

import (
	"net/http"
	usecaseAvatar "retarget/internal/profile-service/usecase/avatar"

	authenticate "retarget/pkg/middleware/auth"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type AvatarController struct {
	avatarUsecase *usecaseAvatar.AvatarUsecase
	logger        *zap.SugaredLogger
}

func NewAvatarController(avatarUsecase *usecaseAvatar.AvatarUsecase, logger *zap.SugaredLogger) *AvatarController {
	return &AvatarController{avatarUsecase: avatarUsecase, logger: logger}
}

func SetupAvatarRoutes(authenticator *authenticate.Authenticator, avatarUsecase *usecaseAvatar.AvatarUsecase, logger *zap.SugaredLogger) http.Handler {
	muxRouter := mux.NewRouter()
	avatarController := NewAvatarController(avatarUsecase, logger)

	muxRouter.Handle("/api/v1/avatar/download", authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(avatarController.DownloadAvatarHandler)))
	muxRouter.Handle("/api/v1/avatar/upload", authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(avatarController.UploadAvatarHandler)))

	return muxRouter
}
