package http

import (
	handlerAvatar "retarget/internal/profile-service/controller/http/avatar"
	handlerProfile "retarget/internal/profile-service/controller/http/profile"
	usecaseAvatar "retarget/internal/profile-service/usecase/avatar"
	usecaseProfile "retarget/internal/profile-service/usecase/profile"
	middleware "retarget/pkg/middleware"
	authenticate "retarget/pkg/middleware/auth"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func SetupRoutes(authenticator *authenticate.Authenticator, profileUsecase *usecaseProfile.ProfileUsecase, avatarUsecase *usecaseAvatar.AvatarUsecase, logger *zap.SugaredLogger) *mux.Router {
	r := mux.NewRouter()

	r.Use(middleware.LogMiddleware)

	profileRoutes := handlerProfile.SetupProfileRoutes(authenticator, profileUsecase, logger)
	r.PathPrefix("/api/v1/profile/").Handler(profileRoutes)

	avatarRoutes := handlerAvatar.SetupAvatarRoutes(authenticator, avatarUsecase, logger)
	r.PathPrefix("/api/v1/avatar/").Handler(avatarRoutes)

	return r
}
