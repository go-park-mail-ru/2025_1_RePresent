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
	ID            int       `json:"id"`
	TransactionID string    `json:"transactionId"`
	UserID        int       `json:"user_id"`
	Amount        float64   `json:"amount"`
	Type          string    `json:"type"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}
