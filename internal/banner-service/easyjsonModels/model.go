//go:generate easyjson -all model.go

package model

import "retarget/pkg/entity"

type CreateUpdateBannerRequest struct {
	Title       string         `json:"title" validate:"required,min=3,max=30"`
	Description string         `json:"description" validate:"max=100"`
	Content     string         `json:"content" validate:"required,len=32"`
	Link        string         `json:"link" validate:"required,max=100"`
	Status      int            `json:"status"`
	MaxPrice    entity.Decimal `json:"max_price" validate:"gt_decimal_01"`
}

type Banner struct {
	ID          int            `json:"id"`
	OwnerID     int            `json:"owner"`
	Title       string         `json:"title"`
	Content     string         `json:"content"`
	Description string         `json:"description"`
	Status      int            `json:"status"`
	Balance     int            `json:"balance"`
	Link        string         `json:"link"`
	Deleted     bool           `json:"deleted"`
	MaxPrice    entity.Decimal `json:"max_price"`
}
