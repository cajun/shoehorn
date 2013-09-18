package command

import (
	"fmt"
	"github.com/mgutz/ansi"
	"os"
	"os/user"
	"strconv"
	"strings"
)

// outOpts will colorize the opts as well as print the docker command
// that is about to execute.
func outOpts(opts []string) {
	lime := ansi.ColorCode("green:black")
	reset := ansi.ColorCode("reset")
	msg := lime + "docker %s" + reset + "\n"
	fmt.Println(" ")
	fmt.Printf(msg, strings.Join(opts, " "))
	fmt.Println(" ")
}

// pidFileName returns the path of the pid file
func pidFileName(instance int) string {
	return fmt.Sprintf("tmp/pids/%s.%d.pid", process, instance)
}

// settingsToParams converts the parameters in the configuration file
// to params that will be passed into docker.
func settingsToParams(instance int, withPid bool) (opts []string) {

	if withPid {
		opts = append(opts, "-cidfile", pidFileName(instance))
	}

	if cfg.WorkingDir != "" {
		opts = append(opts, "-w", cfg.WorkingDir)
	}

	for _, env := range cfg.Env {
		opts = append(opts, "-e", env)
	}

	if cfg.Bytes != 0 {
		opts = append(opts, limitOpts()...)
	}

	if cfg.Options != "" {
		opts = append(opts, cfg.Options)
	}

	if cfg.Port != 0 {
		opts = append(opts, portOpts()...)
	}

	if cfg.Dns != "" {
		opts = append(opts, dnsOpts()...)
	}

	if len(volumnsOpts()) != 0 {
		opts = append(opts, volumnsOpts()...)
	}

	opts = append(opts, cfg.Container)
	opts = append(opts, strings.Split(cfg.StartCmd, " ")...)
	if cfg.QuotedOpts != "" {
		opts = append(opts, fmt.Sprintf("'%s'", cfg.QuotedOpts))
	}

	return
}

// limitOpts converts the memory limits into docker settings
func limitOpts() []string {
	return []string{"-m", strconv.Itoa(cfg.Bytes)}
}

func portOpts() []string {
	return []string{"-p", strconv.Itoa(cfg.Port)}
}

// dnsOpts converts the dns settings into docker settings
func dnsOpts() []string {
	return []string{"-dns", cfg.Dns}
}

// volumnOpts can have multiple settings.  It will convert each one into
// the volume setting
func volumnsOpts() (volumns []string) {
	for _, volumn := range cfg.Volumn {
		usr, _ := user.Current()
		vol := strings.Replace(volumn, "~", usr.HomeDir, -1)
		path, _ := os.Getwd()
		vol = strings.Replace(vol, ".", path, -1)
		volumns = append(volumns, "-v", vol)
	}
	return volumns
}
