package cmd

import (
	"fmt"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:   "htb-cli",
	Short: "CLI enhancing the HackTheBox user experience.",
	Long:  `This software, engineered using the Go programming language, serves to streamline and automate various tasks for the HackTheBox platform, enhancing user efficiency and productivity.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config.ConfigureLogger()
		config.GlobalConfig.Logger.Debug(fmt.Sprintf("Verbosity level : %v", config.GlobalConfig.Verbose))
		defer config.GlobalConfig.Logger.Sync()
		err := config.Init()
		if err != nil {
			fmt.Println(err)
			return
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		config.GlobalConfig.Logger.Error("Erreur lors de l'exécution de la commande", zap.Error(err))
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().CountVarP(&config.GlobalConfig.Verbose, "verbose", "v", "Verbose level")
	rootCmd.PersistentFlags().StringVarP(&config.GlobalConfig.ProxyParam, "proxy", "p", "", "Configure a URL for an HTTP proxy")
	rootCmd.PersistentFlags().BoolVarP(&config.GlobalConfig.BatchParam, "batch", "b", false, "Don't ask questions")
}
