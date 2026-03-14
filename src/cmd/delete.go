package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"jumper/src/services"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [name...]",
	Short: "Delete one or more bookmarks",
	Long: `Delete bookmarks by name.
If no name is provided, the current directory's base name is used.`,
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
				fmt.Fprintln(os.Stderr, "error: cannot get current directory:", err)
				os.Exit(1)
			}
			args = []string{filepath.Base(dir)}
		}

		for _, name := range args {
			if !services.Delete(name) {
				fmt.Fprintf(os.Stderr, "error: no bookmark named %q\n", name)
				os.Exit(1)
			}
			fmt.Printf("Deleted bookmark %q\n", name)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
