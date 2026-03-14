package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	gap "github.com/muesli/go-app-paths"
)

// BookmarksPath returns the path to the bookmarks file.
func BookmarksPath() string {
	scope := gap.NewScope(gap.User, "jumper")
	path, err := scope.DataPath("bookmarks.json")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: cannot determine data directory:", err)
		os.Exit(1)
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
		fmt.Fprintln(os.Stderr, "error: cannot read bookmarks:", err)
		os.Exit(1)
	}

	var bookmarks []Bookmark
	if err := json.Unmarshal(data, &bookmarks); err != nil {
		fmt.Fprintln(os.Stderr, "error: cannot parse bookmarks:", err)
		os.Exit(1)
	}
	return bookmarks
}

func saveBookmarks(bookmarks []Bookmark) {
	path := BookmarksPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		fmt.Fprintln(os.Stderr, "error: cannot create data directory:", err)
		os.Exit(1)
	}

	data, err := json.MarshalIndent(bookmarks, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: cannot serialize bookmarks:", err)
		os.Exit(1)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "error: cannot write bookmarks:", err)
		os.Exit(1)
	}
}
