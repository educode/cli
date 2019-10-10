package cmd

import (
	"github.com/hhu-educode/cli/cmd/submissions"
	"github.com/hhu-educode/cli/cmd/users"
	"github.com/hhu-educode/cli/cmd/whitelist"
	"github.com/spf13/cobra"
	"os"
)

var rootCommand = &cobra.Command{
	Use:   "educode",
	Short: "A simple command line interface for managing educode",
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

func init() {
	rootCommand.AddCommand(submissions.RootCommand)
	rootCommand.AddCommand(users.RootCommand)
	rootCommand.AddCommand(whitelist.RootCommand)
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
