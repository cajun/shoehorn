package command

import (
	"fmt"
	"github.com/cajun/shoehorn/config"
	"os"
	"os/exec"
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
func settingsToParams() {

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
