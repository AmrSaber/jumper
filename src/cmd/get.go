package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AmrSaber/jumper/src/services"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:     "resolve <bookmark>[/path]",
	Short:   "Print the resolved path for a bookmark",
	Aliases: []string{"get"},
	Long: `Resolve resolves a bookmark and prints the full path.

The argument is a bookmark name, optionally followed by a relative path to a
subdirectory or file within the bookmarked directory.

Examples:
  # Print the bookmarked path
  jumper resolve myproject

  # Print the path for a subdirectory of a bookmark
  jumper resolve myproject/src/internal

  # Print the path for a file inside a bookmark
  jumper resolve myproject/README.md

  # Use the resolved path in a command
  ls $(jumper resolve myproject/docs)

  # Open a file from a bookmark in an editor
  $EDITOR $(jumper resolve myproject/src/index.ts)`,
	Args: cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		bookmarkTitle, subPrefix, hasSlash := strings.Cut(toComplete, "/")
		if !hasSlash {
			var completions []cobra.Completion
			for _, c := range bookmarkCompletions(toComplete) {
				completions = append(completions, c+"/")
			}
			return completions, cobra.ShellCompDirectiveNoFileComp | cobra.ShellCompDirectiveNoSpace
		}

		bookmark, found := services.Get(bookmarkTitle)
		if !found {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		dir, namePrefix := filepath.Split(subPrefix)
		entries, err := os.ReadDir(filepath.Join(bookmark.Path, dir))
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		lower := strings.ToLower(namePrefix)
		var completions []cobra.Completion
		for _, entry := range entries {
			if !strings.Contains(strings.ToLower(entry.Name()), lower) {
				continue
			}

			entryPath := filepath.Join(bookmark.Path, dir, entry.Name())
			info, err := os.Stat(entryPath)
			if err != nil {
				continue
			}

			suffix := bookmarkTitle + "/" + filepath.Join(dir, entry.Name())
			if info.IsDir() {
				suffix += "/"
			}

			completions = append(completions, suffix)
		}

		return completions, cobra.ShellCompDirectiveNoFileComp | cobra.ShellCompDirectiveNoSpace
	},

	Run: func(cmd *cobra.Command, args []string) {
		title, subDir, _ := strings.Cut(args[0], "/")

		bookmark, ok := services.Get(title)
		if !ok {
			services.Fatal("error: no bookmark named %q", title)
		}

		fmt.Println(filepath.Join(bookmark.Path, subDir))
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
