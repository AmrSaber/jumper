package main

import "jumper/src/cmd"

var version string

func main() {
	if version != "" {
		cmd.SetVersion(version)
	}

	cmd.Execute()
}
