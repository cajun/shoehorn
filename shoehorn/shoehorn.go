package main

import (
	"flag"
	"fmt"
	"github.com/cajun/shoehorn/command"
	"github.com/cajun/shoehorn/config"
)

func isCommand(cmd string) bool {
	for _, val := range []string{"start", "stop", "kill", "restart"} {
		if val == cmd {
			return true
		}
	}

	return false
}

// handleParam takes in the given parameters and decides what to do with them.
func handleParam(args []string) {
	name := args[0]

	if name == "list" {
		config.PrintProcesses()
	} else if isCommand(name) {
		for _, process := range config.List() {
			handleCommand([]string{process, name})
		}
	} else if config.Process(name) != nil && len(args) == 1 {
		config.PrintConfig(name)
	} else if config.Process(name) != nil && len(args) > 1 {
		handleCommand(args)
	} else {
		fmt.Printf("Process Name: (%v) doesn't exists\n", name)
	}
}

// handleCommand will take in the argument for the process and run it
func handleCommand(args []string) {
	command.SetProcess(args[0])
	command.SetConfig(config.Process(args[0]))

	switch args[1] {
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
	flag.Parse()
	config.LoadConfigs()

	args := flag.Args()

	if len(args) >= 1 {
		handleParam(args)
	} else {
		flag.PrintDefaults()
	}

}
