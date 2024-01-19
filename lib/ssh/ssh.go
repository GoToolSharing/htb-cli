package ssh

import (
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"
)

func Connect(username, password, host string) {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	cmd := "ls /home"
	output, err := session.CombinedOutput(cmd)
	if err != nil {
		panic(err)
	}

	users := strings.Split(string(output), "\n")
	for _, user := range users {
		if user == "" {
			continue
		}
		checkAndReadFile(client, "/home/"+user+"/user.txt")
	}
}
func checkAndReadFile(client *ssh.Client, filePath string) {
	session, err := client.NewSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	cmd := fmt.Sprintf(`if [ -f %s ]; then 
                            if [ -r %s ]; then 
                                cat %s; 
                            else 
                                echo "File '%s' exists, but no read permission"; 
                            fi; 
                        else 
                            echo "File '%s' does not exist"; 
                        fi`, filePath, filePath, filePath, filePath, filePath)

	output, err := session.CombinedOutput(cmd)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(output))
}
