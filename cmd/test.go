package cmd

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		app := tview.NewApplication()

		userInformationsFlex := tview.NewFlex().SetDirection(tview.FlexRow)
		userInformationsFlex.SetBorder(true).SetTitle("Profile").SetTitleAlign(tview.AlignLeft)
		userInformationsFlex.AddItem(tview.NewTextView().SetText("[::b]Name         : QU35T3190[-]").SetDynamicColors(true), 1, 0, false)
		userInformationsFlex.AddItem(tview.NewTextView().SetText("[::b]Team         : SimianSec[-]").SetDynamicColors(true), 1, 0, false)
		userInformationsFlex.AddItem(tview.NewTextView().SetText("[::b]University   : Epitech[-]").SetDynamicColors(true), 1, 0, false)
		userInformationsFlex.AddItem(tview.NewTextView().SetText("[::b]Country      : France[-]").SetDynamicColors(true), 1, 0, false)
		userInformationsFlex.AddItem(tview.NewTextView().SetText("[::b]Subscription : VIP[-]").SetDynamicColors(true), 1, 0, false)
		userInformationsFlex.AddItem(tview.NewTextView().SetText("[::b]Respects     : 12[-]").SetDynamicColors(true), 1, 0, false)
		userInformationsFlex.AddItem(tview.NewTextView().SetText("").SetDynamicColors(true), 4, 0, false)

		userRankingInformationsFlex := tview.NewFlex().SetDirection(tview.FlexRow)
		userRankingInformationsFlex.SetBorder(true).SetTitle("Ranking Informations").SetTitleAlign(tview.AlignLeft)
		userRankingInformationsFlex.AddItem(tview.NewTextView().SetText("[::b]Global       : 737[-]\U0001F3C6").SetDynamicColors(true), 1, 0, false)
		userRankingInformationsFlex.AddItem(tview.NewTextView().SetText("[::b]Team         : 737[-]\U0001F3C6").SetDynamicColors(true), 1, 0, false)
		userRankingInformationsFlex.AddItem(tview.NewTextView().SetText("[::b]University   : 737[-]\U0001F3C6").SetDynamicColors(true), 1, 0, false)
		userRankingInformationsFlex.AddItem(tview.NewTextView().SetText("[::b]Points       : 370[-]\U0001F396").SetDynamicColors(true), 1, 0, false)
		userRankingInformationsFlex.AddItem(tview.NewTextView().SetText("[::b]Rank         : Pro Hacker[-]").SetDynamicColors(true), 1, 0, false)
		userRankingInformationsFlex.AddItem(tview.NewTextView().SetText("[::b]Ownership    : 4.12%[-]").SetDynamicColors(true), 1, 0, false)
		userRankingInformationsFlex.AddItem(tview.NewTextView().SetText("").SetDynamicColors(true), 4, 0, false)

		userMiscInformationsFlex := tview.NewFlex().SetDirection(tview.FlexRow)
		userMiscInformationsFlex.SetBorder(true).SetTitle("Misc").SetTitleAlign(tview.AlignLeft)
		userMiscInformationsFlex.AddItem(tview.NewTextView().SetText("[::b]User Bloods[-]   : 1\U0001FA78[-]").SetDynamicColors(true), 1, 0, false)
		userMiscInformationsFlex.AddItem(tview.NewTextView().SetText("[::b]System Bloods[-] : 1\U0001FA78[-]").SetDynamicColors(true), 1, 0, false)
		userMiscInformationsFlex.AddItem(tview.NewTextView().SetText("[::b]User Owns[-]     : 0[-]").SetDynamicColors(true), 1, 0, false)
		userMiscInformationsFlex.AddItem(tview.NewTextView().SetText("[::b]System Owns[-]   : 0[-]").SetDynamicColors(true), 1, 0, false)
		userMiscInformationsFlex.AddItem(tview.NewTextView().SetText("[::b]Last User[-]     : Sau[-]").SetDynamicColors(true), 1, 0, false)
		userMiscInformationsFlex.AddItem(tview.NewTextView().SetText("[::b]Last System[-]   : Shoppy[-]").SetDynamicColors(true), 1, 0, false)
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

		// Flex bas gauche history
		fortressesPanel := tview.NewFlex().SetDirection(tview.FlexRow)
		fortressesPanel.SetBorder(true).SetTitle("Fortresses").SetTitleAlign(tview.AlignLeft)
		fortressesPanel.AddItem(tview.NewTextView().SetText("[::b]\U0001F3F0 Jet       : [green]11/11[-]").SetDynamicColors(true), 1, 0, false)
		fortressesPanel.AddItem(tview.NewTextView().SetText("[::b]\U0001F3F0 Akerva    : [green]8/8[-]").SetDynamicColors(true), 1, 0, false)
		fortressesPanel.AddItem(tview.NewTextView().SetText("[::b]\U0001F3F0 Context   : [red]2/7[-]").SetDynamicColors(true), 1, 0, false)
		fortressesPanel.AddItem(tview.NewTextView().SetText("[::b]\U0001F3F0 Synacktiv : [orange]5/7[-]").SetDynamicColors(true), 1, 0, false)
		fortressesPanel.AddItem(tview.NewTextView().SetText("[::b]\U0001F3F0 Faraday   : [orange]4/7[-]").SetDynamicColors(true), 1, 0, false)
		fortressesPanel.AddItem(tview.NewTextView().SetText("[::b]\U0001F3F0 AWS       : [red]0/10[-]").SetDynamicColors(true), 1, 0, false)
		fortressesPanel.AddItem(tview.NewTextView().SetText("").SetDynamicColors(true), 4, 0, false)

		// Flex bas gauche history
		prolabsPanel := tview.NewFlex().SetDirection(tview.FlexRow)
		prolabsPanel.SetBorder(true).SetTitle("Pro Labs").SetTitleAlign(tview.AlignLeft)
		prolabsPanel.AddItem(tview.NewTextView().SetText("[::b]\U0001F47D Dante       : [green]11/11[-]").SetDynamicColors(true), 1, 0, false)
		prolabsPanel.AddItem(tview.NewTextView().SetText("[::b]\U0001F47D Offshore    : [green]11/11[-]").SetDynamicColors(true), 1, 0, false)
		prolabsPanel.AddItem(tview.NewTextView().SetText("[::b]\U0001F47D Zephyr      : [green]11/11[-]").SetDynamicColors(true), 1, 0, false)
		prolabsPanel.AddItem(tview.NewTextView().SetText("[::b]\U0001F47D Rastalabs   : [green]11/11[-]").SetDynamicColors(true), 1, 0, false)
		prolabsPanel.AddItem(tview.NewTextView().SetText("[::b]\U0001F47D Cybernetics : [green]11/11[-]").SetDynamicColors(true), 1, 0, false)
		prolabsPanel.AddItem(tview.NewTextView().SetText("[::b]\U0001F47D APTLabs     : [green]11/11[-]").SetDynamicColors(true), 1, 0, false)
		prolabsPanel.AddItem(tview.NewTextView().SetText("").SetDynamicColors(true), 4, 0, false)

		// Flex bas gauche history
		endgamesPanel := tview.NewFlex().SetDirection(tview.FlexRow)
		endgamesPanel.SetBorder(true).SetTitle("Endgames").SetTitleAlign(tview.AlignLeft)
		endgamesPanel.AddItem(tview.NewTextView().SetText("[::b]\U0001F3AE Solar     : [red]0/10[-]").SetDynamicColors(true), 1, 0, false)
		endgamesPanel.AddItem(tview.NewTextView().SetText("[::b]\U0001F3AE Odyssey   : [red]0/10[-]").SetDynamicColors(true), 1, 0, false)
		endgamesPanel.AddItem(tview.NewTextView().SetText("[::b]\U0001F3AE Ascension : [red]0/10[-]").SetDynamicColors(true), 1, 0, false)
		endgamesPanel.AddItem(tview.NewTextView().SetText("[::b]\U0001F3AE RPG       : [red]0/10[-]").SetDynamicColors(true), 1, 0, false)
		endgamesPanel.AddItem(tview.NewTextView().SetText("[::b]\U0001F3AE Hades     : [red]0/10[-]").SetDynamicColors(true), 1, 0, false)
		endgamesPanel.AddItem(tview.NewTextView().SetText("[::b]\U0001F3AE Xen       : [red]0/10[-]").SetDynamicColors(true), 1, 0, false)
		endgamesPanel.AddItem(tview.NewTextView().SetText("[::b]\U0001F3AE P.O.O     : [red]0/10[-]").SetDynamicColors(true), 1, 0, false)
		endgamesPanel.AddItem(tview.NewTextView().SetText("").SetDynamicColors(true), 3, 0, false)

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
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
