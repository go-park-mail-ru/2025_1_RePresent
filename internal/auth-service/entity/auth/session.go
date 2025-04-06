package entity

import (
	"time"
)

type Session struct {
	ID        string    `json:"id"`
	UserID    int       `json:"user_id"`
	Role      int       `json:"role"`
	Expires   time.Time `json:"expires"`
	CreatedAt time.Time `json:"created_at"`
}

var ErrSessionNotFound = &Error{"Session not found"}

type Error struct {
	Msg string
}

func (e *Error) Error() string {
	return e.Msg
}
