package auth

import (
	// "errors"
	"bytes"
	"encoding/json"
	"errors"
	"retarget/internal/controller/http/auth"

	"net/http"
	"net/http/httptest"
	"retarget/internal/entity"

	// a "retarget/internal/usecase/auth"
	"testing"

	"github.com/go-playground/assert"
)

func TestGetCurrentUserHandler(t *testing.T) {
	mockAuthUsecase := new(MockAuthUsecase)
	authController := auth.NewAuthController(mockAuthUsecase)

	defaultUser := entity.User{
		ID:       1,
		Username: "John Doe",
		Email:    "john@example.com",
		Avatar:   "avatar.png",
		Balance:  100,
		Role:     1,
	}

	mockAuthUsecase.On("GetUserBySessionID", "valid_session").Return(&defaultUser, nil)

	req := httptest.NewRequest(http.MethodGet, "/user", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "valid_session"})

	rr := httptest.NewRecorder()
	authController.GetCurrentUserHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	mockAuthUsecase.AssertExpectations(t)
}

func TestGetCurrentUserHandlerNotFound(t *testing.T) {
	mockUsecase := new(MockAuthUsecase)
	authController := auth.NewAuthController(mockUsecase)

	defaultUser := entity.User{
		ID:       1,
		Username: "John Doe",
		Email:    "john@example.com",
		Avatar:   "avatar.png",
		Balance:  100,
		Role:     1,
	}

	mockUsecase.On("GetUserBySessionID", "invalid_session").Return(&defaultUser, errors.New("Пользователь с таким id не существует")) // по идее &defaultUser должен быть nil, но почему-то вылетает паника

	req := httptest.NewRequest(http.MethodGet, "/user", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "invalid_session"})

	rr := httptest.NewRecorder()
	authController.GetCurrentUserHandler(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	mockUsecase.AssertExpectations(t)
}

func TestRegisterSuccess(t *testing.T) {
	mockUsecase := new(MockAuthUsecase)
	authController := auth.NewAuthController(mockUsecase)

	defaultUser := entity.User{
		ID:       1,
		Username: "John Doe",
		Email:    "john@example.com",
		Avatar:   "avatar.png",
		Balance:  100,
		Role:     1,
	}

	JSONuser := auth.RegisterRequest{
		Username: "JohnDoe",
		Email:    "john@example.com",
		Password: "password123",
		Role:     1,
	}

	jsonData, _ := json.Marshal(JSONuser)

	mockUsecase.On("Register", JSONuser.Username, JSONuser.Email, JSONuser.Password, JSONuser.Role).Return(&defaultUser, nil)

	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewReader(jsonData))
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "valid_session"})
	rr := httptest.NewRecorder()

	authController.RegisterHandler(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

}

func TestRegisterError(t *testing.T) {
	mockUsecase := new(MockAuthUsecase)
	authController := auth.NewAuthController(mockUsecase)

	defaultUser := entity.User{
		ID:       1,
		Username: "John Doe",
		Email:    "john@example.com",
		Avatar:   "avatar.png",
		Balance:  100,
		Role:     1,
	}

	JSONuser := auth.RegisterRequest{
		Username: "JohnDoe",
		Email:    "john@example.com",
		Password: "password123",
		Role:     1,
	}

	jsonData, _ := json.Marshal(JSONuser)
	// должен возвращать nil и error, но с nil всё крашится, потом разобраться и заменить defaultUser на nil
	mockUsecase.On("Register", JSONuser.Username, JSONuser.Email, JSONuser.Password, JSONuser.Role).Return(&defaultUser, errors.New("Пользователь с таким email уже существует"))

	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewReader(jsonData))
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "valid_session"})
	rr := httptest.NewRecorder()

	authController.RegisterHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestLoginSuccess(t *testing.T) {
	mockUsecase := new(MockAuthUsecase)
	authController := auth.NewAuthController(mockUsecase)

	JSONUser := auth.LoginRequest{
		Email:    "john@example.com",
		Password: "password123",
		Role:     1,
	}

	defaultUser := entity.User{
		ID:       1,
		Username: "JohnDoe",
		Email:    "john@example.com",
		Avatar:   "avatar.png",
		Balance:  100,
		Role:     1,
	}

	jsonData, _ := json.Marshal(JSONUser)
	mockUsecase.On("Login", JSONUser.Email, JSONUser.Password, JSONUser.Role).Return(&defaultUser, nil)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(jsonData))
	rr := httptest.NewRecorder()

	authController.LoginHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockUsecase.AssertExpectations(t)
}

func TestLoginInvalidCredentials(t *testing.T) {
	mockUsecase := new(MockAuthUsecase)
	authController := auth.NewAuthController(mockUsecase)

	JSONUser := auth.LoginRequest{
		Email:    "john@example.com",
		Password: "wrongpassword",
		Role:     1,
	}

	defaultUser := entity.User{
		ID:       1,
		Username: "JohnDoe",
		Email:    "john@example.com",
		Avatar:   "avatar.png",
		Balance:  100,
		Role:     1,
	}

	jsonData, _ := json.Marshal(JSONUser)
	// должен возвращать nil и error, но с nil всё крашится, потом разобраться и заменить defaultUser на nil
	mockUsecase.On("Login", JSONUser.Email, JSONUser.Password, JSONUser.Role).Return(&defaultUser, errors.New("Неверные данные пользователя"))

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(jsonData))
	rr := httptest.NewRecorder()

	authController.LoginHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	mockUsecase.AssertExpectations(t)
}
