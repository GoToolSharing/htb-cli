/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"
	"strconv"
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

		// User card
		userCard := tview.NewFlex().SetDirection(tview.FlexRow)
		userCard.SetBorder(true).SetTitle("User Informations").SetTitleAlign(tview.AlignLeft)
		// Creating a line with two elements on the same row
		nameAndRank := tview.NewFlex().SetDirection(tview.FlexColumn)
		nameAndRank.AddItem(tview.NewTextView().SetText("[::b]Name[-]: [blue]QU35T3190[-]").SetDynamicColors(true), 0, 1, false)
		nameAndRank.AddItem(tview.NewTextView().SetText("[::b]Rank[-]: [blue]Pro Hacker[-]").SetDynamicColors(true), 0, 1, false)
		nameAndRank.AddItem(tview.NewTextView().SetText("[::b]Ranking[-]: [blue]737[-]\U0001F3C6").SetDynamicColors(true), 0, 1, false)
		nameAndRank.AddItem(tview.NewTextView().SetText("[::b]Points[-]: [blue]370[-]\U0001F396").SetDynamicColors(true), 0, 1, false)
		userCard.AddItem(nameAndRank, 1, 0, false)

		owns := tview.NewFlex().SetDirection(tview.FlexColumn)
		owns.AddItem(tview.NewTextView().SetText("[::b]User Owns[-]: [blue]0[-]").SetDynamicColors(true), 0, 1, false)
		owns.AddItem(tview.NewTextView().SetText("[::b]System Owns[-]: [blue]0[-]").SetDynamicColors(true), 0, 1, false)
		owns.AddItem(tview.NewTextView().SetText("[::b]User Bloods[-]: [blue]1\U0001FA78[-]").SetDynamicColors(true), 0, 1, false)
		owns.AddItem(tview.NewTextView().SetText("[::b]System Bloods[-]: [blue]1\U0001FA78[-]").SetDynamicColors(true), 0, 1, false)
		userCard.AddItem(owns, 1, 0, false)

		other := tview.NewFlex().SetDirection(tview.FlexColumn)
		other.AddItem(tview.NewTextView().SetText("[::b]Country[-]: [blue]France[-]").SetDynamicColors(true), 0, 1, false)
		other.AddItem(tview.NewTextView().SetText("[::b]Subscription[-]: [blue]VIP[-]").SetDynamicColors(true), 0, 1, false)
		other.AddItem(tview.NewTextView().SetText("[::b]Respects[-]: [blue]12[-]").SetDynamicColors(true), 0, 1, false)
		other.AddItem(tview.NewTextView().SetText("[::b]Rank Ownership[-]: [blue]4.12[-]").SetDynamicColors(true), 0, 1, false)
		userCard.AddItem(other, 2, 0, false)

		userCard.AddItem(tview.NewTextView().SetText("TEAM :").SetDynamicColors(true), 1, 0, false)

		team := tview.NewFlex().SetDirection(tview.FlexColumn)
		team.AddItem(tview.NewTextView().SetText("[::b]Name[-]: [blue]SimianSec[-]").SetDynamicColors(true), 0, 1, false)
		team.AddItem(tview.NewTextView().SetText("[::b]Rank[-]: [blue]4[-]").SetDynamicColors(true), 0, 1, false)
		userCard.AddItem(team, 2, 0, false)

		userCard.AddItem(tview.NewTextView().SetText("UNIVERSITY :").SetDynamicColors(true), 1, 0, false)
		
		university := tview.NewFlex().SetDirection(tview.FlexColumn)
		university.AddItem(tview.NewTextView().SetText("[::b]Name[-]: [blue]Epitech[-]").SetDynamicColors(true), 0, 1, false)
		university.AddItem(tview.NewTextView().SetText("[::b]Rank[-]: [blue]40[-]").SetDynamicColors(true), 0, 1, false)
		userCard.AddItem(university, 2, 0, false)

		// TODO: Add challenges owned
		// Add fortress progression
		// Add endgames progression
		// Add last machine pwned
		// Add prolabs progression

		
		// rankingCard.AddItem(tview.NewTextView().SetText("[::b]Player[-]: [blue]xct[-]").SetDynamicColors(true), 0, 1, false)
		// rankingCard.AddItem(tview.NewTextView().SetText("[::b]Points[-]: [blue]3098[-]").SetDynamicColors(true), 0, 1, false)
		// Création des entêtes
		rankingTable := tview.NewTable().SetBorders(true)
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
		playerData := []struct{
			Name   string
			Points string
			Users string
			Systems string
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

		rankingTable.SetSelectable(true, false) // Set to true, false to enable vertical scrolling
		rankingTable.SetFixed(1, 0) // Sets the first row as non-scrollable
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
				if row < len(playerData) {  // assuming playerData is available in this context
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
		historyFlex.AddItem(tview.NewTextView().SetText("Owned [::b]System[-] - Cozyhosting Machine - 1 day ago - +[20pts]").SetDynamicColors(true), 1, 0, false)
		historyFlex.AddItem(tview.NewTextView().SetText("Owned [::b]User[-] - Cozyhosting Machine - 1 day ago - +[10pts]").SetDynamicColors(true), 1, 0, false)
		historyFlex.AddItem(tview.NewTextView().SetText("Owned [::b]Web[-] - WS-Todo Challenge - 1 month ago - +[4pts]").SetDynamicColors(true), 1, 0, false)
		historyFlex.AddItem(tview.NewTextView().SetText("Owned [::b]Web[-] - TwoDots Horror Challenge - 1 month ago - +[3pts]").SetDynamicColors(true), 1, 0, false)
		historyFlex.AddItem(tview.NewTextView().SetText("Owned [::b]Web[-] - TwoDots Horror Challenge - 1 month ago - +[3pts]").SetDynamicColors(true), 1, 0, false)
		historyFlex.AddItem(tview.NewTextView().SetText("Owned [::b]Web[-] - TwoDots Horror Challenge - 1 month ago - +[3pts]").SetDynamicColors(true), 1, 0, false)
		historyFlex.AddItem(tview.NewTextView().SetText("Owned [::b]Web[-] - TwoDots Horror Challenge - 1 month ago - +[3pts]").SetDynamicColors(true), 1, 0, false)
		historyFlex.AddItem(tview.NewTextView().SetText("Owned [::b]Web[-] - TwoDots Horror Challenge - 1 month ago - +[3pts]").SetDynamicColors(true), 1, 0, false)
		historyFlex.AddItem(tview.NewTextView().SetText("Owned [::b]Web[-] - TwoDots Horror Challenge - 1 month ago - +[3pts]").SetDynamicColors(true), 1, 0, false)
		historyFlex.AddItem(tview.NewTextView().SetText("Owned [::b]Web[-] - TwoDots Horror Challenge - 1 month ago - +[3pts]").SetDynamicColors(true), 1, 0, false)
		historyFlex.AddItem(tview.NewTextView().SetText("Owned [::b]Web[-] - TwoDots Horror Challenge - 1 month ago - +[3pts]").SetDynamicColors(true), 1, 0, false)
		historyFlex.AddItem(tview.NewTextView().SetText("Owned [::b]Web[-] - WS-Todo Challenge - 1 month ago - +[4pts]").SetDynamicColors(true), 1, 0, false)
		historyFlex.AddItem(tview.NewTextView().SetText("Owned [::b]Web[-] - WS-Todo Challenge - 1 month ago - +[4pts]").SetDynamicColors(true), 1, 0, false)
		historyFlex.AddItem(tview.NewTextView().SetText("Owned [::b]Web[-] - WS-Todo Challenge - 1 month ago - +[4pts]").SetDynamicColors(true), 1, 0, false)
		historyFlex.AddItem(tview.NewTextView().SetText("Owned [::b]Web[-] - WS-Todo Challenge - 1 month ago - +[4pts]").SetDynamicColors(true), 1, 0, false)
		
		// historyFlex.AddItem(rankingTable, 0, 1, false)

		// Flex bas gauche history
		boxesFlex := tview.NewFlex().SetDirection(tview.FlexRow)
		boxesFlex.SetBorder(true).SetTitle("Advanced Labs").SetTitleAlign(tview.AlignLeft)
		boxesFlex.AddItem(tview.NewTextView().SetText("[::b]Fortresses[-]: ").SetDynamicColors(true), 1, 0, false)
		boxesFlex.AddItem(tview.NewTextView().SetText("[::b]Jet - Akerva - Context - Synacktiv - Faraday - AWS[-]").SetDynamicColors(true), 1, 0, false)
		boxesFlex.AddItem(tview.NewTextView().SetText("[::b]Pro Labs[-]: ").SetDynamicColors(true), 1, 0, false)
		boxesFlex.AddItem(tview.NewTextView().SetText("[::b]Dante - Offshore - Zephyr - Rastalabs - Cybernetics - APTLabs[-]").SetDynamicColors(true), 1, 0, false)
		boxesFlex.AddItem(tview.NewTextView().SetText("[::b]Endgames[-]: ").SetDynamicColors(true), 1, 0, false)
		boxesFlex.AddItem(tview.NewTextView().SetText("[::b]Solar - Odyssey - Ascension - RPG - Hades - Xen - P.O.O[-]").SetDynamicColors(true), 1, 0, false)

		// Flex droite qui contient les flex droite haute et droite basse
		rightFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(rankingCard, 0, 1, false).
		AddItem(rightBottomFlex, 0, 1, false)

		// Flex gauche qui contient les flex gauche haute et droite basse
		leftFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(userCard, 0, 1, false).
		AddItem(boxesFlex, 0, 1, false).
		AddItem(historyFlex, 0, 1, false)


		// Main container
		// mainFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
		// mainFlex.AddItem(userCard, 0, 2, true)
		// mainFlex.AddItem(rankingCard, 0, 1, true)

		// Gestion du focus avec Tab
		app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			if app.GetFocus() == rankingCard {
				app.SetFocus(rightBottomFlex)
			} else {
				app.SetFocus(rankingCard)
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
