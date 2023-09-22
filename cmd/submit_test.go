package cmd

import (
	"testing"

	"github.com/GoToolSharing/htb-cli/utils"
)

func TestSubmitChallengeBadFlag(t *testing.T) {
	r, w := utils.SetOutputTest()
	defer w.Close()
	defer r.Close()
	difficulty := 1
	machineName := ""
	challengeName := "Leet Test"
	flag := "testingflag"
	output, err := coreSubmitCmd(difficulty, machineName, challengeName, flag, "")
	expected := "Incorrect flag"

	if err != nil || output != expected {
		t.Fatalf("Error \"%v\", Expected output: \"%s\", Got \"%v\"", err, expected, output)
	}
}

func TestSubmitMachineBadFlag(t *testing.T) {
	r, w := utils.SetOutputTest()
	defer w.Close()
	defer r.Close()
	difficulty := 1
	machineName := "Blue"
	challengeName := ""
	flag := "testingflag"
	output, err := coreSubmitCmd(difficulty, machineName, challengeName, flag, "")
	expected := "Incorrect flag"

	if err != nil || output != expected {
		t.Fatalf("Error \"%v\", Expected output: \"%s\", Got \"%v\"", err, expected, output)
	}
}
