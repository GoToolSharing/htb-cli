package webhooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/lib/utils"
)

// SendToDiscord sends a message to a Discord channel using a webhook URL.
func SendToDiscord(command string, message string) error {
	if config.ConfigFile["Discord"] == "False" {
		return nil
	}
	embed := Embed{
		Title:       "htb-cli",
		Description: fmt.Sprintf("User **%s** used the **%s** command.\n**Message:** %s", utils.GetCurrentUsername(), command, message),
		Color:       12345,
		Thumbnail: EmbedThumbnail{
			URL: "https://github.com/GoToolSharing/htb-cli/blob/main/assets/logo.png?raw=true",
		},
	}
	payload := map[string]interface{}{
		"embeds": []Embed{embed},
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to create JSON data: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, config.ConfigFile["Discord"], bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	return nil
}
