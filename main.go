package main

import (
	"context"
	"os"
	"runtime/debug"

	"github.com/jonathonwebb/tilde/cmd"
	"github.com/jonathonwebb/x/conf"
)

var version = "0.0.1"

func main() {
	rev, time := getVcsMeta()
	env := conf.DefaultEnv(map[string]any{
		"version":  version,
		"revision": rev,
		"time":     time,
	})
	os.Exit(int(cmd.CLI.Execute(context.Background(), env)))
}

func getVcsMeta() (rev string, time string) {
	info, ok := debug.ReadBuildInfo()
	if ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				rev = setting.Value
			}
			if setting.Key == "vcs.time" {
				time = setting.Value
			}
		}
	}
	return rev, time
}
