package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"jumper/src/services"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a bookmark",
	Long: `Delete a bookmark by name.
If no name is provided, the current directory's base name is used.`,
	Aliases: []string{"del", "rm", "unmark"},
	Args:    cobra.MaximumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return bookmarkCompletions(toComplete), cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		name := ""
		if len(args) == 1 {
			name = args[0]
		} else {
			dir, err := os.Getwd()
			if err != nil {
				fmt.Fprintln(os.Stderr, "error: cannot get current directory:", err)
				os.Exit(1)
			}
			name = filepath.Base(dir)
		}

		if !services.Delete(name) {
			fmt.Fprintf(os.Stderr, "error: no bookmark named %q\n", name)
			os.Exit(1)
		}

		fmt.Printf("Deleted bookmark %q\n", name)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
