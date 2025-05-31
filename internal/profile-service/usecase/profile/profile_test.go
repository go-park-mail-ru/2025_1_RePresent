package profile

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"gopkg.in/inf.v0"
	"testing"

	entityProfile "retarget/internal/profile-service/entity/profile"
	repoProfile "retarget/internal/profile-service/repo/profile"
)

type MockProfileRepository struct {
	mock.Mock
	logger *zap.SugaredLogger
}

func NewMockProfileRepository() *MockProfileRepository {
	logger, _ := zap.NewDevelopment()
	return &MockProfileRepository{
		logger: logger.Sugar(),
	}
}

func (m *MockProfileRepository) GetProfileByID(userID int, requestID string) (*entityProfile.Profile, error) {
	if m.logger != nil {
		m.logger.Debugw("Mock GetProfileByID called", "request_id", requestID, "userID", userID)
	}

	args := m.Called(userID, requestID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entityProfile.Profile), args.Error(1)
}

func (m *MockProfileRepository) UpdateProfileByID(userID int, username, description string, requestID string) error {
	if m.logger != nil {
		m.logger.Debugw("Mock UpdateProfileByID called", "request_id", requestID, "userID", userID)
	}

	args := m.Called(userID, username, description, requestID)
	return args.Error(0)
}

func (m *MockProfileRepository) CloseConnection() error {
	args := m.Called()
	return args.Error(0)
}

type ProfileUsecaseForTest struct {
	mockRepo *MockProfileRepository
}

func (p *ProfileUsecaseForTest) GetProfile(userID int, requestID string) (*entityProfile.ProfileResponse, error) {
	profile, err := p.mockRepo.GetProfileByID(userID, requestID)
	if err != nil {
		if errors.Is(err, entityProfile.ErrProfileNotFound) {
			return nil, nil
		}
		return nil, err
	}

	response := &entityProfile.ProfileResponse{
		ID:          profile.ID,
		Username:    profile.Username,
		Email:       profile.Email,
		Description: profile.Description,
		Balance:     *profile.Balance.Dec,
		Role:        profile.Role,
	}

	return response, nil
}

func (p *ProfileUsecaseForTest) PutProfile(userID int, username, description string, requestID string) error {
	return p.mockRepo.UpdateProfileByID(userID, username, description, requestID)
}

func NewProfileUsecaseForTest() *ProfileUsecaseForTest {
	mockRepo := NewMockProfileRepository()
	return &ProfileUsecaseForTest{
		mockRepo: mockRepo,
	}
}

func decimalToInfDec(d decimal.Decimal) *inf.Dec {
	infDec := new(inf.Dec)

	s := d.String()
	infDec, _ = new(inf.Dec).SetString(s)

	return infDec
}

func TestNewProfileUsecase(t *testing.T) {
	repo := &repoProfile.ProfileRepository{}
	usecase := NewProfileUsecase(repo)
	assert.NotNil(t, usecase)
}

func TestGetProfile(t *testing.T) {
	testUsecase := NewProfileUsecaseForTest()
	mockRepo := testUsecase.mockRepo

	userID := 1
	requestID := "test-request-id"

	balanceDecimal := decimal.NewFromFloat(100.0)
	infDec := decimalToInfDec(balanceDecimal)

	mockProfile := &entityProfile.Profile{
		ID:          userID,
		Username:    "testuser",
		Email:       "test@example.com",
		Description: "Test description",
		Balance: entityProfile.Decimal{
			Dec: infDec,
		},
		Role: 1,
	}

	mockRepo.On("GetProfileByID", userID, requestID).Return(mockProfile, nil).Once()

	profile, err := testUsecase.GetProfile(userID, requestID)
	assert.NoError(t, err)
	assert.NotNil(t, profile)
	assert.Equal(t, userID, profile.ID)
	assert.Equal(t, "testuser", profile.Username)
	assert.Equal(t, "test@example.com", profile.Email)
	assert.Equal(t, "Test description", profile.Description)
	assert.Equal(t, *infDec, profile.Balance)
	assert.Equal(t, 1, profile.Role)

	mockRepo.On("GetProfileByID", 2, requestID).Return(nil, entityProfile.ErrProfileNotFound).Once()

	profile, err = testUsecase.GetProfile(2, requestID)
	assert.NoError(t, err)
	assert.Nil(t, profile)

	mockRepo.On("GetProfileByID", 3, requestID).Return(nil, fmt.Errorf("database error")).Once()

	profile, err = testUsecase.GetProfile(3, requestID)
	assert.Error(t, err)
	assert.Nil(t, profile)

	mockRepo.AssertExpectations(t)
}

