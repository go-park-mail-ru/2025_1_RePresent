package payment

import (
	"errors"
	"fmt"
	"retarget/internal/pay-service/entity"
	"retarget/internal/pay-service/repo"
	"retarget/internal/pay-service/repo/notice"
)

var (
	errTooLittleBalance = errors.New("balance of user is need to top up")
	BalanceLimit        = 100.00
)

type PaymentUsecase struct {
	PaymentRepository *repo.PaymentRepository
	NoticeRepository  *notice.NoticeRepository
}

func NewPayUsecase(payRepository *repo.PaymentRepository, noticeRepository *notice.NoticeRepository) *PaymentUsecase {
	return &PaymentUsecase{PaymentRepository: payRepository, NoticeRepository: noticeRepository}
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
		// Здесь будет вызов Kafka репозитория
		return nil
	}
	if err != nil {
		return err
	}
	return nil
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
