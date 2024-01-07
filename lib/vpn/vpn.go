package vpn

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/lib/utils"
)

func downloadVPN(url string) error {
	resp, err := utils.HtbRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		parts := strings.Split(url, "=")
		productValue := parts[len(parts)-1]
		fmt.Println("You do not have permissions to download the following vpn:", productValue)
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
	if strings.Contains(url, "product=labs") {
		parts := strings.Split(vpnName, "_")

		if len(parts) > 1 {
			parts[1] = "Labs"
		}

		vpnName = strings.Join(parts, "_")
	} else if strings.Contains(url, "product=competitive") {
		parts := strings.Split(vpnName, "_")

		if len(parts) > 1 {
			parts[1] = "Release_Arena"
		}

		vpnName = strings.Join(parts, "_")
	}
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

	return nil
}

// Start starts the VPN connection using an OpenVPN configuration file.
func Start(configPath string) (string, error) {
	config.GlobalConfig.Logger.Debug(fmt.Sprintf("VPN config file : %s", configPath))
	files, err := filepath.Glob(configPath)
	if err != nil {
		return "", fmt.Errorf("search error : %v", err)
	}
	if len(files) == 0 {
		isConfirmed := utils.AskConfirmation("VPN was not found. Would you like to download it ?")
		if isConfirmed {
			err := DownloadAll()
			if err != nil {
				return "", err
			}
		}
	}
	config.GlobalConfig.Logger.Info("VPN is starting...")
	cmd := "pgrep -fa openvpn"
	hacktheboxFound := false
	processes, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		hacktheboxFound = false
	}

	lines := strings.Split(string(processes), "\n")

	config.GlobalConfig.Logger.Debug(fmt.Sprintf("VPN processes: %v", lines))

	uniquePaths := make(map[string]bool)

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.Fields(line)
		processPath := parts[len(parts)-1]

		if _, found := uniquePaths[processPath]; found {
			continue
		}

		uniquePaths[processPath] = true

		file, err := os.Open(processPath)
		if err != nil {
			return "", fmt.Errorf("error reading file %s: %v", processPath, err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if strings.HasPrefix(scanner.Text(), "remote") && strings.Contains(scanner.Text(), "hackthebox.eu") {
				hacktheboxFound = true
				break
			}
		}

		if err := scanner.Err(); err != nil {
			return "", fmt.Errorf("error reading file %s: %v", processPath, err)
		}
	}

	// If no HackTheBox VPN found, start the VPN.
	if !hacktheboxFound {
		cmd = fmt.Sprintf("sudo openvpn %s", configPath)
		vpnProcess := exec.Command("sh", "-c", cmd)
		stdout, _ := vpnProcess.StdoutPipe()
		err := vpnProcess.Start()
		if err != nil {
			return "", err
		}

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "Initialization Sequence Completed") {
				fmt.Println("VPN Started Successfully!")
				return "VPN Started Successfully!", nil
			}
		}

		if err := scanner.Err(); err != nil {
			return "", fmt.Errorf("error reading output from VPN command: %v", err)
		}
		return "", fmt.Errorf("VPN did not start successfully. Initialization Sequence not completed.")
	}

	fmt.Println("VPN shutdown in progress...")
	_, err = Stop()
	if err != nil {
		return "", err
	}
	_, err = Start(configPath)
	if err != nil {
		return "", err
	}
	return "HackTheBox VPN is already running.", nil
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
		return false, fmt.Errorf("error reading response body: %s", err)
	}

	err = json.Unmarshal(jsonBody, &result)
	if err != nil {
		return false, fmt.Errorf("error unmarshalling JSON: %s", err)
	}

	if len(result) == 0 {
		return false, nil
	}

	for _, item := range result {
		connectionData, ok := item["connection"].(map[string]interface{})
		if !ok {
			return false, errors.New("error asserting connection data")
		}

		name := connectionData["name"]
		ip4 := connectionData["ip4"]

		fmt.Printf("Name: %s, IP4: %s\n", name, ip4)
	}
	return true, nil
}

// Stop terminates any active HackTheBox OpenVPN connections.
func Stop() (string, error) {
	fmt.Println("Stopping VPN if any HackTheBox connection is found...")
	cmd := "pgrep -fa openvpn"
	processes, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return "No vpn connection is active", nil
	}

	lines := strings.Split(string(processes), "\n")

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.Fields(line)
		processID := parts[0]
		configPath := parts[len(parts)-1]

		file, err := os.Open(configPath)
		if err != nil {
			continue
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if strings.HasPrefix(scanner.Text(), "remote") && strings.Contains(scanner.Text(), "hackthebox.eu") {
				killCmd := fmt.Sprintf("sudo kill %s", processID)
				if _, err := exec.Command("sh", "-c", killCmd).Output(); err != nil {
					return "", nil
				}
				fmt.Printf("Killed HackTheBox VPN process %s\n", processID)
			}
		}
		return "Completed checking and stopping HackTheBox VPN processes.", nil

	}
	return "", nil
}

func getVPNConfiguration(url string) error {
	parts := strings.Split(url, "=")
	productValue := parts[len(parts)-1]
	resp, err := utils.HtbRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		fmt.Printf("Product: %s\nNo information available\n\n", productValue)
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

	fmt.Printf("Product: %s\nID: %d\nFriendly Name: %s\nCurrent Clients: %d\nLocation: %s\n\n", productValue, response.Data.Assigned.ID, response.Data.Assigned.FriendlyName, response.Data.Assigned.CurrentClients, response.Data.Assigned.LocationFriendly)
	return nil
}

func List() error {
	config.GlobalConfig.Logger.Info("Recovering VPN configurations")
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
			err := getVPNConfiguration(url)
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

	return nil
}
