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

func TestGetUserSubscription(t *testing.T) {
	log.SetOutput(io.Discard)
	output := GetUserSubscription("")
	if output != "vip+" && output != "vip" && output != "free" {
		t.Fatalf("Expected \"%s\" but got \"%s\"", "free / vip / vip+", output)
	}
}

func TestGetActiveMachineID(t *testing.T) {
	log.SetOutput(io.Discard)
	output := GetActiveMachineID("")
	if output != "" {
		t.Fatalf("Expected \"%s\" but got \"%s\"", "", output)
	}
}

// func TestSearchItemIDByNameMachine(t *testing.T) {
// 	log.SetOutput(io.Discard)
// 	output := SearchItemIDByName("Sau", "", "Machine")
// 	if output != "" {
// 		t.Fatalf("Expected \"%s\" but got \"%s\"", "", output)
// 	}
// }

// func TestSearchItemIDByNameChallenge(t *testing.T) {
// 	log.SetOutput(io.Discard)
// 	output := SearchItemIDByName("test", "", "Challenge")
// 	if output != "" {
// 		t.Fatalf("Expected \"%s\" but got \"%s\"", "", output)
// 	}
// }
