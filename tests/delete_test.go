package tests

import (
	"strings"
	"testing"
)

func TestDeleteCommand(t *testing.T) {
	t.Run("delete existing bookmark", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccess(t, "mark", "todelete")
		RunJumperSuccess(t, "delete", "todelete")
		RunJumperFailure(t, "resolve", "todelete")
	})

	t.Run("delete non-existent bookmark fails", func(t *testing.T) {
		SetupTest(t)
		out := RunJumperFailure(t, "delete", "nonexistent")
		if !strings.Contains(out, "nonexistent") {
			t.Errorf("expected error to mention bookmark name, got: %s", out)
		}
	})

	t.Run("rm alias works", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccess(t, "mark", "rm-test")
		RunJumperSuccess(t, "rm", "rm-test")
		RunJumperFailure(t, "resolve", "rm-test")
	})

	t.Run("unmark alias works", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccess(t, "mark", "unmark-test")
		RunJumperSuccess(t, "unmark", "unmark-test")
		RunJumperFailure(t, "resolve", "unmark-test")
	})

	t.Run("delete with no args uses current directory path", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccessIn(t, "/tmp", "mark", "tmp")
		RunJumperSuccessIn(t, "/tmp", "delete")
		RunJumperFailure(t, "resolve", "tmp")
	})

	t.Run("delete by explicit path", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccessIn(t, "/tmp", "mark", "bypath")
		RunJumperSuccess(t, "delete", "/tmp")
		RunJumperFailure(t, "resolve", "bypath")
	})

	t.Run("delete by path removes all bookmarks for that path", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccessIn(t, "/tmp", "mark", "pathdup1")
		RunJumperSuccessIn(t, "/tmp", "mark", "pathdup2")
		RunJumperSuccess(t, "delete", "/tmp")
		RunJumperFailure(t, "resolve", "pathdup1")
		RunJumperFailure(t, "resolve", "pathdup2")
	})

	t.Run("delete by dot arg uses current directory path", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccessIn(t, "/tmp", "mark", "dottest")
		RunJumperSuccessIn(t, "/tmp", "delete", ".")
		RunJumperFailure(t, "resolve", "dottest")
	})

	t.Run("delete by path fails when no bookmarks found", func(t *testing.T) {
		SetupTest(t)
		out := RunJumperFailure(t, "delete", "/no/such/path")
		if !strings.Contains(out, "/no/such/path") {
			t.Errorf("expected error to mention path, got: %s", out)
		}
	})

	t.Run("delete multiple bookmarks", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccess(t, "mark", "multi1")
		RunJumperSuccess(t, "mark", "multi2")
		RunJumperSuccess(t, "mark", "multi3")
		RunJumperSuccess(t, "delete", "multi1", "multi2", "multi3")
		RunJumperFailure(t, "resolve", "multi1")
		RunJumperFailure(t, "resolve", "multi2")
		RunJumperFailure(t, "resolve", "multi3")
	})

	t.Run("delete multiple fails if any non-existent", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccess(t, "mark", "exists")
		RunJumperFailure(t, "delete", "exists", "nonexistent")
	})

	t.Run("deleted bookmark no longer appears in list", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccess(t, "mark", "gone")
		RunJumperSuccess(t, "delete", "gone")
		out := RunJumperSuccess(t, "list", "--output", "json")
		if strings.Contains(out, `"gone"`) {
			t.Error("deleted bookmark should not appear in list")
		}
	})
}
