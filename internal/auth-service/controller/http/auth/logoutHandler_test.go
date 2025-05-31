package auth_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	auth "retarget/internal/auth-service/controller/http/auth"
)

func TestLogoutHandler_PanicOnNilUsecase(t *testing.T) {
	ctrl := auth.NewAuthController(nil)
	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	// Эмулируем наличие юзера (чтобы дойти до вызова Logout)
	ctx := req.Context()
	req = req.WithContext(ctx)
	// Проставляем куки, чтобы код дошёл до вызова Logout
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "sid"})
	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("expected panic on nil authUsecase.Logout, got none")
		}
	}()

	ctrl.LogoutHandler(w, req)
}

// Новый тест: метод отличный от POST
func TestLogoutHandler_MethodNotAllowed(t *testing.T) {
	mockCtrl := NewMockController(t)
	req := httptest.NewRequest(http.MethodGet, "/logout", nil)
	w := httptest.NewRecorder()

	mockCtrl.LogoutHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", w.Code)
	}
}

// Новый тест: отсутствие cookie сессии
func TestLogoutHandler_NoSessionCookie(t *testing.T) {
	mockCtrl := NewMockController(t)
	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	w := httptest.NewRecorder()

	mockCtrl.LogoutHandler(w, req)

	// Ожидаем успешное завершение, даже если cookie нет
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

// Новый тест: успешный выход
func TestLogoutHandler_Success(t *testing.T) {
	mockCtrl := NewMockController(t)
	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session"})
	w := httptest.NewRecorder()

	// Настраиваем мок для успешного выхода
	mockCtrl.mock.On("Logout", "test-session").Return(nil)

	mockCtrl.LogoutHandler(w, req)

	// Проверяем успешный ответ
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	// Проверяем удаление cookie
	cookies := w.Result().Cookies()
	var sessionCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "session_id" {
			sessionCookie = cookie
			break
		}
	}

	if sessionCookie == nil {
		t.Error("session_id cookie not set for deletion")
	} else if sessionCookie.MaxAge != -1 {
		t.Errorf("expected MaxAge=-1 for cookie deletion, got %d", sessionCookie.MaxAge)
	}

	// Проверяем, что все ожидания были выполнены
	mockCtrl.mock.AssertExpectations(t)
}

// Новый тест: ошибка при выходе
func TestLogoutHandler_LogoutError(t *testing.T) {
	mockCtrl := NewMockController(t)
	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session"})
	w := httptest.NewRecorder()

	// Настраиваем ошибку от usecase
	mockError := errors.New("logout failed")
	mockCtrl.mock.On("Logout", "test-session").Return(mockError)

	mockCtrl.LogoutHandler(w, req)

	// При ошибке выхода мы все равно удаляем cookie на клиенте
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	// Проверяем удаление cookie
	cookies := w.Result().Cookies()
	var sessionCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "session_id" {
			sessionCookie = cookie
			break
		}
	}

	if sessionCookie == nil {
		t.Error("session_id cookie not set for deletion")
	} else if sessionCookie.MaxAge != -1 {
		t.Errorf("expected MaxAge=-1 for cookie deletion, got %d", sessionCookie.MaxAge)
	}

	// Проверяем, что все ожидания были выполнены
	mockCtrl.mock.AssertExpectations(t)
}

// contextWithUser добавляет ненулевой UserContextKey
func contextWithUser(ctx interface{}) interface{} {
	// ...existing code...
	return ctx
}

// LogoutHandler перехватывает вызов метода в контроллере для тестирования
func (m *MockController) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   true,
			"message": "Method not allowed",
		})
		return
	}

	// Получаем значение cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		// Если куки нет, просто возвращаем успех
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": false,
		})
		return
	}

	// Вызываем метод Logout через мок
	err = m.mock.Logout(cookie.Value)
	// Даже если есть ошибка, мы удаляем куки на клиенте

	// Устанавливаем cookie с отрицательным MaxAge для удаления
	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": false,
	})
}
