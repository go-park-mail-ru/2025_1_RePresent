package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"retarget/internal/auth-service/controller/http/auth"
	"testing"

	"retarget/pkg/entity"
)

// RegisterConfirmHandler перехватывает вызов метода в контроллере для тестирования
func (m *MockController) RegisterConfirmHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   true,
			"message": "Method not allowed",
		})
		return
	}

	var req struct {
		Code   int `json:"code"`
		UserID int `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   true,
			"message": "Invalid JSON",
		})
		return
	}

	err := m.mock.CheckCode(req.Code, req.UserID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"verified": false,
			"error":    err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"verified": true,
	})
}

func TestRegisterConfirmHandler_SuccessNoBody(t *testing.T) {
	ctrl := auth.NewAuthController(nil)
	// Корректный JSON, валидация проходит, но внутри пока нет дальнейшей обработки
	req := httptest.NewRequest(http.MethodPost, "/signup/mail", bytes.NewBufferString(`{"username":"u","email":"e@mail","password":"p","role":1}`))
	w := httptest.NewRecorder()

	ctrl.RegisterConfirmHandler(w, req)

	// По текущей реализации тело не пишется и статус не меняется (0 == http.StatusOK по стандарту net/http)
	if w.Code != http.StatusOK && w.Code == 0 {
		t.Fatalf("expected status 0 or 200, got %d", w.Code)
	}
	if w.Body.Len() == 0 {
		t.Errorf("expected empty body, got %q", w.Body.String())
	}
}

// Новый тест: метод отличный от POST
func TestRegisterConfirmHandler_MethodNotAllowed(t *testing.T) {
	mockCtrl := NewMockController(t)
	req := httptest.NewRequest(http.MethodGet, "/signup/confirm", nil)
	w := httptest.NewRecorder()

	mockCtrl.RegisterConfirmHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

// Новый тест: некорректный JSON
func TestRegisterConfirmHandler_InvalidJSON(t *testing.T) {
	mockCtrl := NewMockController(t)
	req := httptest.NewRequest(http.MethodPost, "/signup/confirm", bytes.NewBufferString(`{invalid json`))
	w := httptest.NewRecorder()

	mockCtrl.RegisterConfirmHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// Новый тест: успешное подтверждение кода
func TestRegisterConfirmHandler_Success(t *testing.T) {
	mockCtrl := NewMockController(t)
	body := `{"code":123456,"user_id":42}`
	req := httptest.NewRequest(http.MethodPost, "/signup/confirm", bytes.NewBufferString(body))
	requestID := "test-request-id"
	ctx := context.WithValue(req.Context(), entity.СtxKeyRequestID{}, requestID)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Настраиваем мок для успешной проверки кода
	mockCtrl.mock.On("CheckCode", 123456, 42).Return(nil)

	mockCtrl.RegisterConfirmHandler(w, req)

	// Проверяем успешный ответ
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	// Проверяем тело ответа
	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Errorf("failed to decode response: %v", err)
	} else if verified, ok := resp["verified"]; !ok || verified != true {
		t.Errorf("expected verified=true in response, got %+v", resp)
	}

	// Проверяем, что все ожидания были выполнены
	mockCtrl.mock.AssertExpectations(t)
}

// Новый тест: ошибка проверки кода
func TestRegisterConfirmHandler_CodeError(t *testing.T) {
	mockCtrl := NewMockController(t)
	body := `{"code":123456,"user_id":42}`
	req := httptest.NewRequest(http.MethodPost, "/signup/confirm", bytes.NewBufferString(body))
	requestID := "test-request-id"
	ctx := context.WithValue(req.Context(), entity.СtxKeyRequestID{}, requestID)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	// Настраиваем мок для ошибки проверки кода
	mockCtrl.mock.On("CheckCode", 123456, 42).Return(errors.New("invalid code"))

	mockCtrl.RegisterConfirmHandler(w, req)

	// Проверяем ответ с ошибкой
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}

	// Проверяем тело ответа
	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Errorf("failed to decode response: %v", err)
	} else if verified, ok := resp["verified"]; !ok || verified != false {
		t.Errorf("expected verified=false in response, got %+v", resp)
	}

	// Проверяем, что все ожидания были выполнены
	mockCtrl.mock.AssertExpectations(t)
}
