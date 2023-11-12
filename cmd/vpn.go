package cmd

import (
	"fmt"

	"github.com/GoToolSharing/htb-cli/utils"
	"github.com/spf13/cobra"
)

var DownloadVPNParam bool

var vpnCmd = &cobra.Command{
	Use:   "vpn",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		if DownloadVPNParam {
			err := utils.DownloadVPN(proxyParam)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(vpnCmd)
	vpnCmd.Flags().BoolVarP(&DownloadVPNParam, "download", "d", false, "Download VPN")
}
