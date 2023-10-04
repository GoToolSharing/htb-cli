/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/rivo/tview"
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

		// Main container
		mainFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
		mainFlex.AddItem(userCard, 0, 1, true)

		// Run the application
		if err := app.SetRoot(mainFlex, true).SetFocus(mainFlex).Run(); err != nil {
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
