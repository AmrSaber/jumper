package tests

import (
	"strings"
	"testing"
)

func TestGetCommand(t *testing.T) {
	cleanup := SetupTest(t)
	defer cleanup()

	t.Run("get prints path for existing bookmark", func(t *testing.T) {
		RunJumperSuccessIn(t, "/tmp", "mark", "dest")
		if out := RunJumperSuccess(t, "get", "dest"); out != "/tmp" {
			t.Errorf("expected '/tmp', got: %s", out)
		}
	})

	t.Run("get with subdirectory appends subdir to bookmark path", func(t *testing.T) {
		RunJumperSuccessIn(t, "/tmp", "mark", "dest")
		if out := RunJumperSuccess(t, "get", "dest/foo/bar"); out != "/tmp/foo/bar" {
			t.Errorf("expected '/tmp/foo/bar', got: %s", out)
		}
	})

	t.Run("get fails for non-existent bookmark", func(t *testing.T) {
		out := RunJumperFailure(t, "get", "nonexistent")
		if !strings.Contains(out, "nonexistent") {
			t.Errorf("expected error to mention bookmark name, got: %s", out)
		}
	})

	t.Run("get with subdirectory fails for non-existent bookmark", func(t *testing.T) {
		out := RunJumperFailure(t, "get", "nonexistent/subdir")
		if !strings.Contains(out, "nonexistent") {
			t.Errorf("expected error to mention bookmark name, got: %s", out)
		}
	})
}
