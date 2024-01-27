package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/lib/ssh"
	"github.com/GoToolSharing/htb-cli/lib/submit"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var getflagCmd = &cobra.Command{
	Use:   "getflag",
	Short: "Retrieves and submits flags from an SSH connection (linux only)",
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
		connection, err := ssh.Connect(username, password, host, 22)
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}
		userFlag, err := ssh.GetUserFlag(connection)
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}
		userFlag = strings.ReplaceAll(userFlag, "\n", "")
		hostname, err := ssh.GetHostname(connection)
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}
		connection.Close()
		url, payload, err := ssh.BuildSubmitStuff(hostname, userFlag)
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}

		message, err := submit.SubmitFlag(url, payload)
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}
		fmt.Println(message)

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
