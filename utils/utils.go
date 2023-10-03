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
	"text/tabwriter"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"
)

type Machine struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

type Challenge struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

type Root struct {
	Machines   interface{} `json:"machines"`
	Challenges interface{} `json:"challenges"`
}

// Function that removes the output from the console for unit tests.
// TODO: Move into a utils-tests.go file
func SetOutputTest() (*os.File, *os.File) {
	log.SetOutput(io.Discard)
	r, w, _ := os.Pipe()
	os.Stdout = w
	return r, w
}

// SetTabWriterHeader will display the information in an array
func SetTabWriterHeader(header string) *tabwriter.Writer {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)
	fmt.Fprintln(w, header)
	return w
}

// SetTabWriterData will write the contents of each array cell
func SetTabWriterData(w *tabwriter.Writer, data string) {
	fmt.Fprintf(w, data)
}

// AskConfirmation will request confirmation from the user
func AskConfirmation(message string) bool {
	var confirmation bool
	prompt := &survey.Confirm{
		Message: message,
	}
	if err := survey.AskOne(prompt, &confirmation); err != nil {
		return false
	}
	return confirmation
}

// GetHTBToken checks whether the HTB_TOKEN environment variable exists
func GetHTBToken() string {
	var envName = "HTB_TOKEN"
	if os.Getenv(envName) == "" {
		fmt.Printf("Environment variable is not set : %v\n", envName)
		return ""
	}
	return os.Getenv("HTB_TOKEN")
}

// SearchItemIDByName will return the id of an item (machine / challenge) based on its name
func SearchItemIDByName(item string, proxyURL string, element_type string) (string, error) {
	url := "https://www.hackthebox.com/api/v4/search/fetch?query=" + item
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

	if element_type == "Machine" {
		switch root.Machines.(type) {
		case []interface{}:
			// Checking if machines array is empty
			if len(root.Machines.([]interface{})) == 0 {
				fmt.Println("No machine was found")
				return "", fmt.Errorf("error: No machine was found")
			}
			var machines []Machine
			machineData, _ := json.Marshal(root.Machines)
			err := json.Unmarshal(machineData, &machines)
			if err != nil {
				fmt.Println("error:", err)
			}
			log.Println("Machine found :", machines[0].Value)
			var confirmation bool
			confirmation_message := "The following machine was found : " + machines[0].Value
			prompt := &survey.Confirm{
				Message: confirmation_message,
			}
			if err := survey.AskOne(prompt, &confirmation); err != nil {
				log.Fatal(err)
			}
			if !confirmation {
				log.Fatal("Canceled")
			}
			return machines[0].ID, nil
		case map[string]interface{}:
			// Checking if machines array is empty
			if len(root.Machines.(map[string]interface{})) == 0 {
				fmt.Println("No machine was found")
				return "", fmt.Errorf("error: No machine was found")
			}
			var machines map[string]Machine
			machineData, _ := json.Marshal(root.Machines)
			err := json.Unmarshal(machineData, &machines)
			if err != nil {
				fmt.Println("error:", err)
			}
			log.Println("Machine found :", machines["0"].Value)
			var confirmation bool
			confirmation_message := "The following machine was found : " + machines["0"].Value
			prompt := &survey.Confirm{
				Message: confirmation_message,
			}
			if err := survey.AskOne(prompt, &confirmation); err != nil {
				log.Fatal(err)
			}
			if !confirmation {
				log.Fatal("Canceled")
			}
			return machines["0"].ID, nil
		default:
			return "", fmt.Errorf("No machine was found")
		}
	} else if element_type == "Challenge" {
		switch root.Challenges.(type) {
		case []interface{}:
			// Checking if challenges array is empty
			if len(root.Challenges.([]interface{})) == 0 {
				fmt.Println("No challenge was found")
				return "", fmt.Errorf("error: No challenge was found")
			}
			var challenges []Challenge
			challengeData, _ := json.Marshal(root.Challenges)
			err := json.Unmarshal(challengeData, &challenges)
			if err != nil {
				fmt.Println("error:", err)
			}
			log.Println("Challenge found :", challenges[0].Value)
			var confirmation bool
			confirmation_message := "The following challenge was found : " + challenges[0].Value
			prompt := &survey.Confirm{
				Message: confirmation_message,
			}
			if err := survey.AskOne(prompt, &confirmation); err != nil {
				log.Fatal(err)
			}
			if !confirmation {
				log.Fatal("Canceled")
			}
			return challenges[0].ID, nil
		case map[string]interface{}:
			// Checking if challenges array is empty
			if len(root.Challenges.(map[string]interface{})) == 0 {
				fmt.Println("No challenge was found")
				return "", fmt.Errorf("error: No challenge was found")
			}
			var challenges map[string]Challenge
			challengeData, _ := json.Marshal(root.Challenges)
			err := json.Unmarshal(challengeData, &challenges)
			if err != nil {
				fmt.Println("error:", err)
			}
			log.Println("Challenge found :", challenges["0"].Value)
			var confirmation bool
			confirmation_message := "The following challenge was found : " + challenges["0"].Value
			prompt := &survey.Confirm{
				Message: confirmation_message,
			}
			if err := survey.AskOne(prompt, &confirmation); err != nil {
				log.Fatal(err)
			}
			if !confirmation {
				log.Fatal("Canceled")
			}
			return challenges["0"].ID, nil
		default:
			log.Fatal("No challenge found")
		}
	} else {
		log.Fatal("Bad element_type")
	}

	// The HackTheBox API can return either a slice or a map
	return "", nil
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

// GetUserSubscription returns the user's subscription level
func GetUserSubscription(proxyURL string) string {
	url := "https://www.hackthebox.com/api/v4/user/info"
	resp, err := HtbRequest(http.MethodGet, url, proxyURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	info := ParseJsonMessage(resp, "info").(map[string]interface{})
	canAccessVIP := info["canAccessVIP"].(bool)
	isDedicatedVIP := info["isDedicatedVip"].(bool)

	if canAccessVIP {
		if isDedicatedVIP {
			return "vip+"
		}
		return "vip"
	}

	return "free"
}

// GetActiveMachineID returns the id of the active machine
func GetActiveMachineID(proxyURL string) string {
	url := "https://www.hackthebox.com/api/v4/machine/active"
	resp, err := HtbRequest(http.MethodGet, url, proxyURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	info := ParseJsonMessage(resp, "info")
	if info == nil {
		return ""
	}
	return fmt.Sprintf("%.0f", info.(map[string]interface{})["id"].(float64))
}

// HtbRequest makes an HTTP request to the Hackthebox API
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
