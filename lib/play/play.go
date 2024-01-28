package play

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/lib/utils"
	"go.uber.org/zap"
)

func setTmuxSessionName(ip string) error {
	if _, ok := os.LookupEnv("TMUX"); ok {
		cmd := exec.Command("tmux", "rename-session", ip)

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("erreur lors du changement du nom de la session tmux : %w", err)
		}
		fmt.Println("Nom de la session tmux changé avec succès.")
	} else {
		fmt.Println("Pas dans une session tmux.")
	}
	return nil
}

func Configure(releaseParam bool) {
	var machineID string
	var seasonal string
	var host string
	var machineName string
	var machineIP string

	if releaseParam {
		seasonal = "true"
		machineID, err := utils.SearchLastReleaseArenaMachine()
		config.GlobalConfig.Logger.Debug(fmt.Sprintf("Machine ID : %s", machineID))
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}
		machineName, err = utils.GetActiveReleaseArenaMachineName()
		if err != nil {
			config.GlobalConfig.Logger.Error("", zap.Error(err))
			os.Exit(1)
		}
		// output, err := start.CoreStartCmd("", machineID)
		// if err != nil {
		// 	config.GlobalConfig.Logger.Error("", zap.Error(err))
		// 	os.Exit(1)
		// }

		// machineIP, err = utils.GetActiveReleaseArenaMachineIP()
		// if err != nil {
		// 	config.GlobalConfig.Logger.Error("", zap.Error(err))
		// 	os.Exit(1)
		// }

		machineIP = "10.10.11.248"

		host = fmt.Sprintf("%s.htb", strings.ToLower(machineName))

		fmt.Println("Machine ID :", machineID)
		fmt.Println("Machine IP :", machineIP)
		fmt.Println("Machine Name :", machineName)
		fmt.Println("Host :", host)
		fmt.Println("Seasonal :", seasonal)

		err = createConfigFile(machineName, machineIP, machineID, seasonal, host)
		if err != nil {
			fmt.Println("Error creating config file:", err)
		}

		// fmt.Println(output)

	}

	if err := setTmuxSessionName(machineIP); err != nil {
		fmt.Println(err)
	}

	// TODO: add a line to the rc file !

	// fmt.Println("")

	err := createProfileFile(machineName, machineIP, machineID, seasonal, host)
	if err != nil {
		fmt.Println("Error creating config file:", err)
	}
}

func createConfigFile(name, ip, machineID, seasonal, host string) error {
	settingsPath := filepath.Join(config.BaseDirectory, "settings.conf")

	content := []byte(fmt.Sprintf("[Default]\nname=%s\nip=%s\nmachineID=%s\nseasonal=%s\n", name, ip, machineID, seasonal))
	err := os.WriteFile(settingsPath, content, 0644)
	if err != nil {
		return err
	}
	return nil
}

func createProfileFile(name, ip, machineID, seasonal, host string) error {
	profilePath := filepath.Join(config.BaseDirectory, "profile")

	// content := []byte(fmt.Sprintf("export PS1=\"htb-cli - [%s] > \"", ip))
	ipEnv := []byte(fmt.Sprintf("export TARGET=\"%s\"", ip))
	err := os.WriteFile(profilePath, ipEnv, 0644)
	if err != nil {
		return err
	}

	fmt.Println("Use the following command to update your profile.")
	fmt.Println("source ~/.local/htb-cli/profile")

	return nil
}
