package command

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var (
	daemonizedCommands  map[string]string
	infoCommands        map[string]string
	interactiveCommands map[string]string
)

func init() {
	if daemonizedCommands == nil {
		daemonizedCommands = make(map[string]string)
	}
	if infoCommands == nil {
		infoCommands = make(map[string]string)
	}
	if interactiveCommands == nil {
		interactiveCommands = make(map[string]string)
	}
	daemonizedCommands["start"] = "start the given process"
	daemonizedCommands["stop"] = "stop the given process"
	daemonizedCommands["kill"] = "kill the given process"
	daemonizedCommands["restart"] = "restrat the givne process"

	infoCommands["running"] = "check to see if the process is running"
	infoCommands["status"] = "view the status of the process"
	infoCommands["ip"] = "view the private ip"
	infoCommands["logs"] = "see logs for the process"
	infoCommands["pids"] = "view the pids for theses processes"
	infoCommands["port"] = "view the private port"
	infoCommands["params"] = "view the params that will be used in the docker command"
	infoCommands["public_port"] = "view the public port"

	interactiveCommands["console"] = "execute the console command from the config"
	interactiveCommands["bash"] = "execute a bash shell for the process"
	interactiveCommands["ssh"] = "ssh into the container"
}

// DaemonizedCommands are commands that will be daemonized or manage daemonized
// commands
func DaemonizedCommands() map[string]string {
	return daemonizedCommands
}

// InfoCommands are commands the will pull out information about the
// given process
func InfoCommands() map[string]string {
	return infoCommands
}

// InteractiveCommands will turn over some kind of command back to the user
func InteractiveCommands() map[string]string {
	return interactiveCommands
}

// Start will run the standard start command
func Start() {
	runInstances("Start", func(i int, id string) error {
		return runDaemon("run", settingsToParams(i, true)...)
	})
}

// Stop will stop all the process if this type.  If the 'Kill' setting is turned
// on then the stop will kill the process instead
func Stop() {
	if cfg.Kill {
		Kill()
	} else {
		runInstances("Stopping", func(i int, id string) error {
			defer os.Remove(pidFileName(i))
			return run("stop", id)
		})
	}
}

// Restart will call stop then start for this process
func Restart() {
	fmt.Printf("Restarting %v\n", process)
	Stop()
	Start()
}

// Kill will kill the given process
func Kill() {
	runInstances("Killing", func(i int, id string) error {
		defer os.Remove(pidFileName(i))
		return run("kill", id)
	})
}

// Console will run an interactive command for the given console command
func Console() {
	cfg.StartCmd = cfg.Console
	runInteractive("run", settingsToParams(0, false)...)
}

// Bash will execute a bash command against the given container
func Bash() {
	cfg.StartCmd = "/bin/bash"
	runInteractive("run", settingsToParams(0, false)...)
}

func IP() {
}

func Port() int {
	return cfg.Port
}

func PublicPort() int {
	return 0
}

func Ssh() {
}

// Running determines if the given process is running.
func Running() (found bool) {
	found = false
	cmd := exec.Command("docker", "ps")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}

	_, err = cmd.StderrPipe()
	if err != nil {
		fmt.Println(err)
	}

	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(stdout)
	s := buf.String()

	fmt.Printf("%s\n", s)
	cmd.Wait()

	for _, id := range Pids() {
		if !found {
			found = strings.Contains(s, id)
		}
	}
	return
}

// Logs will print out all of the logs for each of the instances
func Logs() {
	runInstances("Logs", func(i int, id string) error {
		return run("log", id)
	})
}

// Status will list out the statuses for the given process
func Status() {
	runInstances("Status", func(i int, id string) error {
		return run("ps", id)
	})
}

// run will execute a command against docker with the given
// options as a daemon. run also sets the pid. run will not
// execute if there is an existing pid
func runDaemon(command string, inOpts ...string) error {
	base := []string{"-d"}
	opts := append(base, inOpts...)

	return run(command, opts...)
}

// run will execute a command against docker with the given
// options as a daemon. run also sets the pid. run will not
// execute if there is an existing pid
func run(command string, inOpts ...string) error {
	base := []string{command}
	opts := append(base, inOpts...)
	outOpts(opts)

	cmd := exec.Command("docker", opts...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// runInteractive will give the user the option for input into the
// docker command. examples would be running bash or ssh
func runInteractive(command string, inOpts ...string) error {
	base := []string{command, "-i", "-t"}
	opts := append(base, inOpts...)
	outOpts(opts)

	cmd := exec.Command("docker", opts...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// runner is an interface to a function that will execute
// the run command.
type runner func(instance int, pid string) error

// runInstances wraps the given function and execute
// the number of instances requested by the config file
// for the command given.
func runInstances(message string, fn runner) {
	fmt.Printf("%s %v\n", message, process)
	for i := 0; i < cfg.Instances; i++ {
		fmt.Printf("...Instance %d of %d %s\n", i, cfg.Instances, process)
		id, err := pid(i)
		if err != nil {
			fmt.Println(err)
		}
		fn(i, id)
	}
}

func runExec(cmd string, args ...string) {
	joined := strings.Join(args, " ")
	cfg.StartCmd = "/bin/bash -c"
	cfg.QuotedOpts = fmt.Sprintf("'cd %s && %s %s'", cfg.RemoteVolumn, cmd, joined)

	runInteractive("run", settingsToParams(0, false)...)
}
