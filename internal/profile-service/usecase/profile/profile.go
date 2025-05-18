package profile

import (
	"errors"
	"fmt"
	entityProfile "retarget/internal/profile-service/entity/profile"
	repoProfile "retarget/internal/profile-service/repo/profile"
)

type ProfileUsecaseInterface interface {
	GetProfile(userID int, requestID string) (*entityProfile.ProfileResponse, error)
	PutProfile(userID int, username, description string, requestID string) error
}

type ProfileUsecase struct {
	profileRepository *repoProfile.ProfileRepository
}

func NewProfileUsecase(profileRepo *repoProfile.ProfileRepository) *ProfileUsecase {
	return &ProfileUsecase{profileRepository: profileRepo}
}

func (r *ProfileUsecase) PutProfile(userID int, username, description string, requestID string) error {
	err := r.profileRepository.UpdateProfileByID(userID, username, description, requestID)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProfileUsecase) GetProfile(userID int, requestID string) (*entityProfile.ProfileResponse, error) {
	profile, err := r.profileRepository.GetProfileByID(userID, requestID)
	if err != nil {
		if errors.Is(err, entityProfile.ErrProfileNotFound) {
			return nil, nil
		}
		return nil, err
	}
	fmt.Println("юзкейс перед сбором респонса")
	response := &entityProfile.ProfileResponse{
		ID:          profile.ID,
		Username:    profile.Username,
		Email:       profile.Email,
		Description: profile.Description,
		Balance:     *profile.Balance.Dec,
		Role:        profile.Role,
	}
	fmt.Println("юзкейс собрал респонс")
	// validationErrors, err := validator.ValidateStruct(profile)
	// if err != nil {
	// 	return nil, errors.New(validationErrors)
	// }
	return response, nil
}
