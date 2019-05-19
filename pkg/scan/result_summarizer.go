package scan

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	gotree "github.com/DiSiqueira/GoTree"
)

func NewResultSummarizer(out io.Writer) *ResultSummarizer {
	return &ResultSummarizer{out: out}
}

type ResultSummarizer struct {
	out             io.Writer
	results         []*Result
	resultsReceived int
	mux             sync.RWMutex
}

func (s *ResultSummarizer) Add(result *Result) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.resultsReceived++

	if result.Response.StatusCode == http.StatusNotFound {
		return
	}

	s.results = append(s.results, result)
}

func (s *ResultSummarizer) Summarize() {
	s.mux.RLock()
	defer s.mux.RUnlock()

	s.printSummary()
	s.printTree()

	for _, r := range s.results {
		fmt.Fprintln(
			s.out,
			fmt.Sprintf(
				"%s [%d] [%s]",
				r.Response.Request.URL,
				r.Response.StatusCode,
				r.Response.Request.Method,
			),
		)
	}
}

func (s *ResultSummarizer) printSummary() {
	fmt.Fprintln(
		s.out,
		fmt.Sprintf("%d requests made, %d results found", s.resultsReceived, len(s.results)),
	)
}

func (s *ResultSummarizer) printTree() {
	root := gotree.New("/")

	treeMap := map[string]gotree.Tree{}

	for _, r := range s.results {
		currentBranch := root

		parts := strings.Split(r.Response.Request.URL.Path, "/")
		for _, p := range parts {
			if len(p) == 0 {
				continue
			}

			t, ok := treeMap[p]
			if !ok {
				currentBranch = currentBranch.Add(p)
				treeMap[p] = currentBranch
				continue
			}

			currentBranch = t
		}
	}

	fmt.Fprintln(s.out, root.Print())
}
