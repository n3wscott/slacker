package options

import (
	"github.com/spf13/cobra"
)

// FileSystemOptions
type FileSystemOptions struct {
	Workspace string
}

func AddFileSystemArgs(cmd *cobra.Command, o *FileSystemOptions) {
	cmd.Flags().StringVar(&o.Workspace, "workspace", "",
		"The directory to turn into a PR.")

	_ = cmd.MarkFlagRequired("workspace")
}
