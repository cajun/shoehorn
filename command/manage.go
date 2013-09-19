package command

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/cajun/shoehorn/config"
	"github.com/cajun/shoehorn/logger"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var (
	available Available
	root      = "."
)

type Ports struct {
	tcp string
	udp string
}

type Network struct {
	Ip      string
	Public  Ports
	Private Ports
}

type Available struct {
	daemonized  map[string]Executor
	info        map[string]Executor
	interactive map[string]Executor
}

func Root() string {
	return root
}

func (a *Available) addDaemon(name string, exe Executor) {
	if a.daemonized == nil {
		a.daemonized = make(map[string]Executor)
	}
	a.daemonized[name] = exe
}

func (a *Available) addInfo(name string, exe Executor) {
	if a.info == nil {
		a.info = make(map[string]Executor)
	}
	a.info[name] = exe
}

func (a *Available) addInteractive(name string, exe Executor) {
	if a.interactive == nil {
		a.interactive = make(map[string]Executor)
	}
	a.interactive[name] = exe
}

func init() {
	flag.StringVar(&root, "root", ".", "which dir at the apps located")

	available.addDaemon("start", Executor{
		description: "start the given process",
		run:         Start})
	available.addDaemon("stop", Executor{
		description: "stop the given process",
		run:         Stop})
	available.addDaemon("kill", Executor{
		description: "kill the given process",
		run:         Kill})
	available.addDaemon("restart", Executor{
		description: "restrat the givne process",
		run:         Restart})

	available.addInfo("status", Executor{
		description: "view the status of the process",
		run:         Status})
	available.addInfo("logs", Executor{
		description: "see logs for the process",
		run:         Logs})
	available.addInfo("params", Executor{
		description: "view the params that will be used in the docker command",
		run:         PrintParams})
	available.addInfo("build", Executor{
		description: "build a container from a file or url",
		run:         Build})

	available.addInteractive("console", Executor{
		description: "execute the console command from the config",
		run:         Console})
	available.addInteractive("bash", Executor{
		description: "execute a bash shell for the process",
		run:         Bash})
	available.addInteractive("get", Executor{
		description: "clone a git repo and then build the images",
		run:         Install})
	available.addInteractive("update", Executor{
		description: "pull and then build the images",
		run:         Update})
	available.addInteractive("attach", Executor{
		description: "attaches to the first instance running",
		run:         Attach})
	//available.addInteractive("ssh", Executor{
	//description: "ssh into the container",
	//run:         Ssh})
}

// DaemonizedCommands are commands that will be daemonized or manage daemonized
// commands
func DaemonizedCommands() map[string]Executor {
	return available.daemonized
}

// InfoCommands are commands the will pull out information about the
// given process
func InfoCommands() map[string]Executor {
	return available.info
}

// InteractiveCommands will turn over some kind of command back to the user
func InteractiveCommands() map[string]Executor {
	return available.interactive
}

// Build will create the container if nessary
func Build(args ...string) {
	if cfg != nil {
		logger.Log(fmt.Sprintf("Building...%s\n", cfg.App))
		cmd := exec.Command("docker", "build", "-t", cfg.Container, cfg.BuildFile)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Run()
	} else {
		for _, process := range config.List() {
			SetProcess(process)
			SetConfig(config.Process(process))
			Build(args...)
		}
	}
}

// Start will run the standard start command
func Start(args ...string) {
	runInstances("Start", func(i int, id string) error {
		return runDaemon("run", settingsToParams(i, true)...)
	})

	if cfg.UseNginx {
		UpdateNginxConf()
	}
}

// Stop will stop all the process if this type.  If the 'Kill' setting is turned
// on then the stop will kill the process instead
func Stop(args ...string) {
	switch {
	case cfg.Kill:
		Kill(args...)
	default:
		runInstances("Stopping", func(i int, id string) error {
			defer os.Remove(pidFileName(i))
			return run("stop", id)
		})
	}
}

// Restart will call stop then start for this process
func Restart(args ...string) {
	logger.Log(fmt.Sprintf("Restarting %v\n", process))
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
	cfg.StartCmd = "/bin/bash -c"
	cfg.QuotedOpts = "'" + cfg.Console + "'"
	runInteractive("run", settingsToParams(0, false)...)
}

