package cmd

import "github.com/stefanoj3/dirstalk/pkg/scan"

type OutputSaver interface {
	Save(scan.Result) error
	Close() error
}
