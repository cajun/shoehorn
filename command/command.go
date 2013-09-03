package command

import (
	"fmt"
	"github.com/cajun/shoehorn/config"
	"github.com/mgutz/ansi"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var (
	cfg     *config.Settings
	process string
)

func outOpts(opts []string) {
	msg := lime() + "docker %s\n" + reset()
	fmt.Printf(msg, strings.Join(opts, " "))
}

func lime() string {
	return ansi.ColorCode("green+h:black")
}

func reset() string {
	return ansi.ColorCode("reset")
}

func run(command string, inOpts ...string) {
	base := []string{command, "-d"}
	opts := append(base, inOpts...)
	outOpts(opts)

	cmd := exec.Command("docker", opts...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func runInteractive(command string, inOpts ...string) {
	base := []string{"run", "-i", "-t"}
	opts := append(base, inOpts...)
	outOpts(opts)

	cmd := exec.Command("docker", opts...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()
}

func SetConfig(settings *config.Settings) {
	cfg = settings
}

func SetProcess(proc string) {
	process = proc
}

func PrintParams() {
	fmt.Println(settingsToParams())
}

// settingsToParams converts the parameters in the configuration file
// to params that will be passed into docker.
func settingsToParams() (opts []string) {
	if cfg.Options != "" {
		opts = append(opts, cfg.Options)
	}

	opts = append(opts, limitOpts()...)

	if cfg.Dns != "" {
		opts = append(opts, dnsOpts()...)
	}

	if len(volumnsOpts()) != 0 {
		opts = append(opts, volumnsOpts()...)
	}

	opts = append(opts, cfg.Container)
	opts = append(opts, strings.Split(cfg.StartCmd, " ")...)

	return
}

func limitOpts() []string {
	return []string{"-m", strconv.Itoa(cfg.Bytes)}
}

func dnsOpts() []string {
	return []string{"-dns", cfg.Dns}
}

func volumnsOpts() (volumns []string) {
	for _, volumn := range cfg.Volumn {
		volumns = append(volumns, "-v", volumn)
	}
	return volumns
}

func Start() {
	fmt.Printf("Starting %v\n", process)
	for i := 0; i < cfg.Instances; i++ {
		fmt.Printf("...Instance %d of %d %s\n", i, cfg.Instances, process)
		run("run", settingsToParams()...)
	}
}

func Stop() {
	fmt.Printf("Stopping %v\n", process)
	run("stop", "PID_GOES_HERE")
}

func Restart() {
	fmt.Printf("Restarting %v\n", process)
	Start()
	Stop()
}

func Kill() {
	fmt.Printf("Killing %v\n", process)
	run("kill", "PID_GOES_HERE")
}

func Console() {
	cfg.StartCmd = cfg.Console
	runInteractive("run", settingsToParams()...)
}

func Bash() {
	cfg.StartCmd = "/bin/bash"
	runInteractive("run", settingsToParams()...)
}

//func IP() {
//}

//func Port() {
//}

//func PublicPort() {
//}

//func Ssh() {
//}

func Logs() {
	for i := 0; i < cfg.Instances; i++ {
		run("log", "PID_GOES_HERE[i]")
	}
}

func Status() {
	for i := 0; i < cfg.Instances; i++ {
		run("ps", "PID_GOES_HERE[i]")
	}
}
