package repo

import (
	"database/sql"
	"fmt"
	"time"

	"go.uber.org/zap"

	"log"
	"retarget/internal/pay-service/entity"

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
	db         *sql.DB
	logger     *zap.SugaredLogger
	clickhouse *sql.DB
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
		"request_id", requestID,
		"query", query,
		"userID", id,
	)

	var balance float64
	err := r.db.QueryRow(query, id).Scan(&balance)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		r.logger.Debugw("User not found",
			"request_id", requestID,
			"userID", id,
			"timeTakenMs", time.Since(startTime).Milliseconds(),
		)
		return 0, fmt.Errorf("user with id %d not found", id)
	case err != nil:
		r.logger.Debugw("Failed to get balance",
			"request_id", requestID,
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
		"request_id", requestID,
		"query", query,
		"userID", userID,
		"amount", amount,
	)

	var newBalance float64
	err := r.db.QueryRow(query, amount, userID).Scan(&newBalance)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Debugw("User not found for balance update",
				"request_id", requestID,
				"userID", userID,
				"timeTakenMs", time.Since(startTime).Milliseconds(),
			)
			return 0, ErrUserNotFound
		}

		r.logger.Debugw("Balance update failed",
			"request_id", requestID,
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

func (r *PaymentRepository) CreateTransaction(trx entity.Transaction) error {
	const q = `
    INSERT INTO transaction (
        transaction_id, user_id, amount, type, status
    ) VALUES ($1, $2, $3, $4, $5)
    `
	_, err := r.db.Exec(q,
		trx.TransactionID,
		trx.UserID,
		trx.Amount,
		trx.Type,
		trx.Status,
	)
	return err
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

func (r *PaymentRepository) RegUserActivity(user_banner_id, user_slot_id int, amount entity.Decimal) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return -1, fmt.Errorf("failed to begin transaction: %w", err)
	}

	_, err = tx.Exec(`
        UPDATE auth_user 
        SET balance = balance - $1 
        WHERE id = $2`,
		amount,
		user_slot_id)
	if err != nil {
		tx.Rollback()
		return -1, fmt.Errorf("failed to update first user balance: %w", err)
	}

	_, err = tx.Exec(`
        UPDATE auth_user 
        SET balance = balance + $1 
        WHERE id = $2`,
		amount,
		user_banner_id)
	if err != nil {
		tx.Rollback()
		return -1, fmt.Errorf("failed to update second user balance: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return -1, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return user_banner_id, nil
}

func (r *PaymentRepository) GetPendingTransactions(userID int) ([]entity.Transaction, error) {
	const q = `
    	SELECT id, transaction_id, user_id, amount, type, status, created_at
    	FROM transaction
    	WHERE user_id = $1 AND status = 0
    `
	rows, err := r.db.Query(q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []entity.Transaction
	for rows.Next() {
		var tx entity.Transaction
		if err := rows.Scan(&tx.ID, &tx.TransactionID, &tx.UserID, &tx.Amount, &tx.Type, &tx.Status, &tx.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, tx)
	}
	return list, rows.Err()
}

func (r *PaymentRepository) UpdateTransactionStatus(transactionID string, status int) error {
	const q = `
    	UPDATE transaction
    	SET status = $1
    	WHERE transaction_id = $2
    `
	_, err := r.db.Exec(q, status, transactionID)
	return err
}

func (r *PaymentRepository) CloseConnection() error {
	return r.db.Close()
}
