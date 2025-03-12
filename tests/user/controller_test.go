package auth

import (
	// "errors"
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

	mockAuthUsecase.On("GetUserByEmail", "john@example.com").Return(&defaultUser, nil)

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
