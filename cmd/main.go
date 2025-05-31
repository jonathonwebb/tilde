package main

import (
	"context"
	"os"

	"github.com/jonathonwebb/tilde/cmd/root"
	"github.com/jonathonwebb/tilde/internal/cli"
	"github.com/jonathonwebb/tilde/internal/core"
)

func main() {
	var cfg core.Config

	ctx := context.Background()
	env := cli.DefaultEnv(map[string]any{
		"version": version,
		"rev":     rev,
		"time":    revTime,
	})

	os.Exit(int(root.Cmd.Execute(ctx, env, &cfg)))
}
