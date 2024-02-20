package main

import (
	"os"
	"runtime/debug"

	"github.com/tklab-group/forge/cli"
)

var version string

func main() {
	// Version setting for installing with `go install`
	if version == "" {
		buildInfo, ok := debug.ReadBuildInfo()
		if ok {
			version = buildInfo.Main.Version
		} else {
			version = "No version provided"
		}
	}

	cli.SetVersion(version)

	err := cli.Execute()
	if err != nil {
		os.Exit(1)
	}
}
