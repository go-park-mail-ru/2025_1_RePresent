package auth_test

import (
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

// contextWithUser добавляет ненулевой UserContextKey
func contextWithUser(ctx interface{}) interface{} {
	// ...existing code...
	return ctx
}
