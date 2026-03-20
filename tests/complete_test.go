package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func mustMkdir(t *testing.T, path string) {
	t.Helper()
	if err := os.Mkdir(path, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
}

func mustMkdirAll(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("mkdirall %s: %v", path, err)
	}
}

func mustSymlink(t *testing.T, target, link string) {
	t.Helper()
	if err := os.Symlink(target, link); err != nil {
		t.Fatalf("symlink %s -> %s: %v", link, target, err)
	}
}

func mustWriteFile(t *testing.T, path string) {
	t.Helper()
	if err := os.WriteFile(path, []byte{}, 0o644); err != nil {
		t.Fatalf("writefile %s: %v", path, err)
	}
}

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "jumper-complete-*")
	if err != nil {
		t.Fatalf("mkdirtemp: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(dir) })
	return dir
}

func TestCompleteCommand(t *testing.T) {
	t.Run("no prefix returns all bookmarks", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccess(t, "mark", "alpha")
		RunJumperSuccess(t, "mark", "beta")
		RunJumperSuccess(t, "mark", "gamma")

		out := RunJumperSuccess(t, "complete")
		for _, name := range []string{"alpha/", "beta/", "gamma/"} {
			if !strings.Contains(out, name) {
				t.Errorf("expected %q in completions, got:\n%s", name, out)
			}
		}
	})

	t.Run("prefix filters bookmark titles case-insensitively", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccess(t, "mark", "projects")
		RunJumperSuccess(t, "mark", "personal")
		RunJumperSuccess(t, "mark", "work")

		out := RunJumperSuccess(t, "complete", "P")
		if !strings.Contains(out, "projects/") {
			t.Errorf("expected 'projects/' in completions, got:\n%s", out)
		}
		if !strings.Contains(out, "personal/") {
			t.Errorf("expected 'personal/' in completions, got:\n%s", out)
		}
		if strings.Contains(out, "work/") {
			t.Errorf("expected 'work/' to be filtered out, got:\n%s", out)
		}
	})

	t.Run("no match returns empty output", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccess(t, "mark", "alpha")

		out := RunJumperSuccess(t, "complete", "zzz")
		if strings.TrimSpace(out) != "" {
			t.Errorf("expected empty output, got:\n%s", out)
		}
	})

	t.Run("subpath lists subdirectories of bookmark", func(t *testing.T) {
		SetupTest(t)
		base := tempDir(t)
		mustMkdir(t, filepath.Join(base, "src"))
		mustMkdir(t, filepath.Join(base, "docs"))

		RunJumperSuccess(t, "mark", "proj", base)
		out := RunJumperSuccess(t, "complete", "proj/")

		if !strings.Contains(out, "proj/src/") {
			t.Errorf("expected 'proj/src/' in completions, got:\n%s", out)
		}
		if !strings.Contains(out, "proj/docs/") {
			t.Errorf("expected 'proj/docs/' in completions, got:\n%s", out)
		}
	})

	t.Run("subpath prefix filters subdirectories", func(t *testing.T) {
		SetupTest(t)
		base := tempDir(t)
		mustMkdir(t, filepath.Join(base, "src"))
		mustMkdir(t, filepath.Join(base, "scripts"))
		mustMkdir(t, filepath.Join(base, "docs"))

		RunJumperSuccess(t, "mark", "proj", base)
		out := RunJumperSuccess(t, "complete", "proj/s")

		if !strings.Contains(out, "proj/src/") {
			t.Errorf("expected 'proj/src/' in completions, got:\n%s", out)
		}
		if !strings.Contains(out, "proj/scripts/") {
			t.Errorf("expected 'proj/scripts/' in completions, got:\n%s", out)
		}
		if strings.Contains(out, "proj/docs/") {
			t.Errorf("expected 'proj/docs/' to be filtered out, got:\n%s", out)
		}
	})

	t.Run("subpath does not include files", func(t *testing.T) {
		SetupTest(t)
		base := tempDir(t)
		mustMkdir(t, filepath.Join(base, "subdir"))
		mustWriteFile(t, filepath.Join(base, "file.txt"))

		RunJumperSuccess(t, "mark", "proj", base)
		out := RunJumperSuccess(t, "complete", "proj/")

		if strings.Contains(out, "file.txt") {
			t.Errorf("expected files to be excluded, got:\n%s", out)
		}
	})

	t.Run("subpath follows symlinks to directories", func(t *testing.T) {
		SetupTest(t)
		base := tempDir(t)
		target := tempDir(t)
		mustMkdir(t, filepath.Join(base, "realdir"))
		mustSymlink(t, target, filepath.Join(base, "linkeddir"))

		RunJumperSuccess(t, "mark", "cfg", base)
		out := RunJumperSuccess(t, "complete", "cfg/")

		if !strings.Contains(out, "cfg/realdir/") {
			t.Errorf("expected 'cfg/realdir/' in completions, got:\n%s", out)
		}
		if !strings.Contains(out, "cfg/linkeddir/") {
			t.Errorf("expected symlinked 'cfg/linkeddir/' in completions, got:\n%s", out)
		}
	})

	t.Run("subpath does not follow symlinks to files", func(t *testing.T) {
		SetupTest(t)
		base := tempDir(t)
		f, err := os.CreateTemp("", "jumper-file-*")
		if err != nil {
			t.Fatalf("createtemp: %v", err)
		}
		if err := f.Close(); err != nil {
			t.Fatalf("close: %v", err)
		}
		t.Cleanup(func() { _ = os.Remove(f.Name()) })
		mustSymlink(t, f.Name(), filepath.Join(base, "linkedfile"))

		RunJumperSuccess(t, "mark", "cfg", base)
		out := RunJumperSuccess(t, "complete", "cfg/")

		if strings.Contains(out, "linkedfile") {
			t.Errorf("expected symlink to file to be excluded, got:\n%s", out)
		}
	})

	t.Run("unknown bookmark in subpath returns empty output", func(t *testing.T) {
		SetupTest(t)
		out := RunJumperSuccess(t, "complete", "nonexistent/sub")
		if strings.TrimSpace(out) != "" {
			t.Errorf("expected empty output for unknown bookmark, got:\n%s", out)
		}
	})

	t.Run("nested subpath lists deeper subdirectories", func(t *testing.T) {
		SetupTest(t)
		base := tempDir(t)
		mustMkdirAll(t, filepath.Join(base, "src", "internal"))
		mustMkdirAll(t, filepath.Join(base, "src", "cmd"))

		RunJumperSuccess(t, "mark", "proj", base)
		out := RunJumperSuccess(t, "complete", "proj/src/")

		if !strings.Contains(out, "proj/src/internal/") {
			t.Errorf("expected 'proj/src/internal/' in completions, got:\n%s", out)
		}
		if !strings.Contains(out, "proj/src/cmd/") {
			t.Errorf("expected 'proj/src/cmd/' in completions, got:\n%s", out)
		}
	})
}
