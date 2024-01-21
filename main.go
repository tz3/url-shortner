package main

import (
	"github.com/tz3/url-shortner/cmd"
	"os"
)

func main() {
	exitCode := cmd.Execute(cmd.Root)
	os.Exit(exitCode)
}
