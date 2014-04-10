package main

import (
  "flag"
  "fmt"
  "github.com/cajun/shoehorn/command"
  "github.com/cajun/shoehorn/config"
  "github.com/cajun/shoehorn/logger"
  "github.com/cajun/shoehorn/server"
  "os"
  "os/signal"
)

var wait = false

func init() {
  flag.BoolVar(&wait, "wait", false, "wait for process signal")
}

// handleParam takes in the given parameters and decides what to do with them.
func handleParam(args []string) {
  name := args[0]

  if name == "list" {
    config.PrintProcesses()
  } else if name == "init" {
    config.Init()
  } else if command.IsCommand(name) {
    for _, process := range config.List() {
      command.ParseCommand([]string{process, name})
    }
  } else if config.Process(name) != nil && len(args) == 1 {
    config.PrintConfig(name)
  } else if config.Process(name) != nil && len(args) > 1 {
    command.ParseCommand(args)
  } else {
    logger.Log(fmt.Sprintf("Process Name: (%v) doesn't exists\n", name))
  }
}

func doIt(args []string) {
  defer logger.Done()
  if server.On() {
    server.Up()
  } else if len(args) >= 1 {
    handleParam(args)
  } else {
    flag.PrintDefaults()
    command.PrintCommands()
  }
}

// main function pulls in the config and flags.
// then passes off the commands to the handleParams method
func main() {
  flag.Parse()

  os.Chdir(command.Root())
  config.LoadConfigs()

  args := flag.Args()

  pipe := logger.New(os.Stdout)

  go doIt(args)

  result := logger.InitStatus()
  for !result.Done {
    select {
    case result = <-pipe:
      logger.Write(result)
    }
  }

  if wait {
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, os.Kill)
    <-c
  }

}
