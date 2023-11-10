package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/utils"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var difficultyParam int
var machineNameParam string
var challengeNameParam string

// coreSubmitCmd handles the submission of flags for machines or challenges, returning a status message or error.
func coreSubmitCmd(difficultyParam int, machineNameParam string, challengeNameParam string, proxyParam string) (string, error) {
	if difficultyParam < 1 || difficultyParam > 10 {
		return "", errors.New("difficulty must be set between 1 and 10")
	}

	// Common payload elements
	difficultyString := strconv.Itoa(difficultyParam * 10)
	payload := map[string]string{
		"difficulty": difficultyString,
	}

	url := ""

	if challengeNameParam != "" {
		log.Println("Challenge submit requested!")
		challengeID, err := utils.SearchItemIDByName(challengeNameParam, proxyParam, "Challenge", batchParam)
		if err != nil {
			return "", err
		}
		url = baseAPIURL + "/challenge/own"
		payload["challenge_id"] = challengeID
	} else if machineNameParam != "" {
		log.Println("Machine submit requested!")
		machineID, err := utils.SearchItemIDByName(machineNameParam, proxyParam, "Machine", batchParam)
		if err != nil {
			return "", err
		}
		machineType := utils.GetMachineType(machineID, proxyParam)
		log.Printf("Machine Type: %s", machineType)

		if machineType == "release" {
			url = baseAPIURL + "/arena/own"
		} else {
			url = baseAPIURL + "/machine/own"

		}
		payload["id"] = machineID
	} else if machineNameParam == "" && challengeNameParam == "" {
		machineID := utils.GetActiveMachineID(proxyParam)
		if machineID == "" {
			return "No machine is running", nil
		}
		machineType := utils.GetMachineType(machineID, proxyParam)
		log.Printf("Machine Type: %s", machineType)

		if machineType == "release" {
			url = baseAPIURL + "/arena/own"
		} else {
			url = baseAPIURL + "/machine/own"

		}
		payload["id"] = machineID
	}

	fmt.Print("Flag : ")
	flagByte, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Error reading flag")
		return "", fmt.Errorf("error reading flag")
	}
	flagOriginal := string(flagByte)
	flag := strings.ReplaceAll(flagOriginal, " ", "")
	payload["flag"] = flag

	log.Println("Flag :", flag)
	log.Println("Difficulty :", difficultyString)

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to create JSON data: %w", err)
	}

	resp, err := utils.HtbRequest(http.MethodPost, url, proxyParam, jsonData)
	if err != nil {
		return "", err
	}

	message, ok := utils.ParseJsonMessage(resp, "message").(string)
	if !ok {
		return "", errors.New("unexpected response format")
	}
	return message, nil
}

// submitCmd defines the "submit" command for submitting flags for machines or challenges.
var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit credentials (machines / challenges / arena)",
	Long:  "This command allows for the submission of user and root flags discovered on vulnerable machines / challenges",
	Run: func(cmd *cobra.Command, args []string) {
		output, err := coreSubmitCmd(difficultyParam, machineNameParam, challengeNameParam, proxyParam)
		if err != nil {
			log.Fatal(err)
		}
		if config.GlobalConf["Discord"] != "False" {
			utils.SendDiscordWebhook("[SUBMIT] - " + output)
		}
		fmt.Println(output)
	},
}

// init adds the submitCmd to rootCmd and sets flags for the "submit" command.
func init() {
	rootCmd.AddCommand(submitCmd)
	submitCmd.Flags().StringVarP(&machineNameParam, "machine_name", "m", "", "Machine Name")
	submitCmd.Flags().StringVarP(&challengeNameParam, "challenge_name", "c", "", "Challenge Name")
	submitCmd.Flags().IntVarP(&difficultyParam, "difficulty", "d", 0, "Difficulty")
	err := submitCmd.MarkFlagRequired("difficulty")
	if err != nil {
		fmt.Println(err)
	}
}

// Fortresses
// POST /api/v4/fortress/1/flag HTTP/1.1 / {"flag":"aaa"} / message
// GET /api/v4/fortresses HTTP/2 / {"status":true,"data":{"5":{"id":7,
// GET /api/v4/fortress/1 / {"status":true,"data":{"id":1,"name":"Jet","ip":"10.13.37.10","image":"https:\/\/www.hackthebox.com\/storage\/fortresses\/c4ca4238a0b923820dcc509a6f75849b_logo.svg","cover_image_url":"https:\/\/www.hackthebox.com\/storage\/fortresses\/c4ca4238a0b923820dcc509a6f75849b_cover_full.png","company":{"id":3,"name":"Jet.com","description":"<p>Jet.com is currently looking for Security Engineers in the USA.<\/p>\n\n<p><span>Jet\u2019s mission is to\nbecome the smartest way to shop and save on pretty much anything. Combining a\nrevolutionary pricing engine, a world-class technology and fulfillment\nplatform, and incredible customer service, we\u2019ve set out to create a new kind\nof e-commerce.\u00a0 At Jet, we\u2019re passionate about empowering people to live\nand work brilliant.<\/span><br><\/p>\n\n<p><span>We need super smart engineers from all levels to help us\nbuild one of the best engineered e-commerce platforms in the world (big talk we\nknow, but that is our goal!). Our engineers combine creativity, curiosity, and drive\nto continuously perfect and revolutionize Jet from the inside out. We are\nlooking to bring more intellectually curious engineers who are passionate about\ntechnology in general (Jet is a technology first company and prides itself on\nits culture of learning and knowledge sharing and we want all our engineers to\nbe as passionate as we are!) \u00a0\u00a0<\/span><br><\/p>\n\n<p>The Environment\u00a0\u00a0<\/p>\n\n<p>Our infrastructure is largely built on Microsoft Windows. We\nhave a hybrid configuration with on premise servers and cloud based servers\nusing Microsoft Azure with a large number of additional technologies and\nmiddleware. We support three warehouses, a call center, corporate headquarters,\nand the development environment in the cloud. Our team uses a mix of Windows,\nApple, and some Linux for our systems management platforms and cutting edge\nnetwork equipment. About 50% of the development platform runs on Linux and the\nrest Windows. \u00a0<\/p>\n\n<p><br><\/p>\n\n<p>Tour the office at\u00a0<a href=\"http:\/\/www.interiordesign.net\/slideshows\/detail\/9120-start-up-and-away\/\">http:\/\/www.interiordesign.net\/slideshows\/detail\/9120-start-up-and-away\/<\/a><\/p>","url":"https:\/\/jet.com\/careers","image":"https:\/\/www.hackthebox.com\/storage\/companies\/3.png"},"reset_votes":1,"description":"Lift off with this introductory fortress from Jet! Featuring interesting web vectors and challenges, this fortress is perfect for those getting started.","has_completion_message":false,"completion_message":null,"progress_percent":100,"players_completed":2284,"points":"100","user_availability":{"available":true,"code":0,"message":null},"flags":[{"id":1,"title":"Connect","points":5,"owned":true},{"id":2,"title":"Digging in...","points":5,"owned":true},{"id":3,"title":"Going Deeper","points":5,"owned":true},{"id":4,"title":"Bypassing Authentication","points":5,"owned":true},{"id":5,"title":"Command","points":5,"owned":true},{"id":6,"title":"Overflown","points":10,"owned":true},{"id":7,"title":"Secret Message","points":10,"owned":true},{"id":8,"title":"Elasticity","points":10,"owned":true},{"id":9,"title":"Member Manager","points":15,"owned":true},{"id":10,"title":"More Secrets","points":10,"owned":true},{"id":11,"title":"Memo","points":20,"owned":true}]}}
