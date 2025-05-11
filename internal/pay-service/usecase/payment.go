package payment

import (
	"context"
	"errors"
	"fmt"
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
	BalanceLimit        = 100.00
)

type PaymentUsecase struct {
	logger            *zap.SugaredLogger
	PaymentRepository *repo.PaymentRepository
	NoticeRepository  *notice.NoticeRepository
	AttemptRepository *attempt.AttemptRepository
}

func NewPayUsecase(zapLogger *zap.SugaredLogger, payRepository *repo.PaymentRepository, noticeRepository *notice.NoticeRepository, attemptRepository *attempt.AttemptRepository) *PaymentUsecase {
	return &PaymentUsecase{logger: zapLogger, PaymentRepository: payRepository, NoticeRepository: noticeRepository, AttemptRepository: attemptRepository}
}

func (u PaymentUsecase) GetBalanceByUserId(id int, requestID string) (float64, error) {
	balance, err := u.PaymentRepository.GetBalanceByUserId(id, requestID)
	if err != nil {
		return 0, err
	}
	return balance, nil
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
	}()

	return nil
	// return uc.PaymentRepository.GetLastTransaction(userID)
}

func (uc *PaymentUsecase) GetTransactionByID(transactionID string, requestID string) (*entity.Transaction, error) {
	return uc.PaymentRepository.GetTransactionByID(transactionID, requestID)
}

func (uc *PaymentUsecase) RegUserActivity(user_banner_id, user_slot_id, amount int) error {
	user_id, err := uc.PaymentRepository.RegUserActivity(user_banner_id, user_slot_id, amount)
	if err != nil {
		return err
	}
	balance, err := uc.CheckBalance(user_id)
	if err == errTooLittleBalance {
		fmt.Println(balance)
		go uc.requireSend(user_id, strconv.FormatFloat(balance, 'f', 2, 64))
		return nil
	}
	if err != nil {
		return err
	}
	return nil
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
			return uc.NoticeRepository.SendLowBalanceNotification(userID, message)
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
