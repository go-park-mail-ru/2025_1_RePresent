package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"retarget/internal/pay-service/entity"
	"retarget/internal/pay-service/repo"
	"retarget/internal/pay-service/repo/attempt"
	"retarget/internal/pay-service/repo/notice"
	"strconv"
	"time"

	"github.com/cenkalti/backoff"

	"go.uber.org/zap"
)

var (
	errTooLittleBalance = errors.New("balance of user is need to top up")
)

const (
	BalanceLimit     = 100.00
	MinActiveBalance = 10.00
)

type PaymentUsecase struct {
	logger            *zap.SugaredLogger
	PaymentRepository *repo.PaymentRepository
	NoticeRepository  *notice.NoticeRepository
	AttemptRepository *attempt.AttemptRepository
	shopID            string
	secretKey         string
	uuid              string // идентификатор мерчанта для payouts
	accountNumber     string // кошелёк для автоматических выплат
	httpClient        *http.Client
}

func NewPayUsecase(
	zapLogger *zap.SugaredLogger,
	payRepository *repo.PaymentRepository,
	noticeRepository *notice.NoticeRepository,
	attemptRepository *attempt.AttemptRepository,
	shopID, secretKey, uuid, accountNumber string,
	httpClient *http.Client,
) *PaymentUsecase {
	return &PaymentUsecase{
		logger:            zapLogger,
		PaymentRepository: payRepository,
		NoticeRepository:  noticeRepository,
		AttemptRepository: attemptRepository,
		shopID:            shopID,
		secretKey:         secretKey,
		uuid:              uuid, // инициализируем UUID
		accountNumber:     accountNumber,
		httpClient:        httpClient,
	}
}

func (u *PaymentUsecase) GetBalanceByUserId(userID int, requestID string) (float64, error) {
	pending, err := u.PaymentRepository.GetPendingTransactions(userID)

	if err != nil {
		return 0, fmt.Errorf("error while loaing pending tx: %w", err)
	}

	for _, tx := range pending {
		statusStr, err := u.getYooPaymentStatus(tx.TransactionID)

		if err != nil {
			continue
		}
		statusInt := mapYooStatus(statusStr)
		if statusInt == 1 && tx.Status != 1 {
			if err := u.TopUpBalance(userID, tx.Amount, requestID); err == nil {
				_ = u.PaymentRepository.UpdateTransactionStatus(tx.TransactionID, statusInt)
			}
		}
	}

	balik, err := u.PaymentRepository.GetBalanceByUserId(userID, requestID)
	if err != nil {
		return 0, err
	}

	return balik, nil
	// return u.PaymentRepository.GetBalanceByUserId(userID, requestID)
}

func mapYooStatus(s string) int {
	switch s {
	case "pending":
		return 0
	case "waiting_for_capture":
		return 1
	case "canceled":
		return 2
	default:
		return -1
	}
}

func (u *PaymentUsecase) getYooPaymentStatus(paymentID string) (string, error) {
	url := fmt.Sprintf("https://api.yookassa.ru/v3/payments/%s", paymentID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.SetBasicAuth(u.shopID, u.secretKey)

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("call YooKassa: %w", err)
	}
	defer resp.Body.Close()

	var out struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", fmt.Errorf("decode status: %w", err)
	}
	return out.Status, nil
}

func (uc *PaymentUsecase) TopUpBalance(userID int, amount float64, requestID string) error {
	if amount <= 0 {
		return repo.ErrInvalidAmount
	}

	_, err := uc.PaymentRepository.UpdateBalance(userID, float64(amount), requestID)
	if err != nil {
		return err
	}

	go func() {
		if uc.AttemptRepository != nil {
			if err := uc.AttemptRepository.ResetAttemptsByUserID(userID); err != nil {
				uc.logger.Errorw("failed to reset attempts after top up",
					"user_id", userID,
					"error", err)
			}
		}
		if uc.NoticeRepository != nil {
			if err = uc.NoticeRepository.SendTopUpBalanceEvent(userID, float64(amount)); err != nil {
				uc.logger.Errorw("failed to send topUp message after top up",
					"user_id", userID,
					"error", err)
			}
		}
	}()

	return nil
}

