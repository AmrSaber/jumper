package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AmrSaber/jumper/assets"
	"github.com/AmrSaber/jumper/src/common"
	"github.com/AmrSaber/jumper/src/services"

	"github.com/spf13/cobra"
)

var agentSetupCmd = &cobra.Command{
	Use:       "setup <agent-name>",
	Short:     "Install jumper's integration for a supported agent",
	Long:      `Write jumper's skill and plugin into the agent's config so the agent always has the live bookmark list and knows how to manage bookmarks.`,
	Args:      cobra.ExactArgs(1),
	ValidArgs: common.SupportedAgents,
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
		return common.SupportedAgents, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		agentName := args[0]

		switch agentName {
		case "opencode":
			skillPath, pluginPath := common.OpencodePaths()
			_ = os.RemoveAll(filepath.Dir(skillPath))
			_ = os.Remove(pluginPath)
			writeAsset(pluginPath, assets.OpencodePlugin)
		default:
			services.Fatal("error: unsupported agent %q", agentName)
		}
	},
}

func writeAsset(path, content string) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		services.Fatal("error: cannot create directory for %s: %v", path, err)
	}

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		services.Fatal("error: cannot write %s: %v", path, err)
	}

	fmt.Printf("wrote %s\n", path)
}

func init() {
	agentCmd.AddCommand(agentSetupCmd)
}
