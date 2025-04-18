package payment

import (
	"retarget/internal/pay-service/entity"
	"retarget/internal/pay-service/repo"
)

type PaymentUsecase struct {
	PaymentRepository *repo.PaymentRepository
}

func NewPayUsecase(payRepository *repo.PaymentRepository) *PaymentUsecase {
	return &PaymentUsecase{PaymentRepository: payRepository}
}

func (u PaymentUsecase) GetBalanceByUserId(id int) (float64, error) {
	balance, err := u.PaymentRepository.GetBalanceByUserId(id)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

func (uc *PaymentUsecase) TopUpBalance(userID int, amount int64) error {
	if amount <= 0 {
		return repo.ErrInvalidAmount
	}

	_, err := uc.PaymentRepository.UpdateBalance(userID, float64(amount))
	if err != nil {
		return err
	}

	return nil
	// return uc.PaymentRepository.GetLastTransaction(userID)
}

func (uc *PaymentUsecase) GetTransactionByID(transactionID string) (*entity.Transaction, error) {
	return uc.PaymentRepository.GetTransactionByID(transactionID)
}
