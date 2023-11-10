package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/sahilm/fuzzy"
	"golang.org/x/term"
)

type SherlockTask struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type SherlockDataTasks struct {
	Tasks []SherlockTask `json:"data"`
}

type SherlockElement struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type SherlockData struct {
	Data []SherlockElement `json:"data"`
}

type SherlockNameID struct {
	Name string
	ID   int
}

func GetSherlockTasks(proxyURL string, sherlockID int) {
	url := "https://www.hackthebox.com/api/v4/sherlocks/" + strconv.Itoa(sherlockID) + "/tasks"
	resp, err := HtbRequest(http.MethodGet, url, proxyURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	jsonData, _ := io.ReadAll(resp.Body)

	var parsedData SherlockDataTasks
	err = json.Unmarshal([]byte(jsonData), &parsedData)
	if err != nil {
		fmt.Println("Error parsing JSON: ", err)
		return
	}

	for _, task := range parsedData.Tasks {
		if !task.Completed {
			fmt.Printf("ID: %d, Title: %s, Description: %s\n", task.ID, task.Title, task.Description)
			fmt.Print("Flag : ")
			flagByte, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				fmt.Println("Error reading flag")
				// return "", fmt.Errorf("error reading flag")
			}
			flagOriginal := string(flagByte)
			fmt.Println(flagOriginal)
		}
	}
}

func GetSherlockGeneralInformations(proxyURL string, sherlockID int) {
	url := "https://www.hackthebox.com/api/v4/sherlocks/" + strconv.Itoa(sherlockID) + "/play"
	resp, err := HtbRequest(http.MethodGet, url, proxyURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	info := ParseJsonMessage(resp, "data").(map[string]interface{})

	log.Println(info)
	fmt.Println("Scenario :", info["scenario"])
	fmt.Println("File :", info["file_name"])
	fmt.Println("File Size :", info["file_size"])
}

func SearchSherlockIDByName(proxyURL string, sherlockSearch string) (int, error) {
	url := "https://www.hackthebox.com/api/v4/sherlocks"
	resp, err := HtbRequest(http.MethodGet, url, proxyURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// info := ParseJsonMessage(resp, "data")
	// log.Println("Info :", info)

	jsonData, _ := io.ReadAll(resp.Body)

	// fmt.Println(jsonData)

	var parsedData SherlockData
	err = json.Unmarshal([]byte(jsonData), &parsedData)
	if err != nil {
		return 1, fmt.Errorf("error parsing JSON: %s", err)
	}

	var nameIDs []SherlockNameID
	for _, challenge := range parsedData.Data {
		nameIDs = append(nameIDs, SherlockNameID{challenge.Name, challenge.ID})
	}

	var names []string
	for _, ni := range nameIDs {
		names = append(names, ni.Name)
	}

	matches := fuzzy.Find(sherlockSearch, names)

	for _, match := range matches {
		matchedNameID := nameIDs[match.Index]
		log.Printf("Found: %s with ID: %d\n", matchedNameID.Name, matchedNameID.ID)
		return matchedNameID.ID, nil
	}

	return 1, fmt.Errorf("error: Nothing was found")
}
