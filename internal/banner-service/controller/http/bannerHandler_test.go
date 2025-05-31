package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	usecase "retarget/internal/banner-service/usecase"
	response "retarget/pkg/entity"
)

func muxSetVars(r *http.Request, vars map[string]string) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, "vars", vars)
	ctx = context.WithValue(ctx, response.СtxKeyRequestID{}, "test-request-id")
	ctx = context.WithValue(ctx, response.UserContextKey, response.UserContext{UserID: 1})
	return r.WithContext(ctx)
}

func TestReadBanner_InvalidID(t *testing.T) {
	h := &BannerController{BannerUsecase: &usecase.BannerUsecase{}}
	req := httptest.NewRequest("GET", "/api/v1/banner/abc", nil)
	req = muxSetVars(req, map[string]string{"banner_id": "abc"})
	w := httptest.NewRecorder()

	h.ReadBanner(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("got %d, want %d", w.Code, http.StatusInternalServerError)
	}
	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	if errVal, _ := resp["error"].(bool); errVal {
		t.Error("expected error=true, got false")
	}
}

func TestCreateBanner_InvalidJSON(t *testing.T) {
	h := &BannerController{BannerUsecase: &usecase.BannerUsecase{}}
	req := httptest.NewRequest("POST", "/api/v1/banner/create", strings.NewReader("invalid json"))
	req = muxSetVars(req, nil)
	w := httptest.NewRecorder()

	h.CreateBanner(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("got %d, want %d", w.Code, http.StatusUnprocessableEntity)
	}
}

func TestCreateBanner_EmptyBody(t *testing.T) {
	h := &BannerController{BannerUsecase: &usecase.BannerUsecase{}}
	req := httptest.NewRequest("POST", "/api/v1/banner/create", bytes.NewReader([]byte("{}")))
	req = muxSetVars(req, nil)
	w := httptest.NewRecorder()

	h.CreateBanner(w, req)

	// Из-за валидации структуры ожидается ошибка
	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d, want %d", w.Code, http.StatusBadRequest)
	}

	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	if errVal, _ := resp["error"].(bool); errVal {
		t.Error("expected error=true, got false")
	}
}

func TestUpdateBanner_InvalidJSON(t *testing.T) {
	h := &BannerController{BannerUsecase: &usecase.BannerUsecase{}}
	req := httptest.NewRequest("PUT", "/api/v1/banner/1", strings.NewReader("invalid json"))
	req = muxSetVars(req, map[string]string{"banner_id": "1"})
	w := httptest.NewRecorder()

	h.UpdateBanner(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("got %d, want %d", w.Code, http.StatusUnprocessableEntity)
	}
}

func TestUpdateBanner_InvalidID(t *testing.T) {
	h := &BannerController{BannerUsecase: &usecase.BannerUsecase{}}
	req := httptest.NewRequest("PUT", "/api/v1/banner/invalid", strings.NewReader("{}"))
	req = muxSetVars(req, map[string]string{"banner_id": "invalid"})
	w := httptest.NewRecorder()

	h.UpdateBanner(w, req)

	if w.Code != 400 {
		t.Errorf("got %d, want %d", w.Code, 400)
	}
}

func TestDeleteBanner_InvalidID(t *testing.T) {
	h := &BannerController{BannerUsecase: &usecase.BannerUsecase{}}
	req := httptest.NewRequest("DELETE", "/api/v1/banner/invalid", nil)
	req = muxSetVars(req, map[string]string{"banner_id": "invalid"})
	w := httptest.NewRecorder()

	h.DeleteBanner(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("got %d, want %d", w.Code, http.StatusInternalServerError)
	}
}
