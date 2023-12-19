package webhooks

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/lib/utils"
)

// SendToDiscord sends a message to a Discord channel using a webhook URL.
func SendToDiscord(message string) error {
	payload := map[string]string{
		"content": message,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to create JSON data: %w", err)
	}
	if config.ConfigFile["Discord"] != "False" {
		_, err = utils.HTTPRequest(http.MethodPost, config.ConfigFile["Discord"], jsonData)
		if err != nil {
			return err
		}
	}
	return nil
}
