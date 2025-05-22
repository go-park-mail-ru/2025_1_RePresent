package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	httpClient        *http.Client
}

func NewPayUsecase(
	zapLogger *zap.SugaredLogger,
	payRepository *repo.PaymentRepository,
	noticeRepository *notice.NoticeRepository,
	attemptRepository *attempt.AttemptRepository,
	shopID, secretKey string,
	httpClient *http.Client,
) *PaymentUsecase {
	return &PaymentUsecase{
		logger:            zapLogger,
		PaymentRepository: payRepository,
		NoticeRepository:  noticeRepository,
		AttemptRepository: attemptRepository,
		shopID:            shopID,
		secretKey:         secretKey,
		httpClient:        httpClient,
	}
}

func (u *PaymentUsecase) GetBalanceByUserId(userID int, requestID string) (float64, error) {
	balik, err := u.PaymentRepository.GetBalanceByUserId(userID, requestID)
	if err != nil {
		return 0, err
	}

	// pending, err := u.PaymentRepository.GetPendingTransactions(userID)
	// if err != nil {
	// 	return 0, fmt.Errorf("error while loaing pending tx: %w", err)
	// }

	// for _, tx := range pending {
	// 	statusStr, err := u.getYooPaymentStatus(tx.TransactionID)
	// 	if err != nil {
	// 		continue
	// 	}
	// 	statusInt := mapYooStatus(statusStr)
	// 	if statusInt == 1 && tx.Status != 1 {
	// 		if _, err := u.PaymentRepository.UpdateBalance(userID, tx.Amount, requestID); err == nil {
	// 			_ = u.PaymentRepository.UpdateTransactionStatus(tx.TransactionID, statusInt)
	// 		}
	// 	}
	// }

	return balik, nil
	// return u.PaymentRepository.GetBalanceByUserId(userID, requestID)
}

func mapYooStatus(s string) int {
	switch s {
	case "pending":
		return 0
	case "succeeded":
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

func (uc *PaymentUsecase) TopUpBalance(userID int, amount int64, requestID string) error {
	if amount <= 0 {
		return repo.ErrInvalidAmount
	}

	_, err := uc.PaymentRepository.UpdateBalance(userID, float64(amount), requestID)
	if err != nil {
		return err
	}

	go func() {
		if err := uc.AttemptRepository.ResetAttemptsByUserID(userID); err != nil {
			uc.logger.Errorw("failed to reset attempts after top up",
				"user_id", userID,
				"error", err)
		}

		err = uc.NoticeRepository.SendTopUpBalanceEvent(userID, float64(amount))
		uc.logger.Errorw("failed to send topUp message after top up",
			"user_id", userID,
			"error", err)
	}()

	return nil
	// return uc.PaymentRepository.GetLastTransaction(userID)
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
	PaymentMethodData struct {
		Type string `json:"type"`
	} `json:"payment_method_data"`
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
	reqBody := YooPaymentRequest{Description: description}
	reqBody.Amount.Value = value
	reqBody.Amount.Currency = currency
	reqBody.PaymentMethodData.Type = "yoo_money"
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
		Type:          "yoo_money",
		Status:        statusInt,
	}
	if err := u.PaymentRepository.CreateTransaction(trx); err != nil {
		return "", fmt.Errorf("save transaction: %w", err)
	}

	return out.Confirmation.ConfirmationURL, nil
}
