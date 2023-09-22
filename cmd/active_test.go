package cmd

import (
	"testing"

	"github.com/GoToolSharing/htb-cli/utils"
)

func TestListActiveMachines(t *testing.T) {
	r, w := utils.SetOutputTest()
	output := core_activeCmd()
	if output != "success" {
		t.Fatalf("Expected 'Easy / Medium / Hard / Insane' but got '%v'", output)
	}
	w.Close()
	r.Close()
}
