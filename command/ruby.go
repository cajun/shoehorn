package command

import (
	"strings"
)

func init() {
	if interactiveCommands == nil {
		interactiveCommands = make(map[string]string)
	}

	interactiveCommands["irb"] = "open an irb session against this container (NOTE: may not have ruby)"
	interactiveCommands["ruby"] = "execute ruby against this container"
	interactiveCommands["rake"] = "execute rake against this container"
	interactiveCommands["bundle"] = "execute bundle against this container"
	interactiveCommands["bundle_install"] = "execute bundle install --path .gems against this container"
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

func BundleInstall() {
	runExec("bundle install --path .gems", "")
}
