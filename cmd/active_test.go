package cmd

import (
	"strings"
	"testing"
)

func TestListActiveMachines(t *testing.T) {
	output := runCmdAndCaptureOutput(t, activeCmd, []string{})

	if !strings.Contains(output, "Medium") && !strings.Contains(output, "Easy") && !strings.Contains(output, "Hard") && !strings.Contains(output, "Insane") {
		t.Fatalf("Expected 'Easy / Medium / Hard / Insane' but got '%v'", output)
	}
}
