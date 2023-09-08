package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"net/http"

	"github.com/QU35T-code/htb-cli/utils"
	"github.com/spf13/cobra"
)

var proxyParam string

var activeCmd = &cobra.Command{
	Use:   "active",
	Short: "List of active machines",
	Run: func(cmd *cobra.Command, args []string) {
		url := "https://www.hackthebox.com/api/v4/machine/list"
		resp, err := utils.HtbRequest(http.MethodGet, url, proxyParam, nil)
		if err != nil {
			log.Fatal(err)
		}
		info := utils.ParseJsonMessage(resp, "info")
		log.Println(info)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)
		fmt.Fprintln(w, "Name\tOS\tDifficulty\tUser Owns\tSystem Owns\tStars\tStatus\tRelease")

		// red_color := color.New(color.FgRed).SprintFunc()
		status := "Not defined"
		for _, value := range info.([]interface{}) {
			data := value.(map[string]interface{})
			if data["authUserInUserOwns"] == nil && data["authUserInRootOwns"] == nil {
				status = "❌ User - ❌ Root"
			} else if data["authUserInUserOwns"] == true && data["authUserInRootOwns"] == nil {
				status = "✅ User - ❌ Root"
			} else if data["authUserInUserOwns"] == nil && data["authUserInRootOwns"] == true {
				status = "❌ User - ✅ Root"
			} else if data["authUserInUserOwns"] == true && data["authUserInRootOwns"] == true {
				status = "✅ User - ✅ Root"
			}
			fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", data["name"], data["os"], data["difficultyText"], data["user_owns_count"], data["root_owns_count"], data["stars"], status, data["release"])
		}
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(activeCmd)
	activeCmd.Flags().StringVarP(&proxyParam, "proxy", "p", "", "Configure a URL for an HTTP proxy")
}
