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
	Register(username string, email string, password string, role int) (*entityAuth.User, error)

	GetUser(userId int) (*entityAuth.User, error)
	CheckCode(code int, userId int) error
	CreateCode(userId int) (int, error)
}

type AuthUsecase struct {
	authRepository    *repoAuth.AuthRepository
	sessionRepository *repoAuth.SessionRepository
}

func NewAuthUsecase(userRepo *repoAuth.AuthRepository, sessionRepo *repoAuth.SessionRepository) *AuthUsecase {
	return &AuthUsecase{authRepository: userRepo, sessionRepository: sessionRepo}
}

func (a *AuthUsecase) Login(email string, password string, role int) (*entityAuth.User, error) {
	user, err := a.authRepository.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if user.Role != role {
		return nil, errors.New("Incorrect user data")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (a *AuthUsecase) Logout(sessionId string) error {
	// err := // Удалить сессию
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthUsecase) GetUser(user_id int) (*entityAuth.User, error) {
	return &entityAuth.User{}, nil
}

func (a *AuthUsecase) Register(username string, email string, password string, role int) (*entityAuth.User, error) {
	user, err := a.authRepository.GetUserByEmail(email)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if user != nil {
		return nil, errors.New("Пользователь с таким email уже существует")
	}

	user, err = a.authRepository.GetUserByUsername(username)
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
		Username:    username,
		Email:       email,
		Password:    hashedPassword,
		Description: "",
		Balance:     0,
		Role:        role,
	}

	err = entityAuth.ValidateUser(user)
	if err != nil {
		return nil, err
	}

	err = a.authRepository.CreateNewUser(user)
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
