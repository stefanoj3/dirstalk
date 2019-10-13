package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

var (
	Version   string
	BuildTime string
)

//nolint:gochecknoinits
func init() {
	if Version == "" {
		Version = "dev"
	}

	if BuildTime == "" {
		BuildTime = "now"
	}
}

func NewVersionCommand(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the current version",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _ = fmt.Fprintln(
				out,
				fmt.Sprintf("Version: %s", Version),
			)

			_, _ = fmt.Fprintln(
				out,
				fmt.Sprintf("Built at: %s", BuildTime),
			)

			return nil
		},
	}

	return cmd
}
