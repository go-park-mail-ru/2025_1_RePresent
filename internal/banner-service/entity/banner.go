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

func (d *Decimal) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		return d.parseFromString(s)
	}

	var f float64
	if err := json.Unmarshal(b, &f); err == nil {
		return d.parseFromString(fmt.Sprintf("%f", f))
	}

	return fmt.Errorf("invalid decimal json: %s", string(b))
}

type Banner struct {
	ID          int     `json:"id"`
	OwnerID     int     `json:"owner"`
	Title       string  `json:"title"`
	Content     string  `json:"content"`
	Description string  `json:"description"`
	Status      int     `json:"status"`
	Balance     int     `json:"balance"`
	Link        string  `json:"link"`
	Deleted     bool    `json:"deleted"`
	MaxPrice    Decimal `json:"max_price"`
}

var DefaultBanner = Banner{
	ID:          -1,
	OwnerID:     -1,
	Title:       "Здесь ничего нет",
	Content:     "0d93f25e17810c6fc8666718f37a5203", //хардкоднуто
	Description: "Но тут могла быть ваша реклама",
	Status:      0,
	Balance:     0,
	Link:        "http://re-target.ru",
	Deleted:     false,
}
