package command

import (
	"bufio"
	"fmt"
	"github.com/cajun/shoehorn/config"
	"github.com/cajun/shoehorn/logger"
	"os"
)

var (
	cfg     *config.Settings
	process string
)

// init Creates the nessary directories that
// will be required for the commands to work.
// The command will expect that there are the
// following directories.
//
// tmp/pids # This is where the pids are stored
// log      # where any logs files should go
// config   # This is the location of the nginx config
//
func init() {
	MkDirs()
}

func MkDirs() {
	os.MkdirAll("tmp/pids", os.ModeDir|0700)
	os.MkdirAll("log", os.ModeDir|0700)
	os.MkdirAll("config", os.ModeDir|0700)
}

type Runner func(...string)

type Executor struct {
	description string
	run         Runner
}

// SetConfig will set an var for the settings for the given process
// that will be executing.
func SetConfig(settings *config.Settings) {
	cfg = settings
}

// SetProcess is setting the name of the process that this command
// will be executing against. The pid will have the process name in
// it.
func SetProcess(proc string) {
	process = proc
}

// PrintParams will print all of the settings that will be passed.
// into docker It assumes the first instance                     .
func PrintParams(args ...string) {
	logger.Log(fmt.Sprintln(settingsToParams(0, false)))
}

// PrintCommands will list out all of the commands to the end user.
func PrintCommands() {
	logger.Log(fmt.Sprintln("** Daemonized Commands **"))
	for cmd, desc := range DaemonizedCommands() {
		logger.Log(fmt.Sprintf("%15s: %s\n", cmd, desc.description))
	}

	logger.Log(fmt.Sprintln("** Information Commands **"))
	for cmd, desc := range InfoCommands() {
		logger.Log(fmt.Sprintf("%15s: %s\n", cmd, desc.description))
	}

	logger.Log(fmt.Sprintln("** Interactive Commands **"))
	for cmd, desc := range InteractiveCommands() {
		logger.Log(fmt.Sprintf("%15s: %s\n", cmd, desc.description))
	}
}

// isCommand will take in a given command and check to see if it is available to
// be executed against ALL of the commands at one time.
func IsCommand(cmd string) bool {
	for val := range DaemonizedCommands() {
		if val == cmd {
			return true
		}
	}
	for val := range InfoCommands() {
		if val == cmd {
			return true
		}
	}

	return false
}

func commandOpts(args []string) (opts []string) {
	if len(args) >= 2 {
		opts = args[2:len(args)]
	}

	return
}

// handleCommand will take in the argument for the process and run it
func ParseCommand(args []string) {
	SetProcess(args[0])
	SetConfig(config.Process(args[0]))
	opts := commandOpts(args)

	if message, ok := cfg.Valid(); ok {

		name := args[1]
		daemonCmd, daemonOk := DaemonizedCommands()[name]
		infoCmd, infoOk := InfoCommands()[name]
		interactiveCmd, interactiveOk := InteractiveCommands()[name]

		switch {
		case daemonOk:
			daemonCmd.run(opts...)
		case infoOk:
			infoCmd.run(opts...)
		case interactiveOk:
			interactiveCmd.run(opts...)
		default:
			logger.Log(fmt.Sprintf("Running Command: (%v) doesn't exists\n", args[1]))
		}
	} else {
		logger.Log(message)
	}

}

// Pids pulls all docker ids for each of the instances
func pids(args ...string) (pids []string) {
	for i := 0; i < cfg.Instances; i++ {
		id, err := pid(i)

		if err != nil {
			logger.Log(fmt.Sprintln(err))
		} else {
			pids = append(pids, id)
		}

	}
	return pids
}

// pid pulls the docker pid for the given instance
func pid(instance int) (pid string, err error) {
	file, err := os.Open(pidFileName(instance))
	if err != nil {
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	pid = scanner.Text()
	return
}
