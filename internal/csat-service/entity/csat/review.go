package csat

import "time"

type Review struct {
	ID        string    `json:"id,omitempty"`
	UserID    int       `json:"user_id" validate:"required,gte=0"`
	Question  string    `json:"question" validate:"required"`
	Page      string    `json:"page" validate:"required"`
	Comment   string    `json:"comment" validate:"lte=200"`
	Rating    int       `json:"rating" validate:"required,gte=0,lte=10"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
