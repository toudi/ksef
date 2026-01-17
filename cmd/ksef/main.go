package main

import (
	"fmt"
	"ksef/cmd/ksef/commands"
	"ksef/internal/logging"
	"os"
)

func main() {
	defer logging.FinishLogging()
	if err := commands.RootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
