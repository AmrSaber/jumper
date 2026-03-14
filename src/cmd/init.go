package cmd

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"jumper/src/services"

	"github.com/spf13/cobra"
)

//go:embed shells/jumper.zsh
var zshInit string

//go:embed shells/jumper.bash
var bashInit string

var initCmd = &cobra.Command{
	Use:   "init [bash|zsh]",
	Short: "Print the shell init script",
	Long:  `Print the shell integration script. Source it in your shell's rc file.`,
	Example: `  # Auto-detect shell (recommended):
  eval "$(jumper init)"

  # Explicit shell:
  eval "$(jumper init zsh)"
  eval "$(jumper init bash)"`,
	ValidArgs: []string{"bash", "zsh"},
	Args:      cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		shell := ""
		if len(args) == 1 {
			shell = args[0]
		} else {
			shell = filepath.Base(os.Getenv("SHELL"))
		}

		switch shell {
		case "zsh":
			fmt.Print(zshInit)
		case "bash":
			fmt.Print(bashInit)
		default:
			if shell == "" {
				services.Fatal("error: could not detect shell, please pass it explicitly: jumper init [bash|zsh]")
			} else {
				services.Fatal("error: unsupported shell %q, supported: bash, zsh", shell)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
