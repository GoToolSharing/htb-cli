package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GoToolSharing/htb-cli/lib/sherlocks"
	"github.com/GoToolSharing/htb-cli/lib/utils"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

var sherlocksCmd = &cobra.Command{
	Use:   "sherlocks",
	Short: "Displays active / retired / to be released sherlocks",
	Run: func(cmd *cobra.Command, args []string) {
		sherlockNameParam, err := cmd.Flags().GetString("sherlock_name")
		if err != nil {
			fmt.Println(err)
			return
		}

		sherlockDownloadPath, err := cmd.Flags().GetString("download")
		if err != nil {
			fmt.Println(err)
			return
		}

		sherlockTaskID, err := cmd.Flags().GetInt("task")
		if err != nil {
			fmt.Println(err)
			return
		}

		sherlockHint, err := cmd.Flags().GetBool("hint")
		if err != nil {
			fmt.Println(err)
			return
		}

		if sherlockNameParam != "" {
			sherlockID, err := sherlocks.SearchIDByName(sherlockNameParam)
			if err != nil {
				fmt.Println(err)
				return
			}
			log.Println("SherlockID :", sherlockID)

			if sherlockTaskID != 0 {
				err := sherlocks.GetTaskByID(sherlockID, sherlockTaskID, sherlockHint)
				if err != nil {
					fmt.Println(err)
					return
				}
				return
			}

			err = sherlocks.GetGeneralInformations(sherlockID, sherlockDownloadPath)

			if err != nil {
				fmt.Println(err)
				return
			}

			data, err := sherlocks.GetTasks(sherlockID)
			if err != nil {
				fmt.Println(err)
				return
			}
			for _, task := range data.Tasks {
				if task.Completed {
					fmt.Printf("\n%s (DONE) :\n%s\n\n", task.Title, task.Description)
				} else {
					fmt.Printf("\n%s :\n%s\n\n", task.Title, task.Description)
				}
			}
			return
		}
		app := tview.NewApplication()

		getAndDisplayFlex := func(url, title string, isScheduled bool, flex *tview.Flex) error {
			resp, err := utils.HtbRequest(http.MethodGet, url, nil)
			if err != nil {
				return fmt.Errorf("failed to get data from %s: %w", url, err)
			}

			info := utils.ParseJsonMessage(resp, "data")

			machineFlex, err := sherlocks.CreateFlex(info, title, isScheduled)
			if err != nil {
				return fmt.Errorf("failed to create flex for %s: %w", title, err)
			}

			flex.AddItem(machineFlex, 0, 1, false)
			return nil
		}

		leftFlex := tview.NewFlex().SetDirection(tview.FlexRow)
		rightFlex := tview.NewFlex().SetDirection(tview.FlexRow)

		if err := getAndDisplayFlex(sherlocks.SherlocksURL, sherlocks.ActiveSherlocksTitle, false, leftFlex); err != nil {
			log.Fatal(err)
		}

		if err := getAndDisplayFlex(sherlocks.RetiredSherlocksURL, sherlocks.RetiredSherlocksTitle, false, leftFlex); err != nil {
			log.Fatal(err)
		}

		if err := getAndDisplayFlex(sherlocks.ScheduledSherlocksURL, sherlocks.ScheduledSherlocksTitle, true, rightFlex); err != nil {
			log.Fatal(err)
		}

		rightFlex.AddItem(tview.NewTextView().SetText("").SetDynamicColors(true), 0, 0, false)

		mainFlex := tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(leftFlex, 0, 3, false).
			AddItem(rightFlex, 0, 1, false)

		if err := app.SetRoot(mainFlex, true).Run(); err != nil {
			panic(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(sherlocksCmd)
	sherlocksCmd.Flags().StringP("sherlock_name", "s", "", "Sherlock Name")
	sherlocksCmd.Flags().StringP("download", "d", "", "Download Sherlock Resources")
	sherlocksCmd.Flags().IntP("task", "t", 0, "Task ID")
	sherlocksCmd.Flags().BoolP("hint", "", false, "Hint")
}
