// Package common for the common utils and helpers
package common

import (
	"os"
	"path/filepath"

	"github.com/AmrSaber/jumper/src/services"
)

var SupportedAgents = []string{"opencode"}

func OpencodePaths() (string, string) {
	configBase, err := os.UserConfigDir()
	if err != nil {
		services.Fatal("error: cannot determine config directory: %v", err)
	}

	skillPath := filepath.Join(configBase, "opencode", "skills", "jumper", "SKILL.md")
	pluginPath := filepath.Join(configBase, "opencode", "plugins", "jumper-inject.ts")

	return skillPath, pluginPath
}
