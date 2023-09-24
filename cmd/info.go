package cmd

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/GoToolSharing/htb-cli/utils"
	"github.com/spf13/cobra"
)

var machineParam []string
var challengeParam []string

// parseAndFormatDate takes a date string, parses it into a time.Time object, and then formats it to the "2006-01-02" format.
func parseAndFormatDate(input string) (string, error) {
	t, err := time.Parse(time.RFC3339Nano, input)
	if err != nil {
		return "", fmt.Errorf("error: parsing date: %v", err)
	}
	return t.Format("2006-01-02"), nil
}

// setStatus determines the status based on user and root flags.
func setStatus(data map[string]interface{}) string {
	userFlag, userFlagExists := data["authUserInUserOwns"].(bool)
	rootFlag, rootFlagExists := data["authUserInRootOwns"].(bool)

	switch {
	case !userFlagExists && !rootFlagExists:
		return "No flags"
	case userFlag && !rootFlag:
		return "User flag"
	case !userFlag && rootFlag:
		return "Root flag"
	case userFlag && rootFlag:
		return "User & Root"
	default:
		return "No flags"
	}
}

// setRetiredStatus determines whether an item is retired or not.
func setRetiredStatus(data map[string]interface{}) string {
	if retired, exists := data["retired"].(float64); exists && retired == 0 {
		return "Yes"
	}
	return "No"
}

// fetchAndDisplayInfo fetches and displays information based on the specified parameters.
func fetchAndDisplayInfo(url, header string, params []string, elementType string) error {
	log.Println("Params :", params)
	w := utils.SetTabWriterHeader(header)

	for _, param := range params {
		itemID := utils.SearchItemIDByName(param, proxyParam, elementType)

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

		status := setStatus(data)
		retiredStatus := setRetiredStatus(data)

		release_key := ""
		if elementType == "Machine" {
			release_key = "release"
		} else {
			release_key = "release_date"
		}
		datetime, err := parseAndFormatDate(data[release_key].(string))
		if err != nil {
			return err
		}

		var bodyData string
		if strings.Contains(url, "machine") {
			bodyData = fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", data["name"], data["os"], retiredStatus, data["difficultyText"], data["stars"], data["firstUserBloodTime"], data["firstRootBloodTime"], status, datetime)
		} else {
			bodyData = fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", data["name"], data["category_name"], retiredStatus, data["difficulty"], data["stars"], data["solves"], status, datetime)
		}

		utils.SetTabWriterData(w, bodyData)
	}
	return nil
}

// coreInfoCmd is the core of the info command; it checks the parameters and displays corresponding information.
func coreInfoCmd(machineParam []string, challengeParam []string) error {
	if len(machineParam) > 0 && len(challengeParam) > 0 {
		return errors.New("error: You can only specify either -m or -c flags, not both")
	}
	if os.Getenv("TEST") == "" {
		isConfirmed := utils.AskConfirmation("Do you want to check for active machine ?")
		if isConfirmed {
			err := displayActiveMachine()
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	machineHeader := "Name\tOS\tActive\tDifficulty\tStars\tFirstUserBlood\tFirstRootBlood\tStatus\tRelease"
	challengeHeader := "Name\tCategory\tActive\tDifficulty\tStars\tSolves\tStatus\tRelease"
	if len(machineParam) > 0 {
		err := fetchAndDisplayInfo("https://www.hackthebox.com/api/v4/machine/profile/", machineHeader, machineParam, "Machine")
		if err != nil {
			return err
		}
	} else if len(challengeParam) > 0 {
		err := fetchAndDisplayInfo("https://www.hackthebox.com/api/v4/challenge/info/", challengeHeader, challengeParam, "Challenge")
		if err != nil {
			return err
		}
	}
	return nil
}

// displayActiveMachine displays information about the active machine if one is found.
func displayActiveMachine() error {
	machine_id := utils.GetActiveMachineID(proxyParam)
	retired_status := "Not defined"

	if machine_id != "" {
		log.Println("Active machine found !")
		log.Println("Machine ID:", machine_id)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)
		header := "Name\tOS\tActive\tDifficulty\tStars\tIP\tStatus\tRelease"
		w = utils.SetTabWriterHeader(header)

		url := "https://www.hackthebox.com/api/v4/machine/profile/" + machine_id
		resp, err := utils.HtbRequest(http.MethodGet, url, proxyParam, nil)
		if err != nil {
			return err
		}
		info := utils.ParseJsonMessage(resp, "info")
		log.Println(info)
		data := info.(map[string]interface{})

		status := setStatus(data)

		if data["retired"].(float64) == 0 {
			retired_status = "Yes"
		} else {
			retired_status = "No"
		}

		datetime, err := parseAndFormatDate(data["release"].(string))
		if err != nil {
			return err
		}

		bodyData := fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n",
			data["name"], data["os"], retired_status,
			data["difficultyText"], data["stars"],
			data["ip"], status, datetime)

		utils.SetTabWriterData(w, bodyData)
	} else {
		fmt.Print("No machine is running")
	}
	return nil
}

// infoCmd is a Cobra command that serves as an entry point to display detailed information about machines.
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Showcase detailed machine information",
	Long:  "Displays detailed information of the specified machines in a structured table.",
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
