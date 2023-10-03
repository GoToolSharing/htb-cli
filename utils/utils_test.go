package utils

import (
	"testing"
)

func TestCheckHTBToken(t *testing.T) {
	r, w := SetOutputTest()
	defer w.Close()
	defer r.Close()
	token := GetHTBToken()
	if token == "" {
		t.Fatal("Expected token but got empty variable")
	}
}

func TestGetMachineType(t *testing.T) {
	r, w := SetOutputTest()
	defer w.Close()
	defer r.Close()
	var machine_id interface{}
	machine_id = "558"
	output := GetMachineType(machine_id, "")
	if output != "retired" && output != "active" {
		t.Fatalf("Expected \"%s\" but got \"%s\"", "active / retired", output)
	}
}

func TestGetUserSubscription(t *testing.T) {
	r, w := SetOutputTest()
	defer w.Close()
	defer r.Close()
	output := GetUserSubscription("")
	if output != "vip+" && output != "vip" && output != "free" {
		t.Fatalf("Expected \"%s\" but got \"%s\"", "free / vip / vip+", output)
	}
}

// func TestGetActiveMachineID(t *testing.T) {
// 	r, w := SetOutputTest()
// 	defer w.Close()
// 	defer r.Close()
// 	output := GetActiveMachineID("")
// 	if output != "" {
// 		t.Fatalf("Expected \"%s\" but got \"%s\"", "", output)
// 	}
// }

func TestSearchItemIDByNameMachine(t *testing.T) {
	r, w := SetOutputTest()
	defer w.Close()
	defer r.Close()
	output := SearchItemIDByName("Sau", "", "Machine")
	if output != "551" {
		t.Fatalf("Expected \"%s\" but got \"%s\"", "551", output)
	}
}

func TestSearchItemIDByNameChallenge(t *testing.T) {
	r, w := SetOutputTest()
	defer w.Close()
	defer r.Close()
	output := SearchItemIDByName("test", "", "Challenge")
	if output != "173" {
		t.Fatalf("Expected \"%s\" but got \"%s\"", "173", output)
	}
}
