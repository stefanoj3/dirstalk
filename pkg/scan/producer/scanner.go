package producer

import (
	"net/url"

	"github.com/stefanoj3/dirstalk/pkg/scan"
)

type Scanner interface {
	Scan(baseUrl *url.URL, workers int) <-chan scan.Result
}
