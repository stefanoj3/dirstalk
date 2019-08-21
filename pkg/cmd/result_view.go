package cmd

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/stefanoj3/dirstalk/pkg/result"
	"github.com/stefanoj3/dirstalk/pkg/scan/summarizer/tree"
)

func NewResultViewCommand(out io.Writer) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "result.view",
		Short: "Read a scan output file and render the folder tree",
		RunE:  buildResultViewCmd(out),
	}

	cmd.Flags().StringP(
		flagResultFile,
		flagResultFileShort,
		"",
		"result file to read",
	)
	err := cmd.MarkFlagFilename(flagResultFile)
	if err != nil {
		return nil, err
	}

	err = cmd.MarkFlagRequired(flagResultFile)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

func buildResultViewCmd(out io.Writer) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		resultFilePath := cmd.Flag(flagResultFile).Value.String()

		results, err := result.LoadResultsFromFile(resultFilePath)
		if err != nil {
			return errors.Wrapf(err, "failed to load results from %s", resultFilePath)
		}

		treeAsString := tree.NewResultTreeProducer().String(results)

		_, err = fmt.Fprintln(out, treeAsString)

		return errors.Wrap(err, "failed to print result tree")
	}
}
