package auth

import (
	"RE/internal/entity"
	"RE/internal/repo"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	userRepository *repo.UserRepository
}

func NewAuthUsecase(userRepository *repo.UserRepository) *AuthUsecase {
	return &AuthUsecase{userRepository: userRepository}
}

func (a *AuthUsecase) Login(email string, password string, role int) (*entity.User, error) {
	user, err := a.userRepository.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if user.Role != role {
		return nil, errors.New("SQL: Неверные данные пользователя")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (a *AuthUsecase) Logout(sessionId string) error {
	err := DelSession(sessionId)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthUsecase) GetUserBySessionID(sessionId string) (*entity.User, error) {
	session, err := GetSession(sessionId)
	if err != nil {
		return nil, err
	}
	user, err := a.userRepository.GetUserByID(session.UserID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("SQL: пользователь с таким id не существует")
	}
	return user, nil
}

func (a *AuthUsecase) Register(username string, email string, password string, role int) (*entity.User, error) {
	user, err := a.userRepository.GetUserByEmail(email)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if user != nil {
		return nil, errors.New("SQL: пользователь с таким email уже существует")
	}

	user, err = a.userRepository.GetUserByUsername(username)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if user != nil {
		return nil, errors.New("SQL: пользователь с таким username уже существует")
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	user = &entity.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
		Avatar:   "",
		Balance:  0,
		Role:     role,
	}

	err = entity.ValidateUser(user)
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
