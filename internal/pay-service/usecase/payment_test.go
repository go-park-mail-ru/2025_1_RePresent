package payment

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"retarget/internal/pay-service/entity"
	"retarget/internal/pay-service/repo"
)

type roundTripper func(req *http.Request) *http.Response

func (rt roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt(req), nil
}

func Test_mapYooStatus(t *testing.T) {
	cases := []struct {
		s    string
		want int
	}{
		{"pending", 0},
		{"waiting_for_capture", 1},
		{"canceled", 2},
		{"something", -1},
	}
	for _, c := range cases {
		assert.Equal(t, c.want, mapYooStatus(c.s), "input=%s", c.s)
	}
}

func Test_CheckBalance_SuccessAndLow(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	p := repo.NewPaymentRepositoryWithDB(db, zap.NewNop().Sugar())
	uc := &PaymentUsecase{PaymentRepository: p}

	mock.ExpectQuery("SELECT balance FROM auth_user").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(200.0))
	bal, err := uc.CheckBalance(1)
	assert.NoError(t, err)
	assert.Equal(t, 200.0, bal)

	mock.ExpectQuery("SELECT balance FROM auth_user").
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(50.0))
	bal2, err2 := uc.CheckBalance(2)
	assert.ErrorIs(t, err2, errTooLittleBalance)
	assert.Equal(t, 50.0, bal2)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func Test_GetTransactionByID_FoundAndNotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	p := repo.NewPaymentRepositoryWithDB(db, zap.NewNop().Sugar())
	uc := &PaymentUsecase{PaymentRepository: p}

	cols := []string{"id", "transaction_id", "user_id", "amount", "type", "status", "created_at"}
	now := time.Now()

	mock.ExpectQuery("WHERE transaction_id = \\$1").
		WithArgs("tx1").
		WillReturnRows(sqlmock.NewRows(cols).AddRow(1, "tx1", 3, 12.5, "x", 0, now))
	tx, err := uc.GetTransactionByID("tx1", "req")
	assert.NoError(t, err)
	assert.Equal(t, "tx1", tx.TransactionID)

	mock.ExpectQuery("WHERE transaction_id = \\$1").
		WithArgs("nope").
		WillReturnError(sql.ErrNoRows)
	_, err2 := uc.GetTransactionByID("nope", "req")
	assert.Error(t, err2)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func Test_CreateYooMoneyPayment_VariousScenarios(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	p := repo.NewPaymentRepositoryWithDB(db, zap.NewNop().Sugar())

	respOK := `{
		"id":"i1","status":"waiting_for_capture",
		"confirmation":{"confirmation_url":"u1"},
		"amount":{"value":"15.00","currency":"RUB"}
	}`
	clientOK := &http.Client{Transport: roundTripper(func(req *http.Request) *http.Response {
		return &http.Response{StatusCode: http.StatusAccepted, Body: ioutil.NopCloser(strings.NewReader(respOK))}
	})}
	uc := &PaymentUsecase{
		PaymentRepository: p,
		httpClient:        clientOK,
		shopID:            "s",
		secretKey:         "k",
	}

	mock.ExpectQuery("SELECT 1 FROM transaction").
		WithArgs("i1").
		WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO transaction").
		WithArgs("i1", 5, 15.00, "yoo_money", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	url, err := uc.CreateYooMoneyPayment(5, "15.00", "RUB", "ret", "d", "key1")
	assert.NoError(t, err)
	assert.Equal(t, "u1", url)

	clientBadStatus := &http.Client{Transport: roundTripper(func(req *http.Request) *http.Response {
		return &http.Response{StatusCode: 500, Body: ioutil.NopCloser(bytes.NewBufferString(`{}`))}
	})}
	uc.httpClient = clientBadStatus
	_, err = uc.CreateYooMoneyPayment(5, "1", "RUB", "", "", "")
	assert.Error(t, err)

	badJSON := `{"id":`
	clientDec := &http.Client{Transport: roundTripper(func(req *http.Request) *http.Response {
		return &http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(strings.NewReader(badJSON))}
	})}
	uc.httpClient = clientDec
	_, err = uc.CreateYooMoneyPayment(5, "1", "RUB", "", "", "")
	assert.Error(t, err)

	resp2 := `{
		"id":"i3","status":"pending",
		"confirmation":{"confirmation_url":"u3"},
		"amount":{"value":"1.23","currency":"RUB"}
	}`
	client2 := &http.Client{Transport: roundTripper(func(req *http.Request) *http.Response {
		return &http.Response{StatusCode: http.StatusAccepted, Body: ioutil.NopCloser(strings.NewReader(resp2))}
	})}
	uc.httpClient = client2
	mock.ExpectQuery("SELECT 1 FROM transaction").
		WithArgs("i3").
		WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO transaction").
		WithArgs("i3", 6, 1.23, "yoo_money", 0).
		WillReturnError(errors.New("db error"))
	_, err = uc.CreateYooMoneyPayment(6, "1.23", "RUB", "", "", "")
	assert.Error(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func Test_RegUserActivity_ErrorOnDebit(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repoDB := repo.NewPaymentRepositoryWithDB(db, zap.NewNop().Sugar())
	uc := &PaymentUsecase{
		logger:            zap.NewNop().Sugar(),
		PaymentRepository: repoDB,
	}

	amt := entity.Decimal{}
	_ = amt.ParseFromString("5.0")

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE auth_user").
		WithArgs(amt, 2).
		WillReturnError(errors.New("debit fail"))
	mock.ExpectRollback()

	err := uc.RegUserActivity(1, 2, amt)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

type attemptStub struct {
	MaxAttempts int
	OnGet       func(int) (int, error)
	OnInc       func(int) error
	OnDec       func(int) error
}

func (s *attemptStub) GetAttemptsByUserID(id int) (int, error) { return s.OnGet(id) }
func (s *attemptStub) IncrementAttemptsByUserID(id int) error {
	if s.OnInc != nil {
		return s.OnInc(id)
	}
	return nil
}
func (s *attemptStub) ResetAttemptsByUserID(id int) error { return nil }
func (s *attemptStub) DecrementAttemptsByUserID(id int) error {
	if s.OnDec != nil {
		return s.OnDec(id)
	}
	return nil
}

type noticeStub struct {
	OnLow func(int) error
}

func (n *noticeStub) SendLowBalanceNotification(uid int) error { return n.OnLow(uid) }
func (n *noticeStub) SendTopUpBalanceEvent(uid int, a float64) error {
	return nil
}

func Test_offBannersByUserID_SuccessAndErrorLogged(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repoDB := repo.NewPaymentRepositoryWithDB(db, zap.NewNop().Sugar())
	uc := &PaymentUsecase{
		logger:            zap.NewNop().Sugar(),
		PaymentRepository: repoDB,
	}

	mock.ExpectExec("UPDATE banner").
		WithArgs(99).
		WillReturnResult(sqlmock.NewResult(0, 2))
	uc.offBannersByUserID(context.Background(), 99)

	mock.ExpectExec("UPDATE banner").
		WithArgs(100).
		WillReturnError(errors.New("db error"))
	uc.offBannersByUserID(context.Background(), 100)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func Test_GetBalanceByUserId_PendingCanceledTransaction(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repoDB := repo.NewPaymentRepositoryWithDB(db, zap.NewNop().Sugar())
	uc := &PaymentUsecase{
		PaymentRepository: repoDB,
		httpClient: &http.Client{Transport: roundTripper(func(req *http.Request) *http.Response {
			return &http.Response{StatusCode: 200, Body: ioutil.NopCloser(strings.NewReader(`{"status":"canceled"}`))}
		})},
	}

	now := time.Now()
	mock.ExpectQuery("FROM transaction").
		WithArgs(5).
		WillReturnRows(sqlmock.NewRows([]string{"id", "transaction_id", "user_id", "amount", "type", "status", "created_at"}).
			AddRow(1, "tx5", 5, 10.0, "yoo_money", 0, now))
	mock.ExpectQuery("SELECT balance FROM auth_user").
		WithArgs(5, "req5").
		WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(80.0))

	bal, _ := uc.GetBalanceByUserId(5, "req5")
	assert.Equal(t, 0.0, bal)
}

func Test_GetBalanceByUserId_FinalBalanceError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repoDB := repo.NewPaymentRepositoryWithDB(db, zap.NewNop().Sugar())
	uc := &PaymentUsecase{PaymentRepository: repoDB}

	mock.ExpectQuery("FROM transaction").
		WithArgs(6).
		WillReturnRows(sqlmock.NewRows([]string{"id", "transaction_id", "user_id", "amount", "type", "status", "created_at"}))
	mock.ExpectQuery("SELECT balance FROM auth_user").
		WithArgs(6, "req6").
		WillReturnError(errors.New("scan error"))

	_, err := uc.GetBalanceByUserId(6, "req6")
	assert.Error(t, err)
}
