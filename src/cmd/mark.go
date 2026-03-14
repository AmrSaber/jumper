package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"jumper/src/services"

	"github.com/spf13/cobra"
)

var markCmd = &cobra.Command{
	Use:   "mark [name]",
	Short: "Bookmark the current directory",
	Long: `Bookmark the current directory under the given name.
If no name is provided, the directory's base name is used.
If the name already exists, its path is overwritten.`,
	Example: `  # Bookmark current directory as "proj"
  jumper mark proj

  # Bookmark using the directory's base name
  jumper mark`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: cannot get current directory:", err)
			os.Exit(1)
		}

		name := filepath.Base(dir)
		if len(args) == 1 {
			name = args[0]
		}

		services.Upsert(name, dir)

		fmt.Printf("Bookmarked %q -> %s\n", name, dir)
	},
}

func init() {
	rootCmd.AddCommand(markCmd)
}
