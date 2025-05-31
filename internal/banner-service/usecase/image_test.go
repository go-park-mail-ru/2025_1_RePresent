package usecase

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type nopCloser struct {
	io.ReadSeeker
	io.ReaderAt
}

func (nopCloser) Close() error { return nil }

func TestGenerateBannerImageName(t *testing.T) {
	uc := NewBannerImageUsecase(nil)
	n1 := uc.generateBannerImageName()
	n2 := uc.generateBannerImageName()
	assert.Len(t, n1, 32)
	assert.Len(t, n2, 32)
	assert.NotEqual(t, n1, n2)
}

func TestUploadBannerImage_NilFile(t *testing.T) {
	uc := NewBannerImageUsecase(nil)
	name, err := uc.UploadBannerImage(nil)
	assert.Empty(t, name)
	assert.EqualError(t, err, "Uploaded file is nil")
}
