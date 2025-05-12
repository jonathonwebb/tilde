package main

import (
	"context"
	"os"

	"github.com/jonathonwebb/tilde/cmd"
	"github.com/jonathonwebb/tilde/internal/cli"
)

var Build string

func main() {
	ctx := context.Background()
	os.Exit(int(cmd.Execute(ctx, cli.DefaultEnv(Build))))
}
