package auth

import (
	"database/sql"
	"errors"
	entityAuth "retarget/internal/auth-service/entity/auth"
	repoAuth "retarget/internal/auth-service/repo/auth"

	"golang.org/x/crypto/bcrypt"
)

type AuthUsecaseInterface interface {
	Login(email string, password string, role int) (*entityAuth.User, error)
	Logout(sessionId string) error
	GetUserBySessionID(sessionId string) (*entityAuth.User, error)
	Register(username string, email string, password string, role int) (*entityAuth.User, error)
}

type AuthUsecase struct {
	userRepository *repoAuth.UserRepository
}

func NewAuthUsecase(userRepo *repoAuth.UserRepository) *AuthUsecase {
	return &AuthUsecase{userRepository: userRepo}
}

func (a *AuthUsecase) Login(email string, password string, role int) (*entityAuth.User, error) {
	user, err := a.userRepository.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if user.Role != role {
		return nil, errors.New("Неверные данные пользователя")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (a *AuthUsecase) Logout(sessionId string) error {
	err := entityAuth.DelSession(sessionId)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthUsecase) GetUserBySessionID(sessionId string) (*entityAuth.User, error) {
	session, err := entityAuth.GetSession(sessionId)
	if err != nil {
		return nil, err
	}
	user, err := a.userRepository.GetUserByID(session.UserID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("Пользователь с таким id не существует")
	}
	return user, nil
}

func (a *AuthUsecase) Register(username string, email string, password string, role int) (*entityAuth.User, error) {
	user, err := a.userRepository.GetUserByEmail(email)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if user != nil {
		return nil, errors.New("Пользователь с таким email уже существует")
	}

	user, err = a.userRepository.GetUserByUsername(username)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if user != nil {
		return nil, errors.New("Пользователь с таким username уже существует")
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	user = &entityAuth.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
		Avatar:   "",
		Balance:  0,
		Role:     role,
	}

	err = entityAuth.ValidateUser(user)
	if err != nil {
		return nil, err
	}

	err = a.userRepository.CreateNewUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func hashPassword(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hash, nil
}
