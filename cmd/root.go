package cmd

import (
	"github.com/spf13/cobra"
)

// command is the root commaand
var command = &cobra.Command{
	Use:     "simple-backend",
	Short:   "it's a simple/toy backend for testing/studying that is created by Yuki Tang",
	Version: "0.0",
	Run: func(c *cobra.Command, args []string) {
		c.Usage()
	},
}

func Execute() error {

	command.AddCommand(
		httpservice(),
		databaseOperations(),
	) // add cmd here

	return command.Execute()
}
