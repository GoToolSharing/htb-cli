package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Machine struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

type Root struct {
	Machines interface{} `json:"machines"`
}

func SearchMachineIDByName(machine_name string, proxyURL string) string {
	url := "https://www.hackthebox.com/api/v4/search/fetch?query=" + machine_name
	resp, err := HtbRequest(http.MethodGet, url, proxyURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	json_body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

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
	json_body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var result map[string]interface{}
	json.Unmarshal([]byte(json_body), &result)
	return result[key]
}

func GetMachineType(machine_id interface{}) string {
	machine_id = fmt.Sprintf("%v", machine_id)
	url := "https://www.hackthebox.com/api/v4/machine/recommended/"
	resp := HtbGet(url)
	card1 := ParseJsonMessage(resp, "card1").(map[string]interface{})
	fmachine_id, _ := strconv.ParseFloat(machine_id.(string), 64)
	if card1["id"].(float64) == fmachine_id {
		return "release"
	}

	url = "https://www.hackthebox.com/api/v4/machine/profile/" + machine_id.(string)
	resp = HtbGet(url)
	info := ParseJsonMessage(resp, "info").(map[string]interface{})
	if info["active"].(float64) == 1 {
		return "active"
	} else if info["retired"].(float64) == 1 {
		return "retired"
	}
	return "error"
}

func GetUserSubscription() string {
	url := "https://www.hackthebox.com/api/v4/user/info"
	resp := HtbGet(url)
	info := ParseJsonMessage(resp, "info").(map[string]interface{})
	if info["canAccessVIP"].(bool) {
		return "vip"
	}
	return "free"
}

func GetActiveMachineID() interface{} {
	url := "https://www.hackthebox.com/api/v4/machine/active"
	resp := HtbGet(url)
	info := ParseJsonMessage(resp, "info")
	if info == nil {
		url = "https://www.hackthebox.com/api/v4/release_arena/active"
		resp = HtbGet(url)
		info = ParseJsonMessage(resp, "info")
		if info == nil {
			fmt.Println("No machine is running")
			os.Exit(1)
		}
		return info.(map[string]interface{})["id"]
	} else {
		return info.(map[string]interface{})["id"]
	}
}

func GetActiveMachineName(machine_id interface{}) interface{} {
	machine_id = fmt.Sprintf("%v", machine_id)
	url := "https://www.hackthebox.com/api/v4/machine/profile/" + machine_id.(string)
	resp := HtbGet(url)
	info := ParseJsonMessage(resp, "info").(map[string]interface{})
	return info["name"]
}
