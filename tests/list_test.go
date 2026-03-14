package tests

import (
	"strings"
	"testing"
)

func TestListCommand(t *testing.T) {
	cleanup := SetupTest(t)
	defer cleanup()

	t.Run("empty list shows helpful message", func(t *testing.T) {
		out := RunJumperSuccess(t, "list")
		if !strings.Contains(out, "No bookmarks") {
			t.Errorf("expected empty state message, got: %s", out)
		}
	})

	t.Run("list shows all bookmarks", func(t *testing.T) {
		RunJumperSuccess(t, "mark", "alpha")
		RunJumperSuccess(t, "mark", "beta")
		RunJumperSuccess(t, "mark", "gamma")

		out := RunJumperSuccess(t, "list")
		for _, name := range []string{"alpha", "beta", "gamma"} {
			if !strings.Contains(out, name) {
				t.Errorf("expected %q in output", name)
			}
		}
	})

	t.Run("list is sorted by title", func(t *testing.T) {
		RunJumperSuccess(t, "mark", "zebra")
		RunJumperSuccess(t, "mark", "apple")
		RunJumperSuccess(t, "mark", "mango")

		out := RunJumperSuccess(t, "list")
		if !(strings.Index(out, "apple") < strings.Index(out, "mango") &&
			strings.Index(out, "mango") < strings.Index(out, "zebra")) {
			t.Errorf("expected sorted order apple < mango < zebra, got:\n%s", out)
		}
	})

	t.Run("list with JSON output", func(t *testing.T) {
		RunJumperSuccess(t, "mark", "jsonmark")
		out := RunJumperSuccess(t, "list", "--output", "json")

		if !strings.HasPrefix(out, "[") {
			t.Errorf("expected JSON array, got: %s", out)
		}
		for _, field := range []string{"jsonmark", `"title"`, `"path"`} {
			if !strings.Contains(out, field) {
				t.Errorf("expected %q in JSON output", field)
			}
		}
	})

	t.Run("list with YAML output", func(t *testing.T) {
		RunJumperSuccess(t, "mark", "yamlmark")
		out := RunJumperSuccess(t, "list", "--output", "yaml")

		for _, field := range []string{"yamlmark", "title:", "path:"} {
			if !strings.Contains(out, field) {
				t.Errorf("expected %q in YAML output", field)
			}
		}
	})

	t.Run("unsupported output format fails", func(t *testing.T) {
		out := RunJumperFailure(t, "list", "--output", "xml")
		if !strings.Contains(out, "xml") {
			t.Errorf("expected error to mention format, got: %s", out)
		}
	})

	t.Run("ls alias works", func(t *testing.T) {
		RunJumperSuccess(t, "mark", "aliascheck")
		out := RunJumperSuccess(t, "ls")
		if !strings.Contains(out, "aliascheck") {
			t.Errorf("expected ls alias to work, got: %s", out)
		}
	})
}
