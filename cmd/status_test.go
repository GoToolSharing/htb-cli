package cmd

import (
	"testing"

	"github.com/GoToolSharing/htb-cli/utils"
)

func TestCheckStatus(t *testing.T) {
	r, w := utils.SetOutputTest()
	output := core_status("")
	if output != "All Systems Operational" {
		t.Fatalf("Expected 'All Systems Operational' but got '%v'", output)
	}
	w.Close()
	r.Close()
}
