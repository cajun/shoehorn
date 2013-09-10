package command

import (
	"strings"
)

func init() {
	if interactiveCommands == nil {
		interactiveCommands = make(map[string]Executor)
	}

	interactiveCommands["irb"] = Executor{
		description: "open an irb session against this container (NOTE: may not have ruby)",
		run:         Irb}

	interactiveCommands["ruby"] = Executor{
		description: "execute ruby against this container",
		run:         Ruby}

	interactiveCommands["rake"] = Executor{
		description: "execute rake against this container",
		run:         Rake}

	interactiveCommands["bundle"] = Executor{
		description: "execute bundle against this container",
		run:         Bundle}

	interactiveCommands["bundle_install"] = Executor{
		description: "execute bundle install --path .gems against this container",
		run:         BundleInstall}
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
	runExec("bundle install --path .gems", "")
}
