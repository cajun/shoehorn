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

	switch args[2] {
	case "start":
		command.Start()
	case "stop":
		command.Stop()
	case "restart":
		command.Restart()
	case "kill":
		command.Kill()
	case "bash":
		command.Bash()
	case "console":
		command.Console()
	case "params":
		command.PrintParams()
	default:
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
