package entity

import (
	"database/sql/driver"
	"fmt"

	"github.com/go-playground/validator/v10"
	"gopkg.in/inf.v0"
)

type Decimal struct {
	*inf.Dec
}

// Value - преобразует Decimal в тип, понятный PostgreSQL
func (d Decimal) Value() (driver.Value, error) {
	if d.Dec == nil {
		return "0", nil
	}
	return d.String(), nil
}

// Scan - преобразует данные из БД в Decimal
func (d *Decimal) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("cannot convert %T to Decimal", value)
	}

	d.Dec = new(inf.Dec)
	_, ok = d.SetString(str)
	if !ok {
		return fmt.Errorf("invalid decimal format: %s", str)
	}
	return nil
}

type User struct {
	ID          int     `json:"id"`
	Username    string  `json:"username" validate:"required,min=5,max=50"`
	Email       string  `json:"email" validate:"email,required"`
	Password    []byte  `json:"password" validate:"required,min=8"`
	Description string  `json:"description"`
	Balance     Decimal `json:"balance"` // Используем кастомный тип
	Role        int     `json:"role" validate:"required,gte=1,lte=2"`
}

func ValidateUser(user *User) error {
	validate := validator.New()

	if user.Balance.Dec == nil {
		user.Balance.Dec = inf.NewDec(0, 0)
	}

	return validate.Struct(user)
}
