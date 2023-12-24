package vpn

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/lib/utils"
	"github.com/GoToolSharing/htb-cli/lib/webhooks"
)

func downloadVPN(url string) error {
	resp, err := utils.HtbRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: Bad status code : %d", resp.StatusCode)
	}

	jsonData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response Response
	err = json.Unmarshal(jsonData, &response)
	if err != nil {
		return err
	}

	vpnURL := fmt.Sprintf("%s/access/ovpnfile/%d/0", config.BaseHackTheBoxAPIURL, response.Data.Assigned.ID)
	resp, err = utils.HtbRequest(http.MethodGet, vpnURL, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: Bad status code : %d", resp.StatusCode)
	}

	vpnName := strings.ReplaceAll(response.Data.Assigned.FriendlyName, " ", "_")
	downloadPath := fmt.Sprintf("%s/%s-vpn.ovpn", config.BaseDirectory, vpnName)
	outFile, err := os.Create(downloadPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("VPN :", vpnName, "downloaded successfully")
	return nil
}

// DownloadAll downloads VPN configurations from HackTheBox for different server types.
func DownloadAll() error {
	baseURL := fmt.Sprintf("%s/connections/servers?product=", config.BaseHackTheBoxAPIURL)
	urls := []string{
		baseURL + "labs",
		baseURL + "starting_point",
		baseURL + "endgames",
		baseURL + "fortresses",
		baseURL + "competitive",
	}

	var wg sync.WaitGroup
	errors := make(chan error, len(urls))

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			err := downloadVPN(url)
			if err != nil {
				errors <- err
				return
			}
		}(url)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		if err != nil {
			return err
		}
	}

	fmt.Println("")

	message := fmt.Sprintf("VPNs are located at the following path : %s", config.BaseDirectory)

	fmt.Println(message)

	err := webhooks.SendToDiscord("vpn", message)
	if err != nil {
		return err
	}

	return nil
}

// Start starts the VPN connection using an OpenVPN configuration file.
func Start(configPath string) (string, error) {
	fmt.Println("VPN is starting...")
	pidFile := config.BaseDirectory + "/lab-vpn.pid"
	cmd := exec.Command("sudo", "openvpn", "--config", configPath, "--writepid", pidFile)

	err := cmd.Start()
	if err != nil {
		return "", nil
	}

	fmt.Println("Wait 20s and check if the vpn is started...")
	time.Sleep(20 * time.Second)
	isActive, err := Status()
	if err != nil {
		return "", err
	}
	if isActive {
		fmt.Println("The VPN is now active !")
	}

	return "", nil
}

// Status checks the current status of the VPN connection.
func Status() (bool, error) {
	url := fmt.Sprintf("%s/connection/status", config.BaseHackTheBoxAPIURL)
	resp, err := utils.HtbRequest(http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var result []map[string]interface{}
	jsonBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("Error reading response body: %s", err)
	}

	err = json.Unmarshal(jsonBody, &result)
	if err != nil {
		return false, fmt.Errorf("Error unmarshalling JSON: %s", err)
	}

	if len(result) == 0 {
		return false, nil
	}

	for _, item := range result {
		connectionData, ok := item["connection"].(map[string]interface{})
		if !ok {
			return false, errors.New("Error asserting connection data")
		}

		name := connectionData["name"]
		ip4 := connectionData["ip4"]

		fmt.Printf("Name: %s, IP4: %s\n", name, ip4)
	}
	return true, nil
}

// Stop attempts to stop the currently active VPN connection.
func Stop() error {
	fmt.Println("Try to stop the active VPN...")
	pidFile := config.BaseDirectory + "/lab-vpn.pid"
	pidData, err := os.ReadFile(pidFile)
	if err != nil {
		return err
	}
	config.GlobalConfig.Logger.Debug(fmt.Sprintf("VPN PID : %v", pidData))

	cmd := exec.Command("sudo", "kill", string(pidData))

	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