func (uc *PaymentUsecase) GetTransactionByID(transactionID string, requestID string) (*entity.Transaction, error) {
	return uc.PaymentRepository.GetTransactionByID(transactionID, requestID)
}

func (uc *PaymentUsecase) RegUserActivity(user_banner_id, user_slot_id int, amount entity.Decimal) error {
	_, user_from_id, err := uc.PaymentRepository.RegUserActivity(user_banner_id, user_slot_id, amount)
	if err != nil {
		return err
	}
	balance_from, err := uc.CheckBalance(user_from_id)
	if err == errTooLittleBalance {
		if balance_from < MinActiveBalance {
			go uc.offBannersByUserID(context.Background(), user_from_id)
		}
		uc.logger.Infow("balance checked", "user_id", user_from_id, "balance", balance_from, "limit", BalanceLimit)
		go uc.requireSend(user_from_id, strconv.FormatFloat(balance_from, 'f', 2, 64))
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

func (uc *PaymentUsecase) offBannersByUserID(ctx context.Context, userID int) {
	defer func() {
		if r := recover(); r != nil {
			uc.logger.Errorw("error/panic in offing Banners",
				"recovered", r,
				"user_id", userID)
		}
	}()

	if err := uc.PaymentRepository.DeactivateBannersByUserID(ctx, userID); err != nil {
		uc.logger.Errorw("failed to deactivate banners",
			"error", err,
			"user_id", userID)
		return
	}

	uc.logger.Infow("banners deactivated successfully",
		"user_id", userID)
}

func (uc *PaymentUsecase) requireSend(userID int, message string) {
	if uc.AttemptRepository == nil || uc.NoticeRepository == nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			uc.logger.Errorw("error/panic in sendRequired",
				"recovered", r,
				"user_id", userID)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	attempts, err := uc.AttemptRepository.GetAttemptsByUserID(userID)
	if err != nil {
		uc.logger.Errorw("failed to get attempts",
			"user_id", userID,
			"error", err)
		return
	}

	if attempts >= uc.AttemptRepository.MaxAttempts {
		return
	}

	if err := uc.AttemptRepository.IncrementAttemptsByUserID(userID); err != nil {
		uc.logger.Errorw("failed to increment attempts",
			"user_id", userID,
			"error", err)
		return
	}

	backoffConfig := backoff.NewExponentialBackOff()
	backoffConfig.MaxElapsedTime = 1 * time.Minute

	notify := func() error {
		select {
		case <-ctx.Done():
			return backoff.Permanent(ctx.Err())
		default:
			return uc.NoticeRepository.SendLowBalanceNotification(userID)
		}
	}

	if err := backoff.Retry(notify, backoff.WithMaxRetries(backoffConfig, 9)); err != nil {
		uc.logger.Errorw("notification failed after retries",
			"user_id", userID,
			"error", err)

		if err := uc.AttemptRepository.DecrementAttemptsByUserID(userID); err != nil {
			uc.logger.Errorw("failed to decrement attempts",
				"user_id", userID,
				"error", err)
		}
	}
}

func (uc *PaymentUsecase) CheckBalance(user_id int) (float64, error) {
	balance, err := uc.PaymentRepository.GetBalanceByUserId(user_id, "UNIMPLEMENTED request_id")
	if err != nil {
		return balance, err
	}
	if balance <= BalanceLimit {
		return balance, errTooLittleBalance
	}
	return balance, nil
}

type YooPaymentRequest struct {
	Amount struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"amount"`
	Confirmation struct {
		Type      string `json:"type"`
		ReturnURL string `json:"return_url"`
	} `json:"confirmation"`
	Description string `json:"description"`
}

type YooPaymentResponse struct {
	ID           string `json:"id"`
	Status       string `json:"status"`
	Confirmation struct {
		Type            string `json:"type"`
		ConfirmationURL string `json:"confirmation_url"`
	} `json:"confirmation"`
	Amount struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"amount"`
}

func (u *PaymentUsecase) CreateYooMoneyPayment(userID int, value, currency, returnURL, description, idempotenceKey string) (string, error) {
	if len(idempotenceKey) > 40 {
		idempotenceKey = idempotenceKey[:40]
	}

	reqBody := YooPaymentRequest{Description: description}
	reqBody.Amount.Value = value
	reqBody.Amount.Currency = currency
	reqBody.Confirmation.Type = "redirect"
	reqBody.Confirmation.ReturnURL = returnURL

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.yookassa.ru/v3/payments", bytes.NewBuffer(payload))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.SetBasicAuth(u.shopID, u.secretKey)
	req.Header.Set("Idempotence-Key", idempotenceKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return "", fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	var out YooPaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	amt, err := strconv.ParseFloat(out.Amount.Value, 64)
	if err != nil {
		return "", fmt.Errorf("parse amount: %w", err)
	}

	statusInt := mapYooStatus(out.Status)
	trx := entity.Transaction{
		TransactionID: out.ID,
		UserID:        userID,
		Amount:        amt,
		Type:          "yoomoney_payment",
		Status:        statusInt,
	}
	if err := u.PaymentRepository.CreateTransaction(trx); err != nil {
		return "", fmt.Errorf("save transaction: %w", err)
	}

	return out.Confirmation.ConfirmationURL, nil
}

type YooPayoutRequest struct {
	Amount struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"amount"`
	PayoutDestinationData struct {
		Type string `json:"type"`
		Card struct {
			Number string `json:"number"`
		} `json:"card,omitempty"`
		YooMoney struct {
			AccountNumber string `json:"account_number"`
		} `json:"yoo_money,omitempty"`
	} `json:"payoutDestinationData"`
	Confirmation struct {
		Type      string `json:"type"`
		ReturnURL string `json:"return_url"`
	} `json:"confirmation,omitempty"`
	Description string `json:"description"`
	Metadata    struct {
		UserID int `json:"user_id"`
	} `json:"metadata,omitempty"`
}

type YooPayoutResponse struct {
	ID           string `json:"id"`
	Status       string `json:"status"`
	Confirmation struct {
		Type            string `json:"type"`
		ConfirmationURL string `json:"confirmation_url"`
	} `json:"confirmation"`
	Amount struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"amount"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

func (u *PaymentUsecase) CreateYooMoneyPayout(
	userID int,
	amount float64,
	destination string,
	destinationType string,
	description string,
	idempotenceKey string,
) (*entity.Transaction, error) {
	balance, err := u.PaymentRepository.GetBalanceByUserId(userID, idempotenceKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	if balance < amount {
		return nil, errors.New("insufficient funds for withdrawal")
	}

	reqBody := YooPayoutRequest{Description: description}
	reqBody.Amount.Value = fmt.Sprintf("%.2f", amount)
	reqBody.Amount.Currency = "RUB"
	reqBody.PayoutDestinationData.Type = destinationType

	switch destinationType {
	case "bank_card":
		if len(destination) < 16 || len(destination) > 19 {
			return nil, errors.New("invalid card number format")
		}
		reqBody.PayoutDestinationData.Card.Number = destination
	case "yoo_money":
		reqBody.PayoutDestinationData.YooMoney.AccountNumber = destination
	default:
		return nil, errors.New("unsupported destination type")
	}

	reqBody.Metadata.UserID = userID

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	u.logger.Debugw("YooMoney payout request",
		"payload", string(payload),
		"userID", userID,
		"destinationType", destinationType)

	req, err := http.NewRequest(http.MethodPost, "https://api.yookassa.ru/v3/payouts", bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.SetBasicAuth(u.shopID, u.secretKey)
	req.Header.Set("Idempotence-Key", idempotenceKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", idempotenceKey)
	if u.uuid != "" {
		req.Header.Set("X-YooMoney-UUID", u.uuid)
	}

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var out YooPayoutResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	_, err = u.PaymentRepository.UpdateBalance(userID, -amount, idempotenceKey)
	if err != nil {
		u.logger.Errorw("Failed to update balance after payout",
			"user_id", userID,
			"amount", amount,
			"payout_id", out.ID,
			"error", err)
		return nil, fmt.Errorf("update balance after payout: %w", err)
	}

	statusInt := mapYooStatus(out.Status)
	trx := entity.Transaction{
		TransactionID: out.ID,
		UserID:        userID,
		Amount:        amount,
		Type:          "payout_" + destinationType,
		Status:        statusInt,
	}

	if err := u.PaymentRepository.CreateTransaction(trx); err != nil {
		u.logger.Errorw("Failed to save payout transaction",
			"user_id", userID,
			"payout_id", out.ID,
			"error", err)
		return nil, fmt.Errorf("save transaction: %w", err)
	}

	return &trx, nil
}

func (u *PaymentUsecase) CreateYooMoneyPayoutRedirect(
	userID int,
	amount float64,
	description, returnURL, idempotenceKey string,
) (string, error) {
	if len(idempotenceKey) > 40 {
		idempotenceKey = idempotenceKey[:40]
	}

	bal, err := u.PaymentRepository.GetBalanceByUserId(userID, idempotenceKey)
	if err != nil {
		return "", fmt.Errorf("get balance: %w", err)
	}
	if bal < amount {
		return "", errors.New("insufficient funds")
	}

	requestMap := map[string]interface{}{
		"amount": map[string]string{
			"value":    fmt.Sprintf("%.2f", amount),
			"currency": "RUB",
		},
		"description": description,
		"confirmation": map[string]interface{}{
			"type":       "redirect",
			"return_url": returnURL,
		},
		"metadata": map[string]interface{}{
			"user_id": userID,
			"type":    "withdrawal",
		},
	}

	payload, _ := json.Marshal(requestMap)
	u.logger.Debugw("Payment request instead of payout", "payload", string(payload))

	req, err := http.NewRequest(http.MethodPost, "https://api.yookassa.ru/v3/payments", bytes.NewBuffer(payload))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.SetBasicAuth(u.shopID, u.secretKey)
	req.Header.Set("Idempotence-Key", idempotenceKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	u.logger.Debugw("YooMoney API response", "status", resp.StatusCode, "body", bodyString)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return "", fmt.Errorf("unexpected status %d: %s", resp.StatusCode, bodyString)
	}

	var response struct {
		ID           string `json:"id"`
		Status       string `json:"status"`
		Confirmation struct {
			Type            string `json:"type"`
			ConfirmationURL string `json:"confirmation_url"`
		} `json:"confirmation"`
	}

	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	_, err = u.PaymentRepository.UpdateBalance(userID, -amount, idempotenceKey)
	if err != nil {
		u.logger.Errorw("Failed to update balance after successful payment initiation",
			"user_id", userID,
			"amount", amount,
			"error", err)
		return "", fmt.Errorf("update balance after payment: %w", err)
	}

	trx := entity.Transaction{
		TransactionID: response.ID,
		UserID:        userID,
		Amount:        amount,
		Type:          "withdrawal_payment",
		Status:        mapYooStatus(response.Status),
	}
	_ = u.PaymentRepository.CreateTransaction(trx)

	u.logger.Debugw("Payment for withdrawal created",
		"confirmation_url", response.Confirmation.ConfirmationURL,
		"payment_id", response.ID)

	return response.Confirmation.ConfirmationURL, nil
}

func (u *PaymentUsecase) CreateYooMoneyPayoutAuto(
	userID int,
	amount float64,
	description string,
	idempotenceKey string,
) (*entity.Transaction, error) {
	return u.CreateYooMoneyPayout(
		userID,
		amount,
		u.accountNumber,
		"yoo_money",
		description,
		idempotenceKey,
	)
}
