package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/stefanoj3/dirstalk/pkg/dictionary"
)

func NewGenerateDictionaryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dictionary.generate [path]",
		Short: "Generate a dictionary from the given folder",
		RunE:  buildGenerateDictionaryFunc(),
	}

	cmd.Flags().StringP(
		flagOutput,
		flagOutputShort,
		"",
		fmt.Sprintf("where to write the dictionary"),
	)

	cmd.Flags().BoolP(
		flagAbsolutePathOnly,
		"",
		false,
		"determines if the dictionary should contain only the absolute path of the files",
	)

	return cmd
}

func buildGenerateDictionaryFunc() func(cmd *cobra.Command, args []string) error {
	f := func(cmd *cobra.Command, args []string) error {
		p, err := getPath(args)
		if err != nil {
			return err
		}

		out, err := getOutputForDictionaryGenerator(cmd)
		if err != nil {
			return err
		}

		absolutePathOnly, err := cmd.Flags().GetBool(flagAbsolutePathOnly)
		if err != nil {
			return errors.Wrapf(err, "failed to retrieve %s flag", flagAbsolutePathOnly)
		}

		generator := dictionary.NewGenerator(out)

		return generator.GenerateDictionaryFrom(p, absolutePathOnly)
	}

	return f
}

func getOutputForDictionaryGenerator(cmd *cobra.Command) (io.Writer, error) {
	output := cmd.Flag(flagOutput).Value.String()
	if output == "" {
		return os.Stdout, nil
	}

	file, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, errors.Wrap(err, "cannot write on the path provided")
	}

	return file, nil
}

func getPath(args []string) (string, error) {
	if len(args) == 0 {
		return "", errors.New("no path provided")
	}

	path := args[0]

	fileInfo, err := os.Stat(path)
	if err != nil {
		return "", errors.Wrap(err, "unable to use the provided path")
	}

	if !fileInfo.IsDir() {
		return "", errors.New("the path should be a directory")
	}

	return path, nil
}
