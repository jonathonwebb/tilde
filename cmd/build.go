package main

import "runtime/debug"

var (
	version      = "0.0.1"
	rev, revTime = vcs()
)

func vcs() (rev string, time string) {
	info, ok := debug.ReadBuildInfo()
	if ok {
		for _, s := range info.Settings {
			if s.Key == "vcs.revision" {
				rev = s.Value
			}
			if s.Key == "vcs.time" {
				time = s.Value
			}
		}
	}
	return rev, time
}
