package utils

import (
	"fmt"
	"net/http"

	"github.com/GoToolSharing/htb-cli/config"
)

func GetAchievementLink(machineID int) (string, error) {
	resp, err := HtbRequest(http.MethodGet, fmt.Sprintf("%s/user/info", config.BaseHackTheBoxAPIURL), nil)
	if err != nil {
		return "", err
	}
	info := ParseJsonMessage(resp, "info")
	infoMap, _ := info.(map[string]interface{})
	config.GlobalConfig.Logger.Debug(fmt.Sprintf("User ID: %v", infoMap["id"]))
	config.GlobalConfig.Logger.Debug(fmt.Sprintf("Machine ID: %d", machineID))

	resp, err = HtbRequest(http.MethodGet, fmt.Sprintf("%s/user/achievement/machine/%v/%d", config.BaseHackTheBoxAPIURL, infoMap["id"], machineID), nil)
	if err != nil {
		return "", err
	}
	_, ok := ParseJsonMessage(resp, "message").(string)
	if !ok {
		return fmt.Sprintf("\nAchievement link: https://labs.hackthebox.com/achievement/machine/%v/%d", infoMap["id"], machineID), nil
	}
	return "", nil

}
