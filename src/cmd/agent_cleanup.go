package cmd

import (
	"os"
	"path/filepath"

	"github.com/AmrSaber/jumper/src/common"
	"github.com/AmrSaber/jumper/src/services"

	"github.com/spf13/cobra"
)

var agentCleanupCmd = &cobra.Command{
	Use:       "cleanup <agent-name>",
	Short:     "Remove jumper's integration from a supported agent",
	Long:      `Remove the jumper skill and plugin previously written by 'jumper agent setup'. Only the files jumper installed are removed; the rest of the agent's config is left untouched.`,
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
		default:
			services.Fatal("error: unsupported agent %q", agentName)
		}
	},
}

func init() {
	agentCmd.AddCommand(agentCleanupCmd)
}
