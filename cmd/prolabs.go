package cmd

import (
	"fmt"
	"os"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/lib/prolabs"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var prolabsCmd = &cobra.Command{
	Use:   "prolabs",
	Short: "Interact with prolabs",
	Run: func(cmd *cobra.Command, args []string) {
		config.GlobalConfig.Logger.Info("Prolabs command executed")

		// Voir les prolabs
		// Voir le detail d'un prolab (progression, flags, machines, ...)
		// Submit flags
		// Voir les review
		// Voir le changelog
		// Certificate

		p := tea.NewProgram(prolabs.InitialModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
		config.GlobalConfig.Logger.Info("Exit prolabs command correctly")
	},
}

func init() {
	rootCmd.AddCommand(prolabsCmd)
}
