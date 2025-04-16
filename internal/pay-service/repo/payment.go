package repo

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"log"
	entity "retarget/internal/pay-service/entity"

	_ "github.com/lib/pq"

	"errors"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrInvalidAmount = errors.New("invalid amount")
)

// func (t entity.Transaction) Error() string {
// 	//TODO implement me
// 	panic("implement me")
// }

type PaymentRepositoryInterface interface {
	GetPaymentByUserId(id int) ([]*entity.Payment, error)
}

type PaymentRepository struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

func NewPaymentRepository(username, password, dbname, host string, port int, sslmode string, logger *zap.SugaredLogger) *PaymentRepository {
	paymentRepo := &PaymentRepository{}
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=%s",
		username, password, dbname, host, port, sslmode))
	if err != nil {
		log.Fatal(err)
	}
	paymentRepo.logger = logger
	paymentRepo.db = db
	return paymentRepo
}

func (r *PaymentRepository) GetBalanceByUserId(id int, requestID string) (float64, error) {
	startTime := time.Now()
	query := "SELECT balance FROM auth_user WHERE id = $1"

	r.logger.Debugw("Getting user balance",
		"query", query,
		"userID", id,
	)

	var balance float64
	err := r.db.QueryRow(query, id).Scan(&balance)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		r.logger.Debugw("User not found",
			"userID", id,
			"timeTakenMs", time.Since(startTime).Milliseconds(),
		)
		return 0, fmt.Errorf("user with id %d not found", id)
	case err != nil:
		r.logger.Debugw("Failed to get balance",
			"userID", id,
			"error", err.Error(),
			"timeTakenMs", time.Since(startTime).Milliseconds(),
		)
		return 0, fmt.Errorf("error getting balance: %w", err)
	default:
		r.logger.Debugw("Balance retrieved successfully",
			"request_id", requestID,
			"userID", id,
			"balance", balance,
			"timeTakenMs", time.Since(startTime).Milliseconds(),
		)
		return balance, nil
	}
}

func (r *PaymentRepository) UpdateBalance(userID int, amount float64, requestID string) (float64, error) {
	startTime := time.Now()
	query := `
        UPDATE auth_user 
        SET balance = balance + $1
        WHERE id = $2
        RETURNING balance
    `

	r.logger.Debugw("Starting balance update",
		"query", query,
		"userID", userID,
		"amount", amount,
	)

	var newBalance float64
	err := r.db.QueryRow(query, amount, userID).Scan(&newBalance)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Debugw("User not found for balance update",
				"userID", userID,
				"timeTakenMs", time.Since(startTime).Milliseconds(),
			)
			return 0, ErrUserNotFound
		}

		r.logger.Debugw("Balance update failed",
			"userID", userID,
			"amount", amount,
			"error", err.Error(),
			"timeTakenMs", time.Since(startTime).Milliseconds(),
		)
		return 0, fmt.Errorf("update balance failed: %w", err)
	}

	r.logger.Debugw("Balance updated successfully",
		"request_id", requestID,
		"userID", userID,
		"newBalance", newBalance,
		"timeTakenMs", time.Since(startTime).Milliseconds(),
	)

	return newBalance, nil
}

func (r *PaymentRepository) CreateTransaction(tx *entity.Transaction, requestID string) error {
	query := `
        INSERT INTO transaction
            (transaction_id, user_id, amount, type, status, created_at)
        VALUES 
            ($1, $2, $3, $4, $5, $6)

    `

	_, err := r.db.Exec(query, tx.TransactionID, tx.UserID, tx.Amount, tx.Type, tx.Status, tx.CreatedAt)
	return err
}

func (r *PaymentRepository) TopUpAccount(userID int, amount int64, requestID string) error {
	transactionID := uuid.New().String()

	return r.CreateTransaction(&entity.Transaction{
		TransactionID: transactionID,
		UserID:        userID,
		Amount:        amount,
		Type:          "topup",
		Status:        "completed",
		CreatedAt:     time.Now(),
	}, requestID)
}

func (r *PaymentRepository) GetLastTransaction(userID int, requestID string) (*entity.Transaction, error) {
	var tx entity.Transaction

	err := r.db.QueryRow(
		"SELECT * FROM transaction WHERE user_id = $1 ORDER BY created_at DESC LIMIT 1",
		userID,
	).Scan(&tx.ID, &tx.TransactionID, &tx.UserID, &tx.Amount, &tx.Type, &tx.Status, &tx.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &tx, nil
}

func (r *PaymentRepository) GetTransactionByID(transactionID string, requestID string) (*entity.Transaction, error) {
	var tx entity.Transaction

	err := r.db.QueryRow(
		"SELECT * FROM transaction WHERE transaction_id = $1",
		transactionID,
	).Scan(&tx.ID, &tx.TransactionID, &tx.UserID, &tx.Amount, &tx.Type, &tx.Status, &tx.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("transaction with ID %s not found", transactionID)
		}
		return nil, err
	}
	return &tx, nil
}

func (r *PaymentRepository) CloseConnection() error {
	return r.db.Close()
}
