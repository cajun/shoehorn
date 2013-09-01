package main

import (
	"flag"
	"fmt"
	"github.com/cajun/shoehorn/config"
	"os"
)

// handleParam takes in the given parameters and decides what to do with them
func handleParam(command string) {
	if command == "list" {
		config.PrintApps()
	} else if config.App(command) != nil {
		config.PrintConfig(command)
	} else {
		fmt.Println("Command (" + command + ") doesn't exists")
	}
}

// main function pulls in the config and flags.  then passes off the commands to
// the handleParams method
func main() {
	config.LoadConfigs()
	flag.Parse()
	args := os.Args
	if len(args) > 1 {
		handleParam(args[1])
	} else {
		flag.PrintDefaults()
	}
}
