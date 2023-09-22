package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"text/tabwriter"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/GoToolSharing/htb-cli/utils"
	"github.com/kyokomi/emoji/v2"
	"github.com/spf13/cobra"
)

var machineParam []string
var challengeParam []string

func core_infoCmd(machineParam []string, challengeParam []string) string {
	if len(machineParam) > 0 && len(challengeParam) > 0 {
		fmt.Println("Error: You can only specify either -m or -c flags, not both.")
		return "error flags"
	}
	if os.Getenv("TEST") == "" {
		var confirmation bool
		confirmation_message := "Do you want to check for active machine ?"
		prompt := &survey.Confirm{
			Message: confirmation_message,
		}
		if err := survey.AskOne(prompt, &confirmation); err != nil {
			log.Fatal(err)
		}
		if confirmation {
			checkActiveMachine()
		}
	}

	// Machines search
	if len(machineParam) > 0 {
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)
		fmt.Fprintln(w, "Name\tOS\tActive\tDifficulty\tStars\tFirstUserBlood\tFirstRootBlood\tStatus\tRelease")
		status := "Not defined"
		retired_status := "Not defined"
		log.Println(machineParam)
		for index, _ := range machineParam {
			machine_id := utils.SearchItemIDByName(machineParam[index], proxyParam, "Machine")

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
			if data["retired"].(float64) == 0 {
				retired_status = emoji.Sprint(":white_check_mark:")
			} else {
				retired_status = emoji.Sprint(":x:")
			}
			t, err := time.Parse(time.RFC3339Nano, data["release"].(string))
			if err != nil {
				fmt.Println("Erreur when date parsing :", err)
				return "error parsing date"
			}
			datetime := t.Format("2006-01-02")
			fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", data["name"], data["os"], retired_status, data["difficultyText"], data["stars"], data["firstUserBloodTime"], data["firstRootBloodTime"], status, datetime)
			w.Flush()
		}
		return "success"
	}

	// Challenges search
	if len(challengeParam) > 0 {
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)
		fmt.Fprintln(w, "Name\tCategory\tActive\tDifficulty\tStars\tSolves\tStatus\tRelease")
		status := "Not defined"
		retired_status := "Not defined"
		log.Println(challengeParam)
		for index, _ := range challengeParam {
			challenge_id := utils.SearchItemIDByName(challengeParam[index], proxyParam, "Challenge")
			log.Println("Challenge id:", challenge_id)
			url := "https://www.hackthebox.com/api/v4/challenge/info/" + challenge_id
			resp, err := utils.HtbRequest(http.MethodGet, url, proxyParam, nil)
			if err != nil {
				log.Fatal(err)
			}
			info := utils.ParseJsonMessage(resp, "challenge")
			data := info.(map[string]interface{})
			if data["authUserSolve"] == false {
				status = emoji.Sprint(":x:Flag")
			} else {
				status = emoji.Sprint(":white_check_mark:Flag")
			}
			if data["retired"].(float64) == 0 {
				retired_status = emoji.Sprint(":white_check_mark:")
			} else {
				retired_status = emoji.Sprint(":x:")
			}
			t, err := time.Parse(time.RFC3339Nano, data["release_date"].(string))
			if err != nil {
				fmt.Println("Erreur when date parsing :", err)
				return "error parsing date"
			}
			datetime := t.Format("2006-01-02")
			fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", data["name"], data["category_name"], retired_status, data["difficulty"], data["stars"], data["solves"], status, datetime)
			w.Flush()
		}
		return "success"
	}
	return "success"
}

func checkActiveMachine() {
	machine_id := utils.GetActiveMachineID(proxyParam)
	status := "Not defined"
	retired_status := "Not defined"
	if machine_id != "" {
		log.Println("Active machine found !")
		log.Println("Machine ID:", machine_id)
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)
		fmt.Fprintln(w, "Name\tOS\tActive\tDifficulty\tStars\tIP\tStatus\tRelease")
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
		if data["retired"].(float64) == 0 {
			retired_status = emoji.Sprint(":white_check_mark:")
		} else {
			retired_status = emoji.Sprint(":x:")
		}
		t, err := time.Parse(time.RFC3339Nano, data["release"].(string))
		if err != nil {
			fmt.Println("Erreur when date parsing :", err)
			return
		}
		datetime := t.Format("2006-01-02")
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n", data["name"], data["os"], retired_status, data["difficultyText"], data["stars"], data["ip"], status, datetime)
		w.Flush()
	}
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Showcase detailed machine information",
	Long:  "Displays detailed information of the specified machines in a structured table.",
	Run: func(cmd *cobra.Command, args []string) {
		core_infoCmd(machineParam, challengeParam)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
	infoCmd.Flags().StringSliceVarP(&machineParam, "machine", "m", []string{}, "Machine name")
	infoCmd.Flags().StringSliceVarP(&challengeParam, "challenge", "c", []string{}, "Challenge name")
}
