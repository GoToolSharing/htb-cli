package cmd

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/lib/utils"
	"github.com/spf13/cobra"
)

var machineParam []string
var challengeParam []string
var usernameParam []string

// Retrieves data for user profile
func fetchData(itemID string, endpoint string, infoKey string) (map[string]interface{}, error) {
	url := config.BaseHackTheBoxAPIURL + endpoint + itemID
	log.Println("URL :", url)

	resp, err := utils.HtbRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	parsedInfo := utils.ParseJsonMessage(resp, infoKey)
	dataMap, ok := parsedInfo.(map[string]interface{})
	if !ok {
		return nil, errors.New("Could not convert parsedInfo to map[string]interface{}")
	}
	return dataMap, nil
}

// fetchAndDisplayInfo fetches and displays information based on the specified parameters.
func fetchAndDisplayInfo(url, header string, params []string, elementType string) error {
	w := utils.SetTabWriterHeader(header)

	// Iteration on all machines / challenges / users argument
	for _, param := range params {
		itemID, err := utils.SearchItemIDByName(param, elementType)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		resp, err := utils.HtbRequest(http.MethodGet, (url + itemID), nil)
		if err != nil {
			return err
		}

		var infoKey string
		if strings.Contains(url, "machine") {
			infoKey = "info"
		} else if strings.Contains(url, "challenge") {
			infoKey = "challenge"
		} else if strings.Contains(url, "user/profile") {
			infoKey = "profile"
		} else {
			return fmt.Errorf("infoKey not defined")
		}

		log.Println("URL :", url)
		log.Println("InfoKey :", infoKey)

		info := utils.ParseJsonMessage(resp, infoKey)
		data := info.(map[string]interface{})

		endpoints := []struct {
			name string
			url  string
		}{
			{"Fortresses", "/user/profile/progress/fortress/"},
			{"Endgames", "/user/profile/progress/endgame/"},
			{"Prolabs", "/user/profile/progress/prolab/"},
			{"Activity", "/user/profile/activity/"},
		}

		dataMaps := make(map[string]map[string]interface{})

		for _, ep := range endpoints {
			data, err := fetchData(itemID, ep.url, "profile")
			if err != nil {
				fmt.Printf("Error fetching data for %s: %v\n", ep.name, err)
				continue
			}
			dataMaps[ep.name] = data
		}

		var bodyData string
		if elementType == "Machine" {
			status := utils.SetStatus(data)
			retiredStatus := getMachineStatus(data)
			release_key := "release"
			datetime, err := utils.ParseAndFormatDate(data[release_key].(string))
			if err != nil {
				return err
			}
			ip := getIPStatus(data)
			bodyData = fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", data["name"], data["os"], retiredStatus, data["difficultyText"], data["stars"], ip, status, data["last_reset_time"], datetime)
		} else if elementType == "Challenge" {
			status := utils.SetStatus(data)
			retiredStatus := getMachineStatus(data)
			release_key := "release_date"
			datetime, err := utils.ParseAndFormatDate(data[release_key].(string))
			if err != nil {
				return err
			}
			bodyData = fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", data["name"], data["category_name"], retiredStatus, data["difficulty"], data["stars"], data["solves"], status, datetime)
		} else if elementType == "Username" {
			utils.DisplayInformationsGUI(data, dataMaps)
			return nil
		}

		utils.SetTabWriterData(w, bodyData)
		w.Flush()
	}
	return nil
}

// coreInfoCmd is the core of the info command; it checks the parameters and displays corresponding information.
func coreInfoCmd() error {
	machineHeader := "Name\tOS\tRetired\tDifficulty\tStars\tIP\tStatus\tLast Reset\tRelease"
	challengeHeader := "Name\tCategory\tRetired\tDifficulty\tStars\tSolves\tStatus\tRelease"
	usernameHeader := "Name\tUser Owns\tSystem Owns\tUser Bloods\tSystem Bloods\tTeam\tUniversity\tRank\tGlobal Rank\tPoints"

	type infoType struct {
		APIURL string
		Header string
		Params []string
		Name   string
	}

	infos := []infoType{
		{config.BaseHackTheBoxAPIURL + "/machine/profile/", machineHeader, machineParam, "Machine"},
		{config.BaseHackTheBoxAPIURL + "/challenge/info/", challengeHeader, challengeParam, "Challenge"},
		{config.BaseHackTheBoxAPIURL + "/user/profile/basic/", usernameHeader, usernameParam, "Username"},
	}

	// No arguments provided
	if len(machineParam) == 0 && len(usernameParam) == 0 && len(challengeParam) == 0 {
		isConfirmed := utils.AskConfirmation("Do you want to check for active machine ?")
		if isConfirmed {
			err := displayActiveMachine(machineHeader)
			if err != nil {
				log.Fatal(err)
			}
		}
		// Get current account
		resp, err := utils.HtbRequest(http.MethodGet, config.BaseHackTheBoxAPIURL+"/user/info", nil)
		if err != nil {
			return err
		}
		info := utils.ParseJsonMessage(resp, "info")
		infoMap, _ := info.(map[string]interface{})
		newInfo := infoType{
			APIURL: config.BaseHackTheBoxAPIURL + "/user/profile/basic/",
			Header: "",
			Params: []string{infoMap["name"].(string)},
			Name:   "Username",
		}
		infos = append(infos, newInfo)
	}

	for _, info := range infos {
		if len(info.Params) > 0 {
			if info.Name == "Machine" {
				isConfirmed := utils.AskConfirmation("Do you want to check for active " + strings.ToLower(info.Name) + "?")
				if isConfirmed {
					err := displayActiveMachine(info.Header)
					if err != nil {
						log.Fatal(err)
					}
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
	machineID := utils.GetActiveMachineID()

	if machineID != "" {
		log.Println("Active machine found !")
		log.Println("Machine ID:", machineID)

		tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)
		w := utils.SetTabWriterHeader(header)

		url := fmt.Sprintf("%s/machine/profile/%s", config.BaseHackTheBoxAPIURL, machineID)
		resp, err := utils.HtbRequest(http.MethodGet, url, nil)
		if err != nil {
			return err
		}
		info := utils.ParseJsonMessage(resp, "info")

		data := info.(map[string]interface{})
		status := utils.SetStatus(data)
		retiredStatus := getMachineStatus(data)

		datetime, err := utils.ParseAndFormatDate(data["release"].(string))
		if err != nil {
			return err
		}

		machineType := utils.GetMachineType(machineID)
		log.Printf("Machine Type: %s", machineType)

		userSubscription := utils.GetUserSubscription()
		log.Printf("User subscription: %s", userSubscription)

		ip := "Undefined"
		_ = ip
		switch {
		case userSubscription == "vip+" || machineType == "release":
			ip = utils.GetActiveMachineIP()
		default:
			ip = getIPStatus(data).(string)
		}

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
		err := coreInfoCmd()
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
	infoCmd.Flags().StringSliceVarP(&usernameParam, "username", "u", []string{}, "Username")
}
