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
	App        string
	StartCmd   string
	Console    string
	Options    string
	Instances  int
	Port       int
	PublicPort int
	MB         int
	GB         int
	Bytes      int
	Domain     []string
	Allow      []string
	Kill       bool
	Container  string
	Volumn     []string
	WorkingDir string
	QuotedOpts string
	Raw        string
	Dns        string
	AutoStart  bool
	UseNginx   bool
	Env        []string
}

type Config struct {
	Process map[string]*Settings
}

var (
	cfg          Config
	globalFile   string
	configFile   string
	overrideFile string
)

func init() {
	var (
		globalConfig   = "$HOME/.shoehorn.cfg"
		globalUsage    = "configuration file for this app"
		defaultConfig  = "shoehorn.cfg"
		usage          = "configuration file for this app"
		overrideConfig = "shoehorn.override.cfg"
		overrideUsage  = "configuration file for this app"
	)

	flag.StringVar(&globalFile, "global", globalConfig, globalUsage)
	flag.StringVar(&configFile, "config", defaultConfig, usage)
	flag.StringVar(&overrideFile, "override", overrideConfig, overrideUsage)

}

func LoadConfigs() {
	defer setDefaults(&cfg) // Load the defaults at the end of the function call
	cfg = Config{}
	gcfg.ReadFileInto(&cfg, globalFile)
	gcfg.ReadFileInto(&cfg, configFile)
	gcfg.ReadFileInto(&cfg, overrideFile)
}

// setDefaults fills out the configuration with safe defaults to ensure
// no process will try to take over the system by default
func setDefaults(cfg *Config) {
	for section, value := range cfg.Process {
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

		if value.Bytes == 0 {
			value.Bytes = value.MB * 1024 * 1024
		}
	}
}

// List out the available processes
func List() []string {
	list := []string{}
	for k := range cfg.Process {
		list = append(list, k)
	}
	return list
}

func Processes() Config {
	return cfg
}

// Process pulls the settings for the given process
func Process(name string) *Settings {
	return cfg.Process[name]
}

// PrintProcesses prints out the listing of all the sections in the config file
func PrintProcesses() {
	fmt.Println("** List of Apps **")
	fmt.Println(List())
}

func printSetting(name string, value string) {
	if value != "" {
		fmt.Printf("%20v  %v\n", name+":", value)
	}
}

func PrintConfig(name string) {
	settings := Process(name)
	fmt.Println("Process [" + name + "]")

	printSetting("App Name", settings.App)
	printSetting("Start Command", settings.StartCmd)
	printSetting("Number of Instances", strconv.Itoa(settings.Instances))
	printSetting("Private Port", strconv.Itoa(settings.Port))
	printSetting("Public Port", strconv.Itoa(settings.PublicPort))
	printSetting("RAM in GB", strconv.Itoa(settings.GB))
	printSetting("RAM in MB", strconv.Itoa(settings.MB))
	printSetting("RAM in Bytes", strconv.Itoa(settings.Bytes))
	printSetting("Domain", strings.Join(settings.Domain, " "))
	printSetting("Kill Process?", strconv.FormatBool(settings.Kill))
	printSetting("Container Name", settings.Container)
	printSetting("Volumn(s)", strings.Join(settings.Volumn, " "))
	printSetting("Working Directory", settings.WorkingDir)
	printSetting("Raw Command", settings.Raw)
	printSetting("Auto Start", strconv.FormatBool(settings.AutoStart))
}
