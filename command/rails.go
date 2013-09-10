package command

import (
	"strings"
)

func init() {
	if interactiveCommands == nil {
		interactiveCommands = make(map[string]Executor)
	}

	interactiveCommands["rails"] = Executor{
		description: "execute the rails command against this container",
		run:         Rails}

	interactiveCommands["assets"] = Executor{
		description: "build the assets for rails",
		run:         Assets}
}

func Rails(args ...string) {
	runExec("./bin/rails", strings.Join(args, " "))
}

func Assets(args ...string) {
	Rake("assets:clean", "assets:precompile")
}
