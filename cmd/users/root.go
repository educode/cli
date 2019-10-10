package users

import (
	"github.com/spf13/cobra"
	"os"
)

var RootCommand = &cobra.Command{
	Use:   "users",
	Short: "Provides functions for managing users",
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

func init() {
	RootCommand.AddCommand(usersListCommand)
}
