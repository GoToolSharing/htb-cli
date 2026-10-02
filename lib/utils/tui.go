package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/rivo/tview"
)

const (
	maxFortressNameLength = 9
	maxProlabNameLength   = 11
	maxActivityNameLength = 11
)

// Calculates the spacing between keys and values to ensure alignment
func calculateSpacing(baseName string, maxNameLength int) string {
	return fmt.Sprintf("%-*s", maxNameLength-len(baseName)+1, "")
}

// Parse and return the user's subscription level
func parseUserSubscription(profile map[string]interface{}) string {
	isVip := false
	if val, ok := profile["isVip"].(bool); ok {
		isVip = val
	}

	isDedicatedVIP := false
	if val, ok := profile["isDedicatedVip"].(bool); ok {
		isDedicatedVIP = val
	}

	if isDedicatedVIP {
		return "VIP+"
	} else if isVip {
		return "VIP"
	}

	return "Free"
}

// Format output for display in tview
func formatFlagInfo(name string, ownedFlags, totalFlags float64, flagSymbol string, maxNameLength int) string {
	var color string
	if ownedFlags == totalFlags {
		color = "[green]"
	} else if ownedFlags == 0 {
		color = "[red]"
	} else {
		color = "[orange]"
	}

	spacing := calculateSpacing(name, maxNameLength)

	return fmt.Sprintf("[::b]%s %s%s: %s%.0f/%.0f[-]", flagSymbol, name, spacing, color, ownedFlags, totalFlags)
}

// Add items and display in tview
func displayInfoPanel(title string, items []interface{}, formatterFunc func(map[string]interface{}) string, paddingBottom int) *tview.Flex {
	panel := tview.NewFlex().SetDirection(tview.FlexRow)
	panel.SetBorder(true).SetTitle(title).SetTitleAlign(tview.AlignLeft)

	isScrollable := title == "Activity"
	textView := tview.NewTextView().SetDynamicColors(true).SetScrollable(isScrollable)

	for _, itemInterface := range items {
		item, ok := itemInterface.(map[string]interface{})
		if !ok {
			fmt.Fprintln(textView, "Error: couldn't convert item to a map[string]interface{}")
			continue
		}

		text := formatterFunc(item)
		fmt.Fprintln(textView, text)
	}

	panel.AddItem(textView, 0, 1, false)

	return panel
}

// Get the right keys for display
func displayInfo(dataMaps map[string]map[string]interface{}, dataMapKey string, title string, flagSymbol string, maxNameLength int, paddingBottom int) *tview.Flex {
	// Check if the capitalized key exists first
	capitalizedKey := strings.ToUpper(string(dataMapKey[0])) + dataMapKey[1:]
	dataMap, ok := dataMaps[capitalizedKey]
	if !ok {
		// Key doesn't exist in dataMaps
		return nil
	}

	// Check if the inner key exists
	itemsInterface, ok := dataMap[dataMapKey]
	if !ok {
		// Inner key doesn't exist
		return nil
	}

	// Try to convert to slice
	items, ok := itemsInterface.([]interface{})
	if !ok {
		fmt.Printf("Error: couldn't convert data for key '%s'\n", dataMapKey)
		return nil
	}

	var formatterFunc func(item map[string]interface{}) string
	if dataMapKey == "activity" {
		formatterFunc = func(item map[string]interface{}) string {
			// Handle v5 API format
			activityType := "unknown"
			if typeVal, ok := item["type"].(string); ok {
				activityType = typeVal
			}

			name := "N/A"
			if nameVal, ok := item["name"].(string); ok {
				name = nameVal
			}

			points := 0.0
			if pointsVal, ok := item["points"].(float64); ok {
				points = pointsVal
			}

			bloodStatus := ""
			if blood, ok := item["blood"].(bool); ok && blood {
				bloodStatus = " [red]🩸[-]"
			}

			// Parse and format ownDate
			dateStr := ""
			if ownDate, ok := item["ownDate"].(string); ok {
				if parsedDate, err := time.Parse(time.RFC3339, ownDate); err == nil {
					dateStr = parsedDate.Format("2006-01-02 15:04")
				} else {
					dateStr = ownDate
				}
			}

			return fmt.Sprintf("[::b]%s - %s - [green]+[%.0fpts][-] (%s)%s", activityType, name, points, dateStr, bloodStatus)
		}
	} else {
		formatterFunc = func(item map[string]interface{}) string {
			name := "N/A"
			if nameVal, ok := item["name"].(string); ok {
				name = nameVal
			}
			var owned, total float64 = 0, 0
			if ownedVal, ok := item["owned_flags"].(float64); ok {
				owned = ownedVal
			}
			if totalVal, ok := item["total_flags"].(float64); ok {
				total = totalVal
			}
			return formatFlagInfo(name, owned, total, flagSymbol, maxNameLength)
		}
	}

	return displayInfoPanel(title, items, formatterFunc, paddingBottom)
}

