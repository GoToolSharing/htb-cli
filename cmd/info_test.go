package cmd

import (
	"testing"

	"github.com/GoToolSharing/htb-cli/utils"
)

func TestCoreInfoCmdMachine(t *testing.T) {
	r, w := utils.SetOutputTest()
	defer w.Close()
	defer r.Close()
	machineParams := []string{"Sau", "Jupiter"}
	challengeParams := []string{}
	output, err := core_infoCmd(machineParams, challengeParams)
	expected := ""
	if err != nil || output != expected {
		t.Fatalf("Error \"%v\", Expected output: \"%s\", Got \"%v\"", err, expected, output)
	}
}

func TestCoreInfoCmdChallenge(t *testing.T) {
	r, w := utils.SetOutputTest()
	defer w.Close()
	defer r.Close()
	machineParams := []string{}
	challengeParams := []string{"Toxic", "Toxin"}
	output, err := core_infoCmd(machineParams, challengeParams)
	expected := ""
	if err != nil || output != expected {
		t.Fatalf("Error \"%v\", Expected output: \"%s\", Got \"%v\"", err, expected, output)
	}
}
