package cmd

import (
	"fmt"

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
		fmt.Println("test")
		// rankingTable := tview.NewTable().
		// 	SetBorders(true).
		// 	SetSelectable(true, false).
		// 	SetFixed(1, 0)
		// rankingTable.SetCell(0, 0, tview.NewTableCell("Ranking").
		// 	SetTextColor(tcell.ColorWhite).
		// 	SetAttributes(tcell.AttrBold).
		// 	SetAlign(tview.AlignCenter))
		// rankingTable.SetCell(0, 1, tview.NewTableCell("Player").
		// 	SetTextColor(tcell.ColorWhite).
		// 	SetAttributes(tcell.AttrBold).
		// 	SetAlign(tview.AlignCenter))
		// rankingTable.SetCell(0, 2, tview.NewTableCell("Points").
		// 	SetTextColor(tcell.ColorWhite).
		// 	SetAttributes(tcell.AttrBold).
		// 	SetAlign(tview.AlignCenter))
		// rankingTable.SetCell(0, 3, tview.NewTableCell("Users").
		// 	SetTextColor(tcell.ColorWhite).
		// 	SetAttributes(tcell.AttrBold).
		// 	SetAlign(tview.AlignCenter))
		// rankingTable.SetCell(0, 4, tview.NewTableCell("Systems").
		// 	SetTextColor(tcell.ColorWhite).
		// 	SetAttributes(tcell.AttrBold).
		// 	SetAlign(tview.AlignCenter))
		// rankingTable.SetCell(0, 5, tview.NewTableCell("Challenges").
		// 	SetTextColor(tcell.ColorWhite).
		// 	SetAttributes(tcell.AttrBold).
		// 	SetAlign(tview.AlignCenter))

		// Ajout des données de joueurs dans le tableau
		// playerData := []struct {
		// 	Name       string
		// 	Points     string
		// 	Users      string
		// 	Systems    string
		// 	Challenges string
		// }{
		// 	{"Player1", "3000", "12", "13", "343"},
		// 	{"Player2", "2500", "19", "10", "363"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player1", "3000", "12", "13", "343"},
		// 	{"Player2", "2500", "19", "10", "363"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player1", "3000", "12", "13", "343"},
		// 	{"Player2", "2500", "19", "10", "363"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player1", "3000", "12", "13", "343"},
		// 	{"Player2", "2500", "19", "10", "363"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// 	{"Player3", "2000", "22", "14", "353"},
		// }

		// for i, player := range playerData {
		// 	// Ajout du classement du joueur
		// 	rankingTable.SetCell(i+1, 0, tview.NewTableCell(strconv.Itoa(i+1)).
		// 		SetTextColor(tcell.ColorWhite).
		// 		SetAlign(tview.AlignCenter))
		// 	// Ajout du nom du joueur
		// 	rankingTable.SetCell(i+1, 1, tview.NewTableCell(player.Name).
		// 		SetTextColor(tcell.ColorWhite).
		// 		SetAlign(tview.AlignLeft))
		// 	// Ajout des points du joueur
		// 	rankingTable.SetCell(i+1, 2, tview.NewTableCell(string(player.Points)).
		// 		SetTextColor(tcell.ColorWhite).
		// 		SetAlign(tview.AlignRight))
		// 	// Ajout des users machines du joueur
		// 	rankingTable.SetCell(i+1, 3, tview.NewTableCell(string(player.Users)).
		// 		SetTextColor(tcell.ColorWhite).
		// 		SetAlign(tview.AlignRight))
		// 	// Ajout des systems machines du joueur
		// 	rankingTable.SetCell(i+1, 4, tview.NewTableCell(string(player.Systems)).
		// 		SetTextColor(tcell.ColorWhite).
		// 		SetAlign(tview.AlignRight))
		// 	// Ajout des challenges machines du joueur
		// 	rankingTable.SetCell(i+1, 5, tview.NewTableCell(string(player.Challenges)).
		// 		SetTextColor(tcell.ColorWhite).
		// 		SetAlign(tview.AlignRight))
		// }

		// // Scrolling avec les touches fléchées
		// rankingTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// 	switch event.Key() {
		// 	case tcell.KeyUp:
		// 		row, _ := rankingTable.GetSelection()
		// 		if row > 1 {
		// 			rankingTable.Select(row-1, 0)
		// 		}
		// 	case tcell.KeyDown:
		// 		row, _ := rankingTable.GetSelection()
		// 		if row < len(playerData) { // assuming playerData is available in this context
		// 			rankingTable.Select(row+1, 0)
		// 		}
		// 	}
		// 	return event
		// })

		// // Ranking card
		// rankingCard := tview.NewFlex().SetDirection(tview.FlexRow)
		// rankingCard.SetBorder(true).SetTitle("Global Ranking").SetTitleAlign(tview.AlignLeft)
		// rankingCard.AddItem(rankingTable, 0, 1, false)

		// // Flex droite basse
		// rightBottomFlex := tview.NewFlex().SetDirection(tview.FlexRow)
		// rightBottomFlex.SetBorder(true).SetTitle("Country Ranking").SetTitleAlign(tview.AlignLeft)
		// rightBottomFlex.AddItem(rankingTable, 0, 1, false)

		// Get fortresses
		// fortressesPanel := displayFortressesInfo(fortresses)
		// endgamesPanel := displayEndgamesInfo(endgames)
		// prolabsPanel := displayProlabsInfo(prolabs)
		// activityPanel := displayActivityInfo(activity)

		// Flex droite qui contient les flex droite haute et droite basse
		// rightFlex := tview.NewFlex().
		// 	SetDirection(tview.FlexRow).
		// 	AddItem(rankingCard, 0, 1, false).
		// 	AddItem(rightBottomFlex, 0, 1, false)

		// Gestion du focus avec Tab
		// app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// 	if event.Key() == tcell.KeyTab {
		// 		if app.GetFocus() == rankingCard {
		// 			app.SetFocus(rightBottomFlex)
		// 		} else {
		// 			app.SetFocus(rankingTable)
		// 		}
		// 		// Ne pas propager l'événement pour éviter de déplacer le focus inutilement
		// 		return nil
		// 	}
		// 	return event
		// })
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
