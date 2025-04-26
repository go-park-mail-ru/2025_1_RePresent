package csat

type Review struct {
	Page     string `json:"page_id" validate:"required"`
	Question string `json:"question" validate:"required"`
	Rating   int    `json:"rating" validate:"required,gte=0,lte=10"`
	Comment  string `json:"comment" validate:"lte=200"`
	User_id  int    `json:"user_id" validate:"required,gte=0"`
}
