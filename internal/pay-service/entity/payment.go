package entity

import "time"

type Payment struct {
	ID        int       `json:"id"`
	OwnerID   int       `json:"owner"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"timestamp"`
	Status    int       `json:"status"`
	Balance   float64   `json:"balance"`
}

type Transaction struct {
	ID            int       `db:"id"`
	TransactionID string    `db:"transaction_id"`
	UserID        int       `db:"user_id"`
	Amount        float64   `db:"amount"`
	Type          string    `db:"type"`
	Status        int       `db:"status"`
	CreatedAt     time.Time `db:"created_at"`
}
