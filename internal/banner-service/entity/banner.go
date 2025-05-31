package entity

import (
	model "retarget/internal/banner-service/easyjsonModels"
	entity "retarget/pkg/entity"
)

var DefaultBanner = model.Banner{
	ID:          -1,
	OwnerID:     -1,
	Title:       "Здесь ничего нет",
	Content:     "144630119b8384662e7c189c300f7311", //хардкоднуто
	Description: "Но тут могла быть ваша реклама",
	Status:      0,
	Balance:     0,
	Link:        "http://test.re-target.ru",
	Deleted:     false,
	MaxPrice:    *entity.NewDecWithoutErr("0.0"),
}
