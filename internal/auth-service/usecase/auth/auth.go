package auth

import (
	"database/sql"
	"errors"
	entityAuth "retarget/internal/auth-service/entity/auth"
	repoAuth "retarget/internal/auth-service/repo/auth"

	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecaseInterface interface {
	Login(email string, password string, role int) (*entityAuth.User, error)
	Logout(sessionId string) error
	Register(username string, email string, password string, role int) (*entityAuth.User, error)

	GetUser(userId int) (*entityAuth.User, error)
	CheckCode(code int, userId int) error
	CreateCode(userId int) (int, error)

	AddSession(userId int, role int) (*entityAuth.Session, error)
}

type AuthUsecase struct {
	authRepository    *repoAuth.AuthRepository
	sessionRepository *repoAuth.SessionRepository
}

func NewAuthUsecase(userRepo *repoAuth.AuthRepository, sessionRepo *repoAuth.SessionRepository) *AuthUsecase {
	return &AuthUsecase{authRepository: userRepo, sessionRepository: sessionRepo}
}

func (a *AuthUsecase) Login(email string, password string, role int, requestID string) (*entityAuth.User, error) {
	user, err := a.authRepository.GetUserByEmail(email, requestID)
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
	err := a.sessionRepository.DelSession(sessionId)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthUsecase) GetUser(user_id int, requestID string) (*entityAuth.User, error) {
	user, err := a.authRepository.GetUserByID(user_id, requestID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (a *AuthUsecase) Register(username string, email string, password string, role int, requestID string) (*entityAuth.User, error) {
	user, err := a.authRepository.GetUserByEmail(email, requestID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if user != nil {
		return nil, errors.New("Пользователь с таким email уже существует")
	}

	user, err = a.authRepository.GetUserByUsername(username, requestID)
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
		Balance:     decimal.Zero,
		Role:        role,
	}

	err = entityAuth.ValidateUser(user)
	if err != nil {
		return nil, err
	}

	err = a.authRepository.CreateNewUser(user, requestID)
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

func (a *AuthUsecase) AddSession(userId int, role int) (*entityAuth.Session, error) {
	session, err := a.sessionRepository.AddSession(userId, role)
	if err != nil {
		return nil, err
	}
	return session, nil
}
