package dto

import (
	"github.com/google/uuid"
	"gopkg.in/inf.v0"
)

type CreateRequest struct {
	SlotName   string  `json:"slot_name" validate:"required,min=1,max=100"`
	FormatCode int     `json:"format_code" validate:"required,min=1"`
	MinPrice   inf.Dec `json:"min_price" validate:"required"`
	IsActive   bool    `json:"is_active" validate:"required"`
}

type UpdateRequest struct {
	Link       uuid.UUID `json:"link" validate:"required"`
	SlotName   string    `json:"slot_name" validate:"required,min=1,max=100"`
	FormatCode int       `json:"format_code" validate:"required"`
	MinPrice   inf.Dec   `json:"min_price" validate:"required"`
	IsActive   bool      `json:"is_active" validate:"required"`
}
