package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	auth "retarget/internal/auth-service/controller/http/auth"
	usecaseAuth "retarget/internal/auth-service/usecase/auth"
)

func TestSetupAuthRoutes_LogoutMethodNotAllowed(t *testing.T) {
	router := auth.SetupAuthRoutes(nil, &usecaseAuth.AuthUsecase{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/logout", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	if w.Code != 401 {
		t.Errorf("expected 401, got %d", w.Code)
	}
}
