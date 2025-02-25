package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"os/user"
	"strings"
	"sync"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/GoToolSharing/htb-cli/config"
	"github.com/briandowns/spinner"
	"github.com/sahilm/fuzzy"
)

// SetTabWriterHeader will display the information in an array
func SetTabWriterHeader(header string) *tabwriter.Writer {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)
	fmt.Fprintln(w, header)
	return w
}

// SetTabWriterData will write the contents of each array cell
func SetTabWriterData(w *tabwriter.Writer, data string) {
	fmt.Fprint(w, data)
}

// AskConfirmation will request confirmation from the user
func AskConfirmation(message string) bool {
	if !config.GlobalConfig.BatchParam {
		var confirmation bool
		prompt := &survey.Confirm{
			Message: message,
		}
		if err := survey.AskOne(prompt, &confirmation); err != nil {
			return false
		}
		return confirmation
	}
	return true
}

// GetHTBToken checks whether the HTB_TOKEN environment variable exists
func GetHTBToken() (string, error) {
	var envName = "HTB_TOKEN"
	var htbToken = os.Getenv(envName)

	if htbToken == "" {
		return "", fmt.Errorf("environment variable is not set : %v\n", envName)
	}

	parts := strings.Split(htbToken, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("the %s variable must be an app token : https://app.hackthebox.com/profile/settings", envName)
	}

	return htbToken, nil
}

// SearchItemIDByName will return the id of an item (machine / challenge / user) based on its name
func SearchItemIDByName(item string, element_type string) (int, error) {
	url := fmt.Sprintf("%s/search/fetch?query=%s", config.BaseHackTheBoxAPIURL, item)
	resp, err := HtbRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	json_body, _ := io.ReadAll(resp.Body)

	var root Root
	err = json.Unmarshal([]byte(json_body), &root)
	if err != nil {
		fmt.Println("error:", err)
	}

	if element_type == "Machine" {
		switch root.Machines.(type) {
		case []interface{}:
			// Checking if machines array is empty
			if len(root.Machines.([]interface{})) == 0 {
				fmt.Println("No machine was found")
				os.Exit(0)
			}
			var machines []Machine
			machineData, _ := json.Marshal(root.Machines)
			err := json.Unmarshal(machineData, &machines)
			if err != nil {
				fmt.Println("error:", err)
			}
			config.GlobalConfig.Logger.Info("Machine found")
			config.GlobalConfig.Logger.Debug(fmt.Sprintf("Machine name: %s", machines[0].Value))
			isConfirmed := AskConfirmation("The following machine was found : " + machines[0].Value)
			if isConfirmed {
				return machines[0].ID, nil
			}
			os.Exit(0)
		case map[string]interface{}:
			// Checking if machines array is empty
			if len(root.Machines.(map[string]interface{})) == 0 {
				fmt.Println("No machine was found")
				os.Exit(0)
			}
			var machines map[string]Machine
			machineData, _ := json.Marshal(root.Machines)
			err := json.Unmarshal(machineData, &machines)
			if err != nil {
				fmt.Println("error:", err)
			}
			config.GlobalConfig.Logger.Info("Machine found")
			config.GlobalConfig.Logger.Debug(fmt.Sprintf("Machine name: %s", machines["0"].Value))
			isConfirmed := AskConfirmation("The following machine was found : " + machines["0"].Value)
			if isConfirmed {
				return machines["0"].ID, nil
			}
			os.Exit(0)
		default:
			fmt.Println("No machine was found")
			os.Exit(0)
		}
	} else if element_type == "Challenge" {
		switch root.Challenges.(type) {
		case []interface{}:
			// Checking if challenges array is empty
			if len(root.Challenges.([]interface{})) == 0 {
				fmt.Println("No challenge was found")
				os.Exit(0)
			}
			var challenges []Challenge
			challengeData, _ := json.Marshal(root.Challenges)
			err := json.Unmarshal(challengeData, &challenges)
			if err != nil {
				fmt.Println("error:", err)
			}
			config.GlobalConfig.Logger.Info("Challenge found")
			config.GlobalConfig.Logger.Debug(fmt.Sprintf("Challenge name: %s", challenges[0].Value))
			isConfirmed := AskConfirmation("The following challenge was found : " + challenges[0].Value)
			if isConfirmed {
				return challenges[0].ID, nil
			}
			os.Exit(0)
		case map[string]interface{}:
			// Checking if challenges array is empty
			if len(root.Challenges.(map[string]interface{})) == 0 {
				fmt.Println("No challenge was found")
				os.Exit(0)
			}
			var challenges map[string]Challenge
			challengeData, _ := json.Marshal(root.Challenges)
			err := json.Unmarshal(challengeData, &challenges)
			if err != nil {
				fmt.Println("error:", err)
			}
			config.GlobalConfig.Logger.Info("Challenge found")
			config.GlobalConfig.Logger.Debug(fmt.Sprintf("Challenge name: %s", challenges["0"].Value))
			isConfirmed := AskConfirmation("The following challenge was found : " + challenges["0"].Value)
			if isConfirmed {
				return challenges["0"].ID, nil
			}
			os.Exit(0)
		default:
			fmt.Println("No challenge was found")
			os.Exit(0)
		}
	} else if element_type == "Username" {
		switch root.Usernames.(type) {
		case []interface{}:
			// Checking if usernames array is empty
			if len(root.Usernames.([]interface{})) == 0 {
				fmt.Println("No username was found")
				return 0, fmt.Errorf("error: No username was found")
			}
			var usernames []Username
			usernameData, _ := json.Marshal(root.Usernames)
			err := json.Unmarshal(usernameData, &usernames)
			if err != nil {
				fmt.Println("error:", err)
			}
			config.GlobalConfig.Logger.Info("Username found")
			config.GlobalConfig.Logger.Debug(fmt.Sprintf("Username value: %s", usernames[0].Value))
			isConfirmed := AskConfirmation("The following username was found : " + usernames[0].Value)
			if isConfirmed {
				return usernames[0].ID, nil
			}
			os.Exit(0)
		case map[string]interface{}:
			// Checking if usernames array is empty
			if len(root.Usernames.(map[string]interface{})) == 0 {
				fmt.Println("No username was found")
				return 0, fmt.Errorf("error: No username was found")
			}
			var usernames map[string]Username
			usernameData, _ := json.Marshal(root.Usernames)
			err := json.Unmarshal(usernameData, &usernames)
			if err != nil {
				fmt.Println("error:", err)
			}
			config.GlobalConfig.Logger.Info("Username found")
			config.GlobalConfig.Logger.Debug(fmt.Sprintf("Username value: %s", usernames["0"].Value))
			isConfirmed := AskConfirmation("The following username was found : " + usernames["0"].Value)
			if isConfirmed {
				return usernames["0"].ID, nil
			}
			os.Exit(0)
		default:
			fmt.Println("No username found")
		}
	} else {
		return 0, errors.New("bad element_type")
	}

	// The HackTheBox API can return either a slice or a map
	return 0, nil
}

