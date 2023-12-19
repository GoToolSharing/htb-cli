package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/lib/utils"
)

func Check(newVersion string) error {
	// Dev version
	config.GlobalConfig.Logger.Debug(fmt.Sprintf("config.Version: %s", config.Version))
	if len(config.Version) == 40 {
		config.GlobalConfig.Logger.Info("Development version detected")
		githubCommits := "https://api.github.com/repos/GoToolSharing/htb-cli/commits?sha=dev"

		resp, err := utils.HTTPRequest(http.MethodGet, githubCommits, nil)
		if err != nil {
			return err
		}
		body, err := io.ReadAll(resp.Body)
		config.GlobalConfig.Logger.Debug(fmt.Sprintf("Body: %s", utils.TruncateString(string(body), 2000)))
		if err != nil {
			return fmt.Errorf("Error when reading the response: %v", err)
		}
		var commits []Commit
		err = json.Unmarshal(body, &commits)
		if err != nil {
			return fmt.Errorf("Error when decoding JSON: %v", err)
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
		if commitHash != config.Version {
			message := fmt.Sprintf("A new update is now available (dev) ! (%s)", commitHash)
			fmt.Println(message)
			fmt.Println("Update with : git pull")
		} else {
			message := fmt.Sprintf("You're up to date (dev) ! (%s)", commitHash)
			fmt.Println(message)
		}
		return nil
	}

	// Main version
	githubVersion := "https://api.github.com/repos/GoToolSharing/htb-cli/releases/latest"

	resp, err := utils.HTTPRequest(http.MethodGet, githubVersion, nil)
	if err != nil {
		return err
	}
	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return fmt.Errorf("Error when decoding JSON: %v", err)
	}
	config.GlobalConfig.Logger.Debug(fmt.Sprintf("release.TagName : %s", release.TagName))
	config.GlobalConfig.Logger.Debug(fmt.Sprintf("config.Version : %s", config.Version))
	if release.TagName != config.Version {
		message := fmt.Sprintf("A new update is now available ! (%s)", release.TagName)
		fmt.Println(message)
		fmt.Println("Update with : go install github.com/GoToolSharing/htb-cli@latest")
	} else {
		message := fmt.Sprintf("You're up to date ! (%s)", config.Version)
		fmt.Println(message)
	}

	return nil

}
