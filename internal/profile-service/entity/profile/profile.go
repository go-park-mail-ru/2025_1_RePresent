package profile

import "errors"

var ErrProfileNotFound = errors.New("Profile not found")

type Profile struct {
	ID          int    `json:"id"`
	Username    string `json:"username" validate:"required,min=5,max=35"`
	Email       string `json:"email" validate:"required,email"`
	Description string `json:"description" validate:"required,min=0,max=200"`
	Balance     int    `json:"balance" validate:"required,gte=0"`
	Role        int    `json:"role" validate:"required,gte=1,lte=2"`
}
