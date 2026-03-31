package tests

import (
	"os"
	"strings"
	"testing"
)

func TestPruneCommand(t *testing.T) {
	t.Run("no stale bookmarks reports message", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccess(t, "mark", "tmp", "/tmp")
		out := RunJumperSuccess(t, "prune")
		if !strings.Contains(out, "No stale bookmarks found") {
			t.Errorf("expected no-stale message, got: %s", out)
		}
	})

	t.Run("removes stale bookmarks", func(t *testing.T) {
		SetupTest(t)

		dir, err := os.MkdirTemp("", "jumper-prune-*")
		if err != nil {
			t.Fatal(err)
		}

		RunJumperSuccess(t, "mark", "stale", dir)
		_ = os.RemoveAll(dir)

		out := RunJumperSuccess(t, "prune")
		if !strings.Contains(out, "stale") {
			t.Errorf("expected deleted bookmark name in output, got: %s", out)
		}

		RunJumperFailure(t, "resolve", "stale")
	})

	t.Run("keeps valid bookmarks intact", func(t *testing.T) {
		SetupTest(t)

		dir, err := os.MkdirTemp("", "jumper-prune-*")
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = os.RemoveAll(dir) })

		RunJumperSuccess(t, "mark", "valid", dir)
		RunJumperSuccess(t, "prune")
		RunJumperSuccess(t, "resolve", "valid")
	})

	t.Run("removes only stale, keeps valid", func(t *testing.T) {
		SetupTest(t)

		dir, err := os.MkdirTemp("", "jumper-prune-*")
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = os.RemoveAll(dir) })

		staleDir, err := os.MkdirTemp("", "jumper-prune-stale-*")
		if err != nil {
			t.Fatal(err)
		}

		RunJumperSuccess(t, "mark", "valid", dir)
		RunJumperSuccess(t, "mark", "stale", staleDir)
		_ = os.RemoveAll(staleDir)

		RunJumperSuccess(t, "prune")

		RunJumperSuccess(t, "resolve", "valid")
		RunJumperFailure(t, "resolve", "stale")
	})

	t.Run("clean alias works", func(t *testing.T) {
		SetupTest(t)
		out := RunJumperSuccess(t, "clean")
		if !strings.Contains(out, "No stale bookmarks found") {
			t.Errorf("expected no-stale message, got: %s", out)
		}
	})
}
