package usecase

import (
	"crypto/rand"
	// "crypto/md5"
	"encoding/hex"
	"errors"
	"mime/multipart"
	repoBannerImage "retarget/internal/banner-service/repo"

	"github.com/minio/minio-go/v7"
)

type BannerImageUsecaseInterface interface {
	generateBannerImageName() string
	DownloadBannerImage(imageID string) (*minio.Object, error)
	UploadBannerImage(file multipart.File) error
}

type BannerImageUsecase struct {
	BannerImageRepository *repoBannerImage.BannerImageRepository
}

func NewBannerImageUsecase(BannerImageRepo *repoBannerImage.BannerImageRepository) *BannerImageUsecase {
	return &BannerImageUsecase{BannerImageRepository: BannerImageRepo}
}

func (r *BannerImageUsecase) generateBannerImageName() string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	// Превращаем байты в строку
	randomString := hex.EncodeToString(bytes)
	return randomString
}

func (r *BannerImageUsecase) DownloadBannerImage(imageID string) (*minio.Object, error) {
	return r.BannerImageRepository.DownloadFile(imageID)
}

func (r *BannerImageUsecase) UploadBannerImage(file multipart.File) (string, error) {
	if file == nil {
		return "", errors.New("Uploaded file is nil")
	}
	objectName := r.generateBannerImageName()
	err := r.BannerImageRepository.UploadFile(objectName, file)
	if err != nil {
		return "", err
	}
	return objectName, nil
}
