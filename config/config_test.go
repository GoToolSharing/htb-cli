package config

import "testing"

func TestMachineOwnAPIUsesV5(t *testing.T) {
	t.Parallel()

	want := "https://labs.hackthebox.com/api/v5/machine/own"
	if MachineOwnAPIURL != want {
		t.Fatalf("MachineOwnAPIURL = %q, want %q", MachineOwnAPIURL, want)
	}
}
