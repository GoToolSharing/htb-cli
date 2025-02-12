package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/lib/cache"
	"github.com/GoToolSharing/htb-cli/lib/machines"
	"github.com/GoToolSharing/htb-cli/lib/utils"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	machineURL     = config.BaseHackTheBoxAPIURL + "/machine/paginated/?per_page=20"
	retiredURL     = config.BaseHackTheBoxAPIURL + "/machine/list/retired/paginated/?per_page=20&sort_by=release-date"
	scheduledURL   = config.BaseHackTheBoxAPIURL + "/machine/unreleased/"
	activeTitle    = "Active"
	retiredTitle   = "Retired"
	scheduledTitle = "Scheduled"
	CheckMark      = "\U00002705"
	CrossMark      = "\U0000274C"
	Penguin        = "\U0001F427"
	Computer       = "\U0001F5A5 "
)

// getColorFromDifficultyText returns the color corresponding to the given difficulty.
func getColorFromDifficultyText(difficultyText string) string {
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

// getOSEmoji returns an emoji corresponding to the given operating system
func getOSEmoji(os string) string {
	switch os {
	case "Linux":
		return Penguin
	case "Windows":
		return Computer
	default:
		return ""
	}
}

// createFlex creates and returns a Flex view with machine information
func createFlex(info interface{}, title string, isScheduled bool) (*tview.Flex, error) {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.SetBorder(true).SetTitle(title).SetTitleAlign(tview.AlignLeft)

	for _, value := range info.([]interface{}) {
		data := value.(map[string]interface{})

		// Determining the color according to difficulty

		key := "Undefined"
		_ = key
		if title == "Scheduled" {
			key = data["difficulty_text"].(string)
		} else {
			key = data["difficultyText"].(string)
		}
		color := getColorFromDifficultyText(key)
		osEmoji := getOSEmoji(data["os"].(string))

		var formatString string

		// Choice of display format depending on the nature of the information
		if isScheduled {
			formatString = fmt.Sprintf("%-10s %s%-10s %s%-10s[-]",
				data["name"], osEmoji, data["os"], color, data["difficulty_text"])
		} else {
			// Convert and format date
			parsedDate, err := time.Parse(time.RFC3339Nano, data["release"].(string))
			if err != nil {
				return nil, fmt.Errorf("error parsing date: %v", err)
			}
			formattedDate := parsedDate.Format("02 January 2006")

			userEmoji := CrossMark + "User"
			if value, ok := data["authUserInUserOwns"]; ok && value != nil {
				if value.(bool) {
					userEmoji = CheckMark + "User"
				}
			}

			rootEmoji := CrossMark + "Root"
			if value, ok := data["authUserInRootOwns"]; ok && value != nil {
				if value.(bool) {
					rootEmoji = CheckMark + "Root"
				}
			}

			formatString = fmt.Sprintf("%-15s %s%-10s %s%-10s[-] %-5v %-5v %-7v %-30s",
				data["name"], osEmoji, data["os"], color, data["difficultyText"],
				data["star"], userEmoji, rootEmoji, formattedDate)
		}

		flex.AddItem(tview.NewTextView().SetText(formatString).SetDynamicColors(true), 1, 0, false)
	}

	return flex, nil
}

var machinesCmd = &cobra.Command{
	Use:   "machines",
	Short: "Displays active / retired machines and next machines to be released",
	Run: func(cmd *cobra.Command, args []string) {
		app := tview.NewApplication()

		refreshParam, err := cmd.Flags().GetBool("refresh")
		if err != nil {
			fmt.Println(err)
			return
		}

		db, err := sql.Open("sqlite3", config.BaseDirectory+config.Database)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		updateCache, err := cache.CheckCacheDate(db)
		if err != nil {
			log.Fatal(err)
		}

		if updateCache || (refreshParam && utils.AskConfirmation("The cache data will be deleted, do you want to continue ?")) {
			_, err = db.Exec("DELETE FROM machines")
			if err != nil {
				log.Fatalf("Error deleting existing data: %v", err)
			}
			config.GlobalConfig.Logger.Info("Machine cache data has been deleted")
		}

		var info interface{}

		getAndDisplayFlex := func(url, title string, isScheduled bool, flex *tview.Flex) error {
			if refreshParam || updateCache {
				config.GlobalConfig.Logger.Info("Machine data recovery via API")
				resp, err := utils.HtbRequest(http.MethodGet, url, nil)
				if err != nil {
					return fmt.Errorf("failed to get data from %s: %w", url, err)
				}

				info = utils.ParseJsonMessage(resp, "data")

				machines.InsertMachines(db, info, title)
			} else {
				config.GlobalConfig.Logger.Info("Machine data recovery via cache")
				info, err = machines.GetMachinesFromCache(db, title)

				if err != nil {
					log.Fatal(err)
				}
			}

			machineFlex, err := createFlex(info, title, isScheduled)
			if err != nil {
				return fmt.Errorf("failed to create flex for %s: %w", title, err)
			}

			flex.AddItem(machineFlex, 0, 1, false)
			return nil
		}

		leftFlex := tview.NewFlex().SetDirection(tview.FlexRow)
		rightFlex := tview.NewFlex().SetDirection(tview.FlexRow)

		if err := getAndDisplayFlex(machineURL, activeTitle, false, leftFlex); err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}

		if err := getAndDisplayFlex(retiredURL, retiredTitle, false, leftFlex); err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}

		if err := getAndDisplayFlex(scheduledURL, scheduledTitle, true, rightFlex); err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}

		err = cache.UpdateCacheDate(db, "machines_cache_date")

		if err != nil {
			log.Fatal(err)
		}

		rightFlex.AddItem(tview.NewTextView().SetText("").SetDynamicColors(true), 0, 0, false)

		mainFlex := tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(leftFlex, 0, 3, false).
			AddItem(rightFlex, 0, 1, false)

		if err := app.SetRoot(mainFlex, true).Run(); err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(machinesCmd)
	machinesCmd.Flags().BoolP("refresh", "r", false, "Refresh cache")
}
