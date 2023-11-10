package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/briandowns/spinner"
)

func SendDiscordWebhook(message string) error {
	payload := map[string]string{
		"content": message,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to create JSON data: %w", err)
	}
	_, err = HTTPRequest(http.MethodPost, config.GlobalConf["Discord"], "", jsonData)
	if err != nil {
		return err
	}
	return nil
}

// HTTPRequest makes an HTTP request
func HTTPRequest(method string, urlParam string, proxyURL string, jsonData []byte) (bool, error) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		s.Stop()
		os.Exit(0)
	}()

	s.Start()

	req, err := http.NewRequest(method, urlParam, bytes.NewBuffer(jsonData))
	if err != nil {
		s.Stop()
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", "htb-cli")

	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	if proxyURL != "" {
		log.Println("Proxy URL found :", proxyURL)
		proxyURLParsed, err := url.Parse(proxyURL)
		if err != nil {
			s.Stop()
			return false, fmt.Errorf("error parsing proxy url : %v", err)
		}
		transport.Proxy = http.ProxyURL(proxyURLParsed)
	}

	log.Println("HTTP request URL :", req.URL)
	log.Println("HTTP request method :", req.Method)
	log.Println("HTTP request body :", req.Body)

	client := &http.Client{Transport: transport}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	s.Stop()
	return true, nil
}
