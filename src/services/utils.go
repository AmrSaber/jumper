package services

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	gap "github.com/muesli/go-app-paths"
)

// Fatal prints a red error message to stderr and exits with code 1.
func Fatal(format string, args ...any) {
	_, _ = color.New(color.FgRed).Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

// Warn prints a yellow warning message to stderr.
func Warn(format string, args ...any) {
	_, _ = color.New(color.FgYellow).Fprintf(os.Stderr, format+"\n", args...)
}

// BookmarksPath returns the path to the bookmarks file.
func BookmarksPath() string {
	scope := gap.NewScope(gap.User, "jumper")
	path, err := scope.DataPath("bookmarks.json")
	if err != nil {
		Fatal("error: cannot determine data directory: %v", err)
	}
	return path
}

func loadBookmarks() []Bookmark {
	path := BookmarksPath()
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return []Bookmark{}
	}
	if err != nil {
		Fatal("error: cannot read bookmarks: %v", err)
	}

	var bookmarks []Bookmark
	if err := json.Unmarshal(data, &bookmarks); err != nil {
		Fatal("error: cannot parse bookmarks: %v", err)
	}
	return bookmarks
}

func saveBookmarks(bookmarks []Bookmark) {
	path := BookmarksPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		Fatal("error: cannot create data directory: %v", err)
	}

	data, err := json.MarshalIndent(bookmarks, "", "  ")
	if err != nil {
		Fatal("error: cannot serialize bookmarks: %v", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		Fatal("error: cannot write bookmarks: %v", err)
	}
}
