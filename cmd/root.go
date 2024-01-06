package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/lib/update"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:   "htb-cli",
	Short: "CLI enhancing the HackTheBox user experience.",
	Long:  `This software, engineered using the Go programming language, serves to streamline and automate various tasks for the HackTheBox platform, enhancing user efficiency and productivity.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		err := config.ConfigureLogger()
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}
		config.GlobalConfig.Logger.Debug(fmt.Sprintf("Verbosity level : %v", config.GlobalConfig.Verbose))
		err = config.Init()
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}
		message, err := update.Check(config.Version)
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}
		config.GlobalConfig.Logger.Debug(fmt.Sprintf("Message : %s", message))
		if strings.Contains(message, "A new update") {
			fmt.Println(message)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().CountVarP(&config.GlobalConfig.Verbose, "verbose", "v", "Verbose level")
	rootCmd.PersistentFlags().StringVarP(&config.GlobalConfig.ProxyParam, "proxy", "p", "", "Configure a URL for an HTTP proxy")
	rootCmd.PersistentFlags().BoolVarP(&config.GlobalConfig.BatchParam, "batch", "b", false, "Don't ask questions")
}
