package shoutbox

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/gorilla/websocket"
)

func ConnectToWebSocket() error {
	config.GlobalConfig.Logger.Info("Starting the websocket connection")
	u := url.URL{Scheme: "wss", Host: "ws-eu.pusher.com", Path: "/app/97608bf7532e6f0fe898", RawQuery: "protocol=7&client=js&version=5.1.1&flash=false"}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("Websocket connection error: %v", err)
	}
	defer c.Close()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			return fmt.Errorf("Error reading websocket message: %v", err)
		}
		config.GlobalConfig.Logger.Debug(fmt.Sprintf("Message received: %s", message))
		var msgData map[string]interface{}
		if err := json.Unmarshal(message, &msgData); err != nil {
			log.Println("Erreur lors de l'analyse du message :", err)
			continue
		}

		if msgData["event"] == "display-info" {
			if data, ok := msgData["data"].(string); ok {
				var dataContent map[string]string
				if err := json.Unmarshal([]byte(data), &dataContent); err != nil {
					log.Println("Erreur lors de l'analyse des données internes :", err)
					continue
				}

				config.GlobalConfig.Logger.Debug(fmt.Sprintf("Data content: %s", dataContent["text"]))
				extractedMessage, _ := parseOwnsMessages(dataContent["text"])
				fmt.Println(extractedMessage)
			}
		}

		var received map[string]interface{}
		if err := json.Unmarshal(message, &received); err != nil {
			return fmt.Errorf("Message parsing error: %v", err)
		}

		if received["event"] == "pusher:connection_established" {
			subscribeMessage := map[string]interface{}{
				"event": "pusher:subscribe",
				"data": map[string]interface{}{
					"auth":    "",
					"channel": "owns-channel",
				},
			}

			subscribeMessageBytes, err := json.Marshal(subscribeMessage)
			if err != nil {
				return fmt.Errorf("Error creating subscription message: %v", err)
			}

			config.GlobalConfig.Logger.Info("Channel owns subscription")
			if err := c.WriteMessage(websocket.TextMessage, subscribeMessageBytes); err != nil {
				return fmt.Errorf("Error sending subscription message: %v", err)
			}

			_, message, err := c.ReadMessage()
			if err != nil {
				return fmt.Errorf("Error reading websocket message: %v", err)
			}
			config.GlobalConfig.Logger.Debug(fmt.Sprintf("Messagess received: %s", message))
		}
	}
}

func parseOwnsMessages(message string) (string, error) {
	re := regexp.MustCompile(`text=([^&]+)`)
	matches := re.FindStringSubmatch(message)

	config.GlobalConfig.Logger.Debug(fmt.Sprintf("Matches: %s", matches))

	if len(matches) >= 2 {
		return strings.ReplaceAll(matches[1], "\\\"", "\""), nil
	}

	return "", nil
}
