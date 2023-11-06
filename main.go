package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/GoToolSharing/htb-cli/cmd"
)

func LoadConfig(filepath string) (map[string]string, error) {
	config := make(map[string]string)

	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("incorrectly formatted line in configuration file : %s", line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if err := validateConfig(key, value); err != nil {
			return nil, err
		}

		config[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return config, nil
}

func validateConfig(key, value string) error {
	switch key {
	case "Logging", "Batch":
		if value != "True" && value != "False" {
			return fmt.Errorf("the value for '%s' must be 'True' or 'False', got : %s", key, value)
		}
	case "Proxy":
		if value != "False" && !isValidHTTPorHTTPSURL(value) {
			return fmt.Errorf("the URL for '%s' must be a valid URL starting with http or https, got : %s", key, value)
		}
	case "Discord":
		if value != "False" && !isValidDiscordWebhook(value) {
			return fmt.Errorf("the Discord webhook URL is invalid : %s", value)
		}
	}

	return nil
}

func isValidDiscordWebhook(u string) bool {
	parsedURL, err := url.Parse(u)
	return err == nil && parsedURL.Scheme == "https" && strings.Contains(parsedURL.Host, "discord.com") && strings.Contains(parsedURL.Path, "/api/webhooks/")
}

func isValidHTTPorHTTPSURL(u string) bool {
	parsedURL, err := url.Parse(u)
	return err == nil && (parsedURL.Scheme == "http" || parsedURL.Scheme == "https")
}

func main() {
	log.SetOutput(io.Discard)
	if _, err := os.Stat(cmd.BaseDirectory); os.IsNotExist(err) {
		log.Printf("The \"%s\" folder does not exist, creation in progress...\n", cmd.BaseDirectory)
		err := os.MkdirAll(cmd.BaseDirectory, os.ModePerm)
		if err != nil {
			fmt.Printf("Folder creation error: %s\n", err)
			return
		}

		log.Printf("\"%s\" folder created successfully\n\n", cmd.BaseDirectory)
	}

	confFilePath := cmd.BaseDirectory + "/default.conf"
	if _, err := os.Stat(confFilePath); os.IsNotExist(err) {
		file, err := os.Create(confFilePath)
		if err != nil {
			fmt.Printf("Error creating file: %v\n", err)
			return
		}
		defer file.Close()

		configContent := `Logging = False
Batch = False
Proxy = False
Discord = False`

		writer := bufio.NewWriter(file)
		_, err = writer.WriteString(configContent)
		if err != nil {
			fmt.Printf("Error when writing to file: %v\n", err)
			return
		}

		err = writer.Flush()
		if err != nil {
			fmt.Printf("Error clearing buffer: %v\n", err)
			return
		}

		log.Println("Configuration file created successfully.")
	}

	// Load configuration file
	log.Println("Loading configuration file...")
	config, err := LoadConfig(cmd.BaseDirectory + "/default.conf")
	if err != nil {
		fmt.Printf("Error loading configuration file : %v\n", err)
		return
	}

	log.Println("Configuration successfully loaded :", config)

	cmd.Execute()
}

///api/v4/access/ovpnfile/18/0

// remote edge-eu-vip-7.hackthebox.eu 1337
