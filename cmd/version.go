package cmd

import (
	"context"
	"fmt"

	"github.com/jonathonwebb/tilde/internal/cli"
)

const (
	versionUsage = `usage: tilde [-h] version`
	versionHelp  = `usage: tilde [-h] version

Version prints program build information.

flags:
  -h, --help   show this help and exit`
)

var versionCmd = cli.Command{
	Name:  "version",
	Usage: versionUsage,
	Help:  versionHelp,
	Action: func(ctx context.Context, e *cli.Env) cli.ExitStatus {
		if e.Build == "" {
			fmt.Fprint(e.Stdout, "tilde (unknown build)\n")
		} else {
			fmt.Fprintf(e.Stdout, "tilde %s\n", e.Build)
		}
		return cli.ExitSuccess
	},
}
