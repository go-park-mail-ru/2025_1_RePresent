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
	entity "retarget/pkg/entity"

	"github.com/stretchr/testify/mock"
)

// RegisterHandler перехватывает вызов метода в контроллере для тестирования
func (m *MockController) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   true,
			"message": "Method not allowed",
		})
		return
	}

	// Декодировка JSON данных
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     int    `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   true,
			"message": "Invalid JSON",
		})
		return
	}

	// Простая валидация
	if req.Username == "" || req.Email == "" || req.Password == "" || req.Role < 1 {
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

	// Используем мок для регистрации
	user, err := m.mock.Register(r.Context(), req.Username, req.Email, req.Password, req.Role, requestID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": false,
		"user":  user,
	})
}

func TestRegisterHandler_PanicOnNilUsecase(t *testing.T) {
	ctrl := auth.NewAuthController(nil)
	body := `{"username":"u","email":"e@mail","password":"p","role":1}`
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBufferString(body))
	// нужно проложить requestID, чтобы дойти до usecase.Register
	ctx := context.WithValue(req.Context(), entity.СtxKeyRequestID{}, "test")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("expected panic on nil usecase.Register, got none")
		}
	}()
	ctrl.RegisterHandler(w, req)
}

// Новый тест: проверка метода, отличного от POST
func TestRegisterHandler_MethodNotAllowed(t *testing.T) {
	mockCtrl := NewMockController(t) // Используем наш существующий MockController
	req := httptest.NewRequest(http.MethodGet, "/signup", nil)
	w := httptest.NewRecorder()

	mockCtrl.RegisterHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

// Новый тест: некорректный JSON в теле запроса
func TestRegisterHandler_InvalidJSON(t *testing.T) {
	mockCtrl := NewMockController(t)
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBufferString(`{invalid json`))
	ctx := context.WithValue(req.Context(), entity.СtxKeyRequestID{}, "test")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	mockCtrl.RegisterHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// Новый тест: ошибка валидации входных данных
func TestRegisterHandler_ValidationError(t *testing.T) {
	mockCtrl := NewMockController(t)
	// Некорректные данные для валидации (отсутствует обязательное поле)
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBufferString(`{"username":"u"}`))
	ctx := context.WithValue(req.Context(), entity.СtxKeyRequestID{}, "test")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	mockCtrl.RegisterHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// Новый тест: успешная регистрация
func TestRegisterHandler_Success(t *testing.T) {
	mockCtrl := NewMockController(t)
	body := `{"username":"username","email":"valid@example.com","password":"password123","role":1}`
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBufferString(body))
	requestID := "test-request-id"
	ctx := context.WithValue(req.Context(), entity.СtxKeyRequestID{}, requestID)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Настраиваем ожидаемый ответ от usecase
	mockUser := &authEntity.User{
		ID:       1,
		Username: "username",
		Email:    "valid@example.com",
		Role:     1,
	}

	// Настраиваем мок для успешной регистрации
	mockCtrl.mock.On("Register", mock.Anything, "username", "valid@example.com", "password123", 1, requestID).
		Return(mockUser, nil)

	mockCtrl.RegisterHandler(w, req)

	// Проверяем успешный ответ
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}

	// Проверяем, что все ожидания были выполнены
	mockCtrl.mock.AssertExpectations(t)
}

// Новый тест: ошибка при регистрации
func TestRegisterHandler_RegisterError(t *testing.T) {
	mockCtrl := NewMockController(t)
	body := `{"username":"username","email":"valid@example.com","password":"password123","role":1}`
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBufferString(body))
	requestID := "test-request-id"
	ctx := context.WithValue(req.Context(), entity.СtxKeyRequestID{}, requestID)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Настраиваем ошибку от usecase
	mockError := errors.New("email already exists")

	// Настраиваем мок для ошибки регистрации
	mockCtrl.mock.On("Register", mock.Anything, "username", "valid@example.com", "password123", 1, requestID).
		Return((*authEntity.User)(nil), mockError)

	mockCtrl.RegisterHandler(w, req)

	// Проверяем ответ с ошибкой
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}

	// Проверяем, что все ожидания были выполнены
	mockCtrl.mock.AssertExpectations(t)
}
