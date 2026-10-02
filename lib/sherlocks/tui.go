package sherlocks

import (
	"fmt"
	"strconv"
	"time"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/lib/utils"
	"github.com/rivo/tview"
)

const (
	SherlocksURL            = config.BaseHackTheBoxAPIURL + "/sherlocks?state=active&sort_by=release_date&sort_type=desc&per_page=20"
	RetiredSherlocksURL     = config.BaseHackTheBoxAPIURL + "/sherlocks?state=retired&sort_by=release_date&sort_type=desc&per_page=20"
	ScheduledSherlocksURL   = config.BaseHackTheBoxAPIURL + "/sherlocks?state=unreleased&sort_by=release_date&sort_type=desc&per_page=20"
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

func getStringField(data map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		value, ok := data[key]
		if !ok || value == nil {
			continue
		}

		switch typed := value.(type) {
		case string:
			if typed != "" {
				return typed
			}
		case fmt.Stringer:
			text := typed.String()
			if text != "" {
				return text
			}
		}
	}

	return ""
}

func getIntField(data map[string]interface{}, keys ...string) (int, bool) {
	for _, key := range keys {
		value, ok := data[key]
		if !ok || value == nil {
			continue
		}

		switch typed := value.(type) {
		case int:
			return typed, true
		case int32:
			return int(typed), true
		case int64:
			return int(typed), true
		case float32:
			return int(typed), true
		case float64:
			return int(typed), true
		case string:
			parsed, err := strconv.Atoi(typed)
			if err == nil {
				return parsed, true
			}
		}
	}

	return 0, false
}

func getProgressLabel(data map[string]interface{}) string {
	progress, ok := getIntField(data, "progress")
	if ok {
		return fmt.Sprintf("%d%%", progress)
	}

	return "Unknown"
}

func getReleaseDateLabel(data map[string]interface{}) string {
	releaseDate := getStringField(data, "release_date")
	if releaseDate == "" {
		return "Unknown"
	}

	parsedDate, err := time.Parse(time.RFC3339Nano, releaseDate)
	if err != nil {
		parsedDate, err = time.Parse(time.RFC3339, releaseDate)
		if err != nil {
			return releaseDate
		}
	}

	return parsedDate.Format("02 January 2006")
}

func fitColumn(value string, width int) string {
	if width <= 0 {
		return ""
	}

	trimmed := value
	if len(trimmed) > width {
		if width <= 3 {
			trimmed = utils.TruncateString(trimmed, width)
		} else {
			trimmed = utils.TruncateString(trimmed, width-3) + "..."
		}
	}

	return fmt.Sprintf("%-*s", width, trimmed)
}

// CreateFlex creates and returns a Flex view with machine information
func CreateFlex(info interface{}, title string, isScheduled bool) (*tview.Flex, error) {
	config.GlobalConfig.Logger.Debug(fmt.Sprintf("Info: %v", info))
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.SetBorder(true).SetTitle(title).SetTitleAlign(tview.AlignLeft)

	for _, value := range info.([]interface{}) {
		data := value.(map[string]interface{})

		difficulty := getStringField(data, "difficulty", "difficultyText")
		if difficulty == "" {
			difficulty = "Undefined"
		}
		color := GetColorFromDifficultyText(difficulty)

		name := getStringField(data, "name")
		if name == "" {
			name = "Undefined"
		}

		category := getStringField(data, "category_name")
		if category == "" {
			category = "Unknown"
		}

		progress := getProgressLabel(data)
		releaseDate := getReleaseDateLabel(data)
		nameColumn := fitColumn(name, 22)
		categoryColumn := fitColumn(category, 16)
		difficultyColumn := fitColumn(difficulty, 12)
		progressColumn := fitColumn(progress, 10)
		releaseDateColumn := fitColumn(releaseDate, 18)

		formatString := fmt.Sprintf("%s %s %s%s[-] %s %s",
			nameColumn, categoryColumn, color, difficultyColumn, progressColumn, releaseDateColumn)

		if isScheduled {
			formatString = fmt.Sprintf("%s %s %s%s[-] %s %s",
				nameColumn, categoryColumn, color, difficultyColumn, fitColumn("Unreleased", 10), releaseDateColumn)
		}

		flex.AddItem(tview.NewTextView().SetText(formatString).SetDynamicColors(true), 1, 0, false)
	}

	return flex, nil
}
