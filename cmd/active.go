package cmd

import (
	"fmt"
	"log"

	"net/http"

	"github.com/GoToolSharing/htb-cli/utils"
	"github.com/spf13/cobra"
)

const machineURL = "https://www.hackthebox.com/api/v4/machine/list"

// coreActiveCmd fetches the list of active machines and displays their details.
func coreActiveCmd(proxyParam string) error {
	resp, err := utils.HtbRequest(http.MethodGet, machineURL, proxyParam, nil)
	if err != nil {
		return err
	}
	info := utils.ParseJsonMessage(resp, "info")
	log.Println(info)
	err = displayActiveMachines(info)
	if err != nil {
		return err
	}
	return nil
}

// displayActiveMachines takes the machine information, processes it, and displays it in a tabulated format.
func displayActiveMachines(info interface{}) error {
	header := "Name\tOS\tDifficulty\tUser Owns\tSystem Owns\tStars\tStatus\tRelease"
	tabWriter := utils.SetTabWriterHeader(header)

	for _, value := range info.([]interface{}) {
		data := value.(map[string]interface{})

		status := utils.SetStatus(data)
		datetime, err := utils.ParseAndFormatDate(data["release"].(string))
		if err != nil {
			return err
		}

		formattedData := fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n",
			data["name"], data["os"], data["difficultyText"],
			data["user_owns_count"], data["root_owns_count"],
			data["stars"], status, datetime)
		utils.SetTabWriterData(tabWriter, formattedData)
	}

	tabWriter.Flush()
	return nil
}

// activeCmd represents the "active" command that the user can call.
var activeCmd = &cobra.Command{
	Use:   "active",
	Short: "Catalogue of active machines",
	Long:  "This command serves to generate a detailed summary of the currently active machines, providing pertinent information for each.",
	Run: func(cmd *cobra.Command, args []string) {
		err := coreActiveCmd(proxyParam)
		if err != nil {
			log.Fatal(err)
		}
	},
}

// init adds the "active" command to the root command,
func init() {
	rootCmd.AddCommand(activeCmd)
}
