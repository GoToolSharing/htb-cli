package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/GoToolSharing/htb-cli/utils"
	"github.com/spf13/cobra"
)

var machineParam []string
var challengeParam []string

// fetchAndDisplayInfo fetches and displays information based on the specified parameters.
func fetchAndDisplayInfo(url, header string, params []string, elementType string) error {
	log.Println("Params :", params)
	w := utils.SetTabWriterHeader(header)

	// Iteration on all machines / challenges argument
	for _, param := range params {
		itemID, err := utils.SearchItemIDByName(param, proxyParam, elementType)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		fullURL := url + itemID
		resp, err := utils.HtbRequest(http.MethodGet, fullURL, proxyParam, nil)
		if err != nil {
			return err
		}

		var infoKey string
		if strings.Contains(url, "machine") {
			infoKey = "info"
		} else {
			infoKey = "challenge"
		}

		info := utils.ParseJsonMessage(resp, infoKey)
		data := info.(map[string]interface{})

		status := utils.SetStatus(data)
		retiredStatus := getMachineStatus(data)

		release_key := ""
		if elementType == "Machine" {
			release_key = "release"
		} else {
			release_key = "release_date"
		}
		datetime, err := utils.ParseAndFormatDate(data[release_key].(string))
		if err != nil {
			return err
		}

		ip := getIPStatus(data)

		var bodyData string
		if strings.Contains(url, "machine") {
			bodyData = fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", data["name"], data["os"], retiredStatus, data["difficultyText"], data["stars"], ip, status, data["last_reset_time"], datetime)
		} else {
			bodyData = fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", data["name"], data["category_name"], retiredStatus, data["difficulty"], data["stars"], data["solves"], status, datetime)
		}

		utils.SetTabWriterData(w, bodyData)
		w.Flush()
	}
	return nil
}

// coreInfoCmd is the core of the info command; it checks the parameters and displays corresponding information.
func coreInfoCmd(machineParam []string, challengeParam []string) error {
    machineHeader := "Name\tOS\tRetired\tDifficulty\tStars\tIP\tStatus\tLast Reset\tRelease"
    challengeHeader := "Name\tCategory\tRetired\tDifficulty\tStars\tSolves\tStatus\tRelease"

    type infoType struct {
        APIURL string
        Header string
        Params []string
        Name   string
    }

    infos := []infoType{
        {"https://www.hackthebox.com/api/v4/machine/profile/", machineHeader, machineParam, "Machine"},
        {"https://www.hackthebox.com/api/v4/challenge/info/", challengeHeader, challengeParam, "Challenge"},
    }

    for _, info := range infos {
        if len(info.Params) > 0 {
            isConfirmed := utils.AskConfirmation("Do you want to check for active " + strings.ToLower(info.Name) + "?")
            if isConfirmed && info.Name == "Machine" { 
                err := displayActiveMachine(info.Header)
                if err != nil {
                    log.Fatal(err)
                }
            }
            for _, param := range info.Params {
                err := fetchAndDisplayInfo(info.APIURL, info.Header, []string{param}, info.Name)
                if err != nil {
                    return err
                }
            }
        }
    }
    return nil
}

// getMachineStatus returns machine status
func getMachineStatus(data map[string]interface{}) string {
	if data["retired"].(float64) == 0 {
		return "No"
	}
	return "Yes"
}

// getIPStatus returns ip status
func getIPStatus(data map[string]interface{}) interface{} {
	if data["ip"] == nil {
		return "Undefined"
	}
	return data["ip"]
}

// displayActiveMachine displays information about the active machine if one is found.
func displayActiveMachine(header string) error {
	machineID := utils.GetActiveMachineID(proxyParam)

	if machineID != "" {
		log.Println("Active machine found !")
		log.Println("Machine ID:", machineID)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)
		w = utils.SetTabWriterHeader(header)

		url := "https://www.hackthebox.com/api/v4/machine/profile/" + machineID
		resp, err := utils.HtbRequest(http.MethodGet, url, proxyParam, nil)
		if err != nil {
			return err
		}
		info := utils.ParseJsonMessage(resp, "info")
		log.Println(info)
		
		data := info.(map[string]interface{})
		status := utils.SetStatus(data)
		retiredStatus := getMachineStatus(data)
		
		datetime, err := utils.ParseAndFormatDate(data["release"].(string))
		if err != nil {
			return err
		}

		ip := getIPStatus(data)

		bodyData := fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n",
			data["name"], data["os"], retiredStatus,
			data["difficultyText"], data["stars"],
			ip, status, data["last_reset_time"], datetime)

		utils.SetTabWriterData(w, bodyData)
		w.Flush()
	} else {
		fmt.Print("No machine is running")
	}
	return nil
}

// infoCmd is a Cobra command that serves as an entry point to display detailed information about machines.
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Detailed information on challenges and machines",
	Long:  "Displays detailed information on machines and challenges in a structured table",
	Run: func(cmd *cobra.Command, args []string) {
		err := coreInfoCmd(machineParam, challengeParam)
		if err != nil {
			log.Fatal(err)
		}
	},
}

// init adds the info command to the root command and sets flags specific to this command.
func init() {
	rootCmd.AddCommand(infoCmd)
	infoCmd.Flags().StringSliceVarP(&machineParam, "machine", "m", []string{}, "Machine name")
	infoCmd.Flags().StringSliceVarP(&challengeParam, "challenge", "c", []string{}, "Challenge name")
}
