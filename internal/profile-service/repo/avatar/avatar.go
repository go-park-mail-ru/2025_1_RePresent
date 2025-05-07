package avatar

import (
	"context"
	"errors"
	"log"
	"mime/multipart"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

type AvatarRepositoryInterface interface {
	createBacket(bucketName string) error

	DownloadFile(objectName string) (*minio.Object, error)
	UploadFile(objectName string, file multipart.File) error
}

type AvatarRepository struct {
	minioClient *minio.Client
	bucketName  string
	logger      *zap.SugaredLogger
}

func NewAvatarRepository(endpoint, accessKeyID, secretAccessKey, token string, useSSL bool, bucketName string, logger *zap.SugaredLogger) *AvatarRepository {
	avatarRepo := &AvatarRepository{bucketName: bucketName}
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, token),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatal(err)
	}
	avatarRepo.minioClient = minioClient

	err = avatarRepo.createBacket(bucketName)
	if err != nil {
		log.Fatal(err)
	}
	avatarRepo.logger = logger
	return avatarRepo
}

func (r *AvatarRepository) createBacket(bucketName string) error {
	err := r.minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	if err != nil {
		if err, ok := err.(minio.ErrorResponse); ok && err.Code == "BucketAlreadyOwnedByYou" {
			return nil
		}
		return err
	}
	return nil
}

func (r *AvatarRepository) DownloadFile(objectName string, requestID string) (*minio.Object, error) {
	object, err := r.minioClient.GetObject(context.Background(), r.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		if errResp, ok := err.(minio.ErrorResponse); ok && errResp.Code == "NoSuchKey" {
			return nil, errors.New("File not found")
		}
		return nil, err
	}
	return object, nil
}

func (r *AvatarRepository) UploadFile(objectName string, file multipart.File, requestID string) error {
	/*
		_, err := r.minioClient.StatObject(context.Background(), r.bucketName, objectName, minio.StatObjectOptions{})
		if err != nil {
			if errResp, ok := err.(minio.ErrorResponse); ok && errResp.Code == "NoSuchKey" {
				// Объект не найден, продолжаем
			} else {
				return err
			}
		} else {
			err := r.minioClient.RemoveObject(context.Background(), r.bucketName, objectName, minio.RemoveObjectOptions{})
			if err != nil {
				return err
			}
		}

		_, err = r.minioClient.PutObject(context.Background(), r.bucketName, objectName, file, -1, minio.PutObjectOptions{})
		if err != nil {
			return err
		}

		return nil
	*/
	startTime := time.Now()
	r.logger.Debugw("Put avatar in storage",
		"request_id", requestID,
		"bucket", r.bucketName,
		"object", objectName,
	)
	_, err := r.minioClient.PutObject(context.Background(), r.bucketName, objectName, file, -1, minio.PutObjectOptions{})
	if err != nil {
		r.logger.Debugw("MINIO Put failed",
			"request_id", requestID,
			"bucket", r.bucketName,
			"object", objectName,
			"error", err.Error(),
			"timeTakenMs", time.Since(startTime).Milliseconds(),
		)
		return err
	}
	r.logger.Debugw("MINIO Put success",
		"request_id", requestID,
		"bucket", r.bucketName,
		"object", objectName,
		"timeTakenMs", time.Since(startTime).Milliseconds(),
	)
	return nil
}
