package command

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var (
	daemonizedCommands  map[string]Executor
	infoCommands        map[string]Executor
	interactiveCommands map[string]Executor
)

func init() {
	if daemonizedCommands == nil {
		daemonizedCommands = make(map[string]Executor)
	}
	if infoCommands == nil {
		infoCommands = make(map[string]Executor)
	}
	if interactiveCommands == nil {
		interactiveCommands = make(map[string]Executor)
	}

	daemonizedCommands["start"] = Executor{
		description: "start the given process",
		run:         Start}
	daemonizedCommands["stop"] = Executor{
		description: "stop the given process",
		run:         Stop}
	daemonizedCommands["kill"] = Executor{
		description: "kill the given process",
		run:         Kill}
	daemonizedCommands["restart"] = Executor{
		description: "restrat the givne process",
		run:         Restart}

	infoCommands["status"] = Executor{
		description: "view the status of the process",
		run:         Status}
	infoCommands["ip"] = Executor{
		description: "view the private ip",
		run:         IP}
	infoCommands["logs"] = Executor{
		description: "see logs for the process",
		run:         Logs}
	infoCommands["port"] = Executor{
		description: "view the private port",
		run:         Port}
	infoCommands["params"] = Executor{
		description: "view the params that will be used in the docker command",
		run:         PrintParams}

	infoCommands["public_port"] = Executor{
		description: "view the public port",
		run:         PublicPort}

	interactiveCommands["console"] = Executor{
		description: "execute the console command from the config",
		run:         Console}
	interactiveCommands["bash"] = Executor{
		description: "execute a bash shell for the process",
		run:         Bash}
	interactiveCommands["ssh"] = Executor{
		description: "ssh into the container",
		run:         Ssh}
}

// DaemonizedCommands are commands that will be daemonized or manage daemonized
// commands
func DaemonizedCommands() map[string]Executor {
	return daemonizedCommands
}

// InfoCommands are commands the will pull out information about the
// given process
func InfoCommands() map[string]Executor {
	return infoCommands
}

// InteractiveCommands will turn over some kind of command back to the user
func InteractiveCommands() map[string]Executor {
	return interactiveCommands
}

// Start will run the standard start command
func Start(args ...string) {
	runInstances("Start", func(i int, id string) error {
		return runDaemon("run", settingsToParams(i, true)...)
	})
}

// Stop will stop all the process if this type.  If the 'Kill' setting is turned
// on then the stop will kill the process instead
func Stop(args ...string) {
	if cfg.Kill {
		Kill(args...)
	} else {
		runInstances("Stopping", func(i int, id string) error {
			defer os.Remove(pidFileName(i))
			return run("stop", id)
		})
	}
}

// Restart will call stop then start for this process
func Restart(args ...string) {
	fmt.Printf("Restarting %v\n", process)
	Stop(args...)
	Start(args...)
}

// Kill will kill the given process
func Kill(args ...string) {
	runInstances("Killing", func(i int, id string) error {
		defer os.Remove(pidFileName(i))
		return run("kill", id)
	})
}

// Console will run an interactive command for the given console command
func Console(args ...string) {
	cfg.StartCmd = cfg.Console
	runInteractive("run", settingsToParams(0, false)...)
}

// Bash will execute a bash command against the given container
func Bash(args ...string) {
	cfg.StartCmd = "/bin/bash"
	runInteractive("run", settingsToParams(0, false)...)
}

func IP(args ...string) {
}

func Port(args ...string) {
}

func publicPort() int {
	return 1
}

func PublicPort(args ...string) {
}

func Ssh(args ...string) {
}

// Running determines if the given process is running.
func running(args ...string) (found bool) {
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

	for _, id := range pids() {
		if !found {
			found = strings.Contains(s, id)
		}
	}

	return
}

// Logs will print out all of the logs for each of the instances
func Logs(args ...string) {
	runInstances("Logs", func(i int, id string) error {
		return run("log", id)
	})

}

// Status will list out the statuses for the given process
func Status(args ...string) {
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
