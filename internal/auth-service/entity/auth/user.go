package entity

import (
	"github.com/go-playground/validator/v10"
	"gopkg.in/inf.v0"
)

type User struct {
	ID          int      `json:"id"`
	Username    string   `json:"username" validate:"required,min=5,max=50"`
	Email       string   `json:"email" validate:"email,required"`
	Password    []byte   `json:"password" validate:"required,min=8"`
	Description string   `json:"description"`
	Balance     *inf.Dec `json:"balance"` // Убираем валидацию, делаем указателем
	Role        int      `json:"role" validate:"required,gte=1,lte=2"`
}

func ValidateUser(user *User) error {
	validate := validator.New()

	// Проверяем Balance отдельно
	if user.Balance == nil {
		user.Balance = inf.NewDec(0, 0)
	}

	// Валидируем только поддерживаемые поля
	type ValidateStruct struct {
		Username string `validate:"required,min=5,max=50"`
		Email    string `validate:"email,required"`
		Password []byte `validate:"required,min=8"`
		Role     int    `validate:"required,gte=1,lte=2"`
	}

	temp := ValidateStruct{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		Role:     user.Role,
	}

	return validate.Struct(temp)
}
