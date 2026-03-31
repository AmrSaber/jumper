package tests

import (
	"strings"
	"testing"
)

func TestResolveCommand(t *testing.T) {
	t.Run("prints path for existing bookmark", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccessIn(t, "/tmp", "mark", "dest")
		if out := RunJumperSuccess(t, "resolve", "dest"); out != "/tmp" {
			t.Errorf("expected '/tmp', got: %s", out)
		}
	})

	t.Run("with subdirectory appends subdir to bookmark path", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccessIn(t, "/tmp", "mark", "dest")
		if out := RunJumperSuccess(t, "resolve", "dest/foo/bar"); out != "/tmp/foo/bar" {
			t.Errorf("expected '/tmp/foo/bar', got: %s", out)
		}
	})

	t.Run("fails for non-existent bookmark", func(t *testing.T) {
		SetupTest(t)
		out := RunJumperFailure(t, "resolve", "nonexistent")
		if !strings.Contains(out, "nonexistent") {
			t.Errorf("expected error to mention bookmark name, got: %s", out)
		}
	})

	t.Run("with subdirectory fails for non-existent bookmark", func(t *testing.T) {
		SetupTest(t)
		out := RunJumperFailure(t, "resolve", "nonexistent/subdir")
		if !strings.Contains(out, "nonexistent") {
			t.Errorf("expected error to mention bookmark name, got: %s", out)
		}
	})

	t.Run("get alias works", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccessIn(t, "/tmp", "mark", "dest")
		if out := RunJumperSuccess(t, "get", "dest"); out != "/tmp" {
			t.Errorf("expected '/tmp', got: %s", out)
		}
	})
}
