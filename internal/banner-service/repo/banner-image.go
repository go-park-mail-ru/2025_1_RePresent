package repo

import (
	"context"
	"errors"
	"log"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type BannerImageRepositoryInterface interface {
	createBacket(bucketName string) error

	DownloadFile(objectName string) (*minio.Object, error)
	UploadFile(objectName string, file multipart.File) error
}

type BannerImageRepository struct {
	minioClient *minio.Client
	bucketName  string
}

func NewBannerImageRepository(endpoint, accessKeyID, secretAccessKey, token string, useSSL bool, bucketName string) *BannerImageRepository {
	bannerRepo := &BannerImageRepository{bucketName: bucketName}
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, token),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatal(err)
	}
	bannerRepo.minioClient = minioClient

	err = bannerRepo.createBacket(bucketName)
	if err != nil {
		log.Fatal(err)
	}

	return bannerRepo
}

func (r *BannerImageRepository) createBacket(bucketName string) error {
	err := r.minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	if err != nil {
		if err, ok := err.(minio.ErrorResponse); ok && err.Code == "BucketAlreadyOwnedByYou" {
			return nil
		}
		return err
	}
	return nil
}

func (r *BannerImageRepository) DownloadFile(objectName string) (*minio.Object, error) {
	object, err := r.minioClient.GetObject(context.Background(), r.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		if errResp, ok := err.(minio.ErrorResponse); ok && errResp.Code == "NoSuchKey" {
			return nil, errors.New("File not found")
		}
		return nil, err
	}
	return object, nil
}

func (r *BannerImageRepository) UploadFile(objectName string, file multipart.File) error {
	_, err := r.minioClient.PutObject(context.Background(), r.bucketName, objectName, file, -1, minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}
