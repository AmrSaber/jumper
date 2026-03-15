package cmd

import (
	"fmt"
	"strings"

	"jumper/src/services"

	"github.com/spf13/cobra"
)

var renameCmd = &cobra.Command{
	Use:     "rename <old-name> <new-name>",
	Short:   "Rename a bookmark",
	Aliases: []string{"mv"},
	Args:    cobra.ExactArgs(2),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return bookmarkCompletions(toComplete), cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		oldName, newName := args[0], args[1]

		if strings.HasPrefix(newName, ".") || strings.HasPrefix(newName, "~") || strings.HasPrefix(newName, "/") {
			services.Fatal("error: bookmark name cannot start with '.', '~', or '/'")
		}

		if !services.Rename(oldName, newName) {
			services.Fatal("error: no bookmark named %q", oldName)
		}

		fmt.Printf("Renamed bookmark %q to %q\n", oldName, newName)
	},
}

func init() {
	rootCmd.AddCommand(renameCmd)
}