// Bash will execute a bash command against the given container
func Bash(args ...string) {
	cfg.StartCmd = "/bin/bash"
	runInteractive("run", settingsToParams(0, false)...)
}

func Install(args ...string) {
	go func() {
		os.Chdir(root)
		opts := []string{"clone", args[0]}

		logger.Log("Cloning " + args[0] + "\n")
		cmd := exec.Command("git", opts...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()

		if err == nil {
			parts := strings.Split(args[0], "/")
			path := parts[len(parts)-1:][0]
			logger.Log("Building Images here: " + path + "\n")
			os.Chdir(path)
			Build(path)
			BundleInstall("now")
		} else {
			logger.Log(err.Error())
		}
	}()
}

func Update(args ...string) {
	go func() {
		logger.Log("Updaing...")
		cmd := exec.Command("git", "pull")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()

		if err == nil {
			Build(root)
			BundleInstall("now")
		} else {
			logger.Log(err.Error())
		}
	}()

}

func ip(instance int) string {
	return networkSettings(instance).Ip
}

func ports(instance int, settings map[string]interface{}) (public, private Ports) {

	if settings["PortMapping"] != nil {
		s := settings["PortMapping"].(map[string]interface{})

		if s["Tcp"] != nil {
			for private_port, public_port := range s["Tcp"].(map[string]interface{}) {
				private.tcp = private_port
				public.tcp = public_port.(string)
			}

		}

		if s["Udp"] != nil {
			for private_port, public_port := range s["Udp"].(map[string]interface{}) {
				private.tcp = private_port
				public.tcp = public_port.(string)
			}
		}

	}
	return
}

//func Ssh(args ...string) {
//}

func networkSettings(instance int) (net Network) {
	settings, _ := inspect(instance)
	settings = settings["NetworkSettings"].(map[string]interface{})

	net.Ip = settings["IPAddress"].(string)

	net.Public, net.Private = ports(instance, settings)
	return
}

func inspect(instance int) (u map[string]interface{}, err error) {
	id, err := pid(instance)
	out, err := exec.Command("docker", "inspect", id).Output()
	if err != nil {
		return
	}

	all := []map[string]interface{}{}
	err = json.Unmarshal(out, &all)
	if len(all) > 0 {
		u = all[0]
	}
	return
}

// Running determines if the given process is running.
func running(args ...string) (found bool) {
	found = false
	cmd := exec.Command("docker", "ps")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logger.Log(fmt.Sprintln(err))
	}

	_, err = cmd.StderrPipe()
	if err != nil {
		logger.Log(fmt.Sprintln(err))
	}

	err = cmd.Start()
	if err != nil {
		logger.Log(fmt.Sprintln(err))
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(stdout)
	s := buf.String()

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
		return run("logs", id)
	})
}

func Attach(args ...string) {
	runInstances("Attach", func(i int, id string) error {
		return run("attach", id)
	})
}

// Status will list out the statuses for the given process
func Status(args ...string) {
	runInstances("Status", func(i int, id string) error {
		on := running()
		logger.Log(fmt.Sprintf("Container ID: %s\n", id))
		logger.Log(fmt.Sprintf("     Running: %s\n", strconv.FormatBool(on)))

		if on {
			net := networkSettings(i)
			logger.Log(fmt.Sprintf("          IP: %s\n", net.Ip))
			logger.Log(fmt.Sprintf(" Public Port: %s\n", net.Public.tcp))
			logger.Log(fmt.Sprintf("Private Port: %s\n", net.Private.tcp))
		}

		return nil
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
	logger.Log(fmt.Sprintf("%s %v\n", message, process))
	for i := 0; i < cfg.Instances; i++ {
		logger.Log(fmt.Sprintf("...Instance %d of %d %s\n", i, cfg.Instances, process))
		id, _ := pid(i)
		fn(i, id)
	}

}

func runExec(cmd string, args ...string) {
	joined := strings.Join(args, " ")
	cfg.StartCmd = "/bin/bash -c"
	cfg.QuotedOpts = fmt.Sprintf("'%s %s'", cmd, joined)

	runInteractive("run", settingsToParams(0, false)...)
}
