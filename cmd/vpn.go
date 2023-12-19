package cmd

import (
	"os"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/lib/vpn"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var vpnCmd = &cobra.Command{
	Use:   "vpn",
	Short: "Interact with HackTheBox VPNs",
	Run: func(cmd *cobra.Command, args []string) {
		config.GlobalConfig.Logger.Info("VPN command executed")
		downloadVPNParam, err := cmd.Flags().GetBool("download")
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}
		if downloadVPNParam {
			config.GlobalConfig.Logger.Info("--download flag detected")
			err := vpn.DownloadAll()
			if err != nil {
				config.GlobalConfig.Logger.Error("", zap.Error(err))
				os.Exit(1)
			}
		}
		config.GlobalConfig.Logger.Info("Exit vpn command correctly")
	},
}

func init() {
	rootCmd.AddCommand(vpnCmd)
	vpnCmd.Flags().Bool("download", false, "Download All VPNs from HackTheBox")
}
