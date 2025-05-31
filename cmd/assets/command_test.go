package assets_test

import (
	"log/slog"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jonathonwebb/tilde/cmd/assets"
	"github.com/jonathonwebb/tilde/internal/cli"
	"github.com/jonathonwebb/tilde/internal/core"
)

func TestAssetsCommand(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		e, cfg, errBuf, outBuf := setUp(t)

		gotCode := assets.Cmd.Execute(t.Context(), e, cfg)
		wantCode := cli.ExitSuccess
		if wantCode != gotCode {
			t.Errorf("want exit status = %v, but got %v", wantCode, gotCode)
		}

		gotErr := errBuf.String()
		wantErr := ""
		if diff := cmp.Diff(wantErr, gotErr); diff != "" {
			t.Errorf("err output mismatch (-want +got):\n%s", diff)
		}

		gotOut := outBuf.String()
		wantOut := ""
		if diff := cmp.Diff(wantOut, gotOut); diff != "" {
			t.Errorf("out output mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("with -h", func(t *testing.T) {
		e, cfg, errBuf, outBuf := setUp(t, "-h")

		got := assets.Cmd.Execute(t.Context(), e, cfg)
		want := cli.ExitSuccess
		if got != want {
			t.Errorf("want exit status = %v, but got %v", want, got)
		}

		gotErr := errBuf.String()
		wantErr := ""
		if diff := cmp.Diff(wantErr, gotErr); diff != "" {
			t.Errorf("err output mismatch (-want +got):\n%s", diff)
		}

		gotOut := outBuf.String()
		wantOut := `usage: tilde [root flags] assets

compile frontend assets.

flags:
  -h, -help   show this help and exit
`
		if diff := cmp.Diff(wantOut, gotOut); diff != "" {
			t.Errorf("out output mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("with additional args", func(t *testing.T) {
		e, cfg, errBuf, outBuf := setUp(t, "something")

		got := assets.Cmd.Execute(t.Context(), e, cfg)
		want := cli.ExitUsageError
		if got != want {
			t.Errorf("want exit status = %v, but got %v", want, got)
		}

		gotErr := errBuf.String()
		wantErr := `expected 0 args, but got 1
usage: tilde [root flags] assets
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

	t.Run("with additional flags", func(t *testing.T) {
		e, cfg, errBuf, outBuf := setUp(t, "-x")

		got := assets.Cmd.Execute(t.Context(), e, cfg)
		want := cli.ExitUsageError
		if got != want {
			t.Errorf("want exit status = %v, but got %v", want, got)
		}

		gotErr := errBuf.String()
		wantErr := `flag provided but not defined: -x
usage: tilde [root flags] assets
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

	dir := t.TempDir()
	assetsDir := path.Join(dir, "assets")
	entrypointsDir := path.Join(assetsDir, "entrypoints")
	staticDir := path.Join(dir, "static")
	if err := os.MkdirAll(entrypointsDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(staticDir, 0755); err != nil {
		t.Fatal(err)
	}

	var errBuf, outBuf strings.Builder

	return &cli.Env{
			Stderr: &errBuf,
			Stdout: &outBuf,
			Args:   append([]string{"assets"}, args...),
		}, &core.Config{
			Env:       "test",
			Level:     slog.LevelError,
			Format:    core.JSONFormat,
			AssetsDir: assetsDir,
			StaticDir: staticDir,
		},
		&errBuf,
		&outBuf
}
