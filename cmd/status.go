package cmd

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/lib/webhooks"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// PageStatus represents the status structure fetched from the API.
type PageStatus struct {
	Status Status `json:"status"`
}

// Status contains the description of the status.
type Status struct {
	Description string `json:"description"`
}

// setupSignalHandler configures a signal handler to stop the spinner and gracefully exit upon receiving specific signals.
func setupSignalHandler(s *spinner.Spinner) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		s.Stop()
		os.Exit(0)
	}()
}

// createClient creates and returns an HTTP client with optional configurations, such as the proxy parameter.
func createClient() (*http.Client, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	if config.GlobalConfig.ProxyParam != "" {
		proxyURLParsed, err := url.Parse(config.GlobalConfig.ProxyParam)
		if err != nil {
			return nil, fmt.Errorf("error parsing proxy URL: %v", err)
		}
		transport.Proxy = http.ProxyURL(proxyURLParsed)
	}

	return &http.Client{Transport: transport}, nil
}

// fetchStatus makes an HTTP request to fetch the status and returns the status description.
func fetchStatus(client *http.Client) (string, error) {
	req, err := http.NewRequest(http.MethodGet, config.StatusURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "htb-cli")
	req.Header.Set("Host", "status.hackthebox.com")

	config.GlobalConfig.Logger.Info("Sending an HTTP request")
	config.GlobalConfig.Logger.Debug(fmt.Sprintf("Request URL: %v", req.URL))
	config.GlobalConfig.Logger.Debug(fmt.Sprintf("Request method: %v", req.Method))
	config.GlobalConfig.Logger.Debug(fmt.Sprintf("Request body: %v", req.Body))

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	var pageStatus PageStatus
	err = json.Unmarshal(body, &pageStatus)
	if err != nil {
		return "", fmt.Errorf("error decoding JSON: %v", err)
	}

	return pageStatus.Status.Description, nil
}

// coreStatusCmd is the main function that orchestrates client creation, fetching the status, and displaying the status.
func coreStatusCmd() (string, error) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	setupSignalHandler(s)
	s.Start()
	defer s.Stop()

	client, err := createClient()
	if err != nil {
		return "", err
	}
	description, err := fetchStatus(client)
	if err != nil {
		return "", err
	}
	s.Stop()
	return description, nil
}

// statusCmd defines the Cobra command to display the status of HackTheBox servers.
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Displays the status of hackthebox servers",
	Run: func(cmd *cobra.Command, args []string) {
		config.GlobalConfig.Logger.Info("Status command executed")
		output, err := coreStatusCmd()
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}
		fmt.Println(output)
		err = webhooks.SendToDiscord("status", output)
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}
		config.GlobalConfig.Logger.Info("Exit status command correctly")
	},
}

// init adds the status command to the root command during the package initialization.
func init() {
	rootCmd.AddCommand(statusCmd)
}
