package profile

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"gopkg.in/inf.v0"
)

var ErrProfileNotFound = errors.New("Profile not found")

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

func (d Decimal) MarshalJSON() ([]byte, error) {
	if d.Dec == nil {
		return []byte(`"0"`), nil
	}
	return []byte(`"` + d.String() + `"`), nil
}

type Profile struct {
	ID          int     `json:"id"`
	Username    string  `json:"username" validate:"required,min=5,max=50"`
	Email       string  `json:"email" validate:"required,email"`
	Description string  `json:"description" validate:"required,min=0,max=200"`
	Balance     Decimal `json:"balance" validate:"required,gte=0"`
	Role        int     `json:"role" validate:"required,gte=1,lte=2"`
}
