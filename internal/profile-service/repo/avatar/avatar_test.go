package avatar

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockMinioClient struct {
	mock.Mock
}

func (m *MockMinioClient) MakeBucket(ctx context.Context, bucketName string, opts minio.MakeBucketOptions) error {
	args := m.Called(ctx, bucketName, opts)
	return args.Error(0)
}

func (m *MockMinioClient) GetObject(ctx context.Context, bucketName, objectName string, opts minio.GetObjectOptions) (*minio.Object, error) {
	args := m.Called(ctx, bucketName, objectName, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*minio.Object), args.Error(1)
}

func (m *MockMinioClient) PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, size int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
	args := m.Called(ctx, bucketName, objectName, reader, size, opts)
	return args.Get(0).(minio.UploadInfo), args.Error(1)
}

type TestAvatarRepository struct {
	bucketName string
	logger     *zap.SugaredLogger
	mockClient *MockMinioClient
}

func (r *TestAvatarRepository) createBacket(bucketName string) error {
	err := r.mockClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	if err != nil {
		if err, ok := err.(minio.ErrorResponse); ok && err.Code == "BucketAlreadyOwnedByYou" {
			return nil
		}
		return err
	}
	return nil
}

func (r *TestAvatarRepository) DownloadFile(objectName string, requestID string) (*minio.Object, error) {
	object, err := r.mockClient.GetObject(context.Background(), r.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		if errResp, ok := err.(minio.ErrorResponse); ok && errResp.Code == "NoSuchKey" {
			return nil, errors.New("File not found")
		}
		return nil, err
	}
	return object, nil
}

func (r *TestAvatarRepository) UploadFile(objectName string, file multipart.File, requestID string) error {
	_, err := r.mockClient.PutObject(context.Background(), r.bucketName, objectName, file, -1, minio.PutObjectOptions{})
	return err
}

func setupTestLogger() *zap.SugaredLogger {
	logger, _ := zap.NewDevelopment()
	return logger.Sugar()
}

func TestCreateBacket(t *testing.T) {
	logger := setupTestLogger()
	mockMinioClient := new(MockMinioClient)

	repo := &TestAvatarRepository{
		bucketName: "test-bucket",
		logger:     logger,
		mockClient: mockMinioClient,
	}

	bucketName := "test-bucket"

	mockMinioClient.On("MakeBucket", mock.Anything, bucketName, mock.Anything).Return(nil).Once()

	err := repo.createBacket(bucketName)
	assert.NoError(t, err)

	errorResponse := minio.ErrorResponse{Code: "BucketAlreadyOwnedByYou"}
	mockMinioClient.On("MakeBucket", mock.Anything, bucketName, mock.Anything).Return(errorResponse).Once()

	err = repo.createBacket(bucketName)
	assert.NoError(t, err)

	mockMinioClient.On("MakeBucket", mock.Anything, bucketName, mock.Anything).Return(assert.AnError).Once()

	err = repo.createBacket(bucketName)
	assert.Error(t, err)

	mockMinioClient.AssertExpectations(t)
}

func TestDownloadFile(t *testing.T) {
	logger := setupTestLogger()
	mockMinioClient := new(MockMinioClient)

	bucketName := "test-bucket"
	repo := &TestAvatarRepository{
		bucketName: bucketName,
		logger:     logger,
		mockClient: mockMinioClient,
	}

	objectName := "test-object"
	requestID := "test-request-id"

	mockObject := &minio.Object{}
	mockMinioClient.On("GetObject", mock.Anything, bucketName, objectName, mock.Anything).Return(mockObject, nil).Once()

	object, err := repo.DownloadFile(objectName, requestID)
	assert.NoError(t, err)
	assert.Equal(t, mockObject, object)

	errorResponse := minio.ErrorResponse{Code: "NoSuchKey"}
	mockMinioClient.On("GetObject", mock.Anything, bucketName, objectName, mock.Anything).Return(nil, errorResponse).Once()

	object, err = repo.DownloadFile(objectName, requestID)
	assert.Error(t, err)
	assert.Nil(t, object)
	assert.Equal(t, "File not found", err.Error())

	mockMinioClient.On("GetObject", mock.Anything, bucketName, objectName, mock.Anything).Return(nil, assert.AnError).Once()

	object, err = repo.DownloadFile(objectName, requestID)
	assert.Error(t, err)
	assert.Nil(t, object)

	mockMinioClient.AssertExpectations(t)
}

