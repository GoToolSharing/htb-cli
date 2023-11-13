package cmd

import (
	"fmt"

	"github.com/GoToolSharing/htb-cli/lib/vpn"
	"github.com/spf13/cobra"
)

var vpnCmd = &cobra.Command{
	Use:   "vpn",
	Short: "Interact with HackTheBox VPNs",
	Run: func(cmd *cobra.Command, args []string) {
		downloadVPNParam, err := cmd.Flags().GetBool("download")
		if err != nil {
			fmt.Println(err)
			return
		}
		if downloadVPNParam {
			err := vpn.DownloadAll()
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(vpnCmd)
	vpnCmd.Flags().Bool("download", false, "Download All VPNs from HackTheBox")
}
