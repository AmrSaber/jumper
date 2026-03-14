package cmd

import (
	"bytes"
	"os"
	"os/exec"
)

var version string

// SetVersion sets the version string (typically injected by goreleaser at build time).
func SetVersion(v string) {
	version = v
}

func getVersion() string {
	if version != "" {
		return version
	}

	if gitVersion := getVersionFromGit(); gitVersion != "" {
		return gitVersion
	}

	return "??"
}

func getVersionFromGit() string {
	if _, err := os.Stat(".git"); err != nil {
		return ""
	}

	tag, err := exec.Command("git", "describe", "--tags").Output()
	if err != nil {
		return ""
	}

	tag = bytes.TrimSpace(tag)

	if isRepoDirty() {
		tag = append(tag, '+')
	}

	return string(tag)
}

func isRepoDirty() bool {
	output, _ := exec.Command("git", "status", "--porcelain").Output()
	return len(bytes.TrimSpace(output)) != 0
}
