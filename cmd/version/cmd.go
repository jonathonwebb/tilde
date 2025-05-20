package version

import (
	"context"
	"fmt"

	"github.com/jonathonwebb/x/conf"
)

var Cmd = conf.Command{
	Name:  "version",
	Usage: `usage: tilde [-h | -help] version`,
	Help: `usage: tilde [-h | -help] version

Version prints information about the tilde build.

flags:
  -h, --help   show this help and exit`,
	Action: func(ctx context.Context, e *conf.Env, _ any) conf.ExitStatus {
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
		return conf.ExitSuccess
	},
}
