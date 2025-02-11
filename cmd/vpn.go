package cmd

import (
	"fmt"
	"os"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/lib/vpn"
	"github.com/GoToolSharing/htb-cli/lib/webhooks"
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
		listVPNParam, err := cmd.Flags().GetBool("list")
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}

		if listVPNParam {
			err := vpn.List()
			if err != nil {
				config.GlobalConfig.Logger.Error("", zap.Error(err))
				os.Exit(1)
			}
			return
		}

		startVPNParam, err := cmd.Flags().GetBool("start")
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}
		stopVPNParam, err := cmd.Flags().GetBool("stop")
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}
		modeVPNParam, err := cmd.Flags().GetString("mode")
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}

		config.GlobalConfig.Logger.Debug(fmt.Sprintf("Mode: %s", modeVPNParam))

		if downloadVPNParam {
			config.GlobalConfig.Logger.Info("--download flag detected")
			err := vpn.DownloadAll()
			if err != nil {
				config.GlobalConfig.Logger.Error("", zap.Error(err))
				os.Exit(1)
			}
		}
		if startVPNParam && stopVPNParam {
			fmt.Println("--start and --stop cannot be used at the same time")
			os.Exit(1)
		}

		// config.BaseDirectory

		var filename string
		var pattern string
		if startVPNParam {
			switch modeVPNParam {
			case "labs":
				pattern = "Labs"
				filename = config.BaseDirectory + "/*" + pattern + "*"
			case "sp":
				pattern = "StartingPoint"
				filename = config.BaseDirectory + "/*" + pattern + "*"
			case "fortresses":
				pattern = "Fortress"
				filename = config.BaseDirectory + "/*" + pattern + "*"
			case "prolabs":
				// TODO : Get VPN name
				pattern = "Pro"
				filename = config.BaseDirectory + "/*" + pattern + "*"
			case "competitive":
				pattern = "Release_Arena"
				filename = config.BaseDirectory + "/*" + pattern + "*"
			default:
				fmt.Println("Available modes (-m) : labs - sp - fortresses - prolabs - competitive")
				return
			}

			_, err = vpn.Start(filename)
			if err != nil {
				config.GlobalConfig.Logger.Error("", zap.Error(err))
				os.Exit(1)
			}
		} else if stopVPNParam {
			message, err := vpn.Stop()
			if err != nil {
				config.GlobalConfig.Logger.Error("", zap.Error(err))
				os.Exit(1)
			}
			fmt.Println(message)
			err = webhooks.SendToDiscord("vpn", message)
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
	vpnCmd.Flags().BoolP("download", "d", false, "Download All VPNs from HackTheBox")
	vpnCmd.Flags().BoolP("start", "", false, "Start a VPN")
	vpnCmd.Flags().BoolP("stop", "", false, "Stop a VPN")
	vpnCmd.Flags().BoolP("list", "", false, "List VPNs")
	vpnCmd.Flags().StringP("mode", "m", "", "Mode")
}
