package tests

import (
	"strings"
	"testing"
)

func TestMarkCommand(t *testing.T) {
	t.Run("mark with explicit name", func(t *testing.T) {
		SetupTest(t)
		out := RunJumperSuccess(t, "mark", "mydir")
		if !strings.Contains(out, "mydir") {
			t.Errorf("expected output to mention bookmark name, got: %s", out)
		}
	})

	t.Run("mark with default name uses directory base name", func(t *testing.T) {
		SetupTest(t)
		out := RunJumperSuccessIn(t, "/tmp", "mark")
		if !strings.Contains(out, "tmp") {
			t.Errorf("expected bookmark name 'tmp', got: %s", out)
		}
	})

	t.Run("mark records the working directory", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccessIn(t, "/tmp", "mark", "tempdir")
		if path := RunJumperSuccess(t, "get", "tempdir"); path != "/tmp" {
			t.Errorf("expected '/tmp', got: %s", path)
		}
	})

	t.Run("mark overwrites existing bookmark", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccess(t, "mark", "overwrite")
		RunJumperSuccessIn(t, "/tmp", "mark", "overwrite")
		if path := RunJumperSuccess(t, "get", "overwrite"); path != "/tmp" {
			t.Errorf("expected updated path '/tmp', got: %s", path)
		}
	})

	t.Run("mark with explicit name and directory", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccess(t, "mark", "explicit-dir", "/tmp")
		if path := RunJumperSuccess(t, "get", "explicit-dir"); path != "/tmp" {
			t.Errorf("expected '/tmp', got: %s", path)
		}
	})

	t.Run("mark with non-existent directory warns", func(t *testing.T) {
		SetupTest(t)
		out := RunJumperSuccess(t, "mark", "ghost", "/no/such/dir")
		if !strings.Contains(out, "warning") {
			t.Errorf("expected a warning, got: %s", out)
		}
	})

	t.Run("mark with relative path resolves to absolute", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccessIn(t, "/tmp", "mark", "relpath", "./")
		if path := RunJumperSuccess(t, "get", "relpath"); path != "/tmp" {
			t.Errorf("expected '/tmp', got: %s", path)
		}
	})

	t.Run("mark with path-like name fails", func(t *testing.T) {
		SetupTest(t)
		RunJumperFailure(t, "mark", ".hidden")
		RunJumperFailure(t, "mark", "~/something")
		RunJumperFailure(t, "mark", "/absolute")
	})

	t.Run("mark with reserved name '-' fails", func(t *testing.T) {
		SetupTest(t)
		out := RunJumperFailure(t, "mark", "-")
		if !strings.Contains(out, "reserved") {
			t.Errorf("expected error to mention reserved name, got: %s", out)
		}
	})
}
