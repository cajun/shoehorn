package command

import (
  "fmt"
  "github.com/cajun/shoehorn/config"
  "github.com/cajun/shoehorn/logger"
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
  msg := "\n\n" + lime + "docker %s" + reset + "\n\n"
  logger.Log(fmt.Sprintf(msg, strings.Join(opts, " ")))
}

// pidFileName returns the path of the pid file
func pidFileName(instance int) string {
  return fmt.Sprintf("tmp/pids/%s.%d.pid", process, instance)
}

// settingsToParams converts the parameters in the configuration file
// to params that will be passed into docker.
func settingsToParams(instance int, withPid bool) (opts []string) {

  if withPid {
    opts = append(opts, "--cidfile", pidFileName(instance))
  }

  if cfg.WorkingDir != "" {
    opts = append(opts, "-w", cfg.WorkingDir)
  }

  if env := envOpts(); len(env) != 0 {
    opts = append(opts, env...)
  }

  if cfg.Bytes != 0 {
    opts = append(opts, limitOpts()...)
  }

  if cfg.Port != 0 {
    opts = append(opts, portOpts()...)
  }

  if cfg.Dns != "" {
    opts = append(opts, dnsOpts()...)
  }

  if vols := volumnsOpts(); len(vols) != 0 {
    opts = append(opts, vols...)
  }

  opts = append(opts, cfg.Container)

  if cfg.StartCmd != "" {
    opts = append(opts, strings.Split(cfg.StartCmd, " ")...)
  }
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
  return []string{"--dns", cfg.Dns}
}

// volumnOpts can have multiple settings.  It will convert each one into
// the volume setting
func volumnsOpts() (volumns []string) {
  for _, volumn := range cfg.Volumn {
    usr, _ := user.Current()
    vol := strings.Replace(volumn, "~", usr.HomeDir, -1)
    path, _ := os.Getwd()
    vol = strings.Replace(vol, ".", path, 1)
    volumns = append(volumns, "-v", vol)
  }
  return volumns
}

func envOpts() (opts []string) {
  for _, env := range cfg.Env {
    opts = append(opts, "-e", env)
  }

  old_cfg := cfg
  old_process := process

  for _, process := range config.List() {
    settings := config.Process(process)
    if settings.IncludeEnv {
      SetProcess(process)
      SetConfig(config.Process(process))

      for i := 0; i < cfg.Instances; i++ {
        if running() {
          net := networkSettings(i)
          name := fmt.Sprintf("%s_%d_IP=%s", process, i, net.Ip)
          opts = append(opts, "-e", strings.ToUpper(name))
          name = fmt.Sprintf("%s_%d_PORT=%s", process, i, net.Private.tcp)
          opts = append(opts, "-e", strings.ToUpper(name))
        }
      }
    }
  }

  cfg = old_cfg
  process = old_process
  return
}
