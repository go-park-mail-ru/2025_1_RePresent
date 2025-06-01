package entity

import (
	"testing"

	"gopkg.in/inf.v0"
)

func TestDecimalValue_Nil(t *testing.T) {
	var d Decimal
	v, err := d.Value()
	if err != nil || v != "0" {
		t.Errorf("expected '0', got %v, err %v", v, err)
	}
}

func TestDecimalValue_NonNil(t *testing.T) {
	d := Decimal{inf.NewDec(123, 0)}
	v, err := d.Value()
	if err != nil || v != "123" {
		t.Errorf("expected '123', got %v, err %v", v, err)
	}
}

func TestDecimalScan_Types(t *testing.T) {
	var d Decimal
	if err := d.Scan(nil); err != nil {
		t.Errorf("scan nil failed: %v", err)
	}
	if err := d.Scan("45.6"); err != nil || d.String() != "45.6" {
		t.Errorf("scan string failed: %v, got %s", err, d.String())
	}
	if err := d.Scan([]byte("7.89")); err != nil || d.String() != "7.89" {
		t.Errorf("scan bytes failed: %v, got %s", err, d.String())
	}
	if err := d.Scan(123); err == nil {
		t.Error("expected error on unsupported type")
	}
}

func TestParseFromString_Invalid(t *testing.T) {
	var d Decimal
	if err := d.parseFromString("bad"); err == nil {
		t.Error("expected parse error")
	}
}

func TestValidateUser(t *testing.T) {
	valid := &User{
		Username: "tester",
		Email:    "a@b.com",
		Password: []byte("strongpwd"),
		Role:     1,
	}
	if err := ValidateUser(valid); err != nil {
		t.Errorf("expected valid user, got %v", err)
	}
	invalid := &User{Username: "u", Email: "x", Password: []byte("short"), Role: 0}
	if err := ValidateUser(invalid); err == nil {
		t.Error("expected validation failure")
	}
}
