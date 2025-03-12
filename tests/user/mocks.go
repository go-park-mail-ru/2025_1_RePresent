package auth

import (
	"retarget/internal/entity"

	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUserByID(id int) (*entity.User, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(email string) (*entity.User, error) {
	args := m.Called(email)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByUsername(username string) (*entity.User, error) {
	args := m.Called(username)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) CreateNewUser(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

type MockAuthUsecase struct {
	mock.Mock
}

func (m *MockAuthUsecase) Login(email string, password string, role int) (*entity.User, error) {
	args := m.Called(email, password, role)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockAuthUsecase) Logout(sessionId string) error {
	args := m.Called(sessionId)
	return args.Error(0)
}

func (m *MockAuthUsecase) GetUserBySessionID(sessionId string) (*entity.User, error) {
	args := m.Called(sessionId)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockAuthUsecase) Register(username string, email string, password string, role int) (*entity.User, error) {
	args := m.Called(username, email, password, role)
	return args.Get(0).(*entity.User), args.Error(1)
}
