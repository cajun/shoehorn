package command

import (
	"strings"
)

func init() {
	if interactiveCommands == nil {
		interactiveCommands = make(map[string]string)
	}

	interactiveCommands["rails"] = "execute the rails command against this container"
	interactiveCommands["assets"] = "build the assets for rails"
}

func Rails(args ...string) {
	runExec("./bin/rails", strings.Join(args, " "))
}

func Assets(args ...string) {
	Rake("assets:clean", "assets:precompile")
}
