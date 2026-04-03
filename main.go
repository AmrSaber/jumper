package main

import "github.com/AmrSaber/jumper/src/cmd"

var version string

func main() {
	if version != "" {
		cmd.SetVersion(version)
	}

	cmd.Execute()
}
