package cmd

import (
	"fmt"
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
			case "endgames":
				pattern = "Endgame"
				filename = config.BaseDirectory + "/*" + pattern + "*"
			case "competitive":
				pattern = "Release_Arena"
				filename = config.BaseDirectory + "/*" + pattern + "*"
			default:
				fmt.Println("Available modes : labs - sp - fortresses - prolabs - endgames - competitive")
				return
			}
			// fmt.Println("FILENAME")
			// fmt.Println(filename)

			// pattern := regexp.MustCompile(pattern)

			// err := filepath.Walk(config.BaseDirectory, func(path string, info os.FileInfo, err error) error {
			// 	if err != nil {
			// 		fmt.Println("Error accessing path:", path)
			// 		return err
			// 	}

			// 	if !pattern.MatchString(info.Name()) && info.IsDir() {
			// 		vpn.DownloadAll()
			// 	}
			// 	return nil
			// })

			// if err != nil {
			// 	fmt.Println("Erreur lors de la recherche de fichiers:", err)
			// }

			_, err = vpn.Start(filename)
			if err != nil {
				config.GlobalConfig.Logger.Error("", zap.Error(err))
				os.Exit(1)
			}
		} else if stopVPNParam {
			_, err := vpn.Stop()
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
	vpnCmd.Flags().BoolP("start", "", false, "Start VPN")
	vpnCmd.Flags().BoolP("stop", "", false, "Stop VPN")
	vpnCmd.Flags().StringP("mode", "m", "", "Mode")
}
