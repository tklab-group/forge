package main

import (
	"fmt"
	"github.com/tklab-group/forge/cli"
	"os"
)

func main() {
	err := cli.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
