package utils

import (
	"fmt"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func calculateSpacing(baseName string, maxNameLength int) string {
	return fmt.Sprintf("%-*s", maxNameLength-len(baseName)+1, "")
}

func parseUserSubscription(profile map[string]interface{}) string {
	isVip := profile["isVip"].(bool)
	isDedicatedVIP := profile["isDedicatedVip"].(bool)

	if isDedicatedVIP {
		return "VIP+"
	} else if isVip {
		return "VIP"
	}

	return "Free"
}

func displayFortressesInfo(fortresses map[string]interface{}) *tview.Flex {
	fortressesList, ok := fortresses["fortresses"].([]interface{})
	if !ok {
		fmt.Println("Couldn't convert fortresses[\"fortresses\"] to a slice of map")
		return nil
	}
	fortressesPanel := tview.NewFlex().SetDirection(tview.FlexRow)
	fortressesPanel.SetBorder(true).SetTitle("Fortresses").SetTitleAlign(tview.AlignLeft)
	maxNameLength := 9
	for _, fortressInterface := range fortressesList {
		fortress, ok := fortressInterface.(map[string]interface{})
		if !ok {
			fmt.Println("Couldn't convert fortress to a map[string]interface{}")
			continue
		}
		name := fortress["name"].(string)
		ownedFlags := fortress["owned_flags"].(float64)
		totalFlags := fortress["total_flags"].(float64)

		var color string
		if ownedFlags == totalFlags {
			color = "[green]"
		} else if ownedFlags == 0 {
			color = "[red]"
		} else {
			color = "[orange]"
		}

		spacing := calculateSpacing(name, maxNameLength)

		text := fmt.Sprintf("[::b]\U0001F3F0 %s%s: %s%.0f/%.0f[-]", name, spacing, color, ownedFlags, totalFlags)
		fortressesPanel.AddItem(tview.NewTextView().SetText(text).SetDynamicColors(true), 1, 0, false)
	}
	fortressesPanel.AddItem(tview.NewTextView().SetText("").SetDynamicColors(true), 4, 0, false)
	return fortressesPanel
}

func displayEndgamesInfo(endgames map[string]interface{}) *tview.Flex {
	endgamesList, ok := endgames["endgames"].([]interface{})
	if !ok {
		fmt.Println("Couldn't convert endgames[\"endgames\"] to a slice of map")
		return nil
	}
	endgamesPanel := tview.NewFlex().SetDirection(tview.FlexRow)
	endgamesPanel.SetBorder(true).SetTitle("Endgames").SetTitleAlign(tview.AlignLeft)
	maxNameLength := 9
	for _, endgamesInterface := range endgamesList {
		endgame, ok := endgamesInterface.(map[string]interface{})
		if !ok {
			fmt.Println("Couldn't convert endgame to a map[string]interface{}")
			continue
		}
		name := endgame["name"].(string)
		ownedFlags := endgame["owned_flags"].(float64)
		totalFlags := endgame["total_flags"].(float64)

		var color string
		if ownedFlags == totalFlags {
			color = "[green]"
		} else if ownedFlags == 0 {
			color = "[red]"
		} else {
			color = "[orange]"
		}

		spacing := calculateSpacing(name, maxNameLength)

		text := fmt.Sprintf("[::b]\U0001F3AE %s%s: %s%.0f/%.0f[-]", name, spacing, color, ownedFlags, totalFlags)
		endgamesPanel.AddItem(tview.NewTextView().SetText(text).SetDynamicColors(true), 1, 0, false)
	}
	endgamesPanel.AddItem(tview.NewTextView().SetText("").SetDynamicColors(true), 3, 0, false)
	return endgamesPanel
}

func displayProlabsInfo(prolabs map[string]interface{}) *tview.Flex {
	prolabsList, ok := prolabs["prolabs"].([]interface{})
	if !ok {
		fmt.Println("Couldn't convert endgames[\"endgames\"] to a slice of map")
		return nil
	}
	prolabsPanel := tview.NewFlex().SetDirection(tview.FlexRow)
	prolabsPanel.SetBorder(true).SetTitle("Pro Labs").SetTitleAlign(tview.AlignLeft)
	maxNameLength := 11
	for _, prolabsInterface := range prolabsList {
		prolab, ok := prolabsInterface.(map[string]interface{})
		if !ok {
			fmt.Println("Couldn't convert prolab to a map[string]interface{}")
			continue
		}
		name := prolab["name"].(string)
		ownedFlags := prolab["owned_flags"].(float64)
		totalFlags := prolab["total_flags"].(float64)

		var color string
		if ownedFlags == totalFlags {
			color = "[green]"
		} else if ownedFlags == 0 {
			color = "[red]"
		} else {
			color = "[orange]"
		}

		spacing := calculateSpacing(name, maxNameLength)

		text := fmt.Sprintf("[::b]\U0001F47D %s%s: %s%.0f/%.0f[-]", name, spacing, color, ownedFlags, totalFlags)
		prolabsPanel.AddItem(tview.NewTextView().SetText(text).SetDynamicColors(true), 1, 0, false)
	}
	prolabsPanel.AddItem(tview.NewTextView().SetText("").SetDynamicColors(true), 4, 0, false)
	return prolabsPanel
}

func DisplayInformationsGUI(profile map[string]interface{}, fortresses map[string]interface{}, endgames map[string]interface{}, prolabs map[string]interface{}) {

	teamName, teamRank := "N/A", "N/A"
	universityName, universityRank := "N/A", "N/A"

	if teamMap, ok := profile["team"].(map[string]interface{}); ok && teamMap != nil {
		teamName = teamMap["name"].(string)
		teamRank = fmt.Sprintf("%v", teamMap["ranking"].(float64))
	}

	if universityMap, ok := profile["university"].(map[string]interface{}); ok && universityMap != nil {
		universityName = universityMap["name"].(string)
		universityRank = fmt.Sprintf("%v", universityMap["rank"].(float64))
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
	userInformationsFlex.AddItem(tview.NewTextView().SetText(fmt.Sprintf("[::b]ID           : %d[-]", int(profile["id"].(float64)))).SetDynamicColors(true), 1, 0, false)
	userInformationsFlex.AddItem(tview.NewTextView().SetText(fmt.Sprintf("[::b]Name         : %v[-]", profile["name"])).SetDynamicColors(true), 1, 0, false)
	userInformationsFlex.AddItem(tview.NewTextView().SetText(fmt.Sprintf("[::b]Team         : %v[-]", teamName)).SetDynamicColors(true), 1, 0, false)
	userInformationsFlex.AddItem(tview.NewTextView().SetText(fmt.Sprintf("[::b]University   : %v[-]", universityName)).SetDynamicColors(true), 1, 0, false)
	userInformationsFlex.AddItem(tview.NewTextView().SetText(fmt.Sprintf("[::b]Country      : %v[-]", profile["country_name"])).SetDynamicColors(true), 1, 0, false)
	userInformationsFlex.AddItem(tview.NewTextView().SetText(fmt.Sprintf("[::b]Subscription : %v[-]", subscription)).SetDynamicColors(true), 1, 0, false)
	userInformationsFlex.AddItem(tview.NewTextView().SetText(fmt.Sprintf("[::b]Respects     : %v[-]", profile["respects"])).SetDynamicColors(true), 1, 0, false)
	userInformationsFlex.AddItem(tview.NewTextView().SetText("").SetDynamicColors(true), 3, 0, false)

	userRankingInformationsFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	userRankingInformationsFlex.SetBorder(true).SetTitle("Ranking Informations").SetTitleAlign(tview.AlignLeft)
	userRankingInformationsFlex.AddItem(tview.NewTextView().SetText(fmt.Sprintf("[::b]Global     : %v[-]\U0001F3C6", ranking)).SetDynamicColors(true), 1, 0, false)
	userRankingInformationsFlex.AddItem(tview.NewTextView().SetText(fmt.Sprintf("[::b]Team       : %v[-]\U0001F91D", teamRank)).SetDynamicColors(true), 1, 0, false)
	userRankingInformationsFlex.AddItem(tview.NewTextView().SetText(fmt.Sprintf("[::b]University : %v[-]\U0001F3EB", universityRank)).SetDynamicColors(true), 1, 0, false)
	userRankingInformationsFlex.AddItem(tview.NewTextView().SetText(fmt.Sprintf("[::b]Points     : %v[-]\U0001F396", profile["points"])).SetDynamicColors(true), 1, 0, false)
	userRankingInformationsFlex.AddItem(tview.NewTextView().SetText(fmt.Sprintf("[::b]Rank       : %v[-]", profile["rank"])).SetDynamicColors(true), 1, 0, false)
	userRankingInformationsFlex.AddItem(tview.NewTextView().SetText(fmt.Sprintf("[::b]Ownership  : %v%% / %v%%[-]", profile["rank_ownership"], rankRequirement)).SetDynamicColors(true), 1, 0, false)
	userRankingInformationsFlex.AddItem(tview.NewTextView().SetText("").SetDynamicColors(true), 4, 0, false)

	userMiscInformationsFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	userMiscInformationsFlex.SetBorder(true).SetTitle("Misc").SetTitleAlign(tview.AlignLeft)
	userMiscInformationsFlex.AddItem(tview.NewTextView().SetText(fmt.Sprintf("[::b]User Bloods   : %v\U0001FA78[-]", profile["user_bloods"])).SetDynamicColors(true), 1, 0, false)
	userMiscInformationsFlex.AddItem(tview.NewTextView().SetText(fmt.Sprintf("[::b]System Bloods : %v\U0001FA78[-]", profile["system_bloods"])).SetDynamicColors(true), 1, 0, false)
	userMiscInformationsFlex.AddItem(tview.NewTextView().SetText(fmt.Sprintf("[::b]User Owns     : %v[-]", profile["user_owns"])).SetDynamicColors(true), 1, 0, false)
	userMiscInformationsFlex.AddItem(tview.NewTextView().SetText(fmt.Sprintf("[::b]System Owns   : %v[-]", profile["system_owns"])).SetDynamicColors(true), 1, 0, false)
	userMiscInformationsFlex.AddItem(tview.NewTextView().SetText(fmt.Sprintf("[::b]Last User     : N/A[-]")).SetDynamicColors(true), 1, 0, false)
	userMiscInformationsFlex.AddItem(tview.NewTextView().SetText(fmt.Sprintf("[::b]Last System   : N/A[-]")).SetDynamicColors(true), 1, 0, false)
	userMiscInformationsFlex.AddItem(tview.NewTextView().SetText("").SetDynamicColors(true), 4, 0, false)

	userInformationsContainer := tview.NewFlex().
		SetDirection(tview.FlexRow).
		SetDirection(tview.FlexColumn).
		AddItem(userInformationsFlex, 0, 1, false).
		AddItem(userRankingInformationsFlex, 0, 1, false).
		AddItem(userMiscInformationsFlex, 0, 1, false)

	rankingTable := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false).
		SetFixed(1, 0)
	rankingTable.SetCell(0, 0, tview.NewTableCell("Ranking").
		SetTextColor(tcell.ColorWhite).
		SetAttributes(tcell.AttrBold).
		SetAlign(tview.AlignCenter))
	rankingTable.SetCell(0, 1, tview.NewTableCell("Player").
		SetTextColor(tcell.ColorWhite).
		SetAttributes(tcell.AttrBold).
		SetAlign(tview.AlignCenter))
	rankingTable.SetCell(0, 2, tview.NewTableCell("Points").
		SetTextColor(tcell.ColorWhite).
		SetAttributes(tcell.AttrBold).
		SetAlign(tview.AlignCenter))
	rankingTable.SetCell(0, 3, tview.NewTableCell("Users").
		SetTextColor(tcell.ColorWhite).
		SetAttributes(tcell.AttrBold).
		SetAlign(tview.AlignCenter))
	rankingTable.SetCell(0, 4, tview.NewTableCell("Systems").
		SetTextColor(tcell.ColorWhite).
		SetAttributes(tcell.AttrBold).
		SetAlign(tview.AlignCenter))
	rankingTable.SetCell(0, 5, tview.NewTableCell("Challenges").
		SetTextColor(tcell.ColorWhite).
		SetAttributes(tcell.AttrBold).
		SetAlign(tview.AlignCenter))

	// Ajout des données de joueurs dans le tableau
	playerData := []struct {
		Name       string
		Points     string
		Users      string
		Systems    string
		Challenges string
	}{
		{"Player1", "3000", "12", "13", "343"},
		{"Player2", "2500", "19", "10", "363"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player1", "3000", "12", "13", "343"},
		{"Player2", "2500", "19", "10", "363"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player1", "3000", "12", "13", "343"},
		{"Player2", "2500", "19", "10", "363"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player1", "3000", "12", "13", "343"},
		{"Player2", "2500", "19", "10", "363"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player3", "2000", "22", "14", "353"},
		{"Player3", "2000", "22", "14", "353"},
	}

	for i, player := range playerData {
		// Ajout du classement du joueur
		rankingTable.SetCell(i+1, 0, tview.NewTableCell(strconv.Itoa(i+1)).
			SetTextColor(tcell.ColorWhite).
			SetAlign(tview.AlignCenter))
		// Ajout du nom du joueur
		rankingTable.SetCell(i+1, 1, tview.NewTableCell(player.Name).
			SetTextColor(tcell.ColorWhite).
			SetAlign(tview.AlignLeft))
		// Ajout des points du joueur
		rankingTable.SetCell(i+1, 2, tview.NewTableCell(string(player.Points)).
			SetTextColor(tcell.ColorWhite).
			SetAlign(tview.AlignRight))
		// Ajout des users machines du joueur
		rankingTable.SetCell(i+1, 3, tview.NewTableCell(string(player.Users)).
			SetTextColor(tcell.ColorWhite).
			SetAlign(tview.AlignRight))
		// Ajout des systems machines du joueur
		rankingTable.SetCell(i+1, 4, tview.NewTableCell(string(player.Systems)).
			SetTextColor(tcell.ColorWhite).
			SetAlign(tview.AlignRight))
		// Ajout des challenges machines du joueur
		rankingTable.SetCell(i+1, 5, tview.NewTableCell(string(player.Challenges)).
			SetTextColor(tcell.ColorWhite).
			SetAlign(tview.AlignRight))
	}

	// Scrolling avec les touches fléchées
	rankingTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp:
			row, _ := rankingTable.GetSelection()
			if row > 1 {
				rankingTable.Select(row-1, 0)
			}
		case tcell.KeyDown:
			row, _ := rankingTable.GetSelection()
			if row < len(playerData) { // assuming playerData is available in this context
				rankingTable.Select(row+1, 0)
			}
		}
		return event
	})

	// Ranking card
	rankingCard := tview.NewFlex().SetDirection(tview.FlexRow)
	rankingCard.SetBorder(true).SetTitle("Global Ranking").SetTitleAlign(tview.AlignLeft)
	rankingCard.AddItem(rankingTable, 0, 1, false)

	// Flex droite basse
	rightBottomFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	rightBottomFlex.SetBorder(true).SetTitle("Country Ranking").SetTitleAlign(tview.AlignLeft)
	rightBottomFlex.AddItem(rankingTable, 0, 1, false)

	// Flex bas gauche history
	historyFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	historyFlex.SetBorder(true).SetTitle("History").SetTitleAlign(tview.AlignLeft)
	historyFlex.AddItem(tview.NewTextView().SetText("Owned [::b]System[-] - Cozyhosting Machine - 1 day ago - [green]+[20pts][-]").SetDynamicColors(true), 1, 0, false)
	historyFlex.AddItem(tview.NewTextView().SetText("Owned [::b]User[-] - Cozyhosting Machine - 1 day ago - [green]+[10pts][-]").SetDynamicColors(true), 1, 0, false)
	historyFlex.AddItem(tview.NewTextView().SetText("Owned [::b]Web[-] - WS-Todo Challenge - 1 month ago - [green]+[4pts][-]").SetDynamicColors(true), 1, 0, false)
	historyFlex.AddItem(tview.NewTextView().SetText("Owned [::b]Web[-] - TwoDots Horror Challenge - 1 month ago - [green]+[3pts][-]").SetDynamicColors(true), 1, 0, false)
	historyFlex.AddItem(tview.NewTextView().SetText("Owned [::b]Web[-] - TwoDots Horror Challenge - 1 month ago - [green]+[3pts][-]").SetDynamicColors(true), 1, 0, false)
	historyFlex.AddItem(tview.NewTextView().SetText("Owned [::b]Web[-] - TwoDots Horror Challenge - 1 month ago - [green]+[3pts][-]").SetDynamicColors(true), 1, 0, false)
	historyFlex.AddItem(tview.NewTextView().SetText("Owned [::b]Web[-] - TwoDots Horror Challenge - 1 month ago - [green]+[3pts][-]").SetDynamicColors(true), 1, 0, false)
	historyFlex.AddItem(tview.NewTextView().SetText("Owned [::b]Web[-] - TwoDots Horror Challenge - 1 month ago - [green]+[3pts][-]").SetDynamicColors(true), 1, 0, false)
	historyFlex.AddItem(tview.NewTextView().SetText("Owned [::b]Web[-] - TwoDots Horror Challenge - 1 month ago - [green]+[3pts][-]").SetDynamicColors(true), 1, 0, false)
	historyFlex.AddItem(tview.NewTextView().SetText("Owned [::b]Web[-] - TwoDots Horror Challenge - 1 month ago - [green]+[3pts][-]").SetDynamicColors(true), 1, 0, false)
	historyFlex.AddItem(tview.NewTextView().SetText("Owned [::b]Web[-] - TwoDots Horror Challenge - 1 month ago - [green]+[3pts][-]").SetDynamicColors(true), 1, 0, false)
	historyFlex.AddItem(tview.NewTextView().SetText("Owned [::b]Web[-] - WS-Todo Challenge - 1 month ago - [green]+[4pts][-]").SetDynamicColors(true), 1, 0, false)
	historyFlex.AddItem(tview.NewTextView().SetText("Owned [::b]Web[-] - WS-Todo Challenge - 1 month ago - [green]+[4pts][-]").SetDynamicColors(true), 1, 0, false)
	historyFlex.AddItem(tview.NewTextView().SetText("Owned [::b]Web[-] - WS-Todo Challenge - 1 month ago - [green]+[4pts][-]").SetDynamicColors(true), 1, 0, false)
	historyFlex.AddItem(tview.NewTextView().SetText("Owned [::b]Web[-] - WS-Todo Challenge - 1 month ago - [green]+[4pts][-]").SetDynamicColors(true), 1, 0, false)
	historyFlex.AddItem(tview.NewTextView().SetText("Owned [::b]Web[-] - WS-Todo Challenge - 1 month ago - [green]+[4pts][-]").SetDynamicColors(true), 1, 0, false)
	historyFlex.AddItem(tview.NewTextView().SetText("Owned [::b]Web[-] - WS-Todo Challenge - 1 month ago - [green]+[4pts][-]").SetDynamicColors(true), 1, 0, false)
	historyFlex.AddItem(tview.NewTextView().SetText("Owned [::b]Web[-] - WS-Todo Challenge - 1 month ago - [green]+[4pts][-]").SetDynamicColors(true), 1, 0, false)
	historyFlex.AddItem(tview.NewTextView().SetText("Owned [::b]Web[-] - WS-Todo Challenge - 1 month ago - [green]+[4pts][-]").SetDynamicColors(true), 1, 0, false)
	historyFlex.AddItem(tview.NewTextView().SetText("Owned [::b]Web[-] - WS-Todo Challenge - 1 month ago - [green]+[4pts][-]").SetDynamicColors(true), 1, 0, false)
	historyFlex.AddItem(tview.NewTextView().SetText("Owned [::b]Web[-] - WS-Todo Challenge - 1 month ago - [green]+[4pts][-]").SetDynamicColors(true), 1, 0, false)
	historyFlex.AddItem(tview.NewTextView().SetText("Owned [::b]Web[-] - WS-Todo Challenge - 1 month ago - [green]+[4pts][-]").SetDynamicColors(true), 1, 0, false)
	historyFlex.AddItem(tview.NewTextView().SetText("Owned [::b]Web[-] - WS-Todo Challenge - 1 month ago - [green]+[4pts][-]").SetDynamicColors(true), 1, 0, false)

	// historyFlex.AddItem(rankingTable, 0, 1, false)

	// Fortresses
	// var jsonStr = `{"fortresses":[{"name":"Jet","avatar":"https://www.hackthebox.com/storage/companies/3.png","completion_percentage":64,"owned_flags":7,"total_flags":11},{"name":"Akerva","avatar":"https://www.hackthebox.com/storage/companies/61.png","completion_percentage":100,"owned_flags":8,"total_flags":8},{"name":"Context","avatar":"https://www.hackthebox.com/storage/companies/9.png","completion_percentage":0,"owned_flags":0,"total_flags":7},{"name":"Synacktiv","avatar":"https://www.hackthebox.com/storage/companies/195.png","completion_percentage":71,"owned_flags":5,"total_flags":7},{"name":"Faraday","avatar":"https://www.hackthebox.com/storage/companies/28.png","completion_percentage":0,"owned_flags":0,"total_flags":7},{"name":"AWS","avatar":"https://www.hackthebox.com/storage/companies/206.png","completion_percentage":0,"owned_flags":0,"total_flags":10}]}`

	// var data map[string][]map[string]interface{}

	// err := json.Unmarshal([]byte(jsonStr), &data)
	// if err != nil {
	// 	fmt.Println("error:", err)
	// }

	// Flex bas gauche history
	// prolabsPanel := tview.NewFlex().SetDirection(tview.FlexRow)
	// prolabsPanel.SetBorder(true).SetTitle("Pro Labs").SetTitleAlign(tview.AlignLeft)
	// prolabsPanel.AddItem(tview.NewTextView().SetText("[::b]\U0001F47D Dante       : [green]11/11[-]").SetDynamicColors(true), 1, 0, false)
	// prolabsPanel.AddItem(tview.NewTextView().SetText("[::b]\U0001F47D Offshore    : [green]11/11[-]").SetDynamicColors(true), 1, 0, false)
	// prolabsPanel.AddItem(tview.NewTextView().SetText("[::b]\U0001F47D Zephyr      : [green]11/11[-]").SetDynamicColors(true), 1, 0, false)
	// prolabsPanel.AddItem(tview.NewTextView().SetText("[::b]\U0001F47D Rastalabs   : [green]11/11[-]").SetDynamicColors(true), 1, 0, false)
	// prolabsPanel.AddItem(tview.NewTextView().SetText("[::b]\U0001F47D Cybernetics : [green]11/11[-]").SetDynamicColors(true), 1, 0, false)
	// prolabsPanel.AddItem(tview.NewTextView().SetText("[::b]\U0001F47D APTLabs     : [green]11/11[-]").SetDynamicColors(true), 1, 0, false)
	// prolabsPanel.AddItem(tview.NewTextView().SetText("").SetDynamicColors(true), 4, 0, false)

	// Flex bas gauche history
	// endgamesPanel := tview.NewFlex().SetDirection(tview.FlexRow)
	// endgamesPanel.SetBorder(true).SetTitle("Endgames").SetTitleAlign(tview.AlignLeft)
	// endgamesPanel.AddItem(tview.NewTextView().SetText("[::b]\U0001F3AE Solar     : [red]0/10[-]").SetDynamicColors(true), 1, 0, false)
	// endgamesPanel.AddItem(tview.NewTextView().SetText("[::b]\U0001F3AE Odyssey   : [red]0/10[-]").SetDynamicColors(true), 1, 0, false)
	// endgamesPanel.AddItem(tview.NewTextView().SetText("[::b]\U0001F3AE Ascension : [red]0/10[-]").SetDynamicColors(true), 1, 0, false)
	// endgamesPanel.AddItem(tview.NewTextView().SetText("[::b]\U0001F3AE RPG       : [red]0/10[-]").SetDynamicColors(true), 1, 0, false)
	// endgamesPanel.AddItem(tview.NewTextView().SetText("[::b]\U0001F3AE Hades     : [red]0/10[-]").SetDynamicColors(true), 1, 0, false)
	// endgamesPanel.AddItem(tview.NewTextView().SetText("[::b]\U0001F3AE Xen       : [red]0/10[-]").SetDynamicColors(true), 1, 0, false)
	// endgamesPanel.AddItem(tview.NewTextView().SetText("[::b]\U0001F3AE P.O.O     : [red]0/10[-]").SetDynamicColors(true), 1, 0, false)
	// endgamesPanel.AddItem(tview.NewTextView().SetText("").SetDynamicColors(true), 3, 0, false)

	// Get fortresses
	fortressesPanel := displayFortressesInfo(fortresses)
	endgamesPanel := displayEndgamesInfo(endgames)
	prolabsPanel := displayProlabsInfo(prolabs)

	advancedLabsFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		SetDirection(tview.FlexColumn).
		AddItem(fortressesPanel, 0, 1, false).
		AddItem(prolabsPanel, 0, 1, false).
		AddItem(endgamesPanel, 0, 1, false)

	// Flex droite qui contient les flex droite haute et droite basse
	rightFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(rankingCard, 0, 1, false).
		AddItem(rightBottomFlex, 0, 1, false)

	// Flex gauche qui contient les flex gauche haute et droite basse
	leftFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(userInformationsContainer, 0, 1, false).
		AddItem(advancedLabsFlex, 0, 1, false).
		AddItem(historyFlex, 0, 2, false)

	// Gestion du focus avec Tab
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			if app.GetFocus() == rankingCard {
				app.SetFocus(rightBottomFlex)
			} else {
				app.SetFocus(rankingTable)
			}
			// Ne pas propager l'événement pour éviter de déplacer le focus inutilement
			return nil
		}
		return event
	})

	mainFlex := tview.NewFlex().
		AddItem(leftFlex, 0, 2, false).
		AddItem(rightFlex, 0, 1, false)

	// Run the application
	if err := app.SetRoot(mainFlex, true).Run(); err != nil {
		panic(err)
	}
}
