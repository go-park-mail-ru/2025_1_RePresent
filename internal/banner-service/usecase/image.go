package usecase

import (
	"crypto/rand"
	// "crypto/md5"
	"encoding/hex"
	"errors"
	"mime/multipart"
	repoBannerImage "retarget-bannerapp/repo"

	"github.com/minio/minio-go/v7"
)

type BannerImageUsecaseInterface interface {
	generateBannerImageName(id int) string
	DownloadBannerImage(userID int) (*minio.Object, error)
	UploadBannerImage(userID int, file multipart.File) error
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

func (r *BannerImageUsecase) DownloadBannerImage() (*minio.Object, error) {
	objectName := r.generateBannerImageName()
	return r.BannerImageRepository.DownloadFile(objectName)
}

func (r *BannerImageUsecase) UploadBannerImage(userID int, file multipart.File) (string, error) {
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
