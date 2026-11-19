package cmd

import (
	"fmt"

	"github.com/AmrSaber/jumper/src/services"

	"github.com/spf13/cobra"
)

var commentFlags = struct{ noComment bool }{}

var commentCmd = &cobra.Command{
	Use:   "comment <name> [comment]",
	Short: "Set or remove a bookmark's comment",
	Long: `Set the comment of an existing bookmark. The name is matched case-insensitively.
The comment can only be omitted together with --no-comment, which removes the comment.`,
	Example: `  # Set a comment
  jumper comment proj "Main project"

  # Remove the comment
  jumper comment proj --no-comment`,
	Args: cobra.RangeArgs(1, 2),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return bookmarkCompletions(toComplete), cobra.ShellCompDirectiveNoFileComp
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		var comment string
		switch {
		case commentFlags.noComment && len(args) == 2:
			services.Fatal("error: cannot use --no-comment together with a comment")
		case !commentFlags.noComment && len(args) == 1:
			services.Fatal("error: a comment is required unless --no-comment is provided")
		case len(args) == 2:
			comment = args[1]
		}

		if err := services.UpdateComment(name, comment); err != nil {
			services.Fatal("error: %v", err)
		}

		if commentFlags.noComment {
			fmt.Printf("Removed comment from %q\n", name)
		} else {
			fmt.Printf("Updated comment of %q\n", name)
		}
	},
}

func init() {
	rootCmd.AddCommand(commentCmd)

	commentCmd.Flags().BoolVar(&commentFlags.noComment, "no-comment", false, "Remove the bookmark's comment")
}
