package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/tabwriter"
	"time"

	"github.com/QU35T-code/htb-cli/utils"
	"github.com/kyokomi/emoji/v2"
	"github.com/spf13/cobra"
)

var machineParam []string

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Displays machine information",
	Run: func(cmd *cobra.Command, args []string) {
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)
		fmt.Fprintln(w, "Name\tOS\tActive\tDifficulty\tStars\tFirstUserBlood\tFirstRootBlood\tStatus\tRelease")
		status := "Not defined"
		log.Println(machineParam)
		for index, _ := range machineParam {
			machine_id := utils.SearchMachineIDByName(machineParam[index], proxyParam)

			url := "https://www.hackthebox.com/api/v4/machine/profile/" + machine_id
			resp, err := utils.HtbRequest(http.MethodGet, url, proxyParam, nil)
			if err != nil {
				log.Fatal(err)
			}
			info := utils.ParseJsonMessage(resp, "info")

			data := info.(map[string]interface{})
			if data["authUserInUserOwns"] == nil && data["authUserInRootOwns"] == nil {
				status = emoji.Sprint(":x:User - :x:Root")
			} else if data["authUserInUserOwns"] == true && data["authUserInRootOwns"] == nil {
				status = emoji.Sprint(":white_check_mark:User - :x:Root")
			} else if data["authUserInUserOwns"] == nil && data["authUserInRootOwns"] == true {
				status = emoji.Sprint(":x:User - :white_check_mark:Root")
			} else if data["authUserInUserOwns"] == true && data["authUserInRootOwns"] == true {
				status = emoji.Sprint(":white_check_mark:User - :white_check_mark:Root")
			}
			t, err := time.Parse(time.RFC3339Nano, data["release"].(string))
			if err != nil {
				fmt.Println("Erreur when date parsing :", err)
				return
			}
			datetime := t.Format("2006-01-02")
			fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", data["name"], data["os"], data["active"], data["difficultyText"], data["stars"], data["firstUserBloodTime"], data["firstRootBloodTime"], status, datetime)
		}
		w.Flush()

	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
	infoCmd.Flags().StringSliceVarP(&machineParam, "machine", "m", []string{}, "Machine name")
	// infoCmd.Flags().StringVarP(&proxyParam, "proxy", "p", "", "Configure a URL for an HTTP proxy")
	infoCmd.MarkFlagRequired("machine")
}
