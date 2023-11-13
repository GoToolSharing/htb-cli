package webhooks

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/utils"
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
	_, err = utils.HTTPRequest(http.MethodPost, config.GlobalConf["Discord"], "", jsonData)
	if err != nil {
		return err
	}
	return nil
}
