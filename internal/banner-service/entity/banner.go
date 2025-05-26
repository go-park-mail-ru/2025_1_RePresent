package entity

import (
	entity "retarget/pkg/entity"
)

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

var DefaultBanner = Banner{
	ID:          -1,
	OwnerID:     -1,
	Title:       "Здесь ничего нет",
	Content:     "0a8a2d2cecbddf2d2ae18cb36483277c", //хардкоднуто
	Description: "Но тут могла быть ваша реклама",
	Status:      0,
	Balance:     0,
	Link:        "http://test.re-target.ru",
	Deleted:     false,
	MaxPrice:    *entity.NewDecWithoutErr("0.0"),
}
