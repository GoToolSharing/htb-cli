package cmd

import (
	"testing"

	"github.com/GoToolSharing/htb-cli/utils"
)

func TestStopNoMachine(t *testing.T) {
	r, w := utils.SetOutputTest()
	defer w.Close()
	defer r.Close()
	output, err := core_stopCmd("")
	expected := "No machine is running"

	if err == nil || err.Error() != expected {
		t.Fatalf("Error \"%v\", Expected output: \"%s\", Got \"%v\"", err, expected, output)
	}
}
