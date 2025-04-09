package repo

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"time"

	"log"
	"retarget/internal/pay-service/entity"

	_ "github.com/lib/pq"

	"errors"
)

type Transaction struct {
	ID            int       `json:"id"`
	TransactionID string    `json:"transactionId"`
	UserID        int       `json:"user_id"`
	Amount        int64     `json:"amount"`
	Type          string    `json:"type"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrInvalidAmount = errors.New("invalid amount")
)

func (t Transaction) Error() string {
	//TODO implement me
	panic("implement me")
}

type PaymentRepositoryInterface interface {
	GetPaymentByUserId(id int) ([]*entity.Payment, error)
}

type PaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(username, password, dbname, host string, port int, sslmode string) *PaymentRepository {
	paymentRepo := &PaymentRepository{}
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=%s",
		username, password, dbname, host, port, sslmode))
	if err != nil {
		log.Fatal(err)
	}
	paymentRepo.db = db
	return paymentRepo
}

func (r *PaymentRepository) GetBalanceByUserId(id int) (float64, error) {
	var balance float64

	err := r.db.QueryRow(
		"SELECT balance FROM users WHERE id = $1",
		id,
	).Scan(&balance)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return 0, fmt.Errorf("user with id %d not found", id)
	case err != nil:
		return 0, fmt.Errorf("error getting balance: %w", err)
	default:
		return balance, nil
	}
}

func (r *PaymentRepository) UpdateBalance(userID int, amount float64) (float64, error) {
	query := `
        UPDATE user 
        SET balance = balance + $1
        WHERE id = $2
        RETURNING balance
    `
	var newBalance float64
	err := r.db.QueryRow(query, amount, userID).Scan(&newBalance)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrUserNotFound
		}
		return 0, fmt.Errorf("update balance failed: %w", err)

	}
	return newBalance, nil
}

func (r *PaymentRepository) CreateTransaction(tx *Transaction) error {
	query := `
        INSERT INTO external_transaction
            (transaction_id, user_id, amount, type, status, created_at)
        VALUES 
            ($1, $2, $3, $4, $5, $6)

    `

	_, err := r.db.Exec(query, tx.TransactionID, tx.UserID, tx.Amount, tx.Type, tx.Status, tx.CreatedAt)
	return err
}

func (r *PaymentRepository) TopUpAccount(userID int, amount int64) error {
	transactionID := uuid.New().String()

	return r.CreateTransaction(&Transaction{
		TransactionID: transactionID,
		UserID:        userID,
		Amount:        amount,
		Type:          "topup",
		Status:        "completed",
		CreatedAt:     time.Now(),
	})
}

func (r *PaymentRepository) GetLastTransaction(userID int) (*Transaction, error) {
	var tx Transaction

	err := r.db.QueryRow(
		"SELECT * FROM external_transaction WHERE user_id = $1 ORDER BY created_at DESC LIMIT 1",
		userID,
	).Scan(&tx.ID, &tx.UserID, &tx.Amount, &tx.Type, &tx.Status, &tx.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &tx, nil
}

func (r *PaymentRepository) GetTransactionByID(transactionID string) (*Transaction, error) {
	var tx Transaction

	err := r.db.QueryRow(
		"SELECT * FROM external_transaction WHERE transaction_id = $1",
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
