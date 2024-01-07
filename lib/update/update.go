package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/lib/utils"
)

func Check(newVersion string) (string, error) {
	// Dev version
	config.GlobalConfig.Logger.Debug(fmt.Sprintf("config.Version: %s", config.Version))
	if len(config.Version) == 40 {
		config.GlobalConfig.Logger.Info("Development version detected")
		githubCommits := "https://api.github.com/repos/GoToolSharing/htb-cli/commits?sha=dev"

		resp, err := utils.HTTPRequest(http.MethodGet, githubCommits, nil)
		if err != nil {
			return "", err
		}
		body, err := io.ReadAll(resp.Body)
		config.GlobalConfig.Logger.Debug(fmt.Sprintf("Body: %s", utils.TruncateString(string(body), 500)))
		if strings.Contains(string(body), "API rate limit") {
			return "htb-cli cannot check for new updates at this time. Please try again later", nil
		}
		if err != nil {
			return "", fmt.Errorf("error when reading the response: %v", err)
		}
		var commits []Commit
		err = json.Unmarshal(body, &commits)
		if err != nil {
			return "", fmt.Errorf("error when decoding JSON: %v", err)
		}
		config.GlobalConfig.Logger.Debug(fmt.Sprintf("Commits : %v", commits))

		var commitHash string
		for _, commit := range commits {
			if commit.Commit.Author.Name != "Github Action" {
				config.GlobalConfig.Logger.Debug(fmt.Sprintf("Last commit hash : %s", commit.SHA))
				commitHash = commit.SHA
				break
			}
		}
		var message string
		if commitHash != config.Version {
			message = fmt.Sprintf("A new update is now available (dev) ! (%s)\nUpdate with : git pull", commitHash)
		} else {
			message = fmt.Sprintf("You're up to date (dev) ! (%s)", commitHash)
		}

		return message, nil
	}

	// Main version
	githubVersion := "https://api.github.com/repos/GoToolSharing/htb-cli/releases/latest"

	resp, err := utils.HTTPRequest(http.MethodGet, githubVersion, nil)
	if err != nil {
		return "", err
	}
	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("error when decoding JSON: %v", err)
	}
	config.GlobalConfig.Logger.Debug(fmt.Sprintf("release.TagName : %s", release.TagName))
	config.GlobalConfig.Logger.Debug(fmt.Sprintf("config.Version : %s", config.Version))
	var message string
	if release.TagName != config.Version {
		message = fmt.Sprintf("A new update is now available ! (%s)\nUpdate with : go install github.com/GoToolSharing/htb-cli@latest", release.TagName)
	} else {
		message = fmt.Sprintf("You're up to date ! (%s)", config.Version)
	}

	return message, nil
}
