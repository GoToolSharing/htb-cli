package cmd

import (
	"testing"

	"github.com/GoToolSharing/htb-cli/utils"
)

func TestCheckStatus(t *testing.T) {
	r, w := utils.SetOutputTest()
	defer w.Close()
	defer r.Close()
	output, err := coreStatusCmd("")
	expected := "All Systems Operational"
	if err != nil || output != expected {
		t.Fatalf("Error \"%v\", Expected output: \"%s\", Got \"%v\"", err, expected, output)
	}
}
