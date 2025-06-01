package entity

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/inf.v0"
)

func TestDecimal_Value(t *testing.T) {
	d1 := Decimal{}
	v1, err := d1.Value()
	assert.NoError(t, err)
	assert.Equal(t, "0", v1)

	d2 := Decimal{Dec: inf.NewDec(12345, -2)} // 123.45
	v2, err := d2.Value()
	assert.NoError(t, err)
	assert.Equal(t, "1234500", v2)
}

func TestDecimal_ParseFromString(t *testing.T) {
	var d Decimal
	err := d.ParseFromString("99.99")
	assert.NoError(t, err)
	assert.Equal(t, "99.99", d.String())

	err = d.ParseFromString("bad")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid decimal format")
}

func TestDecimal_Scan(t *testing.T) {
	var d Decimal

	// nil case
	err := d.Scan(nil)
	assert.NoError(t, err)
	assert.Equal(t, "0", d.String())

	// string
	err = d.Scan("1.23")
	assert.NoError(t, err)
	assert.Equal(t, "1.23", d.String())

	// bytes
	err = d.Scan([]byte("4.56"))
	assert.NoError(t, err)
	assert.Equal(t, "4.56", d.String())

	// wrong type
	err = d.Scan(123)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot convert")
}

func TestDecimal_JSON(t *testing.T) {
	// Marshal nil
	d1 := Decimal{}
	b1, err := json.Marshal(d1)
	assert.NoError(t, err)
	assert.JSONEq(t, `"0"`, string(b1))

	// Marshal non-nil
	d2 := Decimal{Dec: inf.NewDec(789, -2)} // 7.89
	b2, err := json.Marshal(d2)
	assert.NoError(t, err)
	assert.JSONEq(t, `"78900"`, string(b2))

	// Unmarshal string
	var d3 Decimal
	err = json.Unmarshal([]byte(`"3.21"`), &d3)
	assert.NoError(t, err)
	assert.Equal(t, "3.21", d3.String())

	// Unmarshal number
	var d4 Decimal
	err = json.Unmarshal([]byte(`4.50`), &d4)
	assert.NoError(t, err)
	assert.Equal(t, "4.500000", d4.String()) // fmt.Sprintf("%f",4.50)

	// Unmarshal invalid
	var d5 Decimal
	err = json.Unmarshal([]byte(`{"foo":1}`), &d5)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid decimal json")
}
