package cmd

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"
	"github.com/stefanoj3/dirstalk/pkg/common"
	"github.com/stefanoj3/dirstalk/pkg/result"
	"github.com/stefanoj3/dirstalk/pkg/scan/summarizer/tree"
)

func NewResultDiffCommand(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "result.diff",
		Short: "Prints differences between 2 result files",
		RunE:  buildResultDiffCmd(out),
	}

	cmd.Flags().StringP(
		flagResultDiffFirstFile,
		flagResultDiffFirstFileShort,
		"",
		"first result file to read",
	)
	common.Must(cmd.MarkFlagFilename(flagResultDiffFirstFile))
	common.Must(cmd.MarkFlagRequired(flagResultDiffFirstFile))

	cmd.Flags().StringP(
		flagResultDiffSecondFile,
		flagResultDiffSecondFileShort,
		"",
		"second result file to read",
	)
	common.Must(cmd.MarkFlagFilename(flagResultDiffSecondFile))
	common.Must(cmd.MarkFlagRequired(flagResultDiffSecondFile))

	return cmd
}

func buildResultDiffCmd(out io.Writer) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		firstResultFilePath := cmd.Flag(flagResultDiffFirstFile).Value.String()

		resultsFirst, err := result.LoadResultsFromFile(firstResultFilePath)
		if err != nil {
			return errors.Wrapf(err, "failed to load results from %s", firstResultFilePath)
		}

		secondResultFilePath := cmd.Flag(flagResultDiffSecondFile).Value.String()

		resultsSecond, err := result.LoadResultsFromFile(secondResultFilePath)
		if err != nil {
			return errors.Wrapf(err, "failed to load results from %s", secondResultFilePath)
		}

		treeProducer := tree.NewResultTreeProducer()

		differ := diffmatchpatch.New()
		diffs := differ.DiffMain(
			treeProducer.String(resultsFirst),
			treeProducer.String(resultsSecond),
			false,
		)

		if isEqual(diffs) {
			return errors.New("no diffs found")
		}

		_, err = fmt.Fprintln(out, differ.DiffPrettyText(diffs))

		return errors.Wrap(err, "failed to print results diff")
	}
}

func isEqual(diffs []diffmatchpatch.Diff) bool {
	if len(diffs) != 1 {
		return false
	}

	return diffs[0].Type == diffmatchpatch.DiffEqual
}
