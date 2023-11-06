package main

import (
	"fmt"
	"os"

	"github.com/GoToolSharing/htb-cli/cmd"
)

func main() {
	if _, err := os.Stat(cmd.BaseDirectory); os.IsNotExist(err) {
		fmt.Printf("[INFO] - The \"%s\" folder does not exist, creation in progress...\n", cmd.BaseDirectory)
		err := os.MkdirAll(cmd.BaseDirectory, os.ModePerm)
		if err != nil {
			fmt.Printf("Folder creation error: %s\n", err)
			return
		}

		fmt.Printf("[INFO] - \"%s\" folder created successfully\n\n", cmd.BaseDirectory)
	}
	cmd.Execute()
}

///api/v4/access/ovpnfile/18/0

// remote edge-eu-vip-7.hackthebox.eu 1337
