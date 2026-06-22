package cmd

import (
	"github.com/spf13/cobra"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Commands for agent integration",
}

func init() {
	rootCmd.AddCommand(agentCmd)
}
