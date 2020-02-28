package commands

import (
	"context"
	"errors"
	"github.com/n3wscott/slacker/pkg/commands/options"
	"github.com/n3wscott/slacker/pkg/send"
	"github.com/spf13/cobra"
	"strings"
)

func addSend(topLevel *cobra.Command) {
	co := &options.ChannelOptions{}
	mo := &options.MessageOptions{}

	cmd := &cobra.Command{
		Use:       "send",
		ValidArgs: []string{},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires a message") // TODO: when we support reactions, this will be true only if no reaction.
			}
			mo.Message = strings.Join(args, " ")

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			s := send.Send{
				ChannelID:      co.ID,
				ChannelName:    co.Name,
				ThreadID:       co.Thread,
				Message:        mo.Message,
				Reaction:       mo.Reaction,
				RemoveReaction: mo.RemoveReaction,
			}
			if oo.JSON {
				s.Output = "json"
			}
			return s.Do(context.Background())
		},
	}

	options.AddChannelArgs(cmd, co)
	options.AddMessageArgs(cmd, mo)
	options.AddOutputArg(cmd, oo)
	topLevel.AddCommand(cmd)
}
