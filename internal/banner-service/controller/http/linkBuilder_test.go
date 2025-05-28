package controller

import (
	"testing"

	"github.com/gorilla/mux"
)

func TestBannerImageURL_Success(t *testing.T) {
	router := mux.NewRouter()
	router.Handle("/api/v1/banner/image/{image_id}", nil).Name("download_image")
	b := NewLinkBuilder(router)

	url, err := b.BannerImageURL("abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := "/api/v1/banner/image/abc"; url != want {
		t.Errorf("got %q, want %q", url, want)
	}
}

func TestBannerImageURL_WrongPlaceholder(t *testing.T) {
	router := mux.NewRouter()
	router.Handle("/api/v1/banner/image", nil).Name("download_image")
	b := NewLinkBuilder(router)
	if _, err := b.BannerImageURL("1"); err != nil {
		t.Error("expected error for missing placeholder, got nil")
	}
}

func TestNewLinkBuilder(t *testing.T) {
	router := mux.NewRouter()
	lb := NewLinkBuilder(router)
	if lb == nil {
		t.Fatal("expected non-nil LinkBuilder")
	}
}

func TestBannerImageURL_EdgeCases(t *testing.T) {
	router := mux.NewRouter()
	router.Handle("/api/v1/banner/image/{image_id}", nil).Name("download_image")
	b := NewLinkBuilder(router)

	url, err := b.BannerImageURL("")
	if want := ""; url != want {
		t.Errorf("empty ID: got %q, want %q", url, want)
	}

	url, err = b.BannerImageURL("test/with spaces")
	if err == nil {
		t.Fatalf("unexpected error for special chars: %v", err)
	}
	if url != "" {
		t.Error("got empty URL for special chars")
	}
}
