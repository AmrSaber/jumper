package tests

import (
	"strings"
	"testing"
)

func TestRenameCommand(t *testing.T) {
	cleanup := SetupTest(t)
	defer cleanup()

	t.Run("rename existing bookmark", func(t *testing.T) {
		RunJumperSuccess(t, "mark", "oldname")
		RunJumperSuccess(t, "rename", "oldname", "newname")
		RunJumperFailure(t, "get", "oldname")
		RunJumperSuccess(t, "get", "newname")
	})

	t.Run("rename non-existent bookmark fails", func(t *testing.T) {
		out := RunJumperFailure(t, "rename", "ghost", "newname")
		if !strings.Contains(out, "ghost") {
			t.Errorf("expected error to mention old name, got: %s", out)
		}
	})

	t.Run("mv alias works", func(t *testing.T) {
		RunJumperSuccess(t, "mark", "mv-test")
		RunJumperSuccess(t, "mv", "mv-test", "mv-renamed")
		RunJumperFailure(t, "get", "mv-test")
		RunJumperSuccess(t, "get", "mv-renamed")
	})

	t.Run("rename is case-insensitive on old name", func(t *testing.T) {
		RunJumperSuccess(t, "mark", "CaseMark")
		RunJumperSuccess(t, "rename", "casemark", "renamed-case")
		RunJumperSuccess(t, "get", "renamed-case")
	})

	t.Run("rename to dot-prefixed name fails", func(t *testing.T) {
		RunJumperSuccess(t, "mark", "dot-rename-test")
		RunJumperFailure(t, "rename", "dot-rename-test", ".hidden")
	})

	t.Run("renamed bookmark appears in list with new name", func(t *testing.T) {
		RunJumperSuccess(t, "mark", "before")
		RunJumperSuccess(t, "rename", "before", "after")
		out := RunJumperSuccess(t, "list", "--output", "json")
		if strings.Contains(out, `"before"`) {
			t.Error("old name should not appear in list after rename")
		}
		if !strings.Contains(out, `"after"`) {
			t.Error("new name should appear in list after rename")
		}
	})
}
