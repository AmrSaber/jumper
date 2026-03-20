package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"jumper/src/services"

	"github.com/spf13/cobra"
)

var completeCmd = &cobra.Command{
	Use:    "complete [prefix]",
	Hidden: true,
	Args:   cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		prefix := ""
		if len(args) == 1 {
			prefix = args[0]
		}

		// If the prefix contains a slash, expand subdirectories under the bookmark.
		// e.g. "proj/src" -> bookmark "proj" + subpath "src"
		if bookmarkTitle, subPrefix, ok := strings.Cut(prefix, "/"); ok {
			bookmark, found := services.Get(bookmarkTitle)
			if !found {
				return
			}

			dir, namePrefix := filepath.Split(subPrefix)
			entries, err := os.ReadDir(filepath.Join(bookmark.Path, dir))
			if err != nil {
				return
			}
			lower := strings.ToLower(namePrefix)
			for _, entry := range entries {
				entryPath := filepath.Join(bookmark.Path, dir, entry.Name())
				info, err := os.Stat(entryPath)
				if err != nil || !info.IsDir() {
					continue
				}
				if strings.HasPrefix(strings.ToLower(entry.Name()), lower) {
					fmt.Println(bookmark.Title + "/" + filepath.Join(dir, entry.Name()) + "/")
				}
			}
			return
		}

		// Otherwise, list bookmark titles matching the prefix (case-insensitive).
		for _, c := range bookmarkCompletions(prefix) {
			fmt.Println(c + "/")
		}
	},
}

// bookmarkCompletions returns bookmark titles that start with the given prefix (case-insensitive).
func bookmarkCompletions(prefix string) []cobra.Completion {
	lower := strings.ToLower(prefix)
	var matches []cobra.Completion
	for _, b := range services.List() {
		if strings.HasPrefix(strings.ToLower(b.Title), lower) {
			matches = append(matches, b.Title)
		}
	}
	return matches
}

func init() {
	rootCmd.AddCommand(completeCmd)
}
