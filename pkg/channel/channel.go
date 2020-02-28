package channel

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/gosuri/uitable"
	"github.com/n3wscott/slacker/pkg/slackbot"
)

type Channel struct {
	Output string
	List   bool
}

func (c *Channel) Do(ctx context.Context) error {
	if c.List {
		return c.list(ctx)
	}
	return nil
}

func (c *Channel) list(ctx context.Context) error {
	sb := slackbot.NewInstance(ctx)

	channels, err := sb.GetChannels(ctx)
	if err != nil {
		return err
	}

	switch c.Output {
	case "json":
		b, err := json.Marshal(channels)
		if err != nil {
			return err
		}
		_, _ = fmt.Fprintln(color.Output, string(b))
	default:
		tbl := uitable.New()
		tbl.Separator = "  "
		tbl.AddRow("Name", "ID", "IsMember")
		for _, v := range channels {
			tbl.AddRow(v.Name, v.ID, v.IsMember)
		}
		_, _ = fmt.Fprintln(color.Output, tbl)
	}

	return nil
}
