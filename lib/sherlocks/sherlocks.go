package sherlocks

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/lib/utils"
	"github.com/chzyer/readline"
	"github.com/sahilm/fuzzy"
)

var metaRefreshURLPattern = regexp.MustCompile(`url=['"]?([^'">]+)['"]?`)

func getDownloadLink(sherlockID string) (string, error) {
	url := fmt.Sprintf("%s/sherlocks/%s/download_link", config.BaseHackTheBoxAPIURL, sherlockID)

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

// getSherlockDownloadLink constructs and returns the download link for a specific Sherlock challenge.
func resolveDownloadURL(redirectEndpointURL string) (string, error) {
	resp, err := utils.HtbRequest(http.MethodGet, redirectEndpointURL, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		location := resp.Header.Get("Location")
		if location == "" {
			return "", fmt.Errorf("error: redirect response (status %d) had no Location header", resp.StatusCode)
		}
		return location, nil
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error: unexpected status code %d while resolving download link", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	matches := metaRefreshURLPattern.FindSubmatch(body)
	if len(matches) < 2 {
		return "", fmt.Errorf("error: could not find redirect target in HTML response")
	}
	return html.UnescapeString(string(matches[1])), nil
}

// downloadFile downloads the Sherlock file from a given URL to a specified download path.
func downloadFile(downloadLinkURL string, downloadPath string) error {
	finalURL, err := resolveDownloadURL(downloadLinkURL)
	if err != nil {
		return err
	}

	config.GlobalConfig.Logger.Debug(fmt.Sprintf("Resolved final download URL: %s", finalURL))

	resp, err := http.Get(finalURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: unexpected status code %d while downloading file", resp.StatusCode)
	}

	outFile, err := os.Create(downloadPath)
	if err != nil {
		return err
	}

	_, err = io.Copy(outFile, resp.Body)
	outFile.Close()
	if err != nil {
		return err
	}

	fmt.Println("Archive downloaded successfully to:", downloadPath)
	fmt.Println("The password for unlock is: hacktheblue")

	return nil
}

// submitTask sends a flag for a specific task of a Sherlock challenge and returns the server's response.
func SubmitTask(sherlockID string, taskID string, flag string) (string, error) {
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

		message, err := SubmitTask(sherlockID, taskID, flag)

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
		downloadURL, err := getDownloadLink(sherlockID)
		if err != nil {
			return err
		}
		err = downloadFile(downloadURL, sherlockDownloadPath)
		if err != nil {
			return err
		}
	}

	config.GlobalConfig.Logger.Debug(fmt.Sprintf("Informations: %v", info))
	fmt.Println("\nScenario :", info["scenario"])
	fmt.Println("\nFile :", info["file_name"])
	fmt.Println("File Size :", info["file_size"])
	return nil
}

// SearchIDByName searches for a Sherlock challenge by name and returns its ID.
func SearchIDByName(sherlockSearch string) (string, error) {
	apiURL := fmt.Sprintf("%s/sherlocks?keyword=%s", config.BaseHackTheBoxAPIURL, url.QueryEscape(strings.ToLower(sherlockSearch)))
	resp, err := utils.HtbRequest(http.MethodGet, apiURL, nil)
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

	if len(matches) == 0 && parsedData.Message != "" {
		return "", fmt.Errorf("API error: %s", parsedData.Message)
	}

	for _, match := range matches {
		matchedNameID := nameIDs[match.Index]
		isConfirmed := utils.AskConfirmation("The following sherlock was found : " + matchedNameID.Name)
		if isConfirmed {
			return strconv.Itoa(matchedNameID.ID), nil
		}
	}

	return "", fmt.Errorf("error: Nothing was found")
}
