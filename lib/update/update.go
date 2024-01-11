package update

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/lib/utils"
)

func Check(newVersion string) (string, error) {
	// Dev version
	config.GlobalConfig.Logger.Debug(fmt.Sprintf("config.Version: %s", config.Version))
	if config.Version == "dev" {
		config.GlobalConfig.Logger.Info("Development version detected")
		return "Development version (git pull to update)", nil
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
