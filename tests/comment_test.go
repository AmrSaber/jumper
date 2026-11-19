package tests

import (
	"strings"
	"testing"
)

func TestCommentCommand(t *testing.T) {
	t.Run("sets comment on bookmark without one", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccess(t, "mark", "proj", "/tmp")
		RunJumperSuccess(t, "comment", "proj", "hello")
		if got := storedComment(t, "proj"); got != "hello" {
			t.Errorf("expected 'hello', got: %q", got)
		}
	})

	t.Run("replaces existing comment and keeps path", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccess(t, "mark", "proj", "/tmp", "-c", "old")
		RunJumperSuccess(t, "comment", "proj", "new")
		if got := storedComment(t, "proj"); got != "new" {
			t.Errorf("expected 'new', got: %q", got)
		}
		if path := RunJumperSuccess(t, "resolve", "proj"); path != "/tmp" {
			t.Errorf("expected path '/tmp', got: %s", path)
		}
	})

	t.Run("is case-insensitive on name", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccess(t, "mark", "CaseMark", "/tmp")
		RunJumperSuccess(t, "comment", "casemark", "x")
		if got := storedComment(t, "CaseMark"); got != "x" {
			t.Errorf("expected 'x', got: %q", got)
		}
	})

	t.Run("--no-comment removes comment", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccess(t, "mark", "proj", "/tmp", "-c", "old")
		RunJumperSuccess(t, "comment", "proj", "--no-comment")
		if got := storedComment(t, "proj"); got != "" {
			t.Errorf("expected empty comment, got: %q", got)
		}
	})

	t.Run("missing comment without --no-comment fails", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccess(t, "mark", "proj", "/tmp", "-c", "keep")
		out := RunJumperFailure(t, "comment", "proj")
		if !strings.Contains(out, "--no-comment") {
			t.Errorf("expected error to mention --no-comment, got: %s", out)
		}
		if got := storedComment(t, "proj"); got != "keep" {
			t.Errorf("comment should be untouched, got: %q", got)
		}
	})

	t.Run("comment together with --no-comment fails", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccess(t, "mark", "proj", "/tmp", "-c", "keep")
		RunJumperFailure(t, "comment", "proj", "text", "--no-comment")
		if got := storedComment(t, "proj"); got != "keep" {
			t.Errorf("comment should be untouched, got: %q", got)
		}
	})

	t.Run("non-existent bookmark fails", func(t *testing.T) {
		SetupTest(t)
		out := RunJumperFailure(t, "comment", "ghost", "text")
		if !strings.Contains(out, "ghost") {
			t.Errorf("expected error to mention name, got: %s", out)
		}
		RunJumperFailure(t, "comment", "ghost", "--no-comment")
	})

	t.Run("newlines in comment are replaced with spaces", func(t *testing.T) {
		SetupTest(t)
		RunJumperSuccess(t, "mark", "proj", "/tmp", "-c", "a\nb")
		if got := storedComment(t, "proj"); got != "a b" {
			t.Errorf("expected 'a b', got: %q", got)
		}
		RunJumperSuccess(t, "comment", "proj", "c\r\nd\n")
		if got := storedComment(t, "proj"); got != "c d" {
			t.Errorf("expected 'c d', got: %q", got)
		}
	})
}
