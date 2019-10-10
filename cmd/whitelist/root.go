package whitelist

import (
"github.com/spf13/cobra"
"os"
)

var RootCommand = &cobra.Command{
	Use:   "whitelist",
	Short: "Provides functions for managing educode's whitelist",
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

func init() {
	RootCommand.AddCommand(whitelistSyncCommand)
}

