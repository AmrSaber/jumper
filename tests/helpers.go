// Package tests contains integration tests for the tool
package tests

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

var BinaryLocation string

func run(t *testing.T, dir string, args ...string) (string, error) {
	t.Helper()
	cmd := exec.Command(BinaryLocation, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	output, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

func RunJumper(t *testing.T, args ...string) (string, error) {
	t.Helper()
	return run(t, "", args...)
}

func RunJumperIn(t *testing.T, dir string, args ...string) (string, error) {
	t.Helper()
	return run(t, dir, args...)
}

func RunJumperSuccess(t *testing.T, args ...string) string {
	t.Helper()
	output, err := RunJumper(t, args...)
	if err != nil {
		t.Fatalf("jumper %v failed: %v\n%s", args, err, output)
	}
	return output
}

func RunJumperSuccessIn(t *testing.T, dir string, args ...string) string {
	t.Helper()
	output, err := RunJumperIn(t, dir, args...)
	if err != nil {
		t.Fatalf("jumper %v (in %s) failed: %v\n%s", args, dir, err, output)
	}
	return output
}

func RunJumperFailure(t *testing.T, args ...string) string {
	t.Helper()
	output, err := RunJumper(t, args...)
	if err == nil {
		t.Fatalf("jumper %v should have failed\n%s", args, output)
	}
	return output
}

// SetupTest redirects the data directory to a temp folder for test isolation.
// Cleanup is registered automatically via t.Cleanup.
func SetupTest(t *testing.T) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "jumper-test-*")
	if err != nil {
		t.Fatal(err)
	}
	_ = os.Setenv("XDG_DATA_HOME", tmpDir)
	t.Cleanup(func() {
		_ = os.Unsetenv("XDG_DATA_HOME")
		_ = os.RemoveAll(tmpDir)
	})
}
