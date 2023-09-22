package cmd

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

type PageStatus struct {
	Status Status `json:"status"`
}

type Status struct {
	Description string `json:"description"`
}

func core_status(proxyParam string) (string, error) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		s.Stop()
		os.Exit(0)
	}()

	s.Start()
	var statusURL = "https://status.hackthebox.com/api/v2/status.json"
	req, err := http.NewRequest(http.MethodGet, statusURL, nil)
	if err != nil {
		s.Stop()
		return "", err
	}

	req.Header.Set("User-Agent", "HTB-Tool")
	req.Header.Set("Host", "status.hackthebox.com")

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	if proxyParam != "" {
		log.Println("Proxy URL found :", proxyParam)
		proxyURLParsed, err := url.Parse(proxyParam)
		if err != nil {
			s.Stop()
			return "", errors.New(fmt.Sprintf("Error: parsing proxy url: %v", err))
		}
		transport.Proxy = http.ProxyURL(proxyURLParsed)
	}

	log.Println("HTTP request URL :", req.URL)
	log.Println("HTTP request method :", req.Method)
	log.Println("HTTP request body :", req.Body)

	client := &http.Client{Transport: transport}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	s.Stop()
	body, _ := io.ReadAll(resp.Body)
	var pageStatus PageStatus
	err = json.Unmarshal([]byte(body), &pageStatus)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error: JSON decode: %v", err))
	}
	description := pageStatus.Status.Description
	return description, nil
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Displays the status of HackTheBox servers",
	Run: func(cmd *cobra.Command, args []string) {
		output, err := core_status(proxyParam)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(output)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
