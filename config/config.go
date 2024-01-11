package config

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Settings struct {
	Verbose    int
	Logger     *zap.Logger
	ProxyParam string
	BatchParam bool
	NoCheck    bool
}

var GlobalConfig Settings

var ConfigFile map[string]string

var homeDir = os.Getenv("HOME")

var BaseDirectory = homeDir + "/.local/htb-cli"

const HostHackTheBox = "labs.hackthebox.com"

const BaseHackTheBoxAPIURL = "https://" + HostHackTheBox + "/api/v4"

const StatusURL = "https://status.hackthebox.com/api/v2/status.json"

const Version = "32501bfab6e93e517a12b7017130f96e6b933311"

func ConfigureLogger() error {
	var logLevel zapcore.Level

	switch GlobalConfig.Verbose {
	case 0:
		logLevel = zap.ErrorLevel
	case 1:
		logLevel = zap.InfoLevel
	case 2:
		logLevel = zap.DebugLevel
	default:
		logLevel = zap.DebugLevel
	}

	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // Ajoute des couleurs pour les niveaux de log

	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(logLevel),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	var err error
	GlobalConfig.Logger, err = cfg.Build()
	if err != nil {
		return fmt.Errorf("logger configuration error: %v", err)
	}
	zap.ReplaceGlobals(GlobalConfig.Logger)
	return nil
}

// LoadConfig reads a configuration file from a specified filepath and returns a map of key-value pairs.
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
			return nil, fmt.Errorf("incorrectly formatted line in configuration file: %s", line)
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

// validateConfig checks if the provided key-value pairs in the configuration are valid.
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

// isValidDiscordWebhook checks if a given URL is a valid Discord webhook.
func isValidDiscordWebhook(u string) bool {
	parsedURL, err := url.Parse(u)
	return err == nil && parsedURL.Scheme == "https" && strings.Contains(parsedURL.Host, "discord.com") && strings.Contains(parsedURL.Path, "/api/webhooks/")
}

// isValidHTTPorHTTPSURL checks if a given URL is valid and uses either the HTTP or HTTPS protocol.
func isValidHTTPorHTTPSURL(u string) bool {
	parsedURL, err := url.Parse(u)
	return err == nil && (parsedURL.Scheme == "http" || parsedURL.Scheme == "https")
}

// Init initializes the application by setting up necessary directories, creating a default configuration file if it doesn't exist, and loading the configuration.
func Init() error {
	if _, err := os.Stat(BaseDirectory); os.IsNotExist(err) {
		GlobalConfig.Logger.Info(fmt.Sprintf("The \"%s\" folder does not exist, creation in progress...\n", BaseDirectory))
		err := os.MkdirAll(BaseDirectory, os.ModePerm)
		if err != nil {
			return fmt.Errorf("error folder creation: %s", err)
		}

		GlobalConfig.Logger.Info(fmt.Sprintf("\"%s\" folder created successfully\n\n", BaseDirectory))
	}

	confFilePath := BaseDirectory + "/default.conf"
	if _, err := os.Stat(confFilePath); os.IsNotExist(err) {
		file, err := os.Create(confFilePath)
		if err != nil {
			return fmt.Errorf("error creating file: %v", err)
		}
		defer file.Close()

		configContent := `Discord = False
		Update = False`

		writer := bufio.NewWriter(file)
		_, err = writer.WriteString(configContent)
		if err != nil {
			return fmt.Errorf("error when writing to file: %v", err)
		}

		err = writer.Flush()
		if err != nil {
			return fmt.Errorf("error clearing buffer: %v", err)
		}

		GlobalConfig.Logger.Info("Configuration file created successfully.")
	}

	GlobalConfig.Logger.Info("Loading configuration file...")
	config, err := LoadConfig(BaseDirectory + "/default.conf")
	if err != nil {
		return fmt.Errorf("error loading configuration file: %v", err)
	}

	GlobalConfig.Logger.Info("Configuration successfully loaded")
	GlobalConfig.Logger.Debug(fmt.Sprintf("%v", config))
	ConfigFile = config
	return nil
}
