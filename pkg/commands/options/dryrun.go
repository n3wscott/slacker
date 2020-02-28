package options

import "github.com/spf13/cobra"

// DryRunOptions
type DryRunOptions struct {
	DryRun bool
}

func AddDryRunArg(cmd *cobra.Command, po *DryRunOptions) {
	cmd.Flags().BoolVar(&po.DryRun, "dry-run", false,
		"Output what would happen.")
}
