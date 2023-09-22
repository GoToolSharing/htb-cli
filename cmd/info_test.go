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
	output := core_infoCmd(machineParams, challengeParams)
	expected := "success"
	if output != expected {
		t.Fatalf("Expected \"%s\" but got \"%v\"", expected, output)
	}
}

func TestCoreInfoCmdChallenge(t *testing.T) {
	r, w := utils.SetOutputTest()
	defer w.Close()
	defer r.Close()
	machineParams := []string{}
	challengeParams := []string{"Toxic", "Toxin"}
	output := core_infoCmd(machineParams, challengeParams)
	expected := "success"
	if output != expected {
		t.Fatalf("Expected \"%s\" but got \"%v\"", expected, output)
	}
}
