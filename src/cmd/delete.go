package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AmrSaber/jumper/src/services"

	"github.com/spf13/cobra"
)

func isPath(s string) bool {
	return filepath.IsAbs(s) ||
		strings.HasPrefix(s, "~") ||
		strings.HasPrefix(s, ".")
}

func resolvePath(s string) (string, error) {
	if after, ok := strings.CutPrefix(s, "~"); ok {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		s = filepath.Join(home, after)
	}

	return filepath.Abs(s)
}

func deleteByPath(path string) {
	deleted := services.DeleteByPath(path)
	if len(deleted) == 0 {
		services.Fatal("error: no bookmarks found for path %q", path)
	}

	for _, bookmark := range deleted {
		fmt.Printf("Deleted bookmark %q\n", bookmark.Title)
	}
}

var deleteCmd = &cobra.Command{
	Use:   "delete [name|path...]",
	Short: "Delete one or more bookmarks",
	Long: `Delete bookmarks by name or path.

If the argument looks like a path (starts with /, ~, or .) it is treated as a
path and all bookmarks pointing to that directory are deleted. Otherwise it is
treated as a bookmark name.

If no argument is provided, all bookmarks pointing to the current directory are deleted.`,
	Aliases: []string{"del", "rm", "unmark"},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
		already := make(map[string]bool, len(args))
		for _, a := range args {
			already[a] = true
		}

		var completions []cobra.Completion
		for _, c := range bookmarkCompletions(toComplete) {
			if !already[c] {
				completions = append(completions, c)
			}
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			dir, err := os.Getwd()
			if err != nil {
				services.Fatal("error: cannot get current directory: %v", err)
			}

			deleteByPath(dir)
			return
		}

		for _, arg := range args {
			if isPath(arg) {
				path, err := resolvePath(arg)
				if err != nil {
					services.Fatal("error: cannot resolve path %q: %v", arg, err)
				}

				deleteByPath(path)
			} else {
				if !services.Delete(arg) {
					services.Fatal("error: no bookmark named %q", arg)
				}

				fmt.Printf("Deleted bookmark %q\n", arg)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
