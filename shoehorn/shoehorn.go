package main

import (
	"flag"
	"fmt"
	"github.com/cajun/shoehorn/command"
	"github.com/cajun/shoehorn/config"
	"github.com/cajun/shoehorn/server"
)

// handleParam takes in the given parameters and decides what to do with them.
func handleParam(args []string) {
	name := args[0]

	if name == "list" {
		config.PrintProcesses()
	} else if command.IsCommand(name) {
		for _, process := range config.List() {
			command.ParseCommand([]string{process, name})
		}
	} else if config.Process(name) != nil && len(args) == 1 {
		config.PrintConfig(name)
	} else if config.Process(name) != nil && len(args) > 1 {
		command.ParseCommand(args)
	} else {
		fmt.Printf("Process Name: (%v) doesn't exists\n", name)
	}
}

// main function pulls in the config and flags.
// then passes off the commands to the handleParams method
func main() {
	flag.Parse()
	config.LoadConfigs()

	args := flag.Args()

	if server.On() {
		server.Up()
	} else if len(args) >= 1 {
		handleParam(args)
	} else {
		flag.PrintDefaults()
		command.PrintCommands()
	}
}
