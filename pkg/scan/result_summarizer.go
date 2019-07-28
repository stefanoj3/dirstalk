package scan

import (
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"

	gotree "github.com/DiSiqueira/GoTree"
)

func NewResultSummarizer(out io.Writer) *ResultSummarizer {
	return &ResultSummarizer{out: out, resultMap: make(map[string]struct{})}
}

type ResultSummarizer struct {
	out             io.Writer
	results         []Result
	resultMap       map[string]struct{}
	resultsReceived int
	mux             sync.RWMutex
}

func (s *ResultSummarizer) Add(result Result) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.resultsReceived++

	key := keyForResult(result)
	_, found := s.resultMap[key]
	if found {
		return
	}

	if result.StatusCode == http.StatusNotFound {
		return
	}

	s.resultMap[key] = struct{}{}

	s.results = append(s.results, result)
}

func (s *ResultSummarizer) Summarize() {
	s.mux.Lock()
	defer s.mux.Unlock()

	sort.Slice(s.results, func(i, j int) bool {
		return s.results[i].Target.Path < s.results[j].Target.Path
	})

	s.printSummary()
	s.printTree()

	for _, r := range s.results {
		fmt.Fprintln(
			s.out,
			fmt.Sprintf(
				"%s [%d] [%s]",
				r.URL.String(),
				r.StatusCode,
				r.Target.Method,
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

	// TODO: improve efficiency
	for _, r := range s.results {
		currentBranch := root

		parts := strings.Split(r.URL.Path, "/")
		for _, p := range parts {
			if len(p) == 0 {
				continue
			}

			found := false

			for _, item := range currentBranch.Items() {
				if item.Text() != p {
					continue
				}

				currentBranch = item
				found = true
				break
			}

			if found {
				continue
			}

			newTree := gotree.New(p)
			currentBranch.AddTree(newTree)
			currentBranch = newTree
		}
	}

	fmt.Fprintln(s.out, root.Print())
}

func keyForResult(result Result) string {
	return fmt.Sprintf("%s~%s", result.URL.String(), result.Target.Method)
}
