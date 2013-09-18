package command

import (
	"strings"
)

func init() {

	available.addInteractive("rails", Executor{
		description: "execute the rails command against this container",
		run:         Rails})

	available.addInteractive("assets", Executor{
		description: "build the assets for rails",
		run:         Assets})
}

func Rails(args ...string) {
	runExec("./bin/rails", strings.Join(args, " "))
}

func Assets(args ...string) {
	Rake("assets:clean", "assets:precompile")
}
