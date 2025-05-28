package mail

import (
	"errors"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	authEntity "retarget/internal/auth-service/entity/auth"
	repoAuth "retarget/internal/auth-service/repo/auth"
	entityMail "retarget/internal/mail-service/entity/mail"
	mocksMail "retarget/internal/mail-service/mocks"
)

// Мок для AuthRepository
type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) GetUserByID(userID int, requestID string) (*authEntity.User, error) {
	args := m.Called(userID, requestID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authEntity.User), args.Error(1)
}

// Тест для функции SendCodeMail
func TestMailUsecase_SendCodeMail(t *testing.T) {
	tests := []struct {
		name      string
		operation int
		to        string
		code      string
		mockSetup func(*mocksMail.MailRepositoryInterface)
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "успешная отправка кода регистрации",
			operation: entityMail.REGISTER,
			to:        "test@example.com",
			code:      "123456",
			mockSetup: func(mailRepo *mocksMail.MailRepositoryInterface) {
				mailRepo.On("Send", "test@example.com", mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "успешная отправка кода сброса пароля",
			operation: entityMail.RESET_PASSWORD,
			to:        "test@example.com",
			code:      "123456",
			mockSetup: func(mailRepo *mocksMail.MailRepositoryInterface) {
				mailRepo.On("Send", "test@example.com", mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "успешная отправка кода изменения пароля",
			operation: entityMail.EDIT_PASSWORD,
			to:        "test@example.com",
			code:      "123456",
			mockSetup: func(mailRepo *mocksMail.MailRepositoryInterface) {
				mailRepo.On("Send", "test@example.com", mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "неопределенная операция",
			operation: 999,
			to:        "test@example.com",
			code:      "123456",
			mockSetup: func(mailRepo *mocksMail.MailRepositoryInterface) {
				// нет вызовов
			},
			wantErr: true,
			errMsg:  "undefined operation",
		},
		{
			name:      "ошибка отправки почты",
			operation: entityMail.REGISTER,
			to:        "test@example.com",
			code:      "123456",
			mockSetup: func(mailRepo *mocksMail.MailRepositoryInterface) {
				mailRepo.On("Send", "test@example.com", mock.Anything).Return(errors.New("send error"))
			},
			wantErr: true,
			errMsg:  "send error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем моки
			mockMailRepo := mocksMail.NewMailRepositoryInterface(t)
			mockAuthRepo := new(MockAuthRepository)

			// Настраиваем ожидаемое поведение
			tt.mockSetup(mockMailRepo)

			// Создаем usecase
			mailUsecase := NewMailUsecase(mockMailRepo, mockAuthRepo)

			// Вызываем тестируемый метод
			err := mailUsecase.SendCodeMail(tt.operation, tt.to, tt.code)

			// Проверяем результат
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}

			// Проверяем вызовы моков
			mockMailRepo.AssertExpectations(t)
		})
	}
}

// Тест для функции SendLowBalanceMail
func TestMailUsecase_SendLowBalanceMail(t *testing.T) {
	tests := []struct {
		name      string
		operation int
		to        string
		username  string
		balance   string
		href      string
		mockSetup func(*mocksMail.MailRepositoryInterface)
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "успешная отправка уведомления о низком балансе",
			operation: entityMail.LOW_BALANCE,
			to:        "test@example.com",
			username:  "testuser",
			balance:   "10.50",
			href:      "https://example.com/balance",
			mockSetup: func(mailRepo *mocksMail.MailRepositoryInterface) {
				mailRepo.On("Send", "test@example.com", mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "неопределенная операция",
			operation: 999,
			to:        "test@example.com",
			username:  "testuser",
			balance:   "10.50",
			href:      "https://example.com/balance",
			mockSetup: func(mailRepo *mocksMail.MailRepositoryInterface) {
				// нет вызовов
			},
			wantErr: true,
			errMsg:  "undefined operation",
		},
		{
			name:      "ошибка отправки почты",
			operation: entityMail.LOW_BALANCE,
			to:        "test@example.com",
			username:  "testuser",
			balance:   "10.50",
			href:      "https://example.com/balance",
			mockSetup: func(mailRepo *mocksMail.MailRepositoryInterface) {
				mailRepo.On("Send", "test@example.com", mock.Anything).Return(errors.New("send error"))
			},
			wantErr: true,
			errMsg:  "send error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем моки
			mockMailRepo := mocksMail.NewMailRepositoryInterface(t)
			mockAuthRepo := new(MockAuthRepository)

			// Настраиваем ожидаемое поведение
			tt.mockSetup(mockMailRepo)

			// Создаем usecase
			mailUsecase := NewMailUsecase(mockMailRepo, mockAuthRepo)

			// Вызываем тестируемый метод
			err := mailUsecase.SendLowBalanceMail(tt.operation, tt.to, tt.username, tt.balance, tt.href)

			// Проверяем результат
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}

			// Проверяем вызовы моков
			mockMailRepo.AssertExpectations(t)
		})
	}
}

// Тест для функции SendTopUpBalanceMail
func TestMailUsecase_SendTopUpBalanceMail(t *testing.T) {
	tests := []struct {
		name      string
		operation int
		to        string
		username  string
		amount    string
		mockSetup func(*mocksMail.MailRepositoryInterface)
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "успешная отправка уведомления о пополнении баланса",
			operation: entityMail.TOPUP_BALANCE,
			to:        "test@example.com",
			username:  "testuser",
			amount:    "100.00",
			mockSetup: func(mailRepo *mocksMail.MailRepositoryInterface) {
				mailRepo.On("Send", "test@example.com", mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "неопределенная операция",
			operation: 999,
			to:        "test@example.com",
			username:  "testuser",
			amount:    "100.00",
			mockSetup: func(mailRepo *mocksMail.MailRepositoryInterface) {
				// нет вызовов
			},
			wantErr: true,
			errMsg:  "undefined operation",
		},
		{
			name:      "ошибка отправки почты",
			operation: entityMail.TOPUP_BALANCE,
			to:        "test@example.com",
			username:  "testuser",
			amount:    "100.00",
			mockSetup: func(mailRepo *mocksMail.MailRepositoryInterface) {
				mailRepo.On("Send", "test@example.com", mock.Anything).Return(errors.New("send error"))
			},
			wantErr: true,
			errMsg:  "send error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем моки
			mockMailRepo := mocksMail.NewMailRepositoryInterface(t)
			mockAuthRepo := new(MockAuthRepository)

			// Настраиваем ожидаемое поведение
			tt.mockSetup(mockMailRepo)

			// Создаем usecase
			mailUsecase := NewMailUsecase(mockMailRepo, mockAuthRepo)

			// Вызываем тестируемый метод
			err := mailUsecase.SendTopUpBalanceMail(tt.operation, tt.to, tt.username, tt.amount)

			// Проверяем результат
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}

			// Проверяем вызовы моков
			mockMailRepo.AssertExpectations(t)
		})
	}
}

// Тест для функции GetUserByID
func TestMailUsecase_GetUserByID(t *testing.T) {
	tests := []struct {
		name        string
		userID      int
		mockSetup   func(*MockAuthRepository)
		wantEmail   string
		wantName    string
		wantBalance string
		wantErr     bool
	}{
		{
			name:   "успешное получение пользователя",
			userID: 123,
			mockSetup: func(authRepo *MockAuthRepository) {
				balance, _ := decimal.NewFromString("500.75")
				authRepo.On("GetUserByID", 123, mock.Anything).Return(&authEntity.User{
					ID:       123,
					Username: "testuser",
					Email:    "test@example.com",
					Balance:  balance,
				}, nil)
			},
			wantEmail:   "test@example.com",
			wantName:    "testuser",
			wantBalance: "500.75",
			wantErr:     false,
		},
		{
			name:   "ошибка получения пользователя",
			userID: 456,
			mockSetup: func(authRepo *MockAuthRepository) {
				authRepo.On("GetUserByID", 456, mock.Anything).Return(nil, errors.New("user not found"))
			},
			wantEmail:   "froloff1830@gmail.com",
			wantName:    "Error if parse Username",
			wantBalance: "Balance",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем моки
			mockMailRepo := mocksMail.NewMailRepositoryInterface(t)
			mockAuthRepo := new(MockAuthRepository)

			// Настраиваем ожидаемое поведение
			tt.mockSetup(mockAuthRepo)

			// Создаем usecase
			mailUsecase := NewMailUsecase(mockMailRepo, mockAuthRepo)

			// Вызываем тестируемый метод
			email, username, balance, err := mailUsecase.GetUserByID(tt.userID)

			// Проверяем результат
			assert.Equal(t, tt.wantEmail, email)
			assert.Equal(t, tt.wantName, username)
			assert.Equal(t, tt.wantBalance, balance)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Проверяем вызовы моков
			mockAuthRepo.AssertExpectations(t)
		})
	}
}
