package config

import (
	"code.google.com/p/gcfg"
	"flag"
	"fmt"
	"github.com/cajun/shoehorn/logger"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Settings struct {
	App        string
	StartCmd   string
	Console    string
	Instances  int
	Port       int
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
	BuildFile  string
	Dns        string
	UseNginx   bool
	Env        []string
	IncludeEnv bool
	UseBundler bool
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

func (s *Settings) Valid() (message string, valid bool) {
	if s.Container == "" {
		message = fmt.Sprintln("You must specifiy a container")
	}

	return message, len(message) == 0
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
	logger.Log(fmt.Sprintln("** List of Apps **"))
	logger.Log(fmt.Sprintln(List()))
}

func printSetting(name string, value string) {
	if value != "" {
		logger.Log(fmt.Sprintf("%20v  %v\n", name+":", value))
	}
}

func Init() {
	_, err := ioutil.ReadFile(configFile)
	if err != nil {
		ioutil.WriteFile(configFile, []byte(sampleFile), 0666)
	}
}

func PrintConfig(name string) {
	settings := Process(name)
	logger.Log(fmt.Sprintln("Process [" + name + "]"))

	printSetting("App Name", settings.App)
	printSetting("Start Command", settings.StartCmd)
	printSetting("Number of Instances", strconv.Itoa(settings.Instances))
	printSetting("Private Port", strconv.Itoa(settings.Port))
	printSetting("RAM in GB", strconv.Itoa(settings.GB))
	printSetting("RAM in MB", strconv.Itoa(settings.MB))
	printSetting("RAM in Bytes", strconv.Itoa(settings.Bytes))
	printSetting("Domain", strings.Join(settings.Domain, " "))
	printSetting("Kill Process?", strconv.FormatBool(settings.Kill))
	printSetting("Container Name", settings.Container)
	printSetting("Volumn(s)", strings.Join(settings.Volumn, " "))
	printSetting("Working Directory", settings.WorkingDir)
	printSetting("Build File", settings.BuildFile)
}

const sampleFile = `
[process "db"]
Container=mongodb_2.4.5 # required
App=MongoDB
StartCmd=mongod
Gb=2
Port=27017
Console=mongo --host $MONGO_0_IP
Dns=10.1.1.7
Volumn=./.data/mongo:/data/db
BuildFile=Docker.file.or.url
IncludeEnv=true


[process "rails"]
StartCmd=./bin/rails server
Container=ruby_2.0.0
Volumn=.:/www
Port=5000
Mb=512
WorkingDir=/www
BuildFile=Docker.file.or.url
Domain=rails.dev   # required for nginx template
Allow=10.12.2.0/24 # required for nginx template
Allow=10.12.1.0/24 # required for nginx template
Dns=10.1.1.7
UseNginx=true
Kill=true

`
