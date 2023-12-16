package cmd

import (
	"fmt"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays the current version",
	Run: func(cmd *cobra.Command, args []string) {
		if len(config.Version) == 40 {
			fmt.Println("Development version (dev branch): " + config.Version)
		} else {
			fmt.Println("Stable version (main branch): " + config.Version)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
