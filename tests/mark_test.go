package tests

import (
	"strings"
	"testing"
)

func TestMarkCommand(t *testing.T) {
	cleanup := SetupTest(t)
	defer cleanup()

	t.Run("mark with explicit name", func(t *testing.T) {
		out := RunJumperSuccess(t, "mark", "mydir")
		if !strings.Contains(out, "mydir") {
			t.Errorf("expected output to mention bookmark name, got: %s", out)
		}
	})

	t.Run("mark with default name uses directory base name", func(t *testing.T) {
		out := RunJumperSuccessIn(t, "/tmp", "mark")
		if !strings.Contains(out, "tmp") {
			t.Errorf("expected bookmark name 'tmp', got: %s", out)
		}
	})

	t.Run("mark records the working directory", func(t *testing.T) {
		RunJumperSuccessIn(t, "/tmp", "mark", "tempdir")
		if path := RunJumperSuccess(t, "get", "tempdir"); path != "/tmp" {
			t.Errorf("expected '/tmp', got: %s", path)
		}
	})

	t.Run("mark overwrites existing bookmark", func(t *testing.T) {
		RunJumperSuccess(t, "mark", "overwrite")
		RunJumperSuccessIn(t, "/tmp", "mark", "overwrite")
		if path := RunJumperSuccess(t, "get", "overwrite"); path != "/tmp" {
			t.Errorf("expected updated path '/tmp', got: %s", path)
		}
	})

	t.Run("mark with explicit name and directory", func(t *testing.T) {
		RunJumperSuccess(t, "mark", "explicit-dir", "/tmp")
		if path := RunJumperSuccess(t, "get", "explicit-dir"); path != "/tmp" {
			t.Errorf("expected '/tmp', got: %s", path)
		}
	})

	t.Run("mark with dot-prefixed name fails", func(t *testing.T) {
		out := RunJumperFailure(t, "mark", ".hidden")
		if !strings.Contains(out, ".") {
			t.Errorf("expected error to mention dot prefix, got: %s", out)
		}
	})

}
