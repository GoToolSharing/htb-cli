package utils

import (
	"encoding/json"
	"errors"
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

func submitTask(proxyURL string, sherlockID string, taskID string, flag string) (string, error) {
	url := "https://www.hackthebox.com/api/v4/sherlocks/" + sherlockID + "/tasks/" + taskID + "/flag"

	body := map[string]string{
		"flag": flag,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("failed to create JSON data: %w", err)
	}
	resp, err := HtbRequest(http.MethodPost, url, proxyURL, jsonBody)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	message, ok := ParseJsonMessage(resp, "message").(string)
	if !ok {
		return "", errors.New("unexpected response format")
	}
	return message, nil
}

func GetSherlockTasks(proxyURL string, sherlockID string) (bool, error) {
	url := "https://www.hackthebox.com/api/v4/sherlocks/" + sherlockID + "/tasks"
	resp, err := HtbRequest(http.MethodGet, url, proxyURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	jsonData, _ := io.ReadAll(resp.Body)

	var parsedData SherlockDataTasks
	err = json.Unmarshal([]byte(jsonData), &parsedData)
	if err != nil {
		return true, fmt.Errorf("error parsing JSON: ", err)
	}

	for _, task := range parsedData.Tasks {
		if !task.Completed {
			fmt.Printf("ID: %d, Title: %s, Description: %s\n", task.ID, task.Title, task.Description)
			fmt.Print("Flag : ")
			taskID := strconv.Itoa(task.ID)
			flagByte, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				fmt.Println("Error reading flag")
				// return "", fmt.Errorf("error reading flag")
			}
			flag := string(flagByte)
			log.Println(flag)

			message, err := submitTask(proxyURL, sherlockID, taskID, flag)

			if err != nil {
				return true, err
			}

			fmt.Println(message)

			return false, nil
		}
	}
	return true, nil
}

func GetSherlockGeneralInformations(proxyURL string, sherlockID string) {
	url := "https://www.hackthebox.com/api/v4/sherlocks/" + sherlockID + "/play"
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

func SearchSherlockIDByName(proxyURL string, sherlockSearch string) (string, error) {
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
		return "", fmt.Errorf("error parsing JSON: %s", err)
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
		return strconv.Itoa(matchedNameID.ID), nil
	}

	return "", fmt.Errorf("error: Nothing was found")
}
