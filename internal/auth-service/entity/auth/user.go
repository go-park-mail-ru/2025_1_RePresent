package entity

import (
	"github.com/go-playground/validator/v10"
	"gopkg.in/inf.v0"
)

func ValidateUser(user *User) error {
	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		return err
	}
	return nil
}

type User struct {
	ID          int     `json:"id"`
	Username    string  `json:"username" validate:"required,min=5,max=50"`
	Email       string  `json:"email" validate:"email,required"`
	Password    []byte  `json:"password" validate:"required,min=8"`
	Description string  `json:"description"`
	Balance     inf.Dec `json:"balance" validate:"gte=0"`
	Role        int     `json:"role" validate:"required,gte=1,lte=2"`
}
