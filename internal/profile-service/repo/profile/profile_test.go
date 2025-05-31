package profile

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	entityProfile "retarget/internal/profile-service/entity/profile"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func setupTestLogger() *zap.SugaredLogger {
	logger, _ := zap.NewDevelopment()
	return logger.Sugar()
}

func TestNewProfileRepository(t *testing.T) {
	logger := setupTestLogger()

	repo := &ProfileRepository{
		logger: logger,
	}
	assert.NotNil(t, repo)
}

func TestGetProfileByID(t *testing.T) {
	logger := setupTestLogger()
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := &ProfileRepository{
		db:     db,
		logger: logger,
	}

	requestID := "test-request-id"
	userID := 1

	rows := sqlmock.NewRows([]string{"id", "username", "email", "description", "balance", "role"}).
		AddRow(1, "testuser", "test@example.com", "test description", 1000.0, "user")

	mock.ExpectQuery("SELECT id, username, email, description, balance, role FROM auth_user WHERE").
		WithArgs(userID).
		WillReturnRows(rows)

	profile, err := repo.GetProfileByID(userID, requestID)

	mock.ExpectQuery("SELECT id, username, email, description, balance, role FROM auth_user WHERE").
		WithArgs(2).
		WillReturnError(sql.ErrNoRows)

	profile, err = repo.GetProfileByID(2, requestID)
	assert.Error(t, err)
	assert.Equal(t, entityProfile.ErrProfileNotFound, err)
	assert.Nil(t, profile)

	mock.ExpectQuery("SELECT id, username, email, description, balance, role FROM auth_user WHERE").
		WithArgs(3).
		WillReturnError(fmt.Errorf("database connection error"))

	profile, err = repo.GetProfileByID(3, requestID)
	assert.Error(t, err)
	assert.Nil(t, profile)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateProfileByID(t *testing.T) {
	logger := setupTestLogger()
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := &ProfileRepository{
		db:     db,
		logger: logger,
	}

	requestID := "test-request-id"
	userID := 1
	username := "newusername"
	description := "new description"

	mock.ExpectExec("UPDATE auth_user SET username = \\$1, description = \\$2 WHERE id = \\$3").
		WithArgs(username, description, userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.UpdateProfileByID(userID, username, description, requestID)
	assert.NoError(t, err)

	mock.ExpectExec("UPDATE auth_user SET username = \\$1, description = \\$2 WHERE id = \\$3").
		WithArgs(username, description, userID).
		WillReturnError(fmt.Errorf("database error"))

	err = repo.UpdateProfileByID(userID, username, description, requestID)
	assert.Error(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCloseConnection(t *testing.T) {
	logger := setupTestLogger()
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	repo := &ProfileRepository{
		db:     db,
		logger: logger,
	}

	mock.ExpectClose()

	err = repo.CloseConnection()
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestGetProfileByID_Timeout(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка при создании мок базы данных: %v", err)
	}
	defer db.Close()

	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()

	repo := &ProfileRepository{
		db:     db,
		logger: sugar,
	}

	userID := 1
	requestID := "timeout-test-request-id"

	mock.ExpectQuery("SELECT id, username, email, description, balance, role FROM auth_user WHERE").
		WithArgs(userID).
		WillReturnError(errors.New("timeout error: connection timed out"))

	profile, err := repo.GetProfileByID(userID, requestID)

	assert.Error(t, err)
	assert.Nil(t, profile)
	assert.Contains(t, err.Error(), "timeout error")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateProfileByID_NoRowsAffected(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка при создании мок базы данных: %v", err)
	}
	defer db.Close()

	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()

	repo := &ProfileRepository{
		db:     db,
		logger: sugar,
	}

	userID := 999
	username := "nonexistent"
	description := "this user does not exist"
	requestID := "no-rows-test-request-id"

	mock.ExpectExec("UPDATE auth_user SET username = \\$1, description = \\$2 WHERE id = \\$3").
		WithArgs(username, description, userID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.UpdateProfileByID(userID, username, description, requestID)

	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCloseConnection_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка при создании мок базы данных: %v", err)
	}

	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()

	repo := &ProfileRepository{
		db:     db,
		logger: sugar,
	}

	mock.ExpectClose().WillReturnError(errors.New("connection close error"))

	err = repo.CloseConnection()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "connection close error")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNewProfileRepository_InvalidConnectionString(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()

	t.Skip("Этот тест будет вызывать log.Fatal, поэтому он пропущен")

	invalidEndpoint := "invalid connection string"
	_ = NewProfileRepository(invalidEndpoint, sugar)
	t.Fail()
}

func TestUpdateProfileByID_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка при создании мок базы данных: %v", err)
	}
	defer db.Close()

	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()

	repo := &ProfileRepository{
		db:     db,
		logger: sugar,
	}

	userID := 1
	username := "testuser"
	description := "test description"
	requestID := "error-test-request-id"

	mock.ExpectExec("UPDATE auth_user SET username = \\$1, description = \\$2 WHERE id = \\$3").
		WithArgs(username, description, userID).
		WillReturnError(sql.ErrConnDone)

	err = repo.UpdateProfileByID(userID, username, description, requestID)

	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetProfileByID_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка при создании мок базы данных: %v", err)
	}
	defer db.Close()

	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()

	repo := &ProfileRepository{
		db:     db,
		logger: sugar,
	}

	userID := 1
	requestID := "scan-error-test-request-id"

	rows := sqlmock.NewRows([]string{"id", "username", "email", "description", "balance", "role"}).
		AddRow("not-an-integer", "testuser", "test@example.com", "test description", "not-a-decimal", "user")

	mock.ExpectQuery("SELECT id, username, email, description, balance, role FROM auth_user WHERE").
		WithArgs(userID).
		WillReturnRows(rows)

	profile, err := repo.GetProfileByID(userID, requestID)

	assert.Error(t, err)
	assert.Nil(t, profile)
	assert.Contains(t, err.Error(), "invalid")

	assert.NoError(t, mock.ExpectationsWereMet())
}
