package options

import (
	"github.com/spf13/cobra"
)

// ChannelOptions
type ChannelOptions struct {
	Name   string
	ID     string
	Thread string
}

func AddChannelArgs(cmd *cobra.Command, o *ChannelOptions) {
	cmd.Flags().StringVar(&o.Name, "name", "",
		"Channel Name to use, will resolve to a Channel ID.")
	cmd.Flags().StringVar(&o.ID, "id", "",
		"Channel ID to use.")
	cmd.Flags().StringVar(&o.Thread, "thread", "",
		"Unique identifier of a thread's parent message.")
}
