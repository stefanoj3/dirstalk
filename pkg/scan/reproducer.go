package scan

import "github.com/chuckpreslar/emission"

var statusCodesToSkip = map[int]bool{
	404: false,
}

type ReProcessor struct {
	eventEmitter *emission.Emitter
	dictionary   []string
}

func (r *ReProcessor) Process(result *Result) {
	if _, ok := statusCodesToSkip[result.Response.StatusCode]; ok {
		return
	}

	// process directories
	// process hidden files (eg .env)
	// process hidden folders (.git)
	// process files with extensions (eg index.php)
}
