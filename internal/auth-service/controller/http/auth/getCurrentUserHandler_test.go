package auth_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	auth "retarget/internal/auth-service/controller/http/auth"
	usecaseAuth "retarget/internal/auth-service/usecase/auth"
	entity "retarget/pkg/entity"
)

func TestGetCurrentUserHandler_MethodNotAllowed(t *testing.T) {
	ctrl := auth.NewAuthController(&usecaseAuth.AuthUsecase{})
	req := httptest.NewRequest(http.MethodPost, "/me", nil)
	w := httptest.NewRecorder()

	ctrl.GetCurrentUserHandler(w, req)
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
	ctrl := auth.NewAuthController(&usecaseAuth.AuthUsecase{})
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	// только requestID
	ctx := context.WithValue(req.Context(), entity.СtxKeyRequestID{}, "test")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	ctrl.GetCurrentUserHandler(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
	var resp struct {
		Error bool `json:"error"`
	}
	_ = json.NewDecoder(w.Body).Decode(&resp)
	if resp.Error {
		t.Errorf("expected error=true on NoContext, got %+v", resp)
	}
}

func TestGetCurrentUserHandler_PanicOnNilUsecase_AfterContext(t *testing.T) {
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
