package migration_test

import (
	"log"
	"log/slog"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/jonathonwebb/tilde/internal/cli"
	"github.com/jonathonwebb/tilde/internal/core"
)

func TestGenMigrationCommand(t *testing.T) {
	// t.Run("")
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
			Log:    log.New(&errBuf, "", 0),
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
