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

func (d Decimal) Value() (driver.Value, error) {
	if d.Dec == nil {
		return "0", nil
	}
	return d.String(), nil
}

func (d *Decimal) parseFromString(s string) error {
	dec := inf.NewDec(0, 0)
	if _, ok := dec.SetString(s); !ok {
		return fmt.Errorf("invalid decimal format: %s", s)
	}
	d.Dec = dec
	return nil
}

func (d *Decimal) Scan(value interface{}) error {
	switch v := value.(type) {
	case nil:
		d.Dec = inf.NewDec(0, 0)
		return nil
	case string:
		return d.parseFromString(v)
	case []byte:
		return d.parseFromString(string(v))
	default:
		return fmt.Errorf("cannot convert %T to Decimal", value)
	}
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
