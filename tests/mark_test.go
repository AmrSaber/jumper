package tests

import (
	"encoding/json"
	"os"
	"path/filepath"
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
		if path := RunJumperSuccess(t, "resolve", "tempdir"); path != "/tmp" {
			t.Errorf("expected '/tmp', got: %s", path)
		}
	})

	t.Run("mark overwrites existing bookmark", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccess(t, "mark", "overwrite")
		RunJumperSuccessIn(t, "/tmp", "mark", "overwrite")
		if path := RunJumperSuccess(t, "resolve", "overwrite"); path != "/tmp" {
			t.Errorf("expected updated path '/tmp', got: %s", path)
		}
	})

	t.Run("mark with explicit name and directory", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccess(t, "mark", "explicit-dir", "/tmp")
		if path := RunJumperSuccess(t, "resolve", "explicit-dir"); path != "/tmp" {
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
		if path := RunJumperSuccess(t, "resolve", "relpath"); path != "/tmp" {
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

	t.Run("mark with comment stores it", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccess(t, "mark", "commented", "/tmp", "--comment", "my note")
		if got := storedComment(t, "commented"); got != "my note" {
			t.Errorf("expected comment 'my note', got: %q", got)
		}
	})

	t.Run("mark with short comment flag", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccess(t, "mark", "short", "/tmp", "-c", "short note")
		if got := storedComment(t, "short"); got != "short note" {
			t.Errorf("expected comment 'short note', got: %q", got)
		}
	})

	t.Run("mark without comment stores none", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccess(t, "mark", "plain", "/tmp")
		if got := storedComment(t, "plain"); got != "" {
			t.Errorf("expected empty comment, got: %q", got)
		}
	})

	t.Run("re-mark replaces comment", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccess(t, "mark", "redo", "/tmp", "-c", "first")
		RunJumperSuccess(t, "mark", "redo", "/tmp", "-c", "second")
		if got := storedComment(t, "redo"); got != "second" {
			t.Errorf("expected comment 'second', got: %q", got)
		}
		RunJumperSuccess(t, "mark", "redo", "/tmp")
		if got := storedComment(t, "redo"); got != "" {
			t.Errorf("expected comment cleared, got: %q", got)
		}
	})
}

// storedComment reads the comment of a bookmark straight from the data file,
// because no command prints comments yet.
func storedComment(t *testing.T, title string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(os.Getenv("XDG_DATA_HOME"), "jumper", "bookmarks.json"))
	if err != nil {
		t.Fatal(err)
	}
	var bookmarks []struct{ Title, Comment string }
	if err := json.Unmarshal(data, &bookmarks); err != nil {
		t.Fatal(err)
	}
	for _, b := range bookmarks {
		if b.Title == title {
			return b.Comment
		}
	}
	t.Fatalf("bookmark %q not found", title)
	return ""
}
