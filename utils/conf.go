package utils

import (
	"fmt"
	"os"
)

func GetHTBToken() string {
	var envName = "HTB_TOKEN"
	if os.Getenv(envName) == "" {
		fmt.Printf("Environment variable is not set : %v", envName)
		os.Exit(84)
	}
	return os.Getenv("HTB_TOKEN")
}
