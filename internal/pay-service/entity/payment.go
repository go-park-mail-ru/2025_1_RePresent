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
