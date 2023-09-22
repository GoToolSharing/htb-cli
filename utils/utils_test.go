package utils

import (
	"testing"
)

func TestCheckHTBToken(t *testing.T) {
	token := GetHTBToken()
	if token == "" {
		t.Fatal("Expected token but got empty variable")
	}
}
