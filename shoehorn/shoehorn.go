package main

import (
	"flag"
	"fmt"
	"github.com/cajun/shoehorn/command"
	"github.com/cajun/shoehorn/config"
	"os"
)

// handleParam takes in the given parameters and decides what to do with them.
func handleParam(args []string) {
	name := args[1]
	if name == "list" {
		config.PrintProcesses()
	} else if config.Process(name) != nil && len(args) == 2 {
		config.PrintConfig(name)
	} else if config.Process(name) != nil && len(args) > 2 {
		handleCommand(args)
	} else {
		fmt.Printf("Process Name: (%v) doesn't exists\n", name)
	}
}

// handleCommand will take in the argument for the process and run it
func handleCommand(args []string) {
	command.SetProcess(args[1])
	command.SetConfig(config.Process(args[1]))

	if args[2] == "start" {
		command.Start()
	} else if args[2] == "stop" {
		command.Stop()
	} else if args[2] == "restart" {
		command.Restart()
	} else if args[2] == "kill" {
		command.Kill()
	} else if args[2] == "bash" {
		command.Bash()
	} else if args[2] == "console" {
		command.Console()
	} else if args[2] == "params" {
		command.PrintParams()
	} else {
		fmt.Printf("Running Command: (%v) doesn't exists\n", args[2])
	}

}

// main function pulls in the config and flags.
// then passes off the commands to the handleParams method
func main() {
	config.LoadConfigs()
	flag.Parse()
	args := os.Args
	if len(args) > 1 {
		handleParam(args)
	} else {
		flag.PrintDefaults()
	}
}
