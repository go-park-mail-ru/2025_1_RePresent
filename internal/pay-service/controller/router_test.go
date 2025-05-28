package http

import (
	"testing"
)

func TestSetupRoutes_NotNil(t *testing.T) {
	r := SetupRoutes(nil, nil)
	if r == nil {
		t.Fatal("expected router, got nil")
	}
}
