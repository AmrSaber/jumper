package tests

import (
	"strings"
	"testing"
)

func TestDeleteCommand(t *testing.T) {
	cleanup := SetupTest(t)
	defer cleanup()

	t.Run("delete existing bookmark", func(t *testing.T) {
		RunJumperSuccess(t, "mark", "todelete")
		RunJumperSuccess(t, "delete", "todelete")
		RunJumperFailure(t, "get", "todelete")
	})

	t.Run("delete non-existent bookmark fails", func(t *testing.T) {
		out := RunJumperFailure(t, "delete", "nonexistent")
		if !strings.Contains(out, "nonexistent") {
			t.Errorf("expected error to mention bookmark name, got: %s", out)
		}
	})

	t.Run("rm alias works", func(t *testing.T) {
		RunJumperSuccess(t, "mark", "rm-test")
		RunJumperSuccess(t, "rm", "rm-test")
		RunJumperFailure(t, "get", "rm-test")
	})

	t.Run("unmark alias works", func(t *testing.T) {
		RunJumperSuccess(t, "mark", "unmark-test")
		RunJumperSuccess(t, "unmark", "unmark-test")
		RunJumperFailure(t, "get", "unmark-test")
	})

	t.Run("delete with no args uses current directory base name", func(t *testing.T) {
		RunJumperSuccessIn(t, "/tmp", "mark", "tmp")
		RunJumperSuccessIn(t, "/tmp", "delete")
		RunJumperFailure(t, "get", "tmp")
	})

	t.Run("deleted bookmark no longer appears in list", func(t *testing.T) {
		RunJumperSuccess(t, "mark", "gone")
		RunJumperSuccess(t, "delete", "gone")
		out := RunJumperSuccess(t, "list", "--output", "json")
		if strings.Contains(out, `"gone"`) {
			t.Error("deleted bookmark should not appear in list")
		}
	})
}