func TestNewAvatarRepository(t *testing.T) {
	logger := setupTestLogger()
	repo := &AvatarRepository{
		bucketName: "test-bucket",
		logger:     logger,
	}
	assert.NotNil(t, repo)
}

func TestCreateBacket_DetailedErrors(t *testing.T) {
	logger := setupTestLogger()
	mockMinioClient := new(MockMinioClient)

	repo := &TestAvatarRepository{
		bucketName: "test-bucket",
		logger:     logger,
		mockClient: mockMinioClient,
	}

	bucketName := "test-bucket"

	mockMinioClient.On("MakeBucket", mock.Anything, bucketName, mock.Anything).
		Return(minio.ErrorResponse{
			Code:    "AccessDenied",
			Message: "Access Denied",
		}).Once()

	err := repo.createBacket(bucketName)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Access Denied")

	mockMinioClient.On("MakeBucket", mock.Anything, bucketName, mock.Anything).
		Return(errors.New("network error: connection refused")).Once()

	err = repo.createBacket(bucketName)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "network error")

	invalidBucketName := "inv@lid"
	mockMinioClient.On("MakeBucket", mock.Anything, invalidBucketName, mock.Anything).
		Return(minio.ErrorResponse{
			Code:    "InvalidBucketName",
			Message: "The specified bucket is not valid",
		}).Once()

	err = repo.createBacket(invalidBucketName)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bucket is not")

	mockMinioClient.AssertExpectations(t)
}

type MockMultipartFile struct {
	mock.Mock
	io.Reader
}

func (m *MockMultipartFile) Read(p []byte) (n int, err error) {
	if m.Reader != nil {
		return m.Reader.Read(p)
	}
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

func (m *MockMultipartFile) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockMultipartFile) ReadAt(p []byte, off int64) (n int, err error) {
	args := m.Called(p, off)
	return args.Int(0), args.Error(1)
}

func (m *MockMultipartFile) Seek(offset int64, whence int) (int64, error) {
	args := m.Called(offset, whence)
	return args.Get(0).(int64), args.Error(1)
}

func TestUploadFile_ReadError(t *testing.T) {
	logger := setupTestLogger()
	mockMinioClient := new(MockMinioClient)

	bucketName := "test-bucket"
	repo := &TestAvatarRepository{
		bucketName: bucketName,
		logger:     logger,
		mockClient: mockMinioClient,
	}

	objectName := "test-object"
	requestID := "read-error-test-request-id"

	mockFile := &MockMultipartFile{
		Reader: bytes.NewReader([]byte("test content")),
	}
	mockFile.On("Close").Return(nil).Maybe()
	mockFile.On("ReadAt", mock.Anything, mock.Anything).Return(0, io.EOF).Maybe()
	mockFile.On("Seek", mock.Anything, mock.Anything).Return(int64(0), nil).Maybe()

	var mpFile multipart.File = mockFile

	mockMinioClient.On("PutObject", mock.Anything, bucketName, objectName, mock.Anything, int64(-1), mock.Anything).
		Return(minio.UploadInfo{}, errors.New("failed to read file")).Once()

	err := repo.UploadFile(objectName, mpFile, requestID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read file")

	mockMinioClient.AssertExpectations(t)
	mockFile.AssertExpectations(t)
}

func TestDownloadFile_DetailedErrors(t *testing.T) {
	logger := setupTestLogger()
	mockMinioClient := new(MockMinioClient)

	bucketName := "test-bucket"
	_ = &TestAvatarRepository{
		bucketName: bucketName,
		logger:     logger,
		mockClient: mockMinioClient,
	}
}

func TestDownloadFile_EmptyObjectName(t *testing.T) {
	logger := setupTestLogger()
	mockMinioClient := new(MockMinioClient)

	bucketName := "test-bucket"
	_ = &TestAvatarRepository{
		bucketName: bucketName,
		logger:     logger,
		mockClient: mockMinioClient,
	}
}

func TestUploadFile_ContextCancellation(t *testing.T) {
	logger := setupTestLogger()
	mockMinioClient := new(MockMinioClient)

	bucketName := "test-bucket"
	_ = &TestAvatarRepository{
		bucketName: bucketName,
		logger:     logger,
		mockClient: mockMinioClient,
	}
}
