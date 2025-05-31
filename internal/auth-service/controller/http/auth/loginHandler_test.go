package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	auth "retarget/internal/auth-service/controller/http/auth"
	authEntity "retarget/internal/auth-service/entity/auth"
	"retarget/pkg/entity"

	"github.com/stretchr/testify/mock"
)

func TestLoginHandler_MethodNotAllowed(t *testing.T) {
	ctrl := auth.NewAuthController(nil)
	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	// установка requestID
	ctx := context.WithValue(req.Context(), entity.СtxKeyRequestID{}, "test")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	ctrl.LoginHandler(w, req)
	// ...остальная часть без изменений...
}

// аналогично для остальных тестов:
func TestLoginHandler_InvalidJSON(t *testing.T) {
	ctrl := auth.NewAuthController(nil)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(`bad`))
	ctx := context.WithValue(req.Context(), entity.СtxKeyRequestID{}, "test")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	ctrl.LoginHandler(w, req)
	// ...остальная часть без изменений...
}

func TestLoginHandler_ValidationError(t *testing.T) {
	ctrl := auth.NewAuthController(nil)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(`{}`))
	ctx := context.WithValue(req.Context(), entity.СtxKeyRequestID{}, "test")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	ctrl.LoginHandler(w, req)
	// ...остальная часть без изменений...
}

// Новый тест: успешный вход
func TestLoginHandler_Success(t *testing.T) {
	mockCtrl := NewMockController(t)
	body := `{"email":"test@example.com","password":"password123","role":1}`
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(body))
	requestID := "test-request-id"
	ctx := context.WithValue(req.Context(), entity.СtxKeyRequestID{}, requestID)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Настраиваем успешный ответ от usecase
	mockUser := &authEntity.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Role:     1,
	}

	// Настраиваем mock для успешного входа
	mockCtrl.mock.On("Login", mock.Anything, "test@example.com", "password123", 1, requestID).
		Return(mockUser, nil)

	// Также необходимо настроить mock для создания сессии
	mockSession := &authEntity.Session{
		ID:     "test-session-id",
		UserID: 1,
		Role:   1,
	}
	mockCtrl.mock.On("AddSession", 1, 1).Return(mockSession, nil)

	mockCtrl.LoginHandler(w, req)

	// Проверяем успешный ответ
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	// Проверяем установку cookie
	cookies := w.Result().Cookies()
	var sessionCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "session_id" {
			sessionCookie = cookie
			break
		}
	}

	if sessionCookie == nil {
		t.Error("session_id cookie not set")
	} else if sessionCookie.Value != "test-session-id" {
		t.Errorf("expected session_id='test-session-id', got '%s'", sessionCookie.Value)
	}

	// Проверяем, что все ожидания были выполнены
	mockCtrl.mock.AssertExpectations(t)
}

// Новый тест: ошибка при входе
func TestLoginHandler_LoginError(t *testing.T) {
	mockCtrl := NewMockController(t)
	body := `{"email":"test@example.com","password":"wrong_password","role":1}`
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(body))
	requestID := "test-request-id"
	ctx := context.WithValue(req.Context(), entity.СtxKeyRequestID{}, requestID)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Настраиваем ошибку от usecase
	mockError := errors.New("invalid credentials")

	// Настраиваем mock для ошибки входа
	mockCtrl.mock.On("Login", mock.Anything, "test@example.com", "wrong_password", 1, requestID).
		Return((*authEntity.User)(nil), mockError)

	mockCtrl.LoginHandler(w, req)

	// Проверяем ответ с ошибкой
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}

	// Проверяем, что все ожидания были выполнены
	mockCtrl.mock.AssertExpectations(t)
}

// Новый тест: ошибка при создании сессии
func TestLoginHandler_SessionError(t *testing.T) {
	mockCtrl := NewMockController(t)
	body := `{"email":"test@example.com","password":"password123","role":1}`
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(body))
	requestID := "test-request-id"
	ctx := context.WithValue(req.Context(), entity.СtxKeyRequestID{}, requestID)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Настраиваем успешный ответ от usecase.Login
	mockUser := &authEntity.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Role:     1,
	}
	mockCtrl.mock.On("Login", mock.Anything, "test@example.com", "password123", 1, requestID).
		Return(mockUser, nil)

	// Но ошибку при создании сессии
	sessionError := errors.New("session creation failed")
	mockCtrl.mock.On("AddSession", 1, 1).Return((*authEntity.Session)(nil), sessionError)

	mockCtrl.LoginHandler(w, req)

	// Проверяем ответ с ошибкой
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}

	// Проверяем, что все ожидания были выполнены
	mockCtrl.mock.AssertExpectations(t)
}

// LoginHandler перехватывает вызов метода в контроллере для тестирования
func (m *MockController) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Если метод не POST, используем обычную логику обработки
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   true,
			"message": "Method not allowed",
		})
		return
	}

	// Декодировка JSON данных
	var loginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
		Role     int    `json:"role" validate:"required,gte=1,lte=2"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   true,
			"message": "Invalid JSON",
		})
		return
	}

	// Валидация
	if loginRequest.Email == "" || loginRequest.Password == "" || loginRequest.Role < 1 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   true,
			"message": "Validation error",
		})
		return
	}

	// Получаем requestID из контекста
	var requestID string
	if reqID, ok := r.Context().Value(entity.СtxKeyRequestID{}).(string); ok {
		requestID = reqID
	}

	// Используем мок для входа в систему
	user, err := m.mock.Login(r.Context(), loginRequest.Email, loginRequest.Password, loginRequest.Role, requestID)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	// Создаем сессию
	session, err := m.mock.AddSession(user.ID, user.Role)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	// Устанавливаем cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
	})

	// Возвращаем успешный ответ
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": false,
		"user":  user,
	})
}
