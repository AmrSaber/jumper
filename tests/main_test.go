package tests

import (
	"os"
	"os/exec"
	"path"
	"testing"
)

func TestMain(m *testing.M) {
	tmpDir, err := os.MkdirTemp("", "jumper-bin")
	if err != nil {
		panic("Could not create temp directory to build project")
	}

	// Variable declared in helpers.go
	BinaryLocation = path.Join(tmpDir, "jumper")

	// Build the jumper binary before running tests
	buildCmd := exec.Command("go", "build", "-o", BinaryLocation)
	buildCmd.Dir = ".."

	if err := buildCmd.Run(); err != nil {
		panic("Failed to build jumper binary: " + err.Error())
	}

	// Run tests
	code := m.Run()

	// Cleanup (must be before os.Exit since defer won't run)
	_ = os.RemoveAll(tmpDir)

	os.Exit(code)
}
