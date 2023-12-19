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
}

var GlobalConfig Settings

var ConfigFile map[string]string

var homeDir = os.Getenv("HOME")

var BaseDirectory = homeDir + "/.local/htb-cli"

const HostHackTheBox = "www.hackthebox.com"

const BaseHackTheBoxAPIURL = "https://" + HostHackTheBox + "/api/v4"

const Version = "08980722978c141314f5fad44a0be92f1ab0e923"

func ConfigureLogger() {
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
		GlobalConfig.Logger.Error(fmt.Sprintf("Logger configuration error: %v\n", err))
		os.Exit(1)
	}
	zap.ReplaceGlobals(GlobalConfig.Logger)
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
			GlobalConfig.Logger.Error(fmt.Sprintf("Incorrectly formatted line in configuration file : %s", line))
			os.Exit(1)
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
			GlobalConfig.Logger.Error(fmt.Sprintf("The value for '%s' must be 'True' or 'False', got : %s", key, value))
			os.Exit(1)
		}
	case "Proxy":
		if value != "False" && !isValidHTTPorHTTPSURL(value) {
			GlobalConfig.Logger.Error(fmt.Sprintf("The URL for '%s' must be a valid URL starting with http or https, got : %s", key, value))
			os.Exit(1)
		}
	case "Discord":
		if value != "False" && !isValidDiscordWebhook(value) {
			GlobalConfig.Logger.Error(fmt.Sprintf("The Discord webhook URL is invalid : %s", value))
			os.Exit(1)
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
			GlobalConfig.Logger.Error(fmt.Sprintf("Error folder creation: %s", err))
			os.Exit(1)
		}

		GlobalConfig.Logger.Info(fmt.Sprintf("\"%s\" folder created successfully\n\n", BaseDirectory))
	}

	confFilePath := BaseDirectory + "/default.conf"
	if _, err := os.Stat(confFilePath); os.IsNotExist(err) {
		file, err := os.Create(confFilePath)
		if err != nil {
			GlobalConfig.Logger.Error(fmt.Sprintf("Error creating file: %w", err))
			os.Exit(1)
		}
		defer file.Close()

		configContent := `Discord = False`

		writer := bufio.NewWriter(file)
		_, err = writer.WriteString(configContent)
		if err != nil {
			GlobalConfig.Logger.Error(fmt.Sprintf("Error when writing to file: %v", err))
			os.Exit(1)
		}

		err = writer.Flush()
		if err != nil {
			GlobalConfig.Logger.Error(fmt.Sprintf("Error clearing buffer: %v", err))
			os.Exit(1)
		}

		GlobalConfig.Logger.Info("Configuration file created successfully.")
	}

	GlobalConfig.Logger.Info("Loading configuration file...")
	config, err := LoadConfig(BaseDirectory + "/default.conf")
	if err != nil {
		GlobalConfig.Logger.Error(fmt.Sprintf("Error loading configuration file : %v", err))
		os.Exit(1)
	}

	GlobalConfig.Logger.Info("Configuration successfully loaded")
	GlobalConfig.Logger.Debug(fmt.Sprintf("%v", config))
	ConfigFile = config
	return nil
}
