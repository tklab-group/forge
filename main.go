package main

import (
	"github.com/tklab-group/forge/cli"
	"os"
)

func main() {
	err := cli.Execute()
	if err != nil {
		os.Exit(1)
	}
}
