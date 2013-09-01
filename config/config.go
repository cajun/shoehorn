package config

import (
	"code.google.com/p/gcfg"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Settings struct {
	App          string
	StartCmd     string
	Instances    int
	Port         int
	PublicPort   int
	MB           int
	GB           int
	Bytes        int
	Domain       []string
	Kill         bool
	Container    string
	Volumn       []string
	RemoteVolumn string
	Raw          string
}
type Config struct {
	App map[string]*Settings
}

var (
	cfg          Config
	globalFile   string
	configFile   string
	overrideFile string
)

func init() {
	const (
		globalConfig   = "$HOME/global.cfg"
		globalUsage    = "configuration file for this app"
		defaultConfig  = "config.cfg"
		usage          = "configuration file for this app"
		overrideConfig = "override.cfg"
		overrideUsage  = "configuration file for this app"
	)

	flag.StringVar(&globalFile, "global", globalConfig, globalUsage)
	flag.StringVar(&globalFile, "g", globalConfig, globalUsage+" (shorthand)")

	flag.StringVar(&configFile, "config", defaultConfig, usage)
	flag.StringVar(&configFile, "c", defaultConfig, usage+" (shorthand)")

	flag.StringVar(&overrideFile, "override", overrideConfig, overrideUsage)
	flag.StringVar(&overrideFile, "o", overrideConfig, overrideUsage+" (shorthand)")

}

func LoadConfigs() {
	gcfg.ReadFileInto(&cfg, globalFile)
	gcfg.ReadFileInto(&cfg, configFile)
	gcfg.ReadFileInto(&cfg, overrideFile)
	setDefaults(&cfg)
}

func setDefaults(cfg *Config) {
	for section, value := range cfg.App {
		if value.App == "" {
			path, _ := os.Getwd()
			dir := filepath.Base(path)
			value.App = dir + "_" + section
		}
		if value.Instances == 0 {
			value.Instances = 1
		}

		if value.GB != 0 && value.Bytes == 0 {
			value.MB = value.GB * 1024
		}

		if value.MB == 0 && value.Bytes == 0 {
			value.MB = 512
		}

		if value.Bytes == 0 {
			value.Bytes = value.MB * 1024 * 1024
		}
	}
}

func List() []string {
	list := []string{}
	for k, _ := range cfg.App {
		list = append(list, k)
	}
	return list
}

func App(name string) *Settings {
	return cfg.App[name]
}

// printApps prints out the listing of all the sections in the config file
func PrintApps() {
	fmt.Println("** List of Apps **")
	fmt.Println(List())
}

func printSetting(name string, value string) {
	if value != "" {
		fmt.Println(name + value)
	}
}

func PrintConfig(name string) {
	settings := App(name)
	fmt.Println("Section [" + name + "]")

	printSetting("App Name: ", settings.App)
	printSetting("Start Command: ", settings.StartCmd)
	printSetting("Number of Instances: ", strconv.Itoa(settings.Instances))
	printSetting("Private Port: ", strconv.Itoa(settings.Port))
	printSetting("Public Port: ", strconv.Itoa(settings.PublicPort))
	printSetting("RAM in GB: ", strconv.Itoa(settings.GB))
	printSetting("RAM in MB: ", strconv.Itoa(settings.MB))
	printSetting("RAM in Bytes: ", strconv.Itoa(settings.Bytes))
	printSetting("Domain: ", strings.Join(settings.Domain, " "))
	printSetting("Kill Process?: ", strconv.FormatBool(settings.Kill))
	printSetting("Container Name: ", settings.Container)
	printSetting("Volumn(s): ", strings.Join(settings.Volumn, " "))
	printSetting("Remote Volumn: ", settings.RemoteVolumn)
	printSetting("Raw Command: ", settings.Raw)
}
