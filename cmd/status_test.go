package cmd

import (
	"bytes"
	"io"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func runCmdAndCaptureOutput(t *testing.T, cmd *cobra.Command, args []string) string {
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	var cmdOutBuf bytes.Buffer
	cmd.SetOut(&cmdOutBuf)

	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd.Run(cmd, args)

	_ = w.Close()
	out, _ := io.ReadAll(r)

	return string(out)
}

func TestCheckStatus(t *testing.T) {
	output := runCmdAndCaptureOutput(t, statusCmd, []string{})

	if !strings.Contains(output, "All Systems Operational") {
		t.Fatalf("Expected 'All Systems Operational' but got '%v'", output)
	}
}
