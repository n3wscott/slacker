package commands

import (
	"github.com/n3wscott/slacker/pkg/commands/options"
	"github.com/spf13/cobra"
)

var (
	oo = &options.OutputOptions{}
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "slacker",
		Short: "Interact with slack from a command line.",
	}

	AddCommands(cmd)
	return cmd
}

func AddCommands(topLevel *cobra.Command) {
	addList(topLevel)
	addSend(topLevel)

	addHelp(topLevel)
}
