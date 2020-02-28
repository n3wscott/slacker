package commands

import (
	"github.com/spf13/cobra"
)

func addHelp(topLevel *cobra.Command) {
	cmd := &cobra.Command{
		Use:       "help",
		Hidden:    true,
		ValidArgs: []string{},
		Short:     "Print the help menu.",

		RunE: func(cmd *cobra.Command, args []string) error {
			return topLevel.Help()
		},
	}
	topLevel.AddCommand(cmd)
}
