package mail

import (
	"crypto/tls"
	"errors"
	"net"
	"net/smtp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSMTPWriter реализует интерфейс SMTPWriter для тестирования
type MockSMTPWriter struct {
	mock.Mock
}

func (m *MockSMTPWriter) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockSMTPWriter) Write(p []byte) (int, error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

// MailRepositoryTestWrapper оборачивает MailRepository для тестирования
type MailRepositoryTestWrapper struct {
	*MailRepository
	mockVerifyEmail     func(email string) (bool, error)
	mockEnsureConnected func() error
}

func (w *MailRepositoryTestWrapper) VerifyEmail(email string) (bool, error) {
	if w.mockVerifyEmail != nil {
		return w.mockVerifyEmail(email)
	}
	return w.MailRepository.VerifyEmail(email)
}

// Переопределяем метод Send чтобы использовать мок для ensureConnected
func (w *MailRepositoryTestWrapper) Send(to, msg string) error {
	// Вызываем наш мок вместо оригинального ensureConnected
	if w.mockEnsureConnected != nil {
		if err := w.mockEnsureConnected(); err != nil {
			return err
		}
	} else {
		if err := w.MailRepository.ensureConnected(); err != nil {
			return err
		}
	}

	// Остальная часть метода Send
	isExist, err := w.VerifyEmail(to)
	if !isExist || err != nil {
		return err
	}

	if err := smtpMail(w.smtpClient, w.from_sender); err != nil {
		return err
	}
	if err := smtpRcpt(w.smtpClient, to); err != nil {
		return err
	}
	wData, err := smtpData(w.smtpClient)
	if err != nil {
		return err
	}
	_, err = wData.Write([]byte(msg))
	if err != nil {
		return err
	}
	err = wData.Close()
	if err != nil {
		return err
	}

	return nil
}

// TestMailRepository_Send проверяет функцию Send
func TestMailRepository_Send(t *testing.T) {
	t.Run("успешная отправка", func(t *testing.T) {
		// Создаем моки
		mockClient := &smtp.Client{}
		mockWriter := &MockSMTPWriter{}

		// Настраиваем ожидаемое поведение
		mockWriter.On("Write", mock.Anything).Return(10, nil)
		mockWriter.On("Close").Return(nil)

		// Создаем патч для smtp функций
		origSmtpMail := smtpMail
		origSmtpRcpt := smtpRcpt
		origSmtpData := smtpData

		smtpMail = func(client *smtp.Client, from string) error {
			return nil
		}

		smtpRcpt = func(client *smtp.Client, to string) error {
			return nil
		}

		smtpData = func(client *smtp.Client) (SMTPWriter, error) {
			return mockWriter, nil
		}

		// Восстанавливаем оригинальные функции после теста
		defer func() {
			smtpMail = origSmtpMail
			smtpRcpt = origSmtpRcpt
			smtpData = origSmtpData
		}()

		// Создаем репозиторий
		repo := &MailRepository{
			smtpClient:  mockClient,
			from_sender: "sender@example.com",
		}

		// Создаем wrapper с мок-функциями
		wrapper := &MailRepositoryTestWrapper{
			MailRepository: repo,
			mockVerifyEmail: func(email string) (bool, error) {
				return true, nil
			},
			mockEnsureConnected: func() error {
				return nil // Пропускаем проверку соединения
			},
		}

		// Вызываем метод
		err := wrapper.Send("test@example.com", "Test message")

		// Проверяем результат
		assert.NoError(t, err)
		mockWriter.AssertExpectations(t)
	})

	t.Run("ошибка при проверке email", func(t *testing.T) {
		// Создаем репозиторий
		repo := &MailRepository{
			from_sender: "sender@example.com",
		}

		// Создаем wrapper с мок-функциями
		wrapper := &MailRepositoryTestWrapper{
			MailRepository: repo,
			mockVerifyEmail: func(email string) (bool, error) {
				return false, ErrEmailNotExists
			},
			mockEnsureConnected: func() error {
				return nil // Пропускаем проверку соединения
			},
		}

		// Вызываем метод
		err := wrapper.Send("invalid@example.com", "Test message")

		// Проверяем результат
		assert.Error(t, err)
		assert.Equal(t, ErrEmailNotExists, err)
	})

	t.Run("ошибка при соединении", func(t *testing.T) {
		// Создаем репозиторий
		repo := &MailRepository{
			from_sender: "sender@example.com",
		}

		// Создаем wrapper с мок-функциями
		wrapper := &MailRepositoryTestWrapper{
			MailRepository: repo,
			mockEnsureConnected: func() error {
				return errors.New("connection failed")
			},
		}

		// Вызываем метод
		err := wrapper.Send("test@example.com", "Test message")

		// Проверяем результат
		assert.Error(t, err)
		assert.Equal(t, "connection failed", err.Error())
	})
}

// TestMailRepository_VerifyEmail проверяет функцию VerifyEmail
func TestMailRepository_VerifyEmail(t *testing.T) {
	// Сохраняем оригинальные функции
	origLookupMX := lookupMXFn
	origDialTimeout := dialTimeoutFn
	origNewClient := newClientFn

	defer func() {
		// Восстанавливаем оригинальные функции
		lookupMXFn = origLookupMX
		dialTimeoutFn = origDialTimeout
		newClientFn = origNewClient
	}()

	t.Run("некорректный email", func(t *testing.T) {
		repo := &MailRepository{from_sender: "sender@example.com"}
		result, err := repo.VerifyEmail("invalid-email")
		assert.False(t, result)
		assert.Equal(t, ErrInvalidEmail, err)
	})

	t.Run("отсутствуют MX записи", func(t *testing.T) {
		repo := &MailRepository{from_sender: "sender@example.com"}

		// Мокируем LookupMX
		lookupMXFn = func(name string) ([]*net.MX, error) {
			return nil, errors.New("no mx records")
		}

		result, err := repo.VerifyEmail("test@example.com")
		assert.False(t, result)
		assert.Equal(t, ErrNoMXRecords, err)
	})
}

// TestMailRepository_OpenConnection проверяет функцию OpenConnection
func TestMailRepository_OpenConnection(t *testing.T) {
	// Сохраняем оригинальные функции
	origDialTLS := dialTLSFn
	origNewClient := newClientFn
	origSmtpAuth := smtpAuth

	defer func() {
		// Восстанавливаем оригинальные функции
		dialTLSFn = origDialTLS
		newClientFn = origNewClient
		smtpAuth = origSmtpAuth
	}()

	t.Run("успешное соединение", func(t *testing.T) {
		repo := &MailRepository{
			smtpServer: "smtp.example.com",
			smtpPort:   "587",
			username:   "user",
			password:   "pass",
		}

		// Мокируем функции
		dialTLSFn = func(network, addr string, config *tls.Config) (*tls.Conn, error) {
			return &tls.Conn{}, nil
		}

		newClientFn = func(conn net.Conn, host string) (*smtp.Client, error) {
			return &smtp.Client{}, nil
		}

		smtpAuth = func(client *smtp.Client, auth smtp.Auth) error {
			return nil
		}

		// Вызываем метод
		client, err := repo.OpenConnection()

		// Проверяем результат
		assert.NoError(t, err)
		assert.NotNil(t, client)
	})

	t.Run("ошибка TLS", func(t *testing.T) {
		repo := &MailRepository{
			smtpServer: "smtp.example.com",
			smtpPort:   "587",
		}

		// Мокируем функции с ошибкой
		dialTLSFn = func(network, addr string, config *tls.Config) (*tls.Conn, error) {
			return nil, errors.New("tls error")
		}

		// Вызываем метод
		client, err := repo.OpenConnection()

		// Проверяем результат
		assert.Error(t, err)
		assert.Nil(t, client)
		assert.Equal(t, "tls error", err.Error())
	})
}

// TestMailRepository_CloseConnection проверяет функцию CloseConnection
func TestMailRepository_CloseConnection(t *testing.T) {
	// Сохраняем оригинальную функцию
	origSmtpQuit := smtpQuit

	defer func() {
		// Восстанавливаем оригинальную функцию
		smtpQuit = origSmtpQuit
	}()

	t.Run("успешное закрытие", func(t *testing.T) {
		repo := &MailRepository{
			smtpClient: &smtp.Client{},
		}

		// Мокируем функцию Quit
		smtpQuit = func(client *smtp.Client) error {
			return nil
		}

		// Вызываем метод
		err := repo.CloseConnection()

		// Проверяем результат
		assert.NoError(t, err)
	})

	t.Run("ошибка при закрытии", func(t *testing.T) {
		repo := &MailRepository{
			smtpClient: &smtp.Client{},
		}

		// Мокируем функцию Quit с ошибкой
		smtpQuit = func(client *smtp.Client) error {
			return errors.New("quit error")
		}

		// Вызываем метод
		err := repo.CloseConnection()

		// Проверяем результат
		assert.Error(t, err)
		assert.Equal(t, "quit error", err.Error())
	})
}

// TestMailRepository_EnsureConnected проверяет функцию ensureConnected
func TestMailRepository_EnsureConnected(t *testing.T) {
	// Сохраняем оригинальные функции
	origSmtpNoop := smtpNoop
	origSmtpQuit := smtpQuit

	defer func() {
		// Восстанавливаем оригинальные функции
		smtpNoop = origSmtpNoop
		smtpQuit = origSmtpQuit
	}()

	t.Run("соединение уже установлено", func(t *testing.T) {
		repo := &MailRepository{
			smtpClient: &smtp.Client{},
		}

		// Мокируем функцию Noop
		smtpNoop = func(client *smtp.Client) error {
			return nil // соединение активно
		}

		// Подменяем метод OpenConnection, чтобы убедиться что он не вызывается
		originalOpenConnection := repo.OpenConnection
		openConnectionCalled := false
		repo.OpenConnection = func() (*smtp.Client, error) {
			openConnectionCalled = true
			return nil, nil
		}
		defer func() { repo.OpenConnection = originalOpenConnection }()

		// Вызываем метод
		err := repo.ensureConnected()

		// Проверяем результат
		assert.NoError(t, err)
		assert.False(t, openConnectionCalled, "OpenConnection не должен вызываться")
	})

	t.Run("соединение закрыто, переподключение", func(t *testing.T) {
		repo := &MailRepository{
			smtpClient: &smtp.Client{},
		}

		// Мокируем функции
		smtpNoop = func(client *smtp.Client) error {
			return errors.New("connection closed") // соединение закрыто
		}

		smtpQuit = func(client *smtp.Client) error {
			return nil // закрытие успешно
		}

		// Подменяем метод OpenConnection
		originalOpenConnection := repo.OpenConnection
		newClient := &smtp.Client{}
		repo.OpenConnection = func() (*smtp.Client, error) {
			return newClient, nil
		}
		defer func() { repo.OpenConnection = originalOpenConnection }()

		// Вызываем метод
		err := repo.ensureConnected()

		// Проверяем результат
		assert.NoError(t, err)
		assert.Equal(t, newClient, repo.smtpClient, "Должен быть установлен новый клиент")
	})
}
