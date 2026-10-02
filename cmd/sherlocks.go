package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/lib/sherlocks"
	"github.com/GoToolSharing/htb-cli/lib/utils"
	"github.com/chzyer/readline"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var sherlocksCmd = &cobra.Command{
	Use:   "sherlocks",
	Short: "Play Sherlocks mode (blue team)",
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

		sherlockFlag, err := cmd.Flags().GetString("flag")
		if err != nil {
			fmt.Println(err)
			return
		}

		// If no name is provided, show the interactive UI
		if sherlockNameParam == "" {
			app := tview.NewApplication()

			getAndDisplayFlex := func(url, title string, isScheduled bool, flex *tview.Flex) error {
				resp, err := utils.HtbRequest(http.MethodGet, url, nil)
				if err != nil {
					return fmt.Errorf("failed to get data from %s: %w", url, err)
				}
				defer resp.Body.Close()

				info := utils.ParseJsonMessage(resp, "data")

				sherlockFlex, err := sherlocks.CreateFlex(info, title, isScheduled)
				if err != nil {
					return fmt.Errorf("failed to create flex for %s: %w", title, err)
				}

				flex.AddItem(sherlockFlex, 0, 1, false)
				return nil
			}

			leftFlex := tview.NewFlex().SetDirection(tview.FlexRow)
			rightFlex := tview.NewFlex().SetDirection(tview.FlexRow)

			if err := getAndDisplayFlex(sherlocks.SherlocksURL, sherlocks.ActiveSherlocksTitle, false, leftFlex); err != nil {
				config.GlobalConfig.Logger.Error("", zap.Error(err))
				os.Exit(1)
			}

			if err := getAndDisplayFlex(sherlocks.RetiredSherlocksURL, sherlocks.RetiredSherlocksTitle, false, leftFlex); err != nil {
				config.GlobalConfig.Logger.Error("", zap.Error(err))
				os.Exit(1)
			}

			if err := getAndDisplayFlex(sherlocks.ScheduledSherlocksURL, sherlocks.ScheduledSherlocksTitle, true, rightFlex); err != nil {
				config.GlobalConfig.Logger.Error("", zap.Error(err))
				os.Exit(1)
			}

			rightFlex.AddItem(tview.NewTextView().SetText("").SetDynamicColors(true), 0, 0, false)

			mainFlex := tview.NewFlex().SetDirection(tview.FlexColumn).
				AddItem(leftFlex, 0, 3, false).
				AddItem(rightFlex, 0, 1, false)

			if err := app.SetRoot(mainFlex, true).Run(); err != nil {
				config.GlobalConfig.Logger.Error("", zap.Error(err))
				os.Exit(1)
			}

			return
		}

		sherlockID, err := sherlocks.SearchIDByName(sherlockNameParam)
		if err != nil {
			fmt.Printf("Error finding Sherlock: %v\n", err)
			return
		}
		config.GlobalConfig.Logger.Debug(fmt.Sprintf("SherlockID: %s", sherlockID))

		if sherlockTaskID != 0 && sherlockFlag == "" {
			err := sherlocks.GetTaskByID(sherlockID, sherlockTaskID, sherlockHint)
			if err != nil {
				fmt.Println(err)
				return
			}
			return
		}

		if sherlockTaskID != 0 {
			url := fmt.Sprintf("%s/sherlocks/%s/tasks", config.BaseHackTheBoxAPIURL, sherlockID)
			resp, err := utils.HtbRequest(http.MethodGet, url, nil)
			if err != nil {
				fmt.Printf("Error fetching tasks: %v\n", err)
				return
			}
			defer resp.Body.Close()

			jsonData, _ := io.ReadAll(resp.Body)
			var sherlockData sherlocks.SherlockDataTasks
			err = json.Unmarshal([]byte(jsonData), &sherlockData)
			if err != nil {
				fmt.Printf("Error parsing JSON: %v\n", err)
				return
			}

			if sherlockTaskID < 1 || sherlockTaskID > len(sherlockData.Tasks) {
				fmt.Printf("Invalid task ID: %d. Valid range: 1-%d\n", sherlockTaskID, len(sherlockData.Tasks))
				return
			}

			actualTaskID := sherlockData.Tasks[sherlockTaskID-1].ID
			taskIDStr := strconv.Itoa(actualTaskID)

			if sherlockFlag != "" {
				config.GlobalConfig.Logger.Debug(fmt.Sprintf("Submitting flag for task %d: %s", actualTaskID, sherlockFlag))

				message, err := sherlocks.SubmitTask(sherlockID, taskIDStr, sherlockFlag)
				if err != nil {
					fmt.Printf("Error submitting flag: %v\n", err)
					return
				}

				fmt.Println(message)
				return
			}

			if sherlockHint && sherlockData.Tasks[sherlockTaskID-1].Hint != "" {
				fmt.Printf("\n%s :\n%s\n\nHint : %s\nMasked Flag : %s\n",
					sherlockData.Tasks[sherlockTaskID-1].Title,
					sherlockData.Tasks[sherlockTaskID-1].Description,
					sherlockData.Tasks[sherlockTaskID-1].Hint,
					sherlockData.Tasks[sherlockTaskID-1].MaskedFlag)
			} else {
				fmt.Printf("\n%s :\n%s\n\nMasked Flag : %s\n",
					sherlockData.Tasks[sherlockTaskID-1].Title,
					sherlockData.Tasks[sherlockTaskID-1].Description,
					sherlockData.Tasks[sherlockTaskID-1].MaskedFlag)
			}

			rl, err := readline.New("Answer: ")
			if err != nil {
				panic(err)
			}
			defer rl.Close()

			flag, err := rl.Readline()
			if err != nil {
				fmt.Printf("Error reading input: %v\n", err)
				return
			}
			flag = strings.TrimSpace(flag)
			config.GlobalConfig.Logger.Debug(fmt.Sprintf("Flag: %s", flag))

			message, err := sherlocks.SubmitTask(sherlockID, taskIDStr, flag)
			if err != nil {
				fmt.Printf("Error submitting flag: %v\n", err)
				return
			}

			fmt.Println(message)
			return
		}

		if sherlockDownloadPath != "" {
			err = sherlocks.GetGeneralInformations(sherlockID, sherlockDownloadPath)
			if err != nil {
				fmt.Println(err)
				return
			}
			return
		}

		data, err := sherlocks.GetTasks(sherlockID)
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, task := range data.Tasks {
			status := ""
			if task.Completed {
				status = " (DONE)"
			}

			fmt.Printf("\n%s%s :\n%s\n", task.Title, status, task.Description)

			if task.MaskedFlag != "" {
				fmt.Printf("Format : %s\n", task.MaskedFlag)
			}
			fmt.Println()
		}
	},
}

func init() {
	rootCmd.AddCommand(sherlocksCmd)
	sherlocksCmd.Flags().StringP("sherlock_name", "s", "", "Sherlock Name")
	sherlocksCmd.Flags().StringP("download", "d", "", "Download Sherlock Resources")
	sherlocksCmd.Flags().IntP("task", "t", 0, "Task ID")
	sherlocksCmd.Flags().StringP("flag", "f", "", "Task Flag (optional - if provided, submits without prompting)")
	sherlocksCmd.Flags().BoolP("hint", "", false, "Show hint for the task")
}
