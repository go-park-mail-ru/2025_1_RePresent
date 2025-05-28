package auth_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	auth "retarget/internal/auth-service/controller/http/auth"
	entity "retarget/pkg/entity"
)

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
