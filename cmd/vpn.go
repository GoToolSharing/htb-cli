package cmd

import (
	"fmt"

	"github.com/GoToolSharing/htb-cli/utils"
	"github.com/spf13/cobra"
)

var DownloadVPNParam bool

// vpnCmd is a Cobra command used to interact with HackTheBox VPNs.
// It defines the "vpn" command for the command-line application.
var vpnCmd = &cobra.Command{
	Use:   "vpn",
	Short: "Interact with HackTheBox VPNs",
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

// init function is called to initialize the vpnCmd command and its flags.
// It adds the vpn command to the root command and sets up its flags.
func init() {
	rootCmd.AddCommand(vpnCmd)
	vpnCmd.Flags().BoolVarP(&DownloadVPNParam, "download", "d", false, "Download VPN")
}
