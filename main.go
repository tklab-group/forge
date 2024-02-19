package main

import (
	"os"

	"github.com/tklab-group/forge/cli"
)

var version = "No version provided"

func main() {
	cli.SetVersion(version)

	err := cli.Execute()
	if err != nil {
		os.Exit(1)
	}
}
