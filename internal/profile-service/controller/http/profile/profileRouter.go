package profile

import (
	"net/http"
	usecaseProfile "retarget/internal/profile-service/usecase/profile"

	authenticate "retarget/pkg/middleware/auth"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type ProfileController struct {
	profileUsecase *usecaseProfile.ProfileUsecase
	logger         *zap.SugaredLogger
}

func NewProfileController(profileUsecase *usecaseProfile.ProfileUsecase, logger *zap.SugaredLogger) *ProfileController {
	return &ProfileController{profileUsecase: profileUsecase, logger: logger}
}

func SetupProfileRoutes(authenticator *authenticate.Authenticator, profileUsecase *usecaseProfile.ProfileUsecase, logger *zap.SugaredLogger) http.Handler {
	muxRouter := mux.NewRouter()
	profileController := NewProfileController(profileUsecase, logger)

	muxRouter.Handle("/api/v1/profile/my", authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(profileController.GetProfileHandler)))
	muxRouter.Handle("/api/v1/profile/edit", authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(profileController.EditProfileHandler)))

	return muxRouter
}
