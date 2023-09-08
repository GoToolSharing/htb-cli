package utils

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/briandowns/spinner"
)

func HtbRequest(method string, urlParam string, proxyURL string, jsonData []byte) (*http.Response, error) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		s.Stop()
		os.Exit(0)
	}()

	s.Start()
	JWT_TOKEN := GetHTBToken()

	req, err := http.NewRequest(method, urlParam, bytes.NewBuffer(jsonData))
	if err != nil {
		s.Stop()
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", "HTB-Tool")
	req.Header.Set("Authorization", "Bearer "+JWT_TOKEN)

	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
	} else if method == http.MethodGet {
		req.Header.Set("Host", "www.hackthebox.com")
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
			return nil, fmt.Errorf("error parsing proxy url : %v", err)
		}
		transport.Proxy = http.ProxyURL(proxyURLParsed)
	}

	log.Println("HTTP request URL :", req.URL)
	log.Println("HTTP request method :", req.Method)
	log.Println("HTTP request body :", req.Body)
	log.Println("HTTP request headers :", req.Header)

	client := &http.Client{Transport: transport}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	s.Stop()
	return resp, nil
}

func HtbGet(url string) *http.Response {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	JWT_TOKEN := GetHTBToken()
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("User-Agent", "HTB-Tool")
	req.Header.Set("Host", "www.hackthebox.com")
	req.Header.Set("Authorization", "Bearer "+JWT_TOKEN)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	return resp
}

func HtbPost(url string, jsonData []byte) *http.Response {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	JWT_TOKEN := GetHTBToken()
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("User-Agent", "HTB-Tool")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+JWT_TOKEN)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	return resp
}
