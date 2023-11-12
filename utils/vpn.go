package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/GoToolSharing/htb-cli/config"
)

type Assigned struct {
	ID           int    `json:"id"`
	FriendlyName string `json:"friendly_name"`
}

type Data struct {
	Disabled bool     `json:"disabled"`
	Assigned Assigned `json:"assigned"`
}

type Response struct {
	Status bool `json:"status"`
	Data   Data `json:"data"`
}

func DownloadVPN(proxyURL string) error {
	baseURL := "https://www.hackthebox.com/api/v4/connections/servers?product="
	urls := []string{
		baseURL + "labs",
		baseURL + "starting_point",
		baseURL + "endgames",
		baseURL + "fortresses",
		baseURL + "competitive",
	}

	for _, url := range urls {

		resp, err := HtbRequest(http.MethodGet, url, proxyURL, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("error:", resp.StatusCode)
			return nil
		}

		jsonData, _ := io.ReadAll(resp.Body)

		var response Response
		err = json.Unmarshal([]byte(jsonData), &response)
		if err != nil {
			panic(err)
		}

		url = "https://www.hackthebox.com/api/v4/access/ovpnfile/" + strconv.Itoa(response.Data.Assigned.ID) + "/0"
		resp, err = HtbRequest(http.MethodGet, url, proxyURL, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("error:", resp.StatusCode)
			return nil
		}

		vpnName := strings.ReplaceAll(response.Data.Assigned.FriendlyName, " ", "_")

		downloadPath := "/home/qu35t/.local/htb-cli/" + vpnName + "-vpn.ovpn"
		outFile, err := os.Create(downloadPath)
		if err != nil {
			log.Fatal(err)
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("VPN :", vpnName, "downloaded successfully")
	}

	fmt.Println("")
	fmt.Println("VPNs are located at the following path :", config.BaseDirectory)

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
