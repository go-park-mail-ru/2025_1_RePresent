package entity

import "testing"

func TestErrorString(t *testing.T) {
	e := &Error{"oops"}
	if e.Error() != "oops" {
		t.Errorf("expected 'oops', got '%s'", e.Error())
	}
}

func TestErrSessionNotFound(t *testing.T) {
	if ErrSessionNotFound.Error() != "Session not found" {
		t.Errorf("unexpected ErrSessionNotFound: %s", ErrSessionNotFound.Error())
	}
}
