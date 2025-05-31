package auth

import (
	// "errors"

	"retarget/internal/usecase/auth"

	"retarget/internal/entity"

	// a "retarget/internal/usecase/auth"
	"testing"

	"github.com/go-playground/assert"
)

func TestAuthUsecaseSessionNotFound(t *testing.T) {
	mockUserRepository := new(MockUserRepository)
	mockAuthUsecase := auth.NewAuthUsecase(mockUserRepository)

	defaultUser := entity.User{
		ID:       1,
		Username: "John Doe",
		Email:    "john@example.com",
		Avatar:   "avatar.png",
		Balance:  100,
		Role:     1,
	}

	mockUserRepository.On("GetUserByID", "1").Return(&defaultUser, nil)
	user, err := mockAuthUsecase.GetUserBySessionID("valid_session")
	assert.Equal(t, user, nil)
	assert.Equal(t, err, entity.ErrSessionNotFound) // Не нашел способа адекватно мокнуть GetSession, поэтому пока так
}
