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

func getSherlockDownloadLink(proxyURL string, sherlockID string) (string, error) {
	// url := "https://www.hackthebox.com/api/v4/sherlocks/" + sherlockID + "/download_link"

	url := "https://www.hackthebox.com/api/v4/challenge/download/196"

	return url, nil

	resp, err := HtbRequest(http.MethodGet, url, proxyURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error: Sherlock is not available for now")
	}

	data := ParseJsonMessage(resp, "data").(map[string]interface{})

	log.Println("Download URL :", data["url"].(string))
	return data["url"].(string), nil
}

func downloadSherlockFile(proxyURL string, url string, downloadPath string) error {
	resp, err := HtbRequest(http.MethodGet, url, proxyURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("error: Status code:", resp.StatusCode)
		return nil
	}

	outFile, err := os.Create(downloadPath)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Archive downloaded successfully. The password for unlock is: hacktheblue")
	fmt.Println("")
	return nil
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

func GetSherlockTaskByID(proxyURL string, sherlockID string, sherlockTaskID int) error {
	url := "https://www.hackthebox.com/api/v4/sherlocks/" + sherlockID + "/tasks"
	resp, err := HtbRequest(http.MethodGet, url, proxyURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	jsonData, _ := io.ReadAll(resp.Body)

	var sherlockData SherlockDataTasks
	err = json.Unmarshal([]byte(jsonData), &sherlockData)
	if err != nil {
		return fmt.Errorf("error parsing JSON: %w", err)
	}

	if sherlockTaskID >= 1 && sherlockTaskID <= len(sherlockData.Tasks) {
		fmt.Printf("\n%s :\n%s\n\n", sherlockData.Tasks[sherlockTaskID-1].Title, sherlockData.Tasks[sherlockTaskID-1].Description)
		fmt.Print("Flag : ")
		taskID := strconv.Itoa(sherlockData.Tasks[sherlockTaskID-1].ID)
		flagByte, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Println("Error reading flag")
		}
		flag := string(flagByte)
		log.Println(flag)

		message, err := submitTask(proxyURL, sherlockID, taskID, flag)

		if err != nil {
			return err
		}

		fmt.Println(message)
	} else {
		fmt.Println("Invalid task ID :", sherlockTaskID)
	}
	return nil
}

func GetSherlockTasks(proxyURL string, sherlockID string) (*SherlockDataTasks, error) {
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
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	return &parsedData, nil
}

func GetSherlockGeneralInformations(proxyURL string, sherlockID string, sherlockDownloadPath string) error {
	url := "https://www.hackthebox.com/api/v4/sherlocks/" + sherlockID + "/play"
	resp, err := HtbRequest(http.MethodGet, url, proxyURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	info := ParseJsonMessage(resp, "data").(map[string]interface{})

	if sherlockDownloadPath != "" {
		url, err := getSherlockDownloadLink(proxyURL, sherlockID)
		if err != nil {
			return err
		}
		err = downloadSherlockFile(proxyURL, url, sherlockDownloadPath)
		if err != nil {
			return err
		}
	}

	log.Println(info)
	fmt.Println("Scenario :", info["scenario"])
	fmt.Println("\nFile :", info["file_name"])
	fmt.Println("File Size :", info["file_size"])
	return nil
}

func SearchSherlockIDByName(proxyURL string, sherlockSearch string, batchParam bool) (string, error) {
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
		isConfirmed := AskConfirmation("The following sherlock was found : "+matchedNameID.Name, batchParam)
		if isConfirmed {
			return strconv.Itoa(matchedNameID.ID), nil
		}
	}

	return "", fmt.Errorf("error: Nothing was found")
}
