package repo

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"go.uber.org/zap"
	"gopkg.in/inf.v0"

	authEntity "retarget/internal/auth-service/entity/auth"
	optiLog "retarget/pkg/utils/optiLog"
)

func setupRepo(t *testing.T) (*AuthRepository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	logger := zap.NewNop().Sugar()
	asyncLogger := optiLog.NewAsyncLogger(logger, 1, 1)
	return &AuthRepository{db: db, asyncLogger: asyncLogger}, mock
}

func TestGetUserByID_Success(t *testing.T) {
	repo, mock := setupRepo(t)
	rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "description", "balance", "role"}).
		AddRow(1, "user1", "u1@example.com", []byte("pass"), "desc", "100.0", 2)
	mock.ExpectQuery("SELECT id, username, email, password, description, balance, role FROM auth_user WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(rows)

	user, err := repo.GetUserByID(1, "req1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.ID != 1 || user.Username != "user1" {
		t.Errorf("unexpected user %+v", user)
	}
}

func TestGetUserByID_NotFound(t *testing.T) {
	repo, mock := setupRepo(t)
	mock.ExpectQuery("SELECT id, username, email, password, description, balance, role FROM auth_user WHERE id = \\$1").
		WithArgs(2).
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetUserByID(2, "req2")
	if !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("expected ErrNoRows, got %v", err)
	}
}

func TestGetUserByEmail_Success(t *testing.T) {
	repo, mock := setupRepo(t)
	rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "description", "balance", "role"}).
		AddRow(3, "user3", "u3@example.com", []byte("pass3"), "", "0.0", 1)
	mock.ExpectQuery("WHERE email = \\$1").
		WithArgs("u3@example.com").
		WillReturnRows(rows)

	user, err := repo.GetUserByEmail("u3@example.com", "req3")
	if err != nil || user.Email != "u3@example.com" {
		t.Fatalf("expected email u3@example.com, got %v, err %v", user, err)
	}
}

func TestGetUserByUsername_Success(t *testing.T) {
	repo, mock := setupRepo(t)
	rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "description", "balance", "role"}).
		AddRow(4, "user4", "u4@example.com", []byte("pass4"), "", "0.0", 1)
	mock.ExpectQuery("WHERE username = \\$1").
		WithArgs("user4").
		WillReturnRows(rows)

	user, err := repo.GetUserByUsername("user4", "req4")
	if err != nil || user.Username != "user4" {
		t.Fatalf("expected username user4, got %v, err %v", user, err)
	}
}

func TestCheckEmailOrUsernameExists(t *testing.T) {
	repo, mock := setupRepo(t)
	// Found
	rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "description", "balance", "role"}).
		AddRow(5, "user5", "u5@example.com", []byte("pass5"), "", "0.0", 2)
	mock.ExpectQuery("SELECT.*FROM auth_user").
		WithArgs("u5@example.com", "user5").
		WillReturnRows(rows)

	user, err := repo.CheckEmailOrUsernameExists("u5@example.com", "user5", "req5")
	if err != nil || user == nil || user.ID != 5 {
		t.Fatalf("expected user ID 5, got %v, err %v", user, err)
	}

	// Not found
	mock.ExpectQuery("SELECT.*FROM auth_user").
		WithArgs("no@no", "noname").
		WillReturnError(sql.ErrNoRows)

	user, err = repo.CheckEmailOrUsernameExists("no@no", "noname", "req6")
	if err != nil || user != nil {
		t.Errorf("expected nil user, got %v, err %v", user, err)
	}
}

func TestCreateNewUser(t *testing.T) {
	repo, mock := setupRepo(t)

	// Validation error
	invalidUser := &authEntity.User{Username: "u", Email: "bad", Password: []byte("short"), Role: 0}
	if err := repo.CreateNewUser(invalidUser, "reqv"); err == nil {
		t.Error("expected validation error")
	}

	// DB error
	validUser := &authEntity.User{
		Username:    "validuser",
		Email:       "v@example.com",
		Password:    []byte("password"),
		Description: "",
		Balance:     authEntity.Decimal{Dec: inf.NewDec(10, 0)},
		Role:        1,
	}
	mock.ExpectQuery("INSERT INTO auth_user").
		WithArgs(validUser.Username, validUser.Email, validUser.Password, validUser.Description, validUser.Balance.String(), validUser.Role).
		WillReturnError(errors.New("db fail"))

	if err := repo.CreateNewUser(validUser, "reqd"); err == nil {
		t.Error("expected DB error")
	}

	// Success
	mock.ExpectQuery("INSERT INTO auth_user").
		WithArgs(validUser.Username, validUser.Email, validUser.Password, validUser.Description, validUser.Balance.String(), validUser.Role).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10))

	if err := repo.CreateNewUser(validUser, "reqs"); err != nil {
		t.Errorf("expected success, got %v", err)
	}
	if validUser.ID != 10 {
		t.Errorf("expected ID set to 10, got %d", validUser.ID)
	}
}

func TestIsTransientErrorAndWithRetry(t *testing.T) {
	// isTransientError
	if !isTransientError(context.DeadlineExceeded) {
		t.Error("expected DeadlineExceeded to be transient")
	}
	if isTransientError(errors.New("permanent")) {
		t.Error("did not expect permanent error to be transient")
	}

	// withRetry
	attempts := 0
	err := withRetry(func() error {
		attempts++
		if attempts < 2 {
			return context.DeadlineExceeded
		}
		return nil
	})
	if err != nil || attempts != 2 {
		t.Errorf("expected 2 attempts and success, got attempts=%d, err=%v", attempts, err)
	}
}
