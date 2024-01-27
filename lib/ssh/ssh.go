package ssh

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/GoToolSharing/htb-cli/config"
	"github.com/GoToolSharing/htb-cli/lib/utils"
	"golang.org/x/crypto/ssh"
)

func Connect(username, password, host string, port int) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	connection, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), config)
	if err != nil {
		return nil, fmt.Errorf("Connection error: %s\n", err)
	}
	fmt.Println("SSH connection established")

	return connection, nil
}

func GetUserFlag(connection *ssh.Client) (string, error) {
	session, err := connection.NewSession()
	if err != nil {
		return "", fmt.Errorf("Session creation error: %s\n", err)
	}
	defer session.Close()

	cmd := "cat /etc/passwd | grep -E '/home|/users' | cut -d: -f6"
	output, err := session.CombinedOutput(cmd)
	if err != nil {
		return "", fmt.Errorf("Error during command execution: %s\n", err)
	}
	homes := strings.Split(string(output), "\n")
	config.GlobalConfig.Logger.Debug(fmt.Sprintf("Users homes : %v", homes))

	fileFound := false
	for _, home := range homes {
		if home == "" {
			continue
		}
		filePath := filepath.Join(home, "user.txt")
		fileSession, err := connection.NewSession()
		if err != nil {
			fmt.Printf("Error creating file session: %s\n", err)
			continue
		}
		cmd := fmt.Sprintf("if [ -f %s ]; then echo found; else echo not found; fi", filePath)
		fileOutput, err := fileSession.CombinedOutput(cmd)
		if err != nil {
			return "", err
		}
		fileSession.Close()

		if strings.TrimSpace(string(fileOutput)) == "found" {
			fileFound = true
			config.GlobalConfig.Logger.Debug(fmt.Sprintf("User flag found: %s\n", filePath))

			contentSession, err := connection.NewSession()
			if err != nil {
				fmt.Printf("Error creating content session: %s\n", err)
				continue
			}
			contentCmd := fmt.Sprintf("cat %s", filePath)
			contentOutput, err := contentSession.CombinedOutput(contentCmd)
			if err != nil {
				fmt.Printf("File read error %s: %s\n", filePath, err)
				permSession, err := connection.NewSession()
				if err != nil {
					fmt.Printf("Error creating permissions session: %s\n", err)
					continue
				}
				permCmd := fmt.Sprintf("ls -la %s", filePath)
				permOutput, err := permSession.CombinedOutput(permCmd)
				if err != nil {
					fmt.Printf("Error obtaining permissions: %s\n", err)
				} else {
					fmt.Printf("Permissions required: %s\n", string(permOutput))
				}
				permSession.Close()
				continue
			}
			fmt.Printf("%s: %s\n", filePath, string(contentOutput))
			if len(contentOutput) == 32 || len(contentOutput) == 33 {
				config.GlobalConfig.Logger.Info("HTB flag detected")
				return (string(contentOutput)), nil
			}
			contentSession.Close()
			break
		}
	}

	if !fileFound {
		fmt.Println("user.txt file not found in home directories")
	}
	return "", nil
}

func GetHostname(connection *ssh.Client) (string, error) {
	hostnameSession, err := connection.NewSession()
	if err != nil {
		return "", fmt.Errorf("Error creating hostname session: %s\n", err)
	}
	cmd := "hostname"
	sessionOutput, err := hostnameSession.CombinedOutput(cmd)
	if err != nil {
		return "", err
	}
	hostnameSession.Close()
	hostname := strings.ReplaceAll(string(sessionOutput), "\n", "")
	config.GlobalConfig.Logger.Debug(fmt.Sprintf("Hotname: %s", hostname))
	return hostname, nil
}

func BuildSubmitStuff(hostname string, userFlag string) (string, map[string]string, error) {
	// Can be release arena or machine
	var payload map[string]string
	var url string

	machineID, err := utils.SearchItemIDByName(hostname, "Machine")
	if err != nil {
		return "", nil, err
	}
	machineType, err := utils.GetMachineType(machineID)
	if err != nil {
		return "", nil, err
	}
	config.GlobalConfig.Logger.Debug(fmt.Sprintf("Machine Type: %s", machineType))

	if machineType == "release" {
		url = config.BaseHackTheBoxAPIURL + "/arena/own"
		payload = map[string]string{
			"flag": userFlag,
		}
	} else {
		url = config.BaseHackTheBoxAPIURL + "/machine/own"
		payload = map[string]string{
			"id":   machineID,
			"flag": userFlag,
		}
	}

	return url, payload, nil
}
