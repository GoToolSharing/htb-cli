package utils

import (
	"testing"
)

func TestCheckHTBToken(t *testing.T) {
	r, w := SetOutputTest()
	token := GetHTBToken()
	if token == "" {
		t.Fatal("Expected token but got empty variable")
	}
	w.Close()
	r.Close()
}

func TestGetMachineType(t *testing.T) {
	r, w := SetOutputTest()
	var machine_id interface{}
	machine_id = "558"
	output := GetMachineType(machine_id, "")
	if output != "retired" && output != "active" {
		t.Fatalf("Expected \"%s\" but got \"%s\"", "active / retired", output)
	}
	w.Close()
	r.Close()
}

func TestGetUserSubscription(t *testing.T) {
	r, w := SetOutputTest()
	output := GetUserSubscription("")
	if output != "vip+" && output != "vip" && output != "free" {
		t.Fatalf("Expected \"%s\" but got \"%s\"", "free / vip / vip+", output)
	}
	w.Close()
	r.Close()
}

func TestGetActiveMachineID(t *testing.T) {
	r, w := SetOutputTest()
	output := GetActiveMachineID("")
	if output != "" {
		t.Fatalf("Expected \"%s\" but got \"%s\"", "", output)
	}
	w.Close()
	r.Close()
}

func TestSearchItemIDByNameMachine(t *testing.T) {
	r, w := SetOutputTest()
	output := SearchItemIDByName("Sau", "", "Machine")
	if output != "551" {
		t.Fatalf("Expected \"%s\" but got \"%s\"", "551", output)
	}
	w.Close()
	r.Close()
}

func TestSearchItemIDByNameChallenge(t *testing.T) {
	r, w := SetOutputTest()
	output := SearchItemIDByName("test", "", "Challenge")
	if output != "173" {
		t.Fatalf("Expected \"%s\" but got \"%s\"", "173", output)
	}
	w.Close()
	r.Close()
}
