package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/briandowns/spinner"
)

type Machine struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

type Root struct {
	Machines interface{} `json:"machines"`
}

func GetHTBToken() string {
	var envName = "HTB_TOKEN"
	if os.Getenv(envName) == "" {
		fmt.Printf("Environment variable is not set : %v\n", envName)
		return ""
	}
	return os.Getenv("HTB_TOKEN")
}

func SearchMachineIDByName(machine_name string, proxyURL string) string {
	url := "https://www.hackthebox.com/api/v4/search/fetch?query=" + machine_name
	resp, err := HtbRequest(http.MethodGet, url, proxyURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	json_body, _ := io.ReadAll(resp.Body)

	var root Root
	err = json.Unmarshal([]byte(json_body), &root)
	if err != nil {
		fmt.Println("error:", err)
	}

	// The HackTheBox API can return either a slice or a map
	switch root.Machines.(type) {
	case []interface{}:
		var machines []Machine
		machineData, _ := json.Marshal(root.Machines)
		json.Unmarshal(machineData, &machines)
		return machines[0].ID
	case map[string]interface{}:
		var machines map[string]Machine
		machineData, _ := json.Marshal(root.Machines)
		json.Unmarshal(machineData, &machines)
		return machines["0"].ID
	default:
		log.Fatal("No machine found")
	}
	return ""
}

func ParseJsonMessage(resp *http.Response, key string) interface{} {
	json_body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal([]byte(json_body), &result)
	return result[key]
}

func GetMachineType(machine_id interface{}, proxyURL string) string {
	// Check if the machine is the latest release
	url := "https://www.hackthebox.com/api/v4/machine/recommended/"
	resp, err := HtbRequest(http.MethodGet, url, proxyURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	card := ParseJsonMessage(resp, "card1").(map[string]interface{})
	fmachine_id, _ := strconv.ParseFloat(machine_id.(string), 64)
	if card["id"].(float64) == fmachine_id {
		return "release"
	}

	// Check if the machine is active or retired
	url = "https://www.hackthebox.com/api/v4/machine/profile/" + machine_id.(string)
	resp, err = HtbRequest(http.MethodGet, url, proxyURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	info := ParseJsonMessage(resp, "info").(map[string]interface{})
	if info["active"].(float64) == 1 {
		return "active"
	} else if info["retired"].(float64) == 1 {
		return "retired"
	}
	return "error"
}

func GetUserSubscription(proxyURL string) string {
	url := "https://www.hackthebox.com/api/v4/user/info"
	resp, err := HtbRequest(http.MethodGet, url, proxyURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	info := ParseJsonMessage(resp, "info").(map[string]interface{})
	if info["canAccessVIP"].(bool) {
		return "vip"
	}
	return "free"
}

func GetActiveMachineID(proxyURL string) string {
	url := "https://www.hackthebox.com/api/v4/machine/active"
	resp, err := HtbRequest(http.MethodGet, url, proxyURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	info := ParseJsonMessage(resp, "info")
	if info == nil {
		fmt.Println("No machine is running")
		os.Exit(0)
	}
	return fmt.Sprintf("%.0f", info.(map[string]interface{})["id"].(float64))
}

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
	if JWT_TOKEN == "" {
		s.Stop()
		os.Exit(1)
	}

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

	client := &http.Client{Transport: transport}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	resp.Body = io.NopCloser(bytes.NewReader(body))
	var i interface{}
	if json.Unmarshal(body, &i) != nil {
		s.Stop()
		fmt.Println("Your token is invalid or expired")
		os.Exit(1)
	}
	s.Stop()
	return resp, nil
}
