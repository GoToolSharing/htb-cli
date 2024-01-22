package ssh

import (
	"fmt"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"
)

func Connect(username, password, host string, port int) error {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	connection, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), config)
	if err != nil {
		return fmt.Errorf("Connection error: %s\n", err)
	}
	defer connection.Close()
	fmt.Println("Connection established")

	session, err := connection.NewSession()
	if err != nil {
		return fmt.Errorf("Session creation error: %s\n", err)
	}
	defer session.Close()

	cmd := "cat /etc/passwd | grep -E '/home|/users' | cut -d: -f6"
	output, err := session.CombinedOutput(cmd)
	if err != nil {
		return fmt.Errorf("Error during command execution: %s\n", err)
	}
	homes := strings.Split(string(output), "\n")
	fmt.Println("Homes :", homes)

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
		fileSession.Close()

		if strings.TrimSpace(string(fileOutput)) == "found" {
			fileFound = true
			fmt.Printf("File found: %s\n", filePath)

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
			if len(contentOutput) == 32 {
				fmt.Println("HTB flag detected")
				// TODO: auto Submission
			}
			contentSession.Close()
			break
		}
	}

	if !fileFound {
		fmt.Println("user.txt file not found in home directories")
	}
	return nil
}
