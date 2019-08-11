package output

import (
	"encoding/json"

	"github.com/stefanoj3/dirstalk/pkg/scan"
)

func convertResultToRawData(r scan.Result) ([]byte, error) {
	return json.Marshal(r)
}