// ParseJsonMessage will parse the result of the API request into a JSON
func ParseJsonMessage(resp *http.Response, key string) interface{} {
	json_body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	err := json.Unmarshal([]byte(json_body), &result)
	if err != nil {
		fmt.Println("error:", err)
	}
	return result[key]
}

// GetMachineType will return the machine type
func GetMachineType(machine_id int) (string, error) {
	// Check if the machine is the latest release
	url := fmt.Sprintf("%s/machine/recommended/", config.BaseHackTheBoxAPIURL)
	resp, err := HtbRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	card := ParseJsonMessage(resp, "card1").(map[string]interface{})
	if card["id"] == machine_id {
		return "release", nil
	}

	// Check if the machine is active or retired
	url = fmt.Sprintf("%s/machine/profile/%d", config.BaseHackTheBoxAPIURL, machine_id)
	resp, err = HtbRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	info := ParseJsonMessage(resp, "info").(map[string]interface{})
	if info["active"] == true {
		return "active", nil
	} else if info["retired"] == true {
		return "retired", nil
	}
	return "", errors.New("error: machine type not found")
}

// GetUserSubscription returns the user's subscription level
func GetUserSubscription() (string, error) {
	url := fmt.Sprintf("%s/user/info", config.BaseHackTheBoxAPIURL)
	resp, err := HtbRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	info := ParseJsonMessage(resp, "info").(map[string]interface{})
	canAccessVIP := info["canAccessVIP"].(bool)
	isDedicatedVIP := info["isDedicatedVip"].(bool)

	if canAccessVIP {
		if isDedicatedVIP {
			return "vip+", nil
		}
		return "vip", nil
	}

	return "free", nil
}

