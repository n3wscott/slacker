package commands

import (
	"context"
	"github.com/n3wscott/slacker/pkg/channel"
	"github.com/n3wscott/slacker/pkg/commands/options"
	"github.com/n3wscott/slacker/pkg/directs"
	"github.com/spf13/cobra"
)

func addList(topLevel *cobra.Command) {
	cmd := &cobra.Command{
		Use:       "get",
		Short:     "Get a list of channels or direct messages.",
		ValidArgs: []string{},
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	addChannelGet(cmd)
	addDirectMessageGet(cmd)

	topLevel.AddCommand(cmd)
}

func addChannelGet(topLevel *cobra.Command) {
	cmd := &cobra.Command{
		Use:       "channel",
		ValidArgs: []string{},
		Short:     "Get a list of channels.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := channel.Channel{
				List: true,
			}
			if oo.JSON {
				c.Output = "json"
			}
			err := c.Do(context.Background())
			return oo.HandleError(err)
		},
	}
	options.AddOutputArg(cmd, oo)

	topLevel.AddCommand(cmd)
}

func addDirectMessageGet(topLevel *cobra.Command) {
	cmd := &cobra.Command{
		Use:       "dm",
		ValidArgs: []string{},
		Short:     "Get a list of direct message channels.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := directs.Directs{
				List: true,
			}
			if oo.JSON {
				c.Output = "json"
			}
			err := c.Do(context.Background())
			return oo.HandleError(err)
		},
	}
	options.AddOutputArg(cmd, oo)

	topLevel.AddCommand(cmd)
}
