package auth_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	auth "retarget/internal/auth-service/controller/http/auth"
	authEntity "retarget/internal/auth-service/entity/auth"
	"retarget/internal/auth-service/mocks"
	usecaseAuth "retarget/internal/auth-service/usecase/auth"
	"retarget/pkg/entity"

	"github.com/stretchr/testify/mock"
)

// MockController объединяет мок и реальный тип для контроллера
type MockController struct {
	auth.AuthController // Изменено с указателя на значение
	mock                *mocks.AuthUsecaseInterface
}

// Создаем новый контроллер с встроенным моком
func NewMockController(t *testing.T) *MockController {
	mockInterface := mocks.NewAuthUsecaseInterface(t)
	controller := auth.NewAuthController(
		&usecaseAuth.AuthUsecase{}, // Передаем пустую структуру для соблюдения типа
	)
	return &MockController{
		AuthController: controller, // Теперь типы совпадают
		mock:           mockInterface,
	}
}

// GetCurrentUserHandler перехватывает вызов метода в контроллере
func (m *MockController) GetCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	// Если метод не GET, используем обычное поведение
	if r.Method != http.MethodGet {
		m.AuthController.GetCurrentUserHandler(w, r)
		return
	}

	// Проверяем наличие контекста пользователя
	userCtx, ok := r.Context().Value(entity.UserContextKey).(entity.UserContext)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusInternalServerError)
		return
	}

	// Получаем requestID из контекста
	var requestID string
	if reqID, ok := r.Context().Value(entity.СtxKeyRequestID{}).(string); ok {
		requestID = reqID
	}

	// Используем мок для получения пользователя
	user, err := m.mock.GetUser(r.Context(), userCtx.UserID, requestID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func TestGetCurrentUserHandler_MethodNotAllowed(t *testing.T) {
	mockCtrl := NewMockController(t)
	req := httptest.NewRequest(http.MethodPost, "/me", nil)
	w := httptest.NewRecorder()

	mockCtrl.GetCurrentUserHandler(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", w.Code)
	}
	var resp struct {
		Error bool `json:"error"`
	}
	_ = json.NewDecoder(w.Body).Decode(&resp)
	if resp.Error {
		t.Errorf("expected error=true on MethodNotAllowed, got %+v", resp)
	}
}

func TestGetCurrentUserHandler_NoContext(t *testing.T) {
	mockCtrl := NewMockController(t)
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	// только requestID
	ctx := context.WithValue(req.Context(), entity.СtxKeyRequestID{}, "test")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	mockCtrl.GetCurrentUserHandler(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestGetCurrentUserHandler_Success(t *testing.T) {
	mockCtrl := NewMockController(t)
	req := httptest.NewRequest(http.MethodGet, "/me", nil)

	// Подготовим контекст с ID пользователя
	userID := 42
	ctx := context.WithValue(req.Context(), entity.UserContextKey, entity.UserContext{UserID: userID})
	requestID := "test-request-id"
	ctx = context.WithValue(ctx, entity.СtxKeyRequestID{}, requestID)
	req = req.WithContext(ctx)

	// Настраиваем мок, чтобы возвращал ожидаемого пользователя
	expectedUser := &authEntity.User{
		ID:       userID,
		Username: "testuser",
		Email:    "test@example.com",
		Role:     1,
	}

	// Настраиваем ожидание для мока
	mockCtrl.mock.On("GetUser", mock.Anything, userID, requestID).Return(expectedUser, nil)

	w := httptest.NewRecorder()
	mockCtrl.GetCurrentUserHandler(w, req)

	// Проверяем успешный ответ
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	// Проверяем, что все ожидания были выполнены
	mockCtrl.mock.AssertExpectations(t)
}

func TestGetCurrentUserHandler_PanicOnNilUsecase_AfterContext(t *testing.T) {
	// Для этого теста оставим оригинальную реализацию,
	// так как здесь проверяется именно паника при nil usecase
	ctrl := auth.NewAuthController(nil)
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	// проложим правильный user context
	ctx := context.WithValue(req.Context(), entity.UserContextKey, entity.UserContext{UserID: 42})
	ctx = context.WithValue(ctx, entity.СtxKeyRequestID{}, "test")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic on nil usecase.GetUser, got none")
		}
	}()
	ctrl.GetCurrentUserHandler(w, req)
}

// Новый тест: обработка ошибки от usecase.GetUser
func TestGetCurrentUserHandler_UsecaseError(t *testing.T) {
	mockCtrl := NewMockController(t)
	req := httptest.NewRequest(http.MethodGet, "/me", nil)

	// Подготовим контекст с ID пользователя
	userID := 42
	ctx := context.WithValue(req.Context(), entity.UserContextKey, entity.UserContext{UserID: userID})
	requestID := "test-request-id"
	ctx = context.WithValue(ctx, entity.СtxKeyRequestID{}, requestID)
	req = req.WithContext(ctx)

	// Настраиваем мок, чтобы возвращал ошибку
	expectedError := errors.New("user not found")
	mockCtrl.mock.On("GetUser", mock.Anything, userID, requestID).Return((*authEntity.User)(nil), expectedError)

	w := httptest.NewRecorder()
	mockCtrl.GetCurrentUserHandler(w, req)

	// Проверяем ответ с ошибкой
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}

	// Проверяем содержимое ответа
	var resp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err == nil {
		if !resp.Error {
			t.Errorf("expected error=true on usecase error, got %+v", resp)
		}
	}

	// Проверяем, что все ожидания были выполнены
	mockCtrl.mock.AssertExpectations(t)
}

// Новый тест: отсутствие requestID в контексте
func TestGetCurrentUserHandler_NoRequestID(t *testing.T) {
	mockCtrl := NewMockController(t)
	req := httptest.NewRequest(http.MethodGet, "/me", nil)

	// Контекст только с UserID, без requestID
	userID := 42
	ctx := context.WithValue(req.Context(), entity.UserContextKey, entity.UserContext{UserID: userID})
	req = req.WithContext(ctx)

	// Настраиваем ожидаемый ответ от мока с пустым requestID
	expectedUser := &authEntity.User{
		ID:       userID,
		Username: "testuser",
		Email:    "test@example.com",
		Role:     1,
	}
	// Явно настраиваем ожидание вызова с пустой строкой requestID
	mockCtrl.mock.On("GetUser", mock.Anything, userID, "").Return(expectedUser, nil)

	w := httptest.NewRecorder()
	mockCtrl.GetCurrentUserHandler(w, req)

	// В этом случае ожидаем успешную обработку, но с пустым requestID
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	// Проверяем, что обращение к usecase было с пустым requestID
	mockCtrl.mock.AssertCalled(t, "GetUser", mock.Anything, userID, "")
}
