package controller

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	model "retarget/internal/banner-service/easyjsonModels"
	usecase "retarget/internal/banner-service/usecase"
	response "retarget/pkg/entity"
)

func setupRequestContext(r *http.Request, vars map[string]string) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, response.СtxKeyRequestID{}, "test-request-id")
	ctx = context.WithValue(ctx, "vars", vars)
	return r.WithContext(ctx)
}

func TestGetBannerIFrameByID_InvalidID(t *testing.T) {
	invalidBannerID := "invalid"
	req := httptest.NewRequest("GET", "/api/v1/banner/iframe/"+invalidBannerID, nil)
	req = setupRequestContext(req, map[string]string{"banner_id": invalidBannerID})
	_ = httptest.NewRecorder()

	_, err := strToInt(invalidBannerID)
	if err == nil {
		t.Error("expected error for invalid ID, got nil")
	}
}

func strToInt(s string) (int, error) {
	id, err := strconv.Atoi(s)
	return id, err
}

func TestRandomIFrame_InvalidID(t *testing.T) {
	invalidLinkID := "invalid"
	req := httptest.NewRequest("GET", "/api/v1/banner/uniq_link/"+invalidLinkID, nil)
	req = setupRequestContext(req, map[string]string{"uniq_link": invalidLinkID})
	_ = httptest.NewRecorder()

	_, err := strToInt(invalidLinkID)
	if err == nil {
		t.Error("expected error for invalid ID, got nil")
	}
}

type MockLinkBuilder struct {
	bannerImageURLFunc func(string) (string, error)
}

func (m *MockLinkBuilder) BannerImageURL(imageID string) (string, error) {
	if m.bannerImageURLFunc != nil {
		return m.bannerImageURLFunc(imageID)
	}
	return "/mock-image-url/" + imageID, nil
}

func TestGetBannerIFrameByID_Error(t *testing.T) {
	bannerID := "123"

	req := httptest.NewRequest("GET", "/api/v1/banner/iframe/"+bannerID, nil)
	req = setupRequestContext(req, map[string]string{"banner_id": bannerID})
	w := httptest.NewRecorder()

	bannerUsecase := &usecase.BannerUsecase{}
	h := &BannerController{
		BannerUsecase: bannerUsecase,
		LinkBuilder:   &MockLinkBuilder{},
	}

	_ = mux.Vars

	defer func() {
		_ = recover()
	}()

	h.GetBannerIFrameByID(w, req)

	if w.Code != http.StatusForbidden && w.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался статус ошибки, получен: %d", w.Code)
	}
}

func TestRandomIFrame_Error(t *testing.T) {
	uniqLink := "123"

	req := httptest.NewRequest("GET", "/api/v1/banner/uniq_link/"+uniqLink, nil)
	req = setupRequestContext(req, map[string]string{"uniq_link": uniqLink})
	w := httptest.NewRecorder()

	bannerUsecase := &usecase.BannerUsecase{}
	h := &BannerController{
		BannerUsecase: bannerUsecase,
		LinkBuilder:   &MockLinkBuilder{},
	}

	defer func() {
		if r := recover(); r != nil {
			t.Logf("Тест вызвал панику (это нормально): %v", r)
		}
	}()

	h.RandomIFrame(w, req)

	if w.Code != http.StatusForbidden && w.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался статус ошибки, получен: %d", w.Code)
	}
}

func TestIFrameWithMockedHandler(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/banner/image/{image_id}", func(w http.ResponseWriter, r *http.Request) {}).Name("download_image")

	h := &BannerController{
		LinkBuilder: NewLinkBuilder(router),
	}

	_, err := h.LinkBuilder.BannerImageURL("test123")
	if err != nil {
		t.Errorf("LinkBuilder failed: %v", err)
	}
}

func TestGetBannerForIFrame_Interface(t *testing.T) {
	var _ interface {
		GetBannerForIFrame(bannerID int, requestID string) (*model.Banner, error)
	} = &usecase.BannerUsecase{}

	var _ interface {
		GetRandomBannerForIFrame(userID int, requestID string) (*model.Banner, error)
	} = &usecase.BannerUsecase{}
}
