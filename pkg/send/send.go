package send

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/gosuri/uitable"
	"strings"

	"github.com/n3wscott/slacker/pkg/slackbot"
)

type Send struct {
	Output         string
	ChannelName    string
	ChannelID      string
	ThreadID       string
	Message        string
	Reaction       bool
	RemoveReaction bool
}

func (c *Send) Do(ctx context.Context) error {
	sb := slackbot.NewInstance(ctx)

	if c.ChannelID == "" && c.ChannelName != "" {
		channels, err := sb.GetChannels(ctx)
		if err != nil {
			return err
		}
		for _, ch := range channels {
			if strings.EqualFold(strings.ToLower(ch.Name), strings.ToLower(c.ChannelName)) {
				c.ChannelID = ch.ID
				break
			}
		}
		ims, err := sb.GetIMs(ctx)
		if err != nil {
			return err
		}
		for _, ch := range ims {
			if strings.EqualFold(strings.ToLower(ch.Name), strings.ToLower(c.ChannelName)) {
				c.ChannelID = ch.ID
				break
			}
		}
	}

	if c.ChannelID == "" {
		return errors.New("unknown channel")
	}

	var err error
	var resp *slackbot.SlackPostResponse
	if c.Reaction || c.RemoveReaction {
		resp, err = sb.RemoveReaction(ctx, c.ChannelID, c.ThreadID, c.Message)
		if err != nil {
			return err
		}
	} else if c.Reaction {
		resp, err = sb.PostReaction(ctx, c.ChannelID, c.ThreadID, c.Message)
		if err != nil {
			return err
		}
	} else if c.ThreadID != "" {
		resp, err = sb.PostResponse(ctx, c.ChannelID, c.ThreadID, c.Message)
		if err != nil {
			return err
		}
	} else {
		resp, err = sb.PostMessage(ctx, c.ChannelID, c.Message)
		if err != nil {
			return err
		}
	}

	switch c.Output {
	case "json":
		b, err := json.Marshal(resp)
		if err != nil {
			return err
		}
		_, _ = fmt.Fprintln(color.Output, string(b))
	default:
		tbl := uitable.New()
		tbl.Separator = "  "
		tbl.AddRow("ChannelID", "Thread", "Parent")
		tbl.AddRow(resp.ChannelID, resp.Thread, resp.Parent)
		_, _ = fmt.Fprintln(color.Output, tbl)
	}
	return nil
}
