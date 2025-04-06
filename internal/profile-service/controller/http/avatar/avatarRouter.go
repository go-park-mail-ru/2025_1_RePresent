package avatar

import (
	"net/http"
	usecaseAvatar "retarget/internal/profile-service/usecase/avatar"

	authenticate "retarget/pkg/middleware/auth"

	"github.com/gorilla/mux"
)

type AvatarController struct {
	avatarUsecase *usecaseAvatar.AvatarUsecase
}

func NewAvatarController(avatarUsecase *usecaseAvatar.AvatarUsecase) *AvatarController {
	return &AvatarController{avatarUsecase: avatarUsecase}
}

func SetupAvatarRoutes(authenticator *authenticate.Authenticator, avatarUsecase *usecaseAvatar.AvatarUsecase) http.Handler {
	muxRouter := mux.NewRouter()
	avatarController := NewAvatarController(avatarUsecase)

	muxRouter.Handle("/api/v1/avatar/download", authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(avatarController.DownloadAvatarHandler)))
	muxRouter.Handle("/api/v1/avatar/upload", authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(avatarController.UploadAvatarHandler)))

	return muxRouter
}
