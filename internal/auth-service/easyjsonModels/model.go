//go:generate easyjson -all model.go

package model

import (
	"retarget/pkg/entity"

	"gopkg.in/inf.v0"
)

//easyjson:json
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Email    string `json:"email" validate:"email,required"`
	// Code     string `json:"code" validate:"required,len=6"`
	Password string `json:"password" validate:"required,min=8"`
	Role     int    `json:"role" validate:"required,gte=1,lte=2"`
}

//easyjson:json
type RegisterConfirmRequest struct {
	Email string `json:"email" validate:"email,required"`
}

//easyjson:json
type LoginRequest struct {
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required,min=8"`
	Role     int    `json:"role" validate:"required,gte=1,lte=2"`
}

//easyjson:json
type UserResponse struct {
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Balance  inf.Dec `json:"balance"`
	Role     int     `json:"role"`
}

//easyjson:json
type UserResponseWithErr struct {
	Service entity.Response `json:"service"`
	Body    UserResponse    `json:"body"`
}

//easyjson:json
type ErrorRequest struct {
	ErrorText string `json:"error"`
}
