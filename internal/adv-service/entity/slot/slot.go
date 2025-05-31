package slot

import (
	"time"

	"gopkg.in/inf.v0"
)

type Slot struct {
	Link       string    `json:"link"`
	UserID     int       `json:"user_id,omitempty"`
	SlotName   string    `json:"slot_name" validate:"required,min=1,max=100"`
	FormatCode int       `json:"format_code" validate:"required,min=1"`
	MinPrice   inf.Dec   `json:"min_price" validate:"required,min=0"`
	IsActive   bool      `json:"is_active" validate:"required"`
	CreatedAt  time.Time `json:"created_at"`
}