func TestPutProfile(t *testing.T) {
	testUsecase := NewProfileUsecaseForTest()
	mockRepo := testUsecase.mockRepo

	userID := 1
	username := "newusername"
	description := "new description"
	requestID := "test-request-id"

	mockRepo.On("UpdateProfileByID", userID, username, description, requestID).Return(nil).Once()

	err := testUsecase.PutProfile(userID, username, description, requestID)
	assert.NoError(t, err)

	mockRepo.On("UpdateProfileByID", 2, username, description, requestID).Return(fmt.Errorf("database error")).Once()

	err = testUsecase.PutProfile(2, username, description, requestID)
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestPutProfile_ValidationErrors(t *testing.T) {
	testUsecase := NewProfileUsecaseForTest()
	mockRepo := testUsecase.mockRepo

	userID := 1
	username := ""
	description := "test description"
	requestID := "validation-test-request-id"

	mockRepo.On("UpdateProfileByID", userID, username, description, requestID).
		Return(errors.New("validation error: username cannot be empty")).Once()

	err := testUsecase.PutProfile(userID, username, description, requestID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation error")

	mockRepo.AssertExpectations(t)
}

func TestGetProfile_NilBalanceField(t *testing.T) {
	testUsecase := NewProfileUsecaseForTest()
	mockRepo := testUsecase.mockRepo

	userID := 1
	requestID := "nil-balance-test-request-id"

	mockProfile := &entityProfile.Profile{
		ID:          userID,
		Username:    "testuser",
		Email:       "test@example.com",
		Description: "Test description",
		Balance: entityProfile.Decimal{
			Dec: nil,
		},
		Role: 1,
	}

	mockRepo.On("GetProfileByID", userID, requestID).Return(mockProfile, nil).Once()

	assert.Panics(t, func() {
		_, _ = testUsecase.GetProfile(userID, requestID)
	})

	mockRepo.AssertExpectations(t)
}

func TestGetProfile_SQLError(t *testing.T) {
	testUsecase := NewProfileUsecaseForTest()
	mockRepo := testUsecase.mockRepo

	userID := 1
	requestID := "sql-error-test-request-id"

	sqlError := sql.ErrConnDone
	mockRepo.On("GetProfileByID", userID, requestID).Return(nil, sqlError).Once()

	profile, err := testUsecase.GetProfile(userID, requestID)

	assert.Error(t, err)
	assert.Equal(t, sqlError, err)
	assert.Nil(t, profile)

	mockRepo.AssertExpectations(t)
}

func TestPutProfile_EmptyFields(t *testing.T) {
	testUsecase := NewProfileUsecaseForTest()
	mockRepo := testUsecase.mockRepo

	userID := 1
	emptyUsername := ""
	emptyDescription := ""
	requestID := "empty-fields-test-request-id"

	mockRepo.On("UpdateProfileByID", userID, emptyUsername, emptyDescription, requestID).
		Return(nil).Once()

	err := testUsecase.PutProfile(userID, emptyUsername, emptyDescription, requestID)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestGetProfile_ComplexProfile(t *testing.T) {
	testUsecase := NewProfileUsecaseForTest()
	mockRepo := testUsecase.mockRepo

	userID := 1
	requestID := "complex-profile-test-request-id"

	negativeBalance := decimal.NewFromFloat(-999999.99)
	infDec := decimalToInfDec(negativeBalance)

	mockProfile := &entityProfile.Profile{
		ID:          userID,
		Username:    "complex-user",
		Email:       "complex@example.com",
		Description: "This is a very long description that contains many words to test the ability to handle large text fields. It might contain special characters like !@#$%^&*(), numbers 12345, and even emojis 😀.",
		Balance: entityProfile.Decimal{
			Dec: infDec,
		},
		Role: 2,
	}

	mockRepo.On("GetProfileByID", userID, requestID).Return(mockProfile, nil).Once()

	profile, err := testUsecase.GetProfile(userID, requestID)

	assert.NoError(t, err)
	assert.NotNil(t, profile)
	assert.Equal(t, userID, profile.ID)
	assert.Equal(t, "complex-user", profile.Username)
	assert.Equal(t, "complex@example.com", profile.Email)
	assert.Equal(t, mockProfile.Description, profile.Description)
	assert.Equal(t, *infDec, profile.Balance)
	assert.Equal(t, 2, profile.Role)

	mockRepo.AssertExpectations(t)
}
