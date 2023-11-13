package sherlocks

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/utils"
	"github.com/sahilm/fuzzy"
)

// getSherlockDownloadLink constructs and returns the download link for a specific Sherlock challenge.
func getDownloadLink(proxyURL string, sherlockID string) (string, error) {
	url := fmt.Sprintf("%s/sherlocks/%s/download_link", config.BaseHackTheBoxAPIURL, sherlockID)

	// url := "https://www.hackthebox.com/api/v4/challenge/download/196"

	// return url, nil

	resp, err := utils.HtbRequest(http.MethodGet, url, proxyURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error: Sherlock is not available for now")
	}

	data := utils.ParseJsonMessage(resp, "data").(map[string]interface{})

	log.Println("Download URL :", data["url"].(string))
	return data["url"].(string), nil
}

// downloadFile downloads the Sherlock file from a given URL to a specified download path.
func downloadFile(proxyURL string, url string, downloadPath string) error {
	resp, err := utils.HtbRequest(http.MethodGet, url, proxyURL, nil)
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

// submitTask sends a flag for a specific task of a Sherlock challenge and returns the server's response.
func submitTask(proxyURL string, sherlockID string, taskID string, flag string) (string, error) {
	url := fmt.Sprintf("%s/sherlocks/%s/tasks/%s/flag", config.BaseHackTheBoxAPIURL, sherlockID, taskID)

	body := map[string]string{
		"flag": flag,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("failed to create JSON data: %w", err)
	}
	resp, err := utils.HtbRequest(http.MethodPost, url, proxyURL, jsonBody)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	message, ok := utils.ParseJsonMessage(resp, "message").(string)
	if !ok {
		return "", errors.New("unexpected response format")
	}
	return message, nil
}

// GetTaskByID retrieves and prints the description of a specific task of a Sherlock challenge.
func GetTaskByID(proxyURL string, sherlockID string, sherlockTaskID int) error {
	// TODO: Add hint
	url := fmt.Sprintf("%s/sherlocks/%s/tasks", config.BaseHackTheBoxAPIURL, sherlockID)
	resp, err := utils.HtbRequest(http.MethodGet, url, proxyURL, nil)
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
		fmt.Print("Answer : ")
		reader := bufio.NewReader(os.Stdin)
		flag, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		flag = strings.TrimSpace(flag)
		log.Println(flag)
		taskID := strconv.Itoa(sherlockData.Tasks[sherlockTaskID-1].ID)

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

// GetTasks retrieves all tasks for a specific Sherlock challenge.
func GetTasks(proxyURL string, sherlockID string) (*SherlockDataTasks, error) {
	url := fmt.Sprintf("%s/sherlocks/%s/tasks", config.BaseHackTheBoxAPIURL, sherlockID)
	resp, err := utils.HtbRequest(http.MethodGet, url, proxyURL, nil)
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

// GetGeneralInformations retrieves and prints general information about a Sherlock challenge.
func GetGeneralInformations(proxyURL string, sherlockID string, sherlockDownloadPath string) error {
	url := fmt.Sprintf("%s/sherlocks/%s/play", config.BaseHackTheBoxAPIURL, sherlockID)
	resp, err := utils.HtbRequest(http.MethodGet, url, proxyURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	info := utils.ParseJsonMessage(resp, "data").(map[string]interface{})

	if sherlockDownloadPath != "" {
		url, err := getDownloadLink(proxyURL, sherlockID)
		if err != nil {
			return err
		}
		err = downloadFile(proxyURL, url, sherlockDownloadPath)
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

// SearchIDByName searches for a Sherlock challenge by name and returns its ID.
func SearchIDByName(proxyURL string, sherlockSearch string, batchParam bool) (string, error) {
	url := fmt.Sprintf("%s/sherlocks", config.BaseHackTheBoxAPIURL)
	resp, err := utils.HtbRequest(http.MethodGet, url, proxyURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	jsonData, _ := io.ReadAll(resp.Body)

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
		isConfirmed := utils.AskConfirmation("The following sherlock was found : "+matchedNameID.Name, batchParam)
		if isConfirmed {
			return strconv.Itoa(matchedNameID.ID), nil
		}
	}

	return "", fmt.Errorf("error: Nothing was found")
}
