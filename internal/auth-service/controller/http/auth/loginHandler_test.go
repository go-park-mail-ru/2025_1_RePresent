package auth_test

import (
	"bytes"
	"context" // добавлено
	"net/http"
	"net/http/httptest"
	"testing"

	auth "retarget/internal/auth-service/controller/http/auth"
	"retarget/pkg/entity" // добавлено
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