// GetActiveMachineID returns the id of the active machine
func GetActiveMachineID() (int, error) {
	url := fmt.Sprintf("%s/machine/active", config.BaseHackTheBoxAPIURL)
	resp, err := HtbRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	info := ParseJsonMessage(resp, "info")
	if info == nil {
		return 0, err
	}
	idFloat, ok := info.(map[string]interface{})["id"].(float64)
	if !ok {
		return 0, fmt.Errorf("unable to parse machine ID")
	}

	return int(idFloat), nil
}

// GetExpiredTime returns the expired date of the machine
func GetExpiredTime(machineType string) (string, error) {
	if machineType == "release" {
		return getReleaseArenaExpiredTime()
	}
	return getActiveExpiredTime()
}

// getActiveExpiredTime returns the expired date of the active machine
func getActiveExpiredTime() (string, error) {
	url := fmt.Sprintf("%s/machine/active", config.BaseHackTheBoxAPIURL)
	resp, err := HtbRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	info := ParseJsonMessage(resp, "info")
	if info == nil {
		return "Undefined", nil
	}
	data := info.(map[string]interface{})
	expiresAt := data["expires_at"]
	if expiresAt == nil {
		return "Undefined", nil
	}
	return expiresAt.(string), nil
}

// getReleaseArenaExpiredTime returns the expired date of the release arena machine
func getReleaseArenaExpiredTime() (string, error) {
	url := fmt.Sprintf("%s/season/machine/active", config.BaseHackTheBoxAPIURL)
	resp, err := HtbRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	info := ParseJsonMessage(resp, "data")
	if info == nil {
		return "Undefined", nil
	}
	data := info.(map[string]interface{})
	playInfo := data["play_info"].(map[string]interface{})
	expiresAt := playInfo["expires_at"]
	if expiresAt == nil {
		return "Undefined", nil
	}
	return expiresAt.(string), nil
}

// GetActiveMachineIP returns the ip of the active machine
func GetActiveMachineIP() (string, error) {
	url := fmt.Sprintf("%s/machine/active", config.BaseHackTheBoxAPIURL)
	resp, err := HtbRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	info := ParseJsonMessage(resp, "info")
	if info == nil {
		return "", err
	}
	config.GlobalConfig.Logger.Debug(fmt.Sprintf("Active machine informations: %v", info))
	if ipValue, ok := info.(map[string]interface{})["ip"].(string); ok {
		return ipValue, nil
	}
	return "Undefined", nil
}

// GetActiveReleaseArenaMachineIP returns the ip of the active release arena machine
func GetActiveReleaseArenaMachineIP() (string, error) {
	url := fmt.Sprintf("%s/season/machine/active", config.BaseHackTheBoxAPIURL)
	resp, err := HtbRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	data := ParseJsonMessage(resp, "data")
	if data == nil {
		return "", err
	}
	config.GlobalConfig.Logger.Debug(fmt.Sprintf("Relase arena active machine informations: %v", data))

	return fmt.Sprintf("%v", data.(map[string]interface{})["ip"].(string)), nil
}

