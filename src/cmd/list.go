package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"jumper/src/services"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var listFlags = struct{ output string }{}

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all bookmarks",
	Aliases: []string{"ls"},
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		bookmarks := services.List()

		if len(bookmarks) == 0 {
			fmt.Println("No bookmarks yet. Use `jumper mark` to add one.")
			return
		}

		sort.Slice(bookmarks, func(i, j int) bool {
			return strings.ToLower(bookmarks[i].Title) < strings.ToLower(bookmarks[j].Title)
		})

		switch listFlags.output {
		case "json":
			output, _ := json.MarshalIndent(bookmarks, "", "  ")
			fmt.Println(string(output))
		case "yaml":
			output, _ := yaml.Marshal(bookmarks)
			fmt.Print(string(output))
		case "table":
			t := table.NewWriter()
			t.SetOutputMirror(os.Stdout)
			t.AppendHeader(table.Row{"Title", "Path"})

			homePath, err := os.UserHomeDir()
			if err != nil {
				services.Fatal("Could not get home directory path: %v", err)
			}

			for _, bookmark := range bookmarks {
				bookmarkPath := bookmark.Path
				if after, ok := strings.CutPrefix(bookmarkPath, homePath); ok {
					bookmarkPath = "~" + after
				}

				t.AppendRow(table.Row{color.New(color.FgBlue).Sprint(bookmark.Title), bookmarkPath})
			}
			t.SetStyle(table.StyleLight)
			t.Render()
		default:
			services.Fatal("error: unsupported format %q", listFlags.output)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&listFlags.output, "output", "o", "table", "Output format: table, json, yaml")
	_ = listCmd.RegisterFlagCompletionFunc(
		"output",
		cobra.FixedCompletions([]string{"table", "json", "yaml"}, cobra.ShellCompDirectiveDefault),
	)
}
