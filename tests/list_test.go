package tests

import (
	"strings"
	"testing"
)

func TestListCommand(t *testing.T) {
	t.Run("empty list shows helpful message", func(t *testing.T) {
		SetupTest(t)
		out := RunJumperSuccess(t, "list")
		if !strings.Contains(out, "No bookmarks") {
			t.Errorf("expected empty state message, got: %s", out)
		}
	})

	t.Run("list shows all bookmarks", func(t *testing.T) {
		SetupTest(t)
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
		SetupTest(t)
		RunJumperSuccess(t, "mark", "zebra")
		RunJumperSuccess(t, "mark", "apple")
		RunJumperSuccess(t, "mark", "mango")

		out := RunJumperSuccess(t, "list")
		if strings.Index(out, "apple") >= strings.Index(out, "mango") ||
			strings.Index(out, "mango") >= strings.Index(out, "zebra") {
			t.Errorf("expected sorted order apple < mango < zebra, got:\n%s", out)
		}
	})

	t.Run("list with JSON output", func(t *testing.T) {
		SetupTest(t)
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
		SetupTest(t)
		RunJumperSuccess(t, "mark", "yamlmark")
		out := RunJumperSuccess(t, "list", "--output", "yaml")

		for _, field := range []string{"yamlmark", "title:", "path:"} {
			if !strings.Contains(out, field) {
				t.Errorf("expected %q in YAML output", field)
			}
		}
	})

	t.Run("unsupported output format fails", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccess(t, "mark", "dummy")
		out := RunJumperFailure(t, "list", "--output", "xml")
		if !strings.Contains(out, "xml") {
			t.Errorf("expected error to mention format, got: %s", out)
		}
	})

	t.Run("list marks missing paths with not found hint", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccess(t, "mark", "ghost", "/no/such/path")
		out := RunJumperSuccess(t, "list")
		if !strings.Contains(out, "[not found]") {
			t.Errorf("expected '[not found]' hint for missing path, got: %s", out)
		}
	})

	t.Run("list does not mark existing paths with not found hint", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccessIn(t, "/tmp", "mark", "existingpath")
		out := RunJumperSuccess(t, "list")
		for _, line := range strings.Split(out, "\n") {
			if strings.Contains(line, "existingpath") && strings.Contains(line, "[not found]") {
				t.Errorf("expected no '[not found]' hint for existing path row, got: %s", line)
			}
		}
	})

	t.Run("JSON output includes missing:true for stale bookmarks", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccess(t, "mark", "ghost", "/no/such/path")
		out := RunJumperSuccess(t, "list", "--output", "json")
		if !strings.Contains(out, `"missing": true`) {
			t.Errorf("expected missing:true in JSON output, got: %s", out)
		}
	})

	t.Run("JSON output omits missing field for valid bookmarks", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccessIn(t, "/tmp", "mark", "valid")
		out := RunJumperSuccess(t, "list", "--output", "json")
		if strings.Contains(out, `"missing"`) {
			t.Errorf("expected missing field to be omitted for valid bookmark, got: %s", out)
		}
	})

	t.Run("YAML output includes missing:true for stale bookmarks", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccess(t, "mark", "ghost", "/no/such/path")
		out := RunJumperSuccess(t, "list", "--output", "yaml")
		if !strings.Contains(out, "missing: true") {
			t.Errorf("expected missing:true in YAML output, got: %s", out)
		}
	})

	t.Run("YAML output omits missing field for valid bookmarks", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccessIn(t, "/tmp", "mark", "valid")
		out := RunJumperSuccess(t, "list", "--output", "yaml")
		if strings.Contains(out, "missing:") {
			t.Errorf("expected missing field to be omitted for valid bookmark, got: %s", out)
		}
	})

	t.Run("ls alias works", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccess(t, "mark", "aliascheck")
		out := RunJumperSuccess(t, "ls")
		if !strings.Contains(out, "aliascheck") {
			t.Errorf("expected ls alias to work, got: %s", out)
		}
	})
}
