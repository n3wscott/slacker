package options

import (
	"github.com/spf13/cobra"
)

// MessageOptions
type MessageOptions struct {
	Message        string
	Reaction       bool
	RemoveReaction bool
	// TODO: reactions.
}

func AddMessageArgs(cmd *cobra.Command, o *MessageOptions) {
	cmd.Flags().BoolVar(&o.Reaction, "reaction", false,
		"Message is treated as a reaction.")
	cmd.Flags().BoolVar(&o.RemoveReaction, "remove-reaction", false,
		"Message is treated as a reaction to remove.")
}
