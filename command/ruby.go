package command

import (
	"github.com/cajun/shoehorn/config"
	"strings"
)

func init() {

	available.addInteractive("irb", Executor{
		description: "open an irb session against this container (NOTE: may not have ruby)",
		run:         Irb})

	available.addInteractive("ruby", Executor{
		description: "execute ruby against this container",
		run:         Ruby})

	available.addInteractive("rake", Executor{
		description: "execute rake against this container",
		run:         Rake})

	available.addInteractive("bundle", Executor{
		description: "execute bundle against this container",
		run:         Bundle})

	available.addInteractive("bundle_install", Executor{
		description: "execute bundle install --path .gems against this container",
		run:         BundleInstall})
}

func Irb(args ...string) {
	runExec("irb", strings.Join(args, " "))
}

func Ruby(args ...string) {
	runExec("ruby", strings.Join(args, " "))
}

func Rake(args ...string) {
	runExec("rake", strings.Join(args, " "))
}

func Bundle(args ...string) {
	runExec("bundle", strings.Join(args, " "))
}

func BundleInstall(args ...string) {
	if cfg == nil {
		config.LoadConfigs()
		for _, process := range config.List() {
			SetProcess(process)
			SetConfig(config.Process(process))
			BundleInstall(args...)
		}
	} else if cfg.UseBundler {
		runExec("bundle", "install", "--path", ".gems")
	}
}
