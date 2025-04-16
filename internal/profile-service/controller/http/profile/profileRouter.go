package profile

import (
	"net/http"
	usecaseProfile "retarget/internal/profile-service/usecase/profile"
	logger "retarget/pkg/middleware"
	authenticate "retarget/pkg/middleware/auth"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type ProfileController struct {
	profileUsecase *usecaseProfile.ProfileUsecase
	logger         *zap.SugaredLogger
}

func NewProfileController(profileUsecase *usecaseProfile.ProfileUsecase, sugar *zap.SugaredLogger) *ProfileController {
	return &ProfileController{profileUsecase: profileUsecase, logger: sugar}
}

func SetupProfileRoutes(authenticator *authenticate.Authenticator, profileUsecase *usecaseProfile.ProfileUsecase, sugar *zap.SugaredLogger) http.Handler {
	muxRouter := mux.NewRouter()
	profileController := NewProfileController(profileUsecase, sugar)

	muxRouter.Handle("/api/v1/profile/my", logger.LogMiddleware(authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(profileController.GetProfileHandler))))
	muxRouter.Handle("/api/v1/profile/edit", logger.LogMiddleware(authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(profileController.EditProfileHandler))))

	return muxRouter
}
