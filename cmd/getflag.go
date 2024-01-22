package cmd

import (
	"os"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/lib/ssh"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var getflagCmd = &cobra.Command{
	Use:   "getflag",
	Short: "Retrieves and submits flags from an SSH connection",
	Run: func(cmd *cobra.Command, args []string) {
		config.GlobalConfig.Logger.Info("Getflag command executed")
		username, err := cmd.Flags().GetString("username")
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}
		password, err := cmd.Flags().GetString("password")
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}
		host, err := cmd.Flags().GetString("host")
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}
		err = ssh.Connect(username, password, host, 22)
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}
		config.GlobalConfig.Logger.Info("Exit getflag command correctly")
	},
}

func init() {
	rootCmd.AddCommand(getflagCmd)
	getflagCmd.Flags().StringP("username", "u", "", "SSH username")
	getflagCmd.Flags().StringP("password", "p", "", "SSH password")
	getflagCmd.Flags().StringP("port", "P", "", "(Optional) SSH Port (Default 22)")
	getflagCmd.Flags().StringP("host", "", "", "(Optional) SSH host")
}
