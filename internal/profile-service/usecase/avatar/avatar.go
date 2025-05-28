package avatar

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"mime/multipart"
	repoAvatar "retarget/internal/profile-service/repo/avatar"
	"strconv"

	"github.com/minio/minio-go/v7"
)

type AvatarUsecaseInterface interface {
	generateAvatarName(id int) string
	DownloadAvatar(userID int, requestID string) (*minio.Object, error)
	UploadAvatar(userID int, file multipart.File, requestID string) error
}

type AvatarUsecase struct {
	avatarRepository *repoAvatar.AvatarRepository
}

func NewAvatarUsecase(avatarRepo *repoAvatar.AvatarRepository) *AvatarUsecase {
	return &AvatarUsecase{avatarRepository: avatarRepo}
}

func (r *AvatarUsecase) generateAvatarName(id int) string {
	hasher := md5.New()
	hasher.Write([]byte(strconv.Itoa(id)))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (r *AvatarUsecase) UploadAvatar(userID int, file multipart.File, requestID string) error {
	if file == nil {
		return errors.New("Uploaded file is nil")
	}
	objectName := r.generateAvatarName(userID)
	err := r.avatarRepository.UploadFile(objectName, file, requestID)
	if err != nil {
		return err
	}
	return nil
}

func (r *AvatarUsecase) DownloadAvatar(userID int, requestID string) (*minio.Object, error) {
	objectName := r.generateAvatarName(userID)
	return r.avatarRepository.DownloadFile(objectName, requestID)
}
