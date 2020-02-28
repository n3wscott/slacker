package directs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/gosuri/uitable"
	"github.com/n3wscott/slacker/pkg/slackbot"
)

type Directs struct {
	Output string
	List   bool
}

func (o *Directs) Do(ctx context.Context) error {
	if o.List {
		return o.list(ctx)
	}
	return nil
}

func (o *Directs) list(ctx context.Context) error {

	sb := slackbot.NewInstance(ctx)

	ims, err := sb.GetIMs(ctx)
	if err != nil {
		return err
	}

	switch o.Output {
	case "json":
		b, err := json.Marshal(ims)
		if err != nil {
			return err
		}
		_, _ = fmt.Fprintln(color.Output, string(b))
	default:
		tbl := uitable.New()
		tbl.Separator = "  "
		tbl.AddRow("ID", "User Name", "User ID")
		for _, v := range ims {
			tbl.AddRow(v.ID, v.Name, v.With)
		}
		_, _ = fmt.Fprintln(color.Output, tbl)
	}

	return nil
}
