package sherlocks

import (
	"fmt"
	"log"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/rivo/tview"
)

const (
	SherlocksURL            = config.BaseHackTheBoxAPIURL + "/sherlocks?state=active"
	RetiredSherlocksURL     = config.BaseHackTheBoxAPIURL + "/sherlocks?state=retired"
	ScheduledSherlocksURL   = config.BaseHackTheBoxAPIURL + "/sherlocks?state=unreleased"
	ActiveSherlocksTitle    = "Active"
	RetiredSherlocksTitle   = "Retired"
	ScheduledSherlocksTitle = "Scheduled"
	SherlocksCheckMark      = "\U00002705"
	SherlocksCrossMark      = "\U0000274C"
	SPenguin                = "\U0001F427"
	SComputer               = "\U0001F5A5 "
)

// GetColorFromDifficultyText returns the color corresponding to the given difficulty.
func GetColorFromDifficultyText(difficultyText string) string {
	switch difficultyText {
	case "Medium":
		return "[orange]"
	case "Easy":
		return "[green]"
	case "Hard":
		return "[red]"
	case "Insane":
		return "[purple]"
	default:
		return "[-]"
	}
}

// CreateFlex creates and returns a Flex view with machine information
func CreateFlex(info interface{}, title string, isScheduled bool) (*tview.Flex, error) {
	log.Println("Info :", info)
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.SetBorder(true).SetTitle(title).SetTitleAlign(tview.AlignLeft)

	for _, value := range info.([]interface{}) {
		data := value.(map[string]interface{})

		// Determining the color according to difficulty

		key := "Undefined"
		if title == "Scheduled" {
			key = data["difficulty"].(string)
		}
		color := GetColorFromDifficultyText(key)

		var formatString string

		// Choice of display format depending on the nature of the information
		if isScheduled {
			formatString = fmt.Sprintf("%-15s %s%-10s[-]",
				data["name"], color, data["difficulty"])
		}
		// else {

		// Convert and format date
		// parsedDate, err := time.Parse(time.RFC3339Nano, data["release"].(string))
		// if err != nil {
		// 	return nil, fmt.Errorf("error parsing date: %v", err)
		// }
		// formattedDate := parsedDate.Format("02 January 2006")

		// userEmoji := SherlocksCrossMark + "User"
		// if value, ok := data["authUserInUserOwns"]; ok && value != nil {
		// 	if value.(bool) {
		// 		userEmoji = SherlocksCheckMark + "User"
		// 	}
		// }

		// rootEmoji := SherlocksCrossMark + "Root"
		// if value, ok := data["authUserInRootOwns"]; ok && value != nil {
		// 	if value.(bool) {
		// 		rootEmoji = SherlocksCheckMark + "Root"
		// 	}
		// }

		// formatString = fmt.Sprintf("%-15s %s%-10s[-] %-5v %-5v %-7v %-30s",
		// 	data["name"], color, data["difficultyText"],
		// 	data["star"], userEmoji, rootEmoji, formattedDate)
		// }

		flex.AddItem(tview.NewTextView().SetText(formatString).SetDynamicColors(true), 1, 0, false)
	}

	return flex, nil
}
