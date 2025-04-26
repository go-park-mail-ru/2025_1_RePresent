package csat

type Review struct {
	Page     string `json:"page_id" validate:"required"`
	Question string `json:"question" validate:"required"`
	Rating   int    `json:"rating" validate:"required,gte=0,lte=10"`
	Comment  string `json:"comment"`
}
