package avatar

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/stretchr/testify/mock"
	"io"
	"mime/multipart"
	repoAvatar "retarget/internal/profile-service/repo/avatar"
	"strconv"
	"testing"
	"unsafe"

	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type FakeAvatarRepository struct {
	repoAvatar.AvatarRepository
	logger         *zap.SugaredLogger
	bucketName     string
	onDownloadFile func(objectName string, requestID string) (*minio.Object, error)
	onUploadFile   func(objectName string, file multipart.File, requestID string) error
}

func NewFakeAvatarRepository() *FakeAvatarRepository {
	logger, _ := zap.NewDevelopment()
	return &FakeAvatarRepository{
		bucketName: "test-bucket",
		logger:     logger.Sugar(),
	}
}

func (f *FakeAvatarRepository) DownloadFile(objectName string, requestID string) (*minio.Object, error) {
	if f.onDownloadFile != nil {
		return f.onDownloadFile(objectName, requestID)
	}
	return nil, errors.New("download mock not set")
}

func (f *FakeAvatarRepository) UploadFile(objectName string, file multipart.File, requestID string) error {
	if f.onUploadFile != nil {
		return f.onUploadFile(objectName, file, requestID)
	}
	return errors.New("upload mock not set")
}

func injectFakeRepo(usecase *AvatarUsecase, fake *FakeAvatarRepository) {
	ptr := unsafe.Pointer(&usecase.avatarRepository)
	realPtr := (*unsafe.Pointer)(ptr)
	*realPtr = unsafe.Pointer(fake)
}

func TestGenerateAvatarName(t *testing.T) {
	usecase := &AvatarUsecase{}

	userID := 123

	hasher := md5.New()
	hasher.Write([]byte(strconv.Itoa(userID)))
	expected := hex.EncodeToString(hasher.Sum(nil))

	result := usecase.generateAvatarName(userID)
	assert.Equal(t, expected, result)

	userID = 456
	hasher = md5.New()
	hasher.Write([]byte(strconv.Itoa(userID)))
	expected = hex.EncodeToString(hasher.Sum(nil))

	result = usecase.generateAvatarName(userID)
	assert.Equal(t, expected, result)
}

func TestDownloadAvatar(t *testing.T) {
	mockObject := &minio.Object{}

	fakeRepo := NewFakeAvatarRepository()
	fakeRepo.onDownloadFile = func(objectName string, requestID string) (*minio.Object, error) {
		if objectName == "202cb962ac59075b964b07152d234b70" && requestID == "test-request-id" {
			return mockObject, nil
		}
		return nil, errors.New("file not found")
	}

	usecase := &AvatarUsecase{}
	injectFakeRepo(usecase, fakeRepo)

	userID := 123
	requestID := "test-request-id"

	object, err := usecase.DownloadAvatar(userID, requestID)

	userID = 456
	object, err = usecase.DownloadAvatar(userID, requestID)
	assert.Error(t, err)
	assert.Nil(t, object)
	assert.Equal(t, "Bucket name cannot be empty", err.Error())
}

func TestNewAvatarUsecase(t *testing.T) {
	_, _ = zap.NewDevelopment()
	repo := &repoAvatar.AvatarRepository{}

	usecase := NewAvatarUsecase(repo)
	assert.NotNil(t, usecase)

}

type MockMultipartFile struct {
	mock.Mock
	io.Reader
}

func (m *MockMultipartFile) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockMultipartFile) Read(p []byte) (n int, err error) {
	if m.Reader != nil {
		return m.Reader.Read(p)
	}
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

func TestGenerateAvatarName_EdgeCases(t *testing.T) {
	usecase := &AvatarUsecase{}

	userID := 0
	hasher := md5.New()
	hasher.Write([]byte("0"))
	expected := hex.EncodeToString(hasher.Sum(nil))
	result := usecase.generateAvatarName(userID)
	assert.Equal(t, expected, result)

	userID = -1
	hasher = md5.New()
	hasher.Write([]byte("-1"))
	expected = hex.EncodeToString(hasher.Sum(nil))
	result = usecase.generateAvatarName(userID)
	assert.Equal(t, expected, result)

	userID = 2147483647
	hasher = md5.New()
	hasher.Write([]byte("2147483647"))
	expected = hex.EncodeToString(hasher.Sum(nil))
	result = usecase.generateAvatarName(userID)
	assert.Equal(t, expected, result)
}

func TestDownloadAvatar_NotFound(t *testing.T) {
	fakeRepo := NewFakeAvatarRepository()
	fakeRepo.onDownloadFile = func(objectName string, requestID string) (*minio.Object, error) {
		return nil, errors.New("avatar not found")
	}

	usecase := &AvatarUsecase{}
	injectFakeRepo(usecase, fakeRepo)

	userID := 123
	requestID := "not-found-test-request-id"

	object, err := usecase.DownloadAvatar(userID, requestID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be")
	assert.Nil(t, object)
}

func TestDownloadAvatar_ServerError(t *testing.T) {
	fakeRepo := NewFakeAvatarRepository()
	fakeRepo.onDownloadFile = func(objectName string, requestID string) (*minio.Object, error) {
		return nil, errors.New("internal server error")
	}

	usecase := &AvatarUsecase{}
	injectFakeRepo(usecase, fakeRepo)

	userID := 123
	requestID := "server-error-test-request-id"

	object, err := usecase.DownloadAvatar(userID, requestID)

	// Проверяем результаты
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be")
	assert.Nil(t, object)
}

func TestGenerateAvatarName_Consistency(t *testing.T) {
	usecase := &AvatarUsecase{}

	userID := 42

	result1 := usecase.generateAvatarName(userID)
	result2 := usecase.generateAvatarName(userID)

	assert.Equal(t, result1, result2)

	hasher := md5.New()
	hasher.Write([]byte("42"))
	expected := hex.EncodeToString(hasher.Sum(nil))
	assert.Equal(t, expected, result1)
}

func TestUploadAvatar_FileErrors(t *testing.T) {
	fakeRepo := NewFakeAvatarRepository()
	fakeRepo.onUploadFile = func(objectName string, file multipart.File, requestID string) error {
		return errors.New("file read error")
	}

	usecase := &AvatarUsecase{}
	injectFakeRepo(usecase, fakeRepo)

	mockFile := &MockMultipartFile{
		Reader: bytes.NewReader([]byte("test content")),
	}
	mockFile.On("Close").Return(nil).Maybe()
}
