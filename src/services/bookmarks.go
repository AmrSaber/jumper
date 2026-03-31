// Package services for io operations
package services

import (
	"os"
	"strings"
)

type Bookmark struct {
	Title string `json:"title"`
	Path  string `json:"path"`
}

// Upsert adds or updates a bookmark with the given title and path.
func Upsert(title, path string) {
	bookmarks := loadBookmarks()

	for i, b := range bookmarks {
		if strings.EqualFold(b.Title, title) {
			bookmarks[i].Path = path
			saveBookmarks(bookmarks)
			return
		}
	}

	saveBookmarks(append(bookmarks, Bookmark{Title: title, Path: path}))
}

// Get returns the bookmark with the given title (case-insensitive), and whether it exists.
func Get(title string) (Bookmark, bool) {
	for _, b := range loadBookmarks() {
		if strings.EqualFold(b.Title, title) {
			return b, true
		}
	}

	return Bookmark{}, false
}

// List returns all bookmarks.
func List() []Bookmark {
	return loadBookmarks()
}

// Rename renames a bookmark from oldTitle to newTitle (case-insensitive match on oldTitle).
// Returns false if oldTitle does not exist.
func Rename(oldTitle, newTitle string) bool {
	bookmarks := loadBookmarks()

	for i, b := range bookmarks {
		if strings.EqualFold(b.Title, oldTitle) {
			bookmarks[i].Title = newTitle
			saveBookmarks(bookmarks)
			return true
		}
	}

	return false
}

// Delete removes a bookmark by title (case-insensitive). Returns true if it existed.
func Delete(title string) bool {
	bookmarks := loadBookmarks()

	for i, b := range bookmarks {
		if strings.EqualFold(b.Title, title) {
			saveBookmarks(append(bookmarks[:i], bookmarks[i+1:]...))
			return true
		}
	}

	return false
}

// Prune removes all bookmarks whose paths no longer exist. Returns the deleted bookmarks.
func Prune() []Bookmark {
	bookmarks := loadBookmarks()

	remaining := make([]Bookmark, 0, len(bookmarks))
	var deleted []Bookmark

	for _, bookmark := range bookmarks {
		if _, err := os.Stat(bookmark.Path); os.IsNotExist(err) {
			deleted = append(deleted, bookmark)
		} else {
			remaining = append(remaining, bookmark)
		}
	}

	if len(deleted) > 0 {
		saveBookmarks(remaining)
	}

	return deleted
}

// DeleteByPath removes all bookmarks pointing to the given path. Returns the deleted bookmarks.
func DeleteByPath(path string) []Bookmark {
	bookmarks := loadBookmarks()

	remaining := make([]Bookmark, 0, len(bookmarks))
	var deleted []Bookmark

	for _, bookmark := range bookmarks {
		if bookmark.Path == path {
			deleted = append(deleted, bookmark)
		} else {
			remaining = append(remaining, bookmark)
		}
	}

	if len(deleted) > 0 {
		saveBookmarks(remaining)
	}

	return deleted
}
