//go:generate easyjson -all model.go

package model

import (
	"retarget/pkg/entity"
	"time"

	"gopkg.in/inf.v0"
)

type IFrame struct {
	ImageSrc    string
	Link        string
	Title       string
	Description string
	Banner      int64
	Slot        string
}

type CreateSlotResponse struct {
	Link       string    `json:"link"`
	SlotName   string    `json:"slot_name"`
	FormatCode int       `json:"format_code"`
	MinPrice   string    `json:"min_price"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
}

type GetSlotResponse struct {
	Link       string    `json:"link"`
	SlotName   string    `json:"slot_name"`
	FormatCode int       `json:"format_code"`
	MinPrice   inf.Dec   `json:"min_price"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
}

type EditSlotResponse struct {
	Link       string `json:"link"`
	SlotName   string `json:"slot_name"`
	FormatCode int    `json:"format_code"`
	MinPrice   string `json:"min_price"`
	IsActive   bool   `json:"is_active"`
}

type ResponseWithSlot struct {
	Service entity.ServiceResponse `json:"service"`
	Body    CreateSlotResponse     `json:"body,omitempty"`
}

type ResponseWithSlots struct {
	Service entity.ServiceResponse `json:"service"`
	Body    []GetSlotResponse      `json:"body,omitempty"`
}
