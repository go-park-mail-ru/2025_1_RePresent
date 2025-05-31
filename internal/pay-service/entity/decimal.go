package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

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

func (d *Decimal) ParseFromString(s string) error {
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
		return d.ParseFromString(v)
	case []byte:
		return d.ParseFromString(string(v))
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

func (d *Decimal) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		return d.ParseFromString(s)
	}

	var f float64
	if err := json.Unmarshal(b, &f); err == nil {
		return d.ParseFromString(fmt.Sprintf("%f", f))
	}

	return fmt.Errorf("invalid decimal json: %s", string(b))
}
