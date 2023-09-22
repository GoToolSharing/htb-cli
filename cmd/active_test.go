package cmd

import (
	"testing"

	"github.com/GoToolSharing/htb-cli/utils"
)

func TestListActiveMachines(t *testing.T) {
	r, w := utils.SetOutputTest()
	defer w.Close()
	defer r.Close()
	output, err := coreActiveCmd("")
	expected := ""
	if err != nil || output != expected {
		t.Fatalf("Error \"%v\", Expected output: \"%s\", Got \"%v\"", err, expected, output)
	}
}
