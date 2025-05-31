package entity

import (
	"database/sql/driver"
	"fmt"
	"strings"

	"gopkg.in/inf.v0"
)

type Decimal struct {
	*inf.Dec
}

func NewDec(s string) (*Decimal, error) {
	d := new(inf.Dec)
	if _, ok := d.SetString(s); !ok {
		return nil, fmt.Errorf("invalid decimal format: %s", s)
	}
	return &Decimal{Dec: d}, nil
}

func NewDecWithoutErr(s string) *Decimal {
	d := new(inf.Dec)
	return &Decimal{Dec: d}
}

func (d *Decimal) Scan(value interface{}) error {
	if value == nil {
		d.Dec = inf.NewDec(0, 0)
		return nil
	}

	var s string
	switch v := value.(type) {
	case []byte:
		s = string(v)
	case string:
		s = v
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}

	s = strings.TrimLeft(s, "+")

	var err error
	d.Dec, err = ParseDecimal(s)
	return err
}

func ParseDecimal(s string) (*inf.Dec, error) {
	d := new(inf.Dec)
	if _, ok := d.SetString(s); !ok {
		return nil, fmt.Errorf("invalid decimal format: %s", s)
	}
	return d, nil
}

func (d Decimal) Value() (driver.Value, error) {
	if d.Dec == nil {
		return "0", nil
	}
	return d.String(), nil
}

func (d *Decimal) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "null" || s == "" {
		d.Dec = inf.NewDec(0, 0)
		return nil
	}

	var err error
	d.Dec, err = ParseDecimal(s)
	return err
}

func (d Decimal) MarshalJSON() ([]byte, error) {
	if d.Dec == nil {
		return []byte(`"0"`), nil
	}
	return []byte(`"` + d.String() + `"`), nil
}
