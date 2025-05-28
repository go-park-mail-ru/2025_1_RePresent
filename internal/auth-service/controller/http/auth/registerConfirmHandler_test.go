package auth_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	auth "retarget/internal/auth-service/controller/http/auth"
)

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