// Initializes and displays user information in tview
func DisplayInformationsGUI(profile map[string]interface{}, advancedLabsMap map[string]map[string]interface{}) {

	teamName, teamRank := "N/A", "N/A"
	universityName, universityRank := "N/A", "N/A"

	if teamMap, ok := profile["team"].(map[string]interface{}); ok && teamMap != nil {
		if name, ok := teamMap["name"].(string); ok {
			teamName = name
		}
		if rank, ok := teamMap["ranking"].(float64); ok {
			teamRank = fmt.Sprintf("%v", rank)
		}
	}

	if universityMap, ok := profile["university"].(map[string]interface{}); ok && universityMap != nil {
		if name, ok := universityMap["name"].(string); ok {
			universityName = name
		}
		if rank, ok := universityMap["rank"].(float64); ok {
			universityRank = fmt.Sprintf("%v", rank)
		}
	}

	subscription := parseUserSubscription(profile)
	rankRequirement := "100"
	if val, ok := profile["rank_requirement"].(float64); ok {
		rankRequirement = fmt.Sprintf("%.0f", val)
	}

	ranking := "N/A"
	if rank, ok := profile["ranking"].(float64); ok {
		ranking = fmt.Sprintf("%.0f", rank)
	}

	app := tview.NewApplication()

	userInformationsFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	userInformationsFlex.SetBorder(true).SetTitle("Profile").SetTitleAlign(tview.AlignLeft)
	id := "N/A"
	if idVal, ok := profile["id"].(float64); ok {
		id = fmt.Sprintf("%d", int(idVal))
	}
	userInformationsFlex.AddItem(tview.NewTextView().SetText(fmt.Sprintf("[::b]ID           : %v[-]", id)).SetDynamicColors(true), 1, 0, false)
	userInformationsFlex.AddItem(tview.NewTextView().SetText(fmt.Sprintf("[::b]Name         : %v[-]", profile["name"])).SetDynamicColors(true), 1, 0, false)
	userInformationsFlex.AddItem(tview.NewTextView().SetText(fmt.Sprintf("[::b]Team         : %v[-]", teamName)).SetDynamicColors(true), 1, 0, false)
	userInformationsFlex.AddItem(tview.NewTextView().SetText(fmt.Sprintf("[::b]University   : %v[-]", universityName)).SetDynamicColors(true), 1, 0, false)
	userInformationsFlex.AddItem(tview.NewTextView().SetText(fmt.Sprintf("[::b]Country      : %v[-]", profile["country_name"])).SetDynamicColors(true), 1, 0, false)
	userInformationsFlex.AddItem(tview.NewTextView().SetText(fmt.Sprintf("[::b]Subscription : %v[-]", subscription)).SetDynamicColors(true), 1, 0, false)

	userRankingInformationsFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	userRankingInformationsFlex.SetBorder(true).SetTitle("Ranking Informations").SetTitleAlign(tview.AlignLeft)
	userRankingInformationsFlex.AddItem(tview.NewTextView().SetText(fmt.Sprintf("[::b]Global     : %v[-]\U0001F3C6", ranking)).SetDynamicColors(true), 1, 0, false)
	userRankingInformationsFlex.AddItem(tview.NewTextView().SetText(fmt.Sprintf("[::b]Team       : %v[-]\U0001F91D", teamRank)).SetDynamicColors(true), 1, 0, false)
	userRankingInformationsFlex.AddItem(tview.NewTextView().SetText(fmt.Sprintf("[::b]University : %v[-]\U0001F3EB", universityRank)).SetDynamicColors(true), 1, 0, false)
	userRankingInformationsFlex.AddItem(tview.NewTextView().SetText(fmt.Sprintf("[::b]Points     : %v[-]\U0001F396", profile["points"])).SetDynamicColors(true), 1, 0, false)
	userRankingInformationsFlex.AddItem(tview.NewTextView().SetText(fmt.Sprintf("[::b]Rank       : %v[-]", profile["rank"])).SetDynamicColors(true), 1, 0, false)
	userRankingInformationsFlex.AddItem(tview.NewTextView().SetText(fmt.Sprintf("[::b]Ownership  : %v%% / %v%%[-]", profile["rank_ownership"], rankRequirement)).SetDynamicColors(true), 1, 0, false)

	userMiscInformationsFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	userMiscInformationsFlex.SetBorder(true).SetTitle("Misc").SetTitleAlign(tview.AlignLeft)
	userMiscInformationsFlex.AddItem(tview.NewTextView().SetText(fmt.Sprintf("[::b]User Bloods   : %v\U0001FA78[-]", profile["user_bloods"])).SetDynamicColors(true), 1, 0, false)
	userMiscInformationsFlex.AddItem(tview.NewTextView().SetText(fmt.Sprintf("[::b]System Bloods : %v\U0001FA78[-]", profile["system_bloods"])).SetDynamicColors(true), 1, 0, false)
	userMiscInformationsFlex.AddItem(tview.NewTextView().SetText(fmt.Sprintf("[::b]User Owns     : %v[-]", profile["user_owns"])).SetDynamicColors(true), 1, 0, false)
	userMiscInformationsFlex.AddItem(tview.NewTextView().SetText(fmt.Sprintf("[::b]System Owns   : %v[-]", profile["system_owns"])).SetDynamicColors(true), 1, 0, false)
	userMiscInformationsFlex.AddItem(tview.NewTextView().SetText(fmt.Sprintf("[::b]Respects      : %v[-]", profile["respects"])).SetDynamicColors(true), 1, 0, false)
	userMiscInformationsFlex.AddItem(tview.NewTextView().SetText(fmt.Sprintf("[::b]Public        : %v[-]", profile["public"])).SetDynamicColors(true), 1, 0, false)

	userInformationsContainer := tview.NewFlex().
		SetDirection(tview.FlexRow).
		SetDirection(tview.FlexColumn).
		AddItem(userInformationsFlex, 0, 1, false).
		AddItem(userRankingInformationsFlex, 0, 1, false).
		AddItem(userMiscInformationsFlex, 0, 1, false)

	fortressesPanel := displayInfo(advancedLabsMap, "fortresses", "Fortresses", "\U0001F3F0", maxFortressNameLength, 4)
	prolabsPanel := displayInfo(advancedLabsMap, "prolabs", "Pro Labs", "\U0001F47D", maxProlabNameLength, 4)
	activityPanel := displayInfo(advancedLabsMap, "activity", "Activity", "", maxActivityNameLength, 3)

	advancedLabsFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		SetDirection(tview.FlexColumn)
	
	// Only add panels that are not nil
	if fortressesPanel != nil {
		advancedLabsFlex.AddItem(fortressesPanel, 0, 1, false)
	}
	if prolabsPanel != nil {
		advancedLabsFlex.AddItem(prolabsPanel, 0, 1, false)
	}

	leftFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(userInformationsContainer, 0, 1, false).
		AddItem(advancedLabsFlex, 0, 1, false)
	
	// Only add activity panel if it's not nil
	if activityPanel != nil {
		leftFlex.AddItem(activityPanel, 0, 2, false)
	}

	mainFlex := tview.NewFlex().
		AddItem(leftFlex, 0, 2, false)

	// Run the application
	if err := app.SetRoot(mainFlex, true).Run(); err != nil {
		panic(err)
	}
}
