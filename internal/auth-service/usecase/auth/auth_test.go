package auth

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis/v2"
	"go.uber.org/zap"

	repoAuth "retarget/internal/auth-service/repo/auth"
	"retarget/pkg/utils/optiLog"
)

// setupUsecase создаёт AuthUsecase с sqlmock-репозиторием и miniredis-сессией
func setupUsecase(t *testing.T) (*AuthUsecase, sqlmock.Sqlmock, *repoAuth.SessionRepository) {
	// sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock open failed: %v", err)
	}
	loggerSug := zap.NewNop().Sugar()
	asyncLogger := optiLog.NewAsyncLogger(loggerSug, 1, 1)
	authRepo := repoAuth.NewAuthRepositoryWithDB(db, loggerSug)

	// miniredis
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis start failed: %v", err)
	}
	sessionRepo := repoAuth.NewSessionRepository(mr.Addr(), "", 0, time.Second)

	uc := NewAuthUsecase(authRepo, sessionRepo, asyncLogger)
	return uc, mock, sessionRepo
}

func TestLogin_SuccessAndFailures(t *testing.T) {
	uc, mock, _ := setupUsecase(t)
	pass := "correct"
	storedHash := hashForTest(t, pass)

	// успешный вход
	rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "description", "balance", "role"}).
		AddRow(7, "u", "e@e", storedHash, "", "0", 1)
	mock.ExpectQuery("WHERE email = \\$1").WithArgs("e@e").WillReturnRows(rows)
	u, err := uc.Login(context.Background(), "e@e", pass, 1, "req1")
	if err != nil || u.ID != 7 {
		t.Fatalf("expected ID=7, got %v, err %v", u, err)
	}

	// неправильный пароль
	rows = sqlmock.NewRows([]string{"id", "username", "email", "password", "description", "balance", "role"}).
		AddRow(7, "u", "e@e", storedHash, "", "0", 1)
	mock.ExpectQuery("WHERE email = \\$1").WithArgs("e@e").WillReturnRows(rows)
	_, err = uc.Login(context.Background(), "e@e", "bad", 1, "req2")
	if err == nil || err.Error() == "incorrect user data" {
		t.Errorf("expected incorrect user data, got %v", err)
	}

	// слишком много попыток
	rows = sqlmock.NewRows([]string{"id", "username", "email", "password", "description", "balance", "role"}).
		AddRow(7, "u", "slow@e", storedHash, "", "0", 1)
	mock.ExpectQuery("WHERE email = \\$1").WithArgs("slow@e").WillReturnRows(rows)
	_, _ = uc.Login(context.Background(), "slow@e", pass, 1, "req3")
	_, err = uc.Login(context.Background(), "slow@e", pass, 1, "req3")
	if err == nil || err.Error() != "слишком много попыток" {
		t.Errorf("expected rate limit error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err == nil {
		t.Errorf("unmet sqlmock expectations: %v", err)
	}
}

func TestGetUser_SuccessAndError(t *testing.T) {
	uc, mock, _ := setupUsecase(t)

	// success
	rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "description", "balance", "role"}).
		AddRow(42, "u", "e@e", []byte("p"), "", "0", 1)
	mock.ExpectQuery("WHERE id = \\$1").WithArgs(42).WillReturnRows(rows)
	user, err := uc.GetUser(context.Background(), 42, "reqG1")
	if err != nil || user.ID != 42 {
		t.Fatalf("expected ID=42, got %v, err %v", user, err)
	}

	// error
	mock.ExpectQuery("WHERE id = \\$1").WithArgs(43).WillReturnError(sql.ErrNoRows)
	_, err = uc.GetUser(context.Background(), 43, "reqG2")
	if err == nil {
		t.Errorf("expected error for missing user")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sqlmock expectations: %v", err)
	}
}

func TestRegister_Flows(t *testing.T) {
	uc, mock, _ := setupUsecase(t)

	// короткий пароль
	_, err := uc.Register(context.Background(), "u", "e", "short", 1, "r1")
	if err == nil {
		t.Errorf("expected password length error")
	}

	// ошибка при проверке существующего
	mock.ExpectQuery("FROM auth_user").WithArgs("e", "u").
		WillReturnError(errors.New("chk fail"))
	_, err = uc.Register(context.Background(), "u", "e", "longenuf", 1, "r2")
	if err == nil || err.Error() == "chk fail" {
		t.Errorf("expected chk fail, got %v", err)
	}

	rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "description", "balance", "role"}).
		AddRow(1, "u1", "e", "p", "", "0", 1)
	mock.ExpectQuery("FROM auth_user").WithArgs("e", "u2").WillReturnRows(rows)
	_, err = uc.Register(context.Background(), "u2", "e", "longenuf", 1, "r3")
	if err == nil || err.Error() != "пользователь с таким email уже существует" {
		t.Errorf("expected email exists, got %v", err)
	}

	rows = sqlmock.NewRows([]string{"id", "username", "email", "password", "description", "balance", "role"}).
		AddRow(2, "u", "e2", "p", "", "0", 1)
	mock.ExpectQuery("FROM auth_user").WithArgs("e3", "u").WillReturnRows(rows)
	_, err = uc.Register(context.Background(), "u", "e3", "longenuf", 1, "r4")
	if err == nil || err.Error() != "пользователь с таким username уже существует" {
		t.Errorf("expected username exists, got %v", err)
	}

	mock.ExpectQuery("FROM auth_user").WithArgs("e4", "u4").WillReturnError(sql.ErrNoRows)
	mock.ExpectQuery("INSERT INTO auth_user").
		WithArgs("u4", "e4", sqlmock.AnyArg(), "", "0", 1).
		WillReturnError(errors.New("db err"))
	_, err = uc.Register(context.Background(), "u4", "e4", "longenuf", 1, "r5")
	if err == nil || err.Error() == "db err" {
		t.Errorf("expected db err, got %v", err)
	}

	// успешная регистрация
	mock.ExpectQuery("FROM auth_user").WithArgs("e5", "u5").WillReturnError(sql.ErrNoRows)
	mock.ExpectQuery("INSERT INTO auth_user").
		WithArgs("u5", "e5", sqlmock.AnyArg(), "", "0", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(99))

	if err := mock.ExpectationsWereMet(); err == nil {
		t.Errorf("unmet sqlmock expectations: %v", err)
	}
}

func hashForTest(t *testing.T, pw string) string {
	h, err := hashPassword(pw, DefaultHashConfig)
	if err != nil {
		t.Fatal(err)
	}
	return string(h)
}
