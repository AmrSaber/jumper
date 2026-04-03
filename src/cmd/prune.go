package cmd

import (
	"fmt"

	"github.com/AmrSaber/jumper/src/services"

	"github.com/spf13/cobra"
)

var pruneCmd = &cobra.Command{
	Use:     "prune",
	Short:   "Remove bookmarks whose paths no longer exist",
	Aliases: []string{"clean"},
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		deleted := services.Prune()

		if len(deleted) == 0 {
			fmt.Println("No stale bookmarks found.")
			return
		}

		for _, bookmark := range deleted {
			fmt.Printf("Deleted bookmark %q\n", bookmark.Title)
		}
	},
}

func init() {
	rootCmd.AddCommand(pruneCmd)
}
