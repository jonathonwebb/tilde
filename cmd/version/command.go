package version

import (
	"context"
	"fmt"

	"github.com/jonathonwebb/tilde/internal/cli"
)

var Cmd = cli.Command{
	Name:  "version",
	Usage: "usage: tilde [-h] version",
	Help: `usage: tilde [-h] version

Version prints information about the app build.

flags:
  -h, -help   show this help and exit`,
	Action: func(ctx context.Context, e *cli.Env, _ any) cli.ExitStatus {
		var (
			version = "?"
			rev     = "?"
		)
		if versionMeta, ok := e.Meta["version"].(string); ok && versionMeta != "" {
			version = versionMeta
		}
		if revMeta, ok := e.Meta["revision"].(string); ok && revMeta != "" {
			rev = revMeta
		}
		if len(rev) > 7 {
			rev = rev[0:7] // shorten long hashes
		}
		fmt.Fprintf(e.Stdout, "tilde v%s (rev %s)\n", version, rev)
		return cli.ExitSuccess
	},
}
