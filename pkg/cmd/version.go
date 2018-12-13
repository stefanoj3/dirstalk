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

func init() {
	if Version == "" {
		Version = "dev"
	}

	if BuildTime == "" {
		BuildTime = "now"
	}
}

func newVersionCommand(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the current version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(
				out,
				fmt.Sprintf("Version: %s", Version),
			)
			fmt.Fprintln(
				out,
				fmt.Sprintf("Built at: %s", BuildTime),
			)
		},
	}

	return cmd
}
