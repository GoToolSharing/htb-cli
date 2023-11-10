package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GoToolSharing/htb-cli/utils"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

var sherlockNameParam string

const (
	sherlocksURL            = baseAPIURL + "/sherlocks?state=active"
	retiredSherlocksURL     = baseAPIURL + "/sherlocks?state=retired"
	scheduledSherlocksURL   = baseAPIURL + "/sherlocks?state=unreleased"
	activeSherlocksTitle    = "Active"
	retiredSherlocksTitle   = "Retired"
	scheduledSherlocksTitle = "Scheduled"
	SherlocksCheckMark      = "\U00002705"
	SherlocksCrossMark      = "\U0000274C"
	SPenguin                = "\U0001F427"
	SComputer               = "\U0001F5A5 "
)

// sgetColorFromDifficultyText returns the color corresponding to the given difficulty.
func sgetColorFromDifficultyText(difficultyText string) string {
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

// sgetOSEmoji returns an emoji corresponding to the given operating system
func sgetOSEmoji(os string) string {
	switch os {
	case "Linux":
		return SPenguin
	case "Windows":
		return SComputer
	default:
		return ""
	}
}

// screateFlex creates and returns a Flex view with machine information
func screateFlex(info interface{}, title string, isScheduled bool) (*tview.Flex, error) {
	log.Println("Info :", info)
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.SetBorder(true).SetTitle(title).SetTitleAlign(tview.AlignLeft)

	for _, value := range info.([]interface{}) {
		data := value.(map[string]interface{})

		// Determining the color according to difficulty

		key := "Undefined"
		if title == "Scheduled" {
			key = data["difficulty"].(string)
		}
		color := sgetColorFromDifficultyText(key)
		// osEmoji := sgetOSEmoji(data["os"].(string))

		var formatString string

		// Choice of display format depending on the nature of the information
		if isScheduled {
			formatString = fmt.Sprintf("%-15s %s%-10s[-]",
				data["name"], color, data["difficulty"])
		} else {
			// Convert and format date
			// parsedDate, err := time.Parse(time.RFC3339Nano, data["release"].(string))
			// if err != nil {
			// 	return nil, fmt.Errorf("error parsing date: %v", err)
			// }
			// formattedDate := parsedDate.Format("02 January 2006")

			// userEmoji := SherlocksCrossMark + "User"
			// if value, ok := data["authUserInUserOwns"]; ok && value != nil {
			// 	if value.(bool) {
			// 		userEmoji = SherlocksCheckMark + "User"
			// 	}
			// }

			// rootEmoji := SherlocksCrossMark + "Root"
			// if value, ok := data["authUserInRootOwns"]; ok && value != nil {
			// 	if value.(bool) {
			// 		rootEmoji = SherlocksCheckMark + "Root"
			// 	}
			// }

			// formatString = fmt.Sprintf("%-15s %s%-10s[-] %-5v %-5v %-7v %-30s",
			// 	data["name"], color, data["difficultyText"],
			// 	data["star"], userEmoji, rootEmoji, formattedDate)
		}

		flex.AddItem(tview.NewTextView().SetText(formatString).SetDynamicColors(true), 1, 0, false)
	}

	return flex, nil
}

var sherlocksCmd = &cobra.Command{
	Use:   "sherlocks",
	Short: "Displays active sherlocks and next sherlocks to be released",
	Run: func(cmd *cobra.Command, args []string) {
		if sherlockNameParam != "" {
			sherlockID, err := utils.SearchSherlockIDByName(proxyParam, sherlockNameParam)
			if err != nil {
				fmt.Println(err)
				return
			}
			log.Println("SherlockID :", sherlockID)

			utils.GetSherlockGeneralInformations(proxyParam, sherlockID)

			ret := false

			for !ret {
				ret, err = utils.GetSherlockTasks(proxyParam, sherlockID)
				if err != nil {
					fmt.Println(err)
					return
				}
			}

			fmt.Println("No tasks left")

			return
		}
		app := tview.NewApplication()

		getAndDisplayFlex := func(url, title string, isScheduled bool, flex *tview.Flex) error {
			resp, err := utils.HtbRequest(http.MethodGet, url, proxyParam, nil)
			if err != nil {
				return fmt.Errorf("failed to get data from %s: %w", url, err)
			}

			info := utils.ParseJsonMessage(resp, "data")

			machineFlex, err := screateFlex(info, title, isScheduled)
			if err != nil {
				return fmt.Errorf("failed to create flex for %s: %w", title, err)
			}

			flex.AddItem(machineFlex, 0, 1, false)
			return nil
		}

		leftFlex := tview.NewFlex().SetDirection(tview.FlexRow)
		rightFlex := tview.NewFlex().SetDirection(tview.FlexRow)

		if err := getAndDisplayFlex(sherlocksURL, activeSherlocksTitle, false, leftFlex); err != nil {
			log.Fatal(err)
		}

		if err := getAndDisplayFlex(retiredSherlocksURL, retiredSherlocksTitle, false, leftFlex); err != nil {
			log.Fatal(err)
		}

		if err := getAndDisplayFlex(scheduledSherlocksURL, scheduledSherlocksTitle, true, rightFlex); err != nil {
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
	sherlocksCmd.Flags().StringVarP(&sherlockNameParam, "sherlock_name", "s", "", "Sherlock Name")
}
