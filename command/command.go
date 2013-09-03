package command

import (
	"fmt"
	"github.com/cajun/shoehorn/config"
	"net"
	"os"
	"os/exec"
	"strings"
)

var (
	cfg     *config.Settings
	process string
)

func run() {
	cmd := exec.Command("ls", "-l")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func runInteractive() {
	cmd := exec.Command("irb")
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

// settingsToParams converts the parameters in the configuration file
// to params that will be passed into docker.
func settingsToParams() string {
	opts := []string{}
	opts = append(opts, cfg.Options)
	opts = append(opts, limitOpts())
	opts = append(opts, dnsOpts())
	opts = append(opts, volumnsOpts())
	opts = append(opts, cfg.Container)
	opts = append(opts, cfg.StartCmd)

	return strings.Join(opts, " ")
}

func limitOpts() string {
	return fmt.Sprintf("-m %v", cfg.Bytes)
}

func dnsOpts() string {
	return fmt.Sprintf("-dns %v", "127.0.0.1")
}

func volumnsOpts() string {
	volumns := []string{}
	for volumn := range cfg.Volumn {
		volumns = append(volumns, fmt.Sprintf("-v %v", volumn))
	}
	return strings.Join(volumns, " ")
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

func Dns() {
	xips, _ := net.InterfaceAddrs()
	for index := range xips {
		fmt.Printf("DNS ip: %v\n", xips[index].Network())
		fmt.Printf("DNS name: %v\n", xips[index].String())
	}

	ips, _ := net.LookupNS("swap")
	fmt.Printf("HOST : %v\n", ips)

}

//func Console() {
//}

//func Bash() {
//}

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

//func Status() {
//}
