package cmd

import (
	"fmt"

	"github.com/GoToolSharing/htb-cli/lib/vpn"
	"github.com/spf13/cobra"
)

// vpnCmd is a Cobra command used to interact with HackTheBox VPNs.
// It defines the "vpn" command for the command-line application.
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

// init function is called to initialize the vpnCmd command and its flags.
// It adds the vpn command to the root command and sets up its flags.
func init() {
	rootCmd.AddCommand(vpnCmd)
	vpnCmd.Flags().Bool("download", false, "Download All VPNs from HackTheBox")
}
