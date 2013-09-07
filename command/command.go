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

func run() {
	base := []string{"run", "-d"}
	opts := append(base, settingsToParams()...)
	outOpts(opts)

	cmd := exec.Command("docker", opts...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func runInteractive() {
	base := []string{"run", "-i", "-t"}
	opts := append(base, settingsToParams()...)
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

func volumnsOpts() []string {
	volumns := []string{}
	for volumn := range cfg.Volumn {
		volumns = append(volumns, fmt.Sprintf("-v %v", volumn))
	}
	return volumns
}

func Start() {
	fmt.Printf("Starting %v\n", process)
	run()
}

func Stop() {
	fmt.Printf("Stopping %v\n", process)
}

func Restart() {
	fmt.Printf("Restarting %v\n", process)
	Start()
	Stop()
}

func Kill() {
	fmt.Printf("Killing %v\n", process)
}

func Console() {
	cfg.StartCmd = cfg.Console
	runInteractive()
}

func Bash() {
	cfg.StartCmd = "/bin/bash"
	runInteractive()
}

//func IP() {
//}

//func Port() {
//}

//func PublicPort() {
//}

//func Ssh() {
//}

//func Logs() {
//}

//func Status() {
//}
