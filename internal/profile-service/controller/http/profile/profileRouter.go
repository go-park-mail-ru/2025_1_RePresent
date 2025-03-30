package profile

import (
	"net/http"
	usecaseProfile "retarget/internal/profile-service/usecase/profile"

	authenticate "retarget/pkg/middleware/auth"

	"github.com/gorilla/mux"
)

type ProfileController struct {
	profileUsecase *usecaseProfile.ProfileUsecase
}

func NewProfileController(profileUsecase *usecaseProfile.ProfileUsecase) *ProfileController {
	return &ProfileController{profileUsecase: profileUsecase}
}

func SetupProfileRoutes(authenticator *authenticate.Authenticator, profileUsecase *usecaseProfile.ProfileUsecase) http.Handler {
	muxRouter := mux.NewRouter()
	profileController := NewProfileController(profileUsecase)

	muxRouter.Handle("/api/v1/profile/my", authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(profileController.GetProfileHandler)))
	muxRouter.Handle("/api/v1/profile/edit", authenticate.AuthMiddleware(authenticator)(http.HandlerFunc(profileController.EditProfileHandler)))

	return muxRouter
}
