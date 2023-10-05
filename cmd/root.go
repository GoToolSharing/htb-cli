package cmd

import (
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var verbose bool
var proxyParam string
var batchParam bool

const baseAPIURL = "https://www.hackthebox.com/api/v4"

var rootCmd = &cobra.Command{
	Use:   "htb-cli",
	Short: "CLI enhancing the HackTheBox user experience.",
	Long:  `This software, engineered using the Go programming language, serves to streamline and automate various tasks for the HackTheBox platform, enhancing user efficiency and productivity.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if verbose {
			log.SetOutput(os.Stdout)
		} else {
			log.SetOutput(io.Discard)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose mode")
	rootCmd.PersistentFlags().StringVarP(&proxyParam, "proxy", "p", "", "Configure a URL for an HTTP proxy")
	rootCmd.PersistentFlags().BoolVarP(&batchParam, "batch", "b", false, "Allows all questions")
}
