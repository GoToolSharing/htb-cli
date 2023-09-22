package utils

import (
	"io"
	"log"
	"testing"
)

func TestCheckHTBToken(t *testing.T) {
	token := GetHTBToken()
	if token == "" {
		t.Fatal("Expected token but got empty variable")
	}
}

func TestGetMachineType(t *testing.T) {
	log.SetOutput(io.Discard)
	var machine_id interface{}
	machine_id = "558"
	output := GetMachineType(machine_id, "")
	if output != "retired" && output != "active" {
		t.Fatalf("Expected \"%s\" but got \"%s\"", "active / retired", output)
	}
}
