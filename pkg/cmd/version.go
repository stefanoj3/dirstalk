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

func NewVersionCommand(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the current version",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := fmt.Fprintln(
				out,
				fmt.Sprintf("Version: %s", Version),
			)
			if err != nil {
				return err
			}

			_, err = fmt.Fprintln(
				out,
				fmt.Sprintf("Built at: %s", BuildTime),
			)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
