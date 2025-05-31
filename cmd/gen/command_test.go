package gen_test

import (
	"log/slog"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jonathonwebb/tilde/cmd/gen"
	"github.com/jonathonwebb/tilde/internal/cli"
	"github.com/jonathonwebb/tilde/internal/core"
)

func TestGenCommand(t *testing.T) {
	t.Run("without command", func(t *testing.T) {
		e, cfg, errBuf, outBuf := setUp(t)

		gotCode := gen.Cmd.Execute(t.Context(), e, cfg)
		wantCode := cli.ExitUsageError
		if wantCode != gotCode {
			t.Errorf("want exit status = %v, but got %v", wantCode, gotCode)
		}

		gotErr := errBuf.String()
		wantErr := `missing command
usage: tilde [root flags] gen <command>
`
		if diff := cmp.Diff(wantErr, gotErr); diff != "" {
			t.Errorf("err output mismatch (-want +got):\n%s", diff)
		}

		gotOut := outBuf.String()
		wantOut := ""
		if diff := cmp.Diff(wantOut, gotOut); diff != "" {
			t.Errorf("out output mismatch (-want +got):\n%s", diff)
		}
	})
}

func setUp(t testing.TB, args ...string) (*cli.Env, *core.Config, *strings.Builder, *strings.Builder) {
	t.Helper()

	var errBuf, outBuf strings.Builder

	return &cli.Env{
			Stderr: &errBuf,
			Stdout: &outBuf,
			Args:   append([]string{"assets"}, args...),
		}, &core.Config{
			Env:    "test",
			Level:  slog.LevelError,
			Format: core.JSONFormat,
		},
		&errBuf,
		&outBuf
}
