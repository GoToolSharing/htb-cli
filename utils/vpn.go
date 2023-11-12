package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func DownloadVPN(proxyURL string) error {
	url := "https://www.hackthebox.com/api/v4/connections/servers?product=labs"
	// url := "https://www.hackthebox.com/api/v4/connections/servers?product=starting_point"
	// url := "https://www.hackthebox.com/api/v4/connections/servers?product=endgames"
	// url := "https://www.hackthebox.com/api/v4/connections/servers?product=fortresses"
	// url := "https://www.hackthebox.com/api/v4/connections/servers?product=competitive"

	// {"status":true,"data":{"disabled":false,"assigned":{"id":18,"friendly_name":"EU VIP 7","current_clients":8,"location":"EU","location_type_friendly":"EU - VIP"},"options":{"EU":{"EU - Free"
	resp, err := HtbRequest(http.MethodGet, url, proxyURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("error:", resp.StatusCode)
		return nil
	}

	return nil
}

func StartVPN(configPath string) string {
	fmt.Println("VPN is starting...")
	pidFile := baseDirectory + "/lab-vpn.pid"
	cmd := exec.Command("sudo", "openvpn", "--config", configPath, "--writepid", pidFile)

	err := cmd.Start()
	if err != nil {
		return ""
	}

	fmt.Println("Wait 20s and check if the vpn is started...")
	time.Sleep(20 * time.Second)
	isActive := CheckVPN("")
	if isActive {
		fmt.Println("The VPN is now active !")
	}

	return ""
}

func CheckVPN(proxyURL string) bool {
	url := "https://www.hackthebox.com/api/v4/connection/status"
	resp, err := HtbRequest(http.MethodGet, url, proxyURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var result []map[string]interface{}
	jsonBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %s\n", err)
		return false
	}

	err = json.Unmarshal(jsonBody, &result)
	if err != nil {
		log.Printf("Error unmarshalling JSON: %s\n", err)
		return false
	}

	if len(result) == 0 {
		return false
	}

	for _, item := range result {
		connectionData, ok := item["connection"].(map[string]interface{})
		if !ok {
			log.Printf("Error asserting connection data\n")
			return false
		}

		name := connectionData["name"]
		ip4 := connectionData["ip4"]

		fmt.Printf("Name: %s, IP4: %s\n", name, ip4)
	}
	return true
}

func StopVPN() error {
	fmt.Println("Try to stop the active VPN...")
	pidFile := baseDirectory + "/lab-vpn.pid"
	pidData, err := os.ReadFile(pidFile)
	if err != nil {
		return err
	}
	log.Println("VPN PID :", string(pidData))

	cmd := exec.Command("sudo", "kill", string(pidData))

	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
