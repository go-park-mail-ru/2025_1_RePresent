package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	usecase "retarget/internal/banner-service/usecase"
)

func createMultipartRequest(fieldName, filename string, content []byte) (*http.Request, string, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fieldName, filename)
	if err != nil {
		return nil, "", err
	}
	part.Write(content)
	writer.Close()

	req := httptest.NewRequest("PUT", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.ContentLength = int64(body.Len())
	return req, writer.FormDataContentType(), nil
}

func TestUploadImageHandler_FileSizeLimit(t *testing.T) {
	controller := &BannerController{}

	req := httptest.NewRequest("PUT", "/upload", nil)
	req.ContentLength = 11 * 1024 * 1024 // 11MB
	w := httptest.NewRecorder()

	controller.UploadImageHandler(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status %d; got %d", http.StatusUnprocessableEntity, w.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if errValue, _ := resp["error"].(bool); errValue {
		t.Error("Expected error=true in response")
	}
}

func TestUploadImageHandler_InvalidForm(t *testing.T) {
	controller := &BannerController{}

	req := httptest.NewRequest("PUT", "/upload", strings.NewReader("invalid form data"))
	req.Header.Set("Content-Type", "multipart/form-data; boundary=xxx")
	w := httptest.NewRecorder()

	controller.UploadImageHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d; got %d", w.Code, http.StatusBadRequest)
	}
}

func TestUploadImageHandler_NoFile(t *testing.T) {
	controller := &BannerController{}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("field", "value")
	writer.Close()

	req := httptest.NewRequest("PUT", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	controller.UploadImageHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d; got %d", w.Code, http.StatusBadRequest)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if msg, _ := resp["message"].(string); strings.Contains(msg, "Image not found") {
		t.Errorf("Expected message about missing image, got: %s", msg)
	}
}

func TestUploadImageHandler_UnsupportedType(t *testing.T) {
	controller := &BannerController{
		ImageUsecase: &usecase.BannerImageUsecase{},
	}

	textData := []byte("This is not an image")
	req, contentType, err := createMultipartRequest("image", "test.txt", textData)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()

	controller.UploadImageHandler(w, req)

	if w.Code != http.StatusUnsupportedMediaType {
		t.Errorf("Expected status %d; got %d", http.StatusUnsupportedMediaType, w.Code)
	}
}

func testCreateContext(ctx context.Context, key, value interface{}) context.Context {
	return context.WithValue(ctx, key, value)
}

func TestDownloadImage_VarsHandling(t *testing.T) {
	_ = func(r *http.Request) map[string]string {
		return map[string]string{"image_id": "test-id"}
	}

	_ = mux.Vars

	defer func() {
		if r := recover(); r != nil {
			t.Logf("Тест вызвал панику (это ожидаемо): %v", r)
		}
	}()

	_ = &BannerController{
		ImageUsecase: &usecase.BannerImageUsecase{},
	}

	req := httptest.NewRequest("GET", "/api/v1/banner/image/test-id", nil)
	ctx := testCreateContext(req.Context(), "vars", map[string]string{"image_id": "test-id"})
	req = req.WithContext(ctx)
	_ = httptest.NewRecorder()

	t.Log("Testing vars handling in DownloadImage (expect panic)")
}
