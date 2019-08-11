package producer

import (
	"net/url"

	"github.com/stefanoj3/dirstalk/pkg/scan"
)

type Scanner interface {
	Scan(baseURL *url.URL, workers int) <-chan scan.Result
}
