package profile

import (
	"errors"
	entityProfile "retarget/internal/profile-service/entity/profile"
	repoProfile "retarget/internal/profile-service/repo/profile"
)

type ProfileUsecaseInterface interface {
	GetProfile(userID int) (*entityProfile.ProfileResponse, error)
	PutProfile(userID int, username, description string) error
}

type ProfileUsecase struct {
	profileRepository *repoProfile.ProfileRepository
}

func NewProfileUsecase(profileRepo *repoProfile.ProfileRepository) *ProfileUsecase {
	return &ProfileUsecase{profileRepository: profileRepo}
}

func (r *ProfileUsecase) PutProfile(userID int, username, description string) error {
	err := r.profileRepository.UpdateProfileByID(userID, username, description)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProfileUsecase) GetProfile(userID int) (*entityProfile.ProfileResponse, error) {
	profile, err := r.profileRepository.GetProfileByID(userID)
	if err != nil {
		if errors.Is(err, entityProfile.ErrProfileNotFound) {
			return nil, nil
		}
		return nil, err
	}
	response := &entityProfile.ProfileResponse{
		Username:    profile.Username,
		Email:       profile.Email,
		Description: profile.Description,
		Balance:     profile.Balance,
		Role:        profile.Role,
	}
	// validationErrors, err := validator.ValidateStruct(response)
	//if err != nil {
	//		return nil, errors.New(validationErrors)
	//}
	return response, nil
}
