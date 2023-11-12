package main

import (
	"fmt"
	"io"
	"log"

	"github.com/GoToolSharing/htb-cli/cmd"
	"github.com/GoToolSharing/htb-cli/config"
)

func main() {
	log.SetOutput(io.Discard)
	err := config.Init()
	if err != nil {
		fmt.Println(err)
		return
	}
	cmd.Execute()
}
