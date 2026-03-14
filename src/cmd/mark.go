package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"jumper/src/services"

	"github.com/spf13/cobra"
)

var markCmd = &cobra.Command{
	Use:   "mark [name] [directory]",
	Short: "Bookmark a directory",
	Long: `Bookmark a directory under the given name.
If no name is provided, the directory's base name is used.
If no directory is provided, the current directory is used.
If the name already exists, its path is overwritten.`,
	Example: `  # Bookmark current directory using its base name
  jumper mark

  # Bookmark current directory as "proj"
  jumper mark proj

  # Bookmark a specific directory as "proj"
  jumper mark proj ~/Projects/my-app`,
	Args: cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := os.Getwd()
		if err != nil {
			services.Fatal("error: cannot get current directory: %v", err)
		}

		var name string
		switch len(args) {
		case 0:
			name = filepath.Base(dir)
		case 1:
			name = args[0]
		case 2:
			name = args[0]
			dir = args[1]
		}

		services.Upsert(name, dir)

		fmt.Printf("Bookmarked %q -> %s\n", name, dir)
	},
}

func init() {
	rootCmd.AddCommand(markCmd)
}
