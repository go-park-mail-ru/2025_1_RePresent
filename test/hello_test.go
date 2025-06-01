package test

import (
	"testing"
)

func TestHelloWorld(t *testing.T) {
	expected := "Hello, world!"
	actual := "Hello, world!"

	if actual != expected {
		t.Errorf("Expected %q but got %q", expected, actual)
	}
}
