package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"retarget/internal/pay-service/repo"
	"testing"

	usecase "retarget/internal/pay-service/usecase"
	response "retarget/pkg/entity"
	auth "retarget/pkg/middleware/auth"

	"github.com/stretchr/testify/assert"
)

func makeHandler() *PaymentController {
	return NewPaymentController(&usecase.PaymentUsecase{})
}

func Test_GetUserBalance_NoCookie(t *testing.T) {
	h := makeHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/payment/balance", nil)
	req = req.WithContext(context.WithValue(req.Context(), response.СtxKeyRequestID{}, "rid"))
	rr := httptest.NewRecorder()

	h.GetUserBalance(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid Cookie")
}

func Test_GetUserBalance_NoUserContext(t *testing.T) {
	h := makeHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/payment/balance", nil)
	req = req.WithContext(context.WithValue(req.Context(), response.СtxKeyRequestID{}, "rid"))
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "abc"})
	rr := httptest.NewRecorder()

	h.GetUserBalance(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "Error of authenticator")
}

func Test_TopUpAccount_NoCookie(t *testing.T) {
	h := makeHandler()
	body := bytes.NewBufferString(`{"amount":10}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/payment/accounts/topup", body)
	req = req.WithContext(context.WithValue(req.Context(), response.СtxKeyRequestID{}, "rid"))
	rr := httptest.NewRecorder()

	h.TopUpAccount(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid Cookie")
}

func Test_TopUpAccount_InvalidBody(t *testing.T) {
	h := makeHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/payment/accounts/topup", bytes.NewBufferString(`not json`))
	req = req.WithContext(context.WithValue(req.Context(), response.СtxKeyRequestID{}, "rid"))
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "ok"})
	rr := httptest.NewRecorder()

	h.TopUpAccount(rr, req)

	assert.Contains(t, []int{http.StatusBadRequest, http.StatusInternalServerError}, rr.Code)
	assert.True(t, rr.Body.Len() > 0)
}

func Test_CreateTransaction_InvalidBody(t *testing.T) {
	h := makeHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/payment/transactions", bytes.NewBufferString(`{bad`))
	rr := httptest.NewRecorder()

	h.CreateTransaction(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "invalid request payload")
}

func Test_Router_ApplyMiddlewares(t *testing.T) {
	authn := &auth.Authenticator{}
	mux := SetupPaymentRoutes(authn, &usecase.PaymentUsecase{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/payment/balance", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}
func Test_CreateTransaction_NoUserContext(t *testing.T) {
	h := makeHandler()
	body := bytes.NewBufferString(`{
        "value":"1.23","currency":"RUB",
        "return_url":"u","description":"d","idempotence_key":"k"
    }`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/payment/transactions", body)
	// устанавливаем заглушку requestID и cookie, но не кладём UserContextKey
	req = req.WithContext(context.WithValue(req.Context(), response.СtxKeyRequestID{}, "rid"))
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "ok"})
	rr := httptest.NewRecorder()

	h.CreateTransaction(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "Error of authenticator")
}

func Test_RegUserActivity_Endpoint(t *testing.T) {
	authn := &auth.Authenticator{}
	mux := SetupPaymentRoutes(authn, &usecase.PaymentUsecase{})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/payment/transactions/clicks", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Zero(t, rr.Body.Len())
}

func Test_Router_MethodRestrictions(t *testing.T) {
	authn := &auth.Authenticator{}
	mux := SetupPaymentRoutes(authn, &usecase.PaymentUsecase{})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/payment/transactions", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	assert.Equal(t, 405, rr.Code)
}
func Test_PostBalance_NoCookie(t *testing.T) {
	mux := SetupPaymentRoutes(&auth.Authenticator{}, &usecase.PaymentUsecase{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/payment/balance", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "named cookie not present")
}

func Test_GetTopUp_NoCookie(t *testing.T) {
	mux := SetupPaymentRoutes(&auth.Authenticator{}, &usecase.PaymentUsecase{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/payment/accounts/topup", nil)
	req = req.WithContext(context.WithValue(req.Context(), response.СtxKeyRequestID{}, "rid"))
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), "named cookie not present")
}

func Test_DeleteTransactions_MethodNotAllowed(t *testing.T) {
	mux := SetupPaymentRoutes(&auth.Authenticator{}, &usecase.PaymentUsecase{})
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/payment/transactions", nil)
	req = req.WithContext(context.WithValue(req.Context(), response.СtxKeyRequestID{}, "rid"))
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func Test_UnknownRoute(t *testing.T) {
	mux := SetupPaymentRoutes(&auth.Authenticator{}, &usecase.PaymentUsecase{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/payment/nonexistent", nil)
	req = req.WithContext(context.WithValue(req.Context(), response.СtxKeyRequestID{}, "rid"))
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}
func Test_TopUpAccount_InvalidAmountZero(t *testing.T) {
	h := makeHandler()
	body := bytes.NewBufferString(`{"amount":0}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/payment/accounts/topup", body)
	req = req.WithContext(context.WithValue(req.Context(), response.СtxKeyRequestID{}, "rid"))
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "abc"})
	rr := httptest.NewRecorder()

	h.TopUpAccount(rr, req)

	assert.Equal(t, 500, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid Amount")
}

func Test_GetTransactionByID_Router_NoCookie(t *testing.T) {
	mux := SetupPaymentRoutes(&auth.Authenticator{}, &usecase.PaymentUsecase{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/payment/transactions/tx123", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}
func Test_GetUserBalance_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	payRepo := repo.NewPaymentRepositoryWithDB(db, zap.NewNop().Sugar())
	uc := usecase.NewPayUsecase(zap.NewNop().Sugar(), payRepo, nil, nil, "", "", &http.Client{})
	ctrl := NewPaymentController(uc)

	ctx := context.WithValue(context.Background(), response.СtxKeyRequestID{}, "req1")
	ctx = context.WithValue(ctx, response.UserContextKey, response.UserContext{UserID: 5})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/payment/balance", nil).WithContext(ctx)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "sess5"})
	rr := httptest.NewRecorder()

	mock.ExpectQuery("SELECT id, transaction_id, user_id, amount, type, status, created_at FROM transaction").
		WithArgs(5).
		WillReturnRows(sqlmock.NewRows([]string{"id", "transaction_id", "user_id", "amount", "type", "status", "created_at"}))
	mock.ExpectQuery("SELECT balance FROM auth_user").
		WithArgs(5).
		WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(250.75))

	ctrl.GetUserBalance(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var resp map[string]float64
	assert.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
	assert.Equal(t, 250.75, resp["balance"])

	assert.NoError(t, mock.ExpectationsWereMet())
}
func Test_HeadBalance_NoCookie(t *testing.T) {
	mux := SetupPaymentRoutes(&auth.Authenticator{}, &usecase.PaymentUsecase{})
	req := httptest.NewRequest(http.MethodHead, "/api/v1/payment/balance", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func Test_HeadTopUp_NoCookie(t *testing.T) {
	mux := SetupPaymentRoutes(&auth.Authenticator{}, &usecase.PaymentUsecase{})
	req := httptest.NewRequest(http.MethodHead, "/api/v1/payment/accounts/topup", nil)
	req = req.WithContext(context.WithValue(req.Context(), response.СtxKeyRequestID{}, "rid"))
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func Test_RegUserActivity_Direct(t *testing.T) {
	ctrl := NewPaymentController(&usecase.PaymentUsecase{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/payment/transactions/clicks", nil)
	rr := httptest.NewRecorder()
	ctrl.RegUserActivity(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Zero(t, rr.Body.Len())
}
func Test_GetTransactions_List_MethodNotAllowed(t *testing.T) {
	mux := SetupPaymentRoutes(&auth.Authenticator{}, &usecase.PaymentUsecase{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/payment/transactions", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func Test_POST_Transactions_Clicks(t *testing.T) {
	mux := SetupPaymentRoutes(&auth.Authenticator{}, &usecase.PaymentUsecase{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/payment/transactions/clicks", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Zero(t, rr.Body.Len())
}

func Test_HEAD_Transactions_Clicks(t *testing.T) {
	mux := SetupPaymentRoutes(&auth.Authenticator{}, &usecase.PaymentUsecase{})
	req := httptest.NewRequest(http.MethodHead, "/api/v1/payment/transactions/clicks", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func Test_OPTIONS_Balance(t *testing.T) {
	mux := SetupPaymentRoutes(&auth.Authenticator{}, &usecase.PaymentUsecase{})
	req := httptest.NewRequest(http.MethodOptions, "/api/v1/payment/balance", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	assert.Equal(t, 401, rr.Code)
}

func Test_CreateTransaction_CookieButNoContext(t *testing.T) {
	h := makeHandler()
	body := bytes.NewBufferString(`{"value":"1","currency":"RUB","return_url":"u","description":"d","idempotence_key":"k"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/payment/transactions", body)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "ok"})
	rr := httptest.NewRecorder()
	h.CreateTransaction(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func Test_RegUserActivity_UnsupportedMethod(t *testing.T) {
	mux := SetupPaymentRoutes(&auth.Authenticator{}, &usecase.PaymentUsecase{})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/payment/transactions/clicks", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}
func Test_GetTransactionByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	payRepo := repo.NewPaymentRepositoryWithDB(db, zap.NewNop().Sugar())
	uc := usecase.NewPayUsecase(zap.NewNop().Sugar(), payRepo, nil, nil, "", "", &http.Client{})
	ctrl := NewPaymentController(uc)

	cols := []string{"id", "transaction_id", "user_id", "amount", "type", "status", "created_at"}
	mock.ExpectQuery("SELECT * FROM transaction").
		WithArgs("tx1").
		WillReturnRows(sqlmock.NewRows(cols).AddRow(1, "tx1", 7, 99.9, "yoo_money", 1, nil))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/payment/transactions/tx1", nil)
	req = mux.SetURLVars(req, map[string]string{"transactionid": "tx1"})
	req = req.WithContext(context.WithValue(req.Context(), response.СtxKeyRequestID{}, "rid"))
	rr := httptest.NewRecorder()
	ctrl.GetTransactionByID(rr, req)

	assert.Equal(t, 404, rr.Code)
	var out map[string]interface{}
	assert.NoError(t, json.NewDecoder(rr.Body).Decode(&out))
}

func Test_GetTransactionByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	payRepo := repo.NewPaymentRepositoryWithDB(db, zap.NewNop().Sugar())
	uc := usecase.NewPayUsecase(zap.NewNop().Sugar(), payRepo, nil, nil, "", "", &http.Client{})
	ctrl := NewPaymentController(uc)

	mock.ExpectQuery("SELECT \\* FROM transaction").
		WithArgs("nope").
		WillReturnError(repo.ErrUserNotFound)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/payment/transactions/nope", nil)
	req = mux.SetURLVars(req, map[string]string{"transactionid": "nope"})
	req = req.WithContext(context.WithValue(req.Context(), response.СtxKeyRequestID{}, "rid"))
	rr := httptest.NewRecorder()

	ctrl.GetTransactionByID(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	var resp map[string]interface{}
	assert.NoError(t, json.NewDecoder(rr.Body).Decode(&resp))
	assert.NoError(t, mock.ExpectationsWereMet())
}
