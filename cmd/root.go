package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "htb-cli",
	Short: "Shortcuts that facilitate Hackthebox",
	Long:  `This is a program developed in Go to facilitate and automate certain tasks for the Hackthebox platform.`,
}

func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
