package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Install the latest update",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running update command...")
		exec.Command("go install github.com/QU35T-code/htb-cli@latest")
		fmt.Println("End of update")
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
