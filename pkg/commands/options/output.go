package options

import "github.com/spf13/cobra"

// OutputOptions
type OutputOptions struct {
	JSON bool
}

func AddOutputArg(cmd *cobra.Command, po *OutputOptions) {
	cmd.Flags().BoolVar(&po.JSON, "json", false,
		"Output as JSON.")
}
