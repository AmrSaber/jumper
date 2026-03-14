package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"jumper/src/services"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:    "get <name>",
	Short:  "Print the path for a bookmark",
	Hidden: true,
	Args:   cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return bookmarkCompletions(toComplete), cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		title, subDir, _ := strings.Cut(args[0], "/")

		bookmark, ok := services.Get(title)
		if !ok {
			fmt.Fprintf(os.Stderr, "error: no bookmark named %q\n", title)
			os.Exit(1)
		}

		fmt.Println(filepath.Join(bookmark.Path, subDir))
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