// HtbRequest makes an HTTP request to the Hackthebox API
func HtbRequest(method string, urlParam string, jsonData []byte) (*http.Response, error) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		s.Stop()
		os.Exit(0)
	}()

	s.Start()
	JWT_TOKEN, err := GetHTBToken()
	if err != nil {
		s.Stop()
		return nil, err
	}

	req, err := http.NewRequest(method, urlParam, bytes.NewBuffer(jsonData))
	if err != nil {
		s.Stop()
		return nil, err
	}

	req.Header.Set("User-Agent", "htb-cli")
	req.Header.Set("Authorization", "Bearer "+JWT_TOKEN)

	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json, text/plain, */*")
	} else if method == http.MethodGet {
		req.Header.Set("Host", config.HostHackTheBox)
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	if config.GlobalConfig.ProxyParam != "" {
		config.GlobalConfig.Logger.Info("Proxy URL found")
		config.GlobalConfig.Logger.Debug(fmt.Sprintf("Proxy value : %s", config.GlobalConfig.ProxyParam))
		proxyURLParsed, err := url.Parse(config.GlobalConfig.ProxyParam)
		if err != nil {
			s.Stop()
			return nil, fmt.Errorf("error parsing proxy url : %v", err)
		}
		transport.Proxy = http.ProxyURL(proxyURLParsed)
	}

	config.GlobalConfig.Logger.Info("Sending an HTTP HTB request")
	config.GlobalConfig.Logger.Debug(fmt.Sprintf("Request URL: %v", req.URL))
	config.GlobalConfig.Logger.Debug(fmt.Sprintf("Request method: %v", req.Method))
	config.GlobalConfig.Logger.Debug(fmt.Sprintf("Request body: %v", req.Body))

	client := &http.Client{Transport: transport, CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	resp.Body = io.NopCloser(bytes.NewReader(body))

	// Check if token is invalid or expired
	if resp.StatusCode == 302 && strings.Contains(resp.Header.Get("Location"), "/login") {
		s.Stop()
		return nil, fmt.Errorf("HTB Token appears invalid or expired")
	}
	s.Stop()
	return resp, nil
}

func TruncateString(str string, maxLength int) string {
	if len(str) > maxLength {
		return str[:maxLength]
	}
	return str
}

func GetInformationsFromActiveMachine() (map[string]interface{}, error) {
	machineID, err := GetActiveMachineID()
	if err != nil {
		return nil, err
	}
	if machineID == 0 {
		fmt.Println("No machine is running")
		return nil, nil
	}
	config.GlobalConfig.Logger.Debug(fmt.Sprintf("Machine ID: %d", machineID))

	url := fmt.Sprintf("%s/machine/profile/%d", config.BaseHackTheBoxAPIURL, machineID)
	resp, err := HtbRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	info := ParseJsonMessage(resp, "info")

	data := info.(map[string]interface{})

	return data, nil
}

// HTTPRequest makes an HTTP request with the specified method, URL, proxy settings, and data.
func HTTPRequest(method string, urlParam string, jsonData []byte) (*http.Response, error) {
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
		return nil, err
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

	if config.GlobalConfig.ProxyParam != "" {
		config.GlobalConfig.Logger.Info("Proxy URL found")
		config.GlobalConfig.Logger.Debug(fmt.Sprintf("Proxy value : %s", config.GlobalConfig.ProxyParam))
		proxyURLParsed, err := url.Parse(config.GlobalConfig.ProxyParam)
		if err != nil {
			s.Stop()
			return nil, fmt.Errorf("error parsing proxy url : %v", err)
		}
		transport.Proxy = http.ProxyURL(proxyURLParsed)
	}

	config.GlobalConfig.Logger.Info("Sending an HTTP request")
	config.GlobalConfig.Logger.Debug(fmt.Sprintf("Request URL: %v", req.URL))
	config.GlobalConfig.Logger.Debug(fmt.Sprintf("Request method: %v", req.Method))
	config.GlobalConfig.Logger.Debug(fmt.Sprintf("Request body: %v", req.Body))

	client := &http.Client{Transport: transport}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	resp.Body = io.NopCloser(bytes.NewReader(body))
	s.Stop()
	return resp, nil
}

// GetCurrentUsername retrieves the current system user's name.
func GetCurrentUsername() string {
	user, _ := user.Current()
	return user.Username
}

func SearchLastReleaseArenaMachine() (int, error) {
	url := fmt.Sprintf("%s/season/machine/active", config.BaseHackTheBoxAPIURL)
	resp, err := HtbRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	info := ParseJsonMessage(resp, "data")
	if info == nil {
		return 0, err
	}
	config.GlobalConfig.Logger.Debug(fmt.Sprintf("Information on the last active machine: %v", info))
	machineID := int(info.(map[string]interface{})["id"].(float64))
	return machineID, nil
}

func extractNamesAndIDs(jsonData string) (map[string]int, error) {
	var response JsonResponse
	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		return nil, err
	}

	namesAndIDs := make(map[string]int)
	for _, item := range response.Data {
		namesAndIDs[item.Name] = item.ID
	}

	return namesAndIDs, nil
}

func SearchFortressID(partialName string) (int, error) {
	url := fmt.Sprintf("%s/fortresses", config.BaseHackTheBoxAPIURL)
	resp, err := HtbRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	jsonData, _ := io.ReadAll(resp.Body)
	namesAndIDs, err := extractNamesAndIDs(string(jsonData))
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return 0, nil
	}

	var names []string
	for name := range namesAndIDs {
		names = append(names, name)
	}

	matches := fuzzy.Find(partialName, names)

	for _, match := range matches {
		matchedName := names[match.Index]
		isConfirmed := AskConfirmation("The following fortress was found : " + matchedName)
		if isConfirmed {
			return namesAndIDs[matchedName], nil
		}
		os.Exit(0)
	}
	return 0, nil
}

func SearchProlabID(partialName string) (int, error) {
	url := fmt.Sprintf("%s/prolabs", config.BaseHackTheBoxAPIURL)
	resp, err := HtbRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	jsonData, _ := io.ReadAll(resp.Body)
	namesAndIDs, err := extractProlabsNamesAndIDs(string(jsonData))
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return 0, nil
	}

	var names []string
	for name := range namesAndIDs {
		names = append(names, name)
	}

	matches := fuzzy.Find(partialName, names)

	for _, match := range matches {
		matchedName := names[match.Index]
		isConfirmed := AskConfirmation("The following prolab was found : " + matchedName)
		if isConfirmed {
			return namesAndIDs[matchedName], nil
		}
		os.Exit(0)
	}
	return 0, nil
}

func extractProlabsNamesAndIDs(jsonData string) (map[string]int, error) {
	var response ProlabJsonResponse
	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		return nil, err
	}

	namesAndIDs := make(map[string]int)
	for _, lab := range response.Data.Labs {
		namesAndIDs[lab.Name] = lab.ID
	}

	return namesAndIDs, nil
}

// Fuzzy finder for challenges

func SearchChallengeByName(partialName string) (ChallengeFinder, error) {
	var wg sync.WaitGroup
	respChan := make(chan *http.Response, 2)
	var allChallenges []ChallengeFinder

	wg.Add(2)
	go sendRequest(fmt.Sprintf("%s/challenge/list", config.BaseHackTheBoxAPIURL), respChan, &wg)
	go sendRequest(fmt.Sprintf("%s/challenge/list/retired", config.BaseHackTheBoxAPIURL), respChan, &wg)

	go func() {
		wg.Wait()
		close(respChan)
	}()

	for resp := range respChan {
		challenges := processResponse(resp)
		allChallenges = append(allChallenges, challenges...)
	}

	challengeNames := make([]string, len(allChallenges))
	for i, ch := range allChallenges {
		challengeNames[i] = ch.Name
	}

	matches := fuzzy.Find(partialName, challengeNames)

	config.GlobalConfig.Logger.Debug(fmt.Sprintf("Challenges matches: %v", matches))

	for _, match := range matches {
		matchedName := challengeNames[match.Index]
		isConfirmed := AskConfirmation("The following challenge was found : " + matchedName)
		if isConfirmed {
			for _, ch := range allChallenges {
				if ch.Name == matchedName {
					return ch, nil
				}
			}
		}
	}

	return ChallengeFinder{}, fmt.Errorf("no matching challenge found")
}

func processResponse(resp *http.Response) []ChallengeFinder {
	defer resp.Body.Close()
	var cr ChallengeResponseFinder
	if err := json.NewDecoder(resp.Body).Decode(&cr); err != nil {
		fmt.Printf("JSON decoding error %v", err)
		os.Exit(1)
	}
	return cr.ChallengesFinder
}

func sendRequest(url string, respChan chan<- *http.Response, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := HtbRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("error sending request : %v", err)
		os.Exit(1)
	}
	respChan <- resp
}

func GetChallengeBlooder(challengeID string) (string, error) {
	url := fmt.Sprintf("%s/challenge/activity/%s", config.BaseHackTheBoxAPIURL, challengeID)
	resp, err := HtbRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	jsonData, _ := io.ReadAll(resp.Body)

	var data DataActivity
	err = json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		panic(err)
	}

	for _, activity := range data.Info.Activities {
		if activity.Type == "blood" {
			return activity.UserName, nil
		}
	}

	return "Not defined", nil
}
