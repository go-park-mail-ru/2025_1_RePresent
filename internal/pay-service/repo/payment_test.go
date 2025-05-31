package repo_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"retarget/internal/pay-service/entity"
	"retarget/internal/pay-service/repo"
)

func setup() (*repo.PaymentRepository, sqlmock.Sqlmock, func()) {
	db, mock, _ := sqlmock.New()
	logger := zap.NewNop().Sugar()
	r := repo.NewPaymentRepositoryWithDB(db, logger)
	return r, mock, func() { db.Close() }
}

func TestGetBalanceByUserId_Success(t *testing.T) {
	r, mock, close := setup()
	defer close()

	mock.ExpectQuery("SELECT balance FROM auth_user").
		WithArgs(42).
		WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(123.45))

	bal, err := r.GetBalanceByUserId(42, "req-1")
	assert.NoError(t, err)
	assert.Equal(t, 123.45, bal)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetBalanceByUserId_NoRows(t *testing.T) {
	r, mock, close := setup()
	defer close()

	mock.ExpectQuery("SELECT balance FROM auth_user").
		WithArgs(7).
		WillReturnError(sql.ErrNoRows)

	bal, err := r.GetBalanceByUserId(7, "req-2")
	assert.Error(t, err)
	assert.Zero(t, bal)
	assert.Contains(t, err.Error(), "user with id 7 not found")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateBalance_UserNotFound(t *testing.T) {
	r, mock, close := setup()
	defer close()

	mock.ExpectQuery("UPDATE auth_user").
		WithArgs(10.0, 99).
		WillReturnError(sql.ErrNoRows)

	bal, err := r.UpdateBalance(99, 10.0, "req-3")
	assert.ErrorIs(t, err, repo.ErrUserNotFound)
	assert.Zero(t, bal)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateTransaction_Exists(t *testing.T) {
	r, mock, close := setup()
	defer close()

	mock.ExpectQuery("SELECT 1 FROM transaction").
		WithArgs("trx-1").
		WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(1))

	err := r.CreateTransaction(entity.Transaction{TransactionID: "trx-1"})
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateTransaction_New(t *testing.T) {
	r, mock, close := setup()
	defer close()

	mock.ExpectQuery("SELECT 1 FROM transaction").
		WithArgs("trx-2").
		WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO transaction").
		WithArgs("trx-2", 0, 0.0, "", 0).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := r.CreateTransaction(entity.Transaction{TransactionID: "trx-2"})
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeactivateBannersByUserID(t *testing.T) {
	r, mock, close := setup()
	defer close()

	mock.ExpectExec("UPDATE banner").
		WithArgs(5).
		WillReturnResult(sqlmock.NewResult(0, 3))

	err := r.DeactivateBannersByUserID(context.Background(), 5)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetLastTransaction_Success(t *testing.T) {
	r, mock, close := setup()
	defer close()

	cols := []string{"id", "transaction_id", "user_id", "amount", "type", "status", "created_at"}
	now := time.Now()
	mock.ExpectQuery("SELECT \\* FROM transaction").
		WithArgs(11).
		WillReturnRows(sqlmock.NewRows(cols).AddRow(1, "trx-id", 11, 50.0, "t", 1, now))

	tx, err := r.GetLastTransaction(11, "req")
	assert.NoError(t, err)
	assert.Equal(t, "trx-id", tx.TransactionID)
	assert.Equal(t, now, tx.CreatedAt)
}

func TestGetLastTransaction_Error(t *testing.T) {
	r, mock, close := setup()
	defer close()

	mock.ExpectQuery("SELECT \\* FROM transaction").
		WithArgs(12).
		WillReturnError(sql.ErrNoRows)

	tx, err := r.GetLastTransaction(12, "req")
	assert.Error(t, err)
	assert.Nil(t, tx)
}

func TestGetTransactionByID_Success(t *testing.T) {
	r, mock, close := setup()
	defer close()

	cols := []string{"id", "transaction_id", "user_id", "amount", "type", "status", "created_at"}
	now := time.Now()
	mock.ExpectQuery("SELECT \\* FROM transaction WHERE transaction_id").
		WithArgs("id-1").
		WillReturnRows(sqlmock.NewRows(cols).AddRow(2, "id-1", 22, 75.5, "t", 0, now))

	tx, err := r.GetTransactionByID("id-1", "req")
	assert.NoError(t, err)
	assert.Equal(t, 22, tx.UserID)
}

func TestGetTransactionByID_NotFound(t *testing.T) {
	r, mock, close := setup()
	defer close()

	mock.ExpectQuery("SELECT \\* FROM transaction WHERE transaction_id").
		WithArgs("nope").
		WillReturnError(sql.ErrNoRows)

	tx, err := r.GetTransactionByID("nope", "req")
	assert.Error(t, err)
	assert.Nil(t, tx)
	assert.Contains(t, err.Error(), "not found")
}

func TestRegUserActivity_Success(t *testing.T) {
	r, mock, close := setup()
	defer close()

	amount := entity.Decimal{Dec: nil} // Value=>"0"
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE auth_user").
		WithArgs(sqlmock.AnyArg(), 2).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("UPDATE auth_user").
		WithArgs(sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	from, to, err := r.RegUserActivity(1, 2, amount)
	assert.NoError(t, err)
	assert.Equal(t, 1, from)
	assert.Equal(t, 2, to)
}

func TestRegUserActivity_FirstUpdateError(t *testing.T) {
	r, mock, close := setup()
	defer close()

	amount := entity.Decimal{Dec: nil}
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE auth_user").
		WithArgs(sqlmock.AnyArg(), 2).
		WillReturnError(fmt.Errorf("first error"))
	mock.ExpectRollback().WillReturnError(fmt.Errorf("rb error"))

	_, _, err := r.RegUserActivity(1, 2, amount)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to rollback")
}

func TestRegUserActivity_SecondUpdateError(t *testing.T) {
	r, mock, close := setup()
	defer close()

	amount := entity.Decimal{Dec: nil}
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE auth_user").
		WithArgs(sqlmock.AnyArg(), 2).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("UPDATE auth_user").
		WithArgs(sqlmock.AnyArg(), 1).
		WillReturnError(fmt.Errorf("second error"))
	mock.ExpectRollback()

	_, _, err := r.RegUserActivity(1, 2, amount)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update second user balance")
}

func TestRegUserActivity_CommitError(t *testing.T) {
	r, mock, close := setup()
	defer close()

	amount := entity.Decimal{Dec: nil}
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE auth_user").
		WithArgs(sqlmock.AnyArg(), 2).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("UPDATE auth_user").
		WithArgs(sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("commit error"))

	_, _, err := r.RegUserActivity(1, 2, amount)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to commit")
}

func TestGetPendingTransactions_Success(t *testing.T) {
	r, mock, close := setup()
	defer close()

	cols := []string{"id", "transaction_id", "user_id", "amount", "type", "status", "created_at"}
	now := time.Now()
	mock.ExpectQuery("SELECT id, transaction_id").
		WithArgs(33).
		WillReturnRows(sqlmock.NewRows(cols).
			AddRow(1, "a", 33, 1.1, "t", 0, now).
			AddRow(2, "b", 33, 2.2, "t", 0, now))

	list, err := r.GetPendingTransactions(33)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestGetPendingTransactions_QueryError(t *testing.T) {
	r, mock, close := setup()
	defer close()

	mock.ExpectQuery("SELECT id, transaction_id").
		WithArgs(34).
		WillReturnError(fmt.Errorf("query error"))

	_, err := r.GetPendingTransactions(34)
	assert.Error(t, err)
}

func TestUpdateTransactionStatus(t *testing.T) {
	r, mock, close := setup()
	defer close()

	mock.ExpectExec("UPDATE transaction").
		WithArgs(2, "trx").
		WillReturnResult(sqlmock.NewResult(1, 1))
	assert.NoError(t, r.UpdateTransactionStatus("trx", 2))

	mock.ExpectExec("UPDATE transaction").
		WithArgs(3, "trx2").
		WillReturnError(fmt.Errorf("exec error"))
	assert.Error(t, r.UpdateTransactionStatus("trx2", 3))
}
