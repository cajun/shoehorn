package command

import (
	"bufio"
	"fmt"
	"github.com/cajun/shoehorn/config"
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
	os.MkdirAll("tmp/pids", os.ModeDir|0700)
	os.MkdirAll("log", os.ModeDir|0700)
	os.MkdirAll("config", os.ModeDir|0700)
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
func PrintParams() {
	fmt.Println(settingsToParams(0, false))
}

// PrintCommands will list out all of the commands to the end user.
func PrintCommands() {
	fmt.Println("** Daemonized Commands **")
	for cmd, desc := range DaemonizedCommands() {
		fmt.Printf("%15s: %s\n", cmd, desc)
	}

	fmt.Println("** Information Commands **")
	for cmd, desc := range InfoCommands() {
		fmt.Printf("%15s: %s\n", cmd, desc)
	}

	fmt.Println("** Interactive Commands **")
	for cmd, desc := range InteractiveCommands() {
		fmt.Printf("%15s: %s\n", cmd, desc)
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

// handleCommand will take in the argument for the process and run it
func ParseCommand(args []string) {
	SetProcess(args[0])
	SetConfig(config.Process(args[0]))

	opts := []string{}
	if len(args) >= 2 {
		opts = args[2:len(args)]
	}

	switch args[1] {
	case "running":
		Running()
	case "pids":
		for _, id := range Pids() {
			fmt.Println(id)
		}
	case "start":
		Start()
	case "stop":
		Stop()
	case "restart":
		Restart()
	case "kill":
		Kill()
	case "bash":
		Bash()
	case "console":
		Console()
	case "params":
		PrintParams()
	case "ip":
		IP()
	case "port":
		Port()
	case "public_port":
		PublicPort()
	case "ssh":
		Ssh()
	case "status":
		Status()
	case "logs":
		Logs()
	case "ruby":
		Ruby(opts...)
	case "irb":
		Irb(opts...)
	case "rake":
		Rake(opts...)
	case "bundle":
		Bundle(opts...)
	case "bundle_install":
		BundleInstall()
	default:
		fmt.Printf("Running Command: (%v) doesn't exists\n", args[2])
	}
}

// Pids pulls all docker ids for each of the instances
func Pids() (pids []string) {
	for i := 0; i < cfg.Instances; i++ {
		id, err := pid(i)

		if err != nil {
			fmt.Println(err)
		}

		pids = append(pids, id)
	}
	return
}

// pid pulls the docker pid for the given instance
func pid(instance int) (pid string, err error) {
	file, err := os.Open(pidFileName(instance))
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	pid = scanner.Text()
	return
}
