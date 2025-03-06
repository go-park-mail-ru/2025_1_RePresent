package auth

import (
	"RE/internal/entity"
	"RE/internal/repo"

	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	userRepository *repo.UserRepository
}

func NewAuthUsecase(userRepository *repo.UserRepository) *AuthUsecase {
	return &AuthUsecase{userRepository: userRepository}
}

func (a *AuthUsecase) Authenticate(username string, password string) (*entity.User, error) {
	user, err := a.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}
	return user, nil
}
