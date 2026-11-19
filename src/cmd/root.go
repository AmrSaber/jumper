// Package cmd contains all the commands used.
package cmd

import (
	"fmt"
	"os"

	"github.com/AmrSaber/jumper/src/services"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "jumper",
	Short: "A bookmark manager for your shell directories",
	Long: fmt.Sprintf(
		`Jumper lets you bookmark directories and jump to them by name from anywhere in your shell.
Bookmarks are saved in %[1]s.`,
		services.BookmarksPath(),
	),
	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	rootCmd.Version = getVersion()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
