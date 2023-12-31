package sherlocks

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/lib/utils"
	"github.com/chzyer/readline"
	"github.com/sahilm/fuzzy"
)

// getSherlockDownloadLink constructs and returns the download link for a specific Sherlock challenge.
func getDownloadLink(sherlockID string) (string, error) {
	url := fmt.Sprintf("%s/sherlocks/%s/download_link", config.BaseHackTheBoxAPIURL, sherlockID)

	// url := "https://www.hackthebox.com/api/v4/challenge/download/196"

	// return url, nil

	resp, err := utils.HtbRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error: Sherlock is not available for now")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var data DownloadFile
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}

	config.GlobalConfig.Logger.Debug(fmt.Sprintf("Download URL: %s", data.URL))
	return data.URL, nil
}

// downloadFile downloads the Sherlock file from a given URL to a specified download path.
func downloadFile(url string, downloadPath string) error {
	resp, err := utils.HtbRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("error: Status code:", resp.StatusCode)
		return nil
	}

	outFile, err := os.Create(downloadPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("Archive downloaded successfully. The password for unlock is: hacktheblue")
	fmt.Println("")
	return nil
}

// submitTask sends a flag for a specific task of a Sherlock challenge and returns the server's response.
func submitTask(sherlockID string, taskID string, flag string) (string, error) {
	url := fmt.Sprintf("%s/sherlocks/%s/tasks/%s/flag", config.BaseHackTheBoxAPIURL, sherlockID, taskID)

	body := map[string]string{
		"flag": flag,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("failed to create JSON data: %w", err)
	}
	resp, err := utils.HtbRequest(http.MethodPost, url, jsonBody)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	message, ok := utils.ParseJsonMessage(resp, "message").(string)
	if !ok {
		return "", errors.New("unexpected response format")
	}
	return message, nil
}

// GetTaskByID retrieves and prints the description of a specific task of a Sherlock challenge.
func GetTaskByID(sherlockID string, sherlockTaskID int, sherlockHint bool) error {
	// TODO: Add hint
	url := fmt.Sprintf("%s/sherlocks/%s/tasks", config.BaseHackTheBoxAPIURL, sherlockID)
	resp, err := utils.HtbRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	jsonData, _ := io.ReadAll(resp.Body)

	var sherlockData SherlockDataTasks
	err = json.Unmarshal([]byte(jsonData), &sherlockData)
	if err != nil {
		return fmt.Errorf("error parsing JSON: %w", err)
	}

	if sherlockTaskID >= 1 && sherlockTaskID <= len(sherlockData.Tasks) {
		if sherlockHint && sherlockData.Tasks[sherlockTaskID-1].Hint != "" {
			fmt.Printf("\n%s :\n%s\n\nHint : %s\nMasked Flag : %s\n", sherlockData.Tasks[sherlockTaskID-1].Title, sherlockData.Tasks[sherlockTaskID-1].Description, sherlockData.Tasks[sherlockTaskID-1].Hint, sherlockData.Tasks[sherlockTaskID-1].MaskedFlag)
		} else {
			fmt.Printf("\n%s :\n%s\n\nMasked Flag : %s\n", sherlockData.Tasks[sherlockTaskID-1].Title, sherlockData.Tasks[sherlockTaskID-1].Description, sherlockData.Tasks[sherlockTaskID-1].MaskedFlag)
		}
		rl, err := readline.New("Answer: ")
		if err != nil {
			panic(err)
		}
		defer rl.Close()

		flag, err := rl.Readline()
		if err != nil {
			return err
		}
		flag = strings.TrimSpace(flag)
		config.GlobalConfig.Logger.Debug(fmt.Sprintf("Flag: %s", flag))
		taskID := strconv.Itoa(sherlockData.Tasks[sherlockTaskID-1].ID)

		message, err := submitTask(sherlockID, taskID, flag)

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
func GetTasks(sherlockID string) (*SherlockDataTasks, error) {
	url := fmt.Sprintf("%s/sherlocks/%s/tasks", config.BaseHackTheBoxAPIURL, sherlockID)
	resp, err := utils.HtbRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
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
func GetGeneralInformations(sherlockID string, sherlockDownloadPath string) error {
	url := fmt.Sprintf("%s/sherlocks/%s/play", config.BaseHackTheBoxAPIURL, sherlockID)
	resp, err := utils.HtbRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	info := utils.ParseJsonMessage(resp, "data").(map[string]interface{})

	if sherlockDownloadPath != "" {
		url, err := getDownloadLink(sherlockID)
		if err != nil {
			return err
		}
		err = downloadFile(url, sherlockDownloadPath)
		if err != nil {
			return err
		}
	}

	config.GlobalConfig.Logger.Debug(fmt.Sprintf("Informations: %v", info))
	fmt.Println("Scenario :", info["scenario"])
	fmt.Println("\nFile :", info["file_name"])
	fmt.Println("File Size :", info["file_size"])
	return nil
}

// SearchIDByName searches for a Sherlock challenge by name and returns its ID.
func SearchIDByName(sherlockSearch string) (string, error) {
	url := fmt.Sprintf("%s/sherlocks", config.BaseHackTheBoxAPIURL)
	resp, err := utils.HtbRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
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
		isConfirmed := utils.AskConfirmation("The following sherlock was found : " + matchedNameID.Name)
		if isConfirmed {
			return strconv.Itoa(matchedNameID.ID), nil
		}
	}

	return "", fmt.Errorf("error: Nothing was found")
}
