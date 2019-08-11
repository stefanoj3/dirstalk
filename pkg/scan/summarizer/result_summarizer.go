package summarizer

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/stefanoj3/dirstalk/pkg/scan"

	gotree "github.com/DiSiqueira/GoTree"
)

const (
	breakingText = "Found something breaking"
	foundText    = "Found"
)

func NewResultSummarizer(logger *logrus.Logger) *ResultSummarizer {
	return &ResultSummarizer{
		logger:    logger,
		resultMap: make(map[string]struct{}),
	}
}

type ResultSummarizer struct {
	logger    *logrus.Logger
	results   []scan.Result
	resultMap map[string]struct{}
	mux       sync.RWMutex
}

func (s *ResultSummarizer) Add(result scan.Result) {
	s.mux.Lock()
	defer s.mux.Unlock()

	key := keyForResult(result)
	_, found := s.resultMap[key]
	if found {
		return
	}

	s.log(result)

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
		_, _ = fmt.Fprintln(
			s.logger.Out,
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
	_, _ = fmt.Fprintln(
		s.logger.Out,
		fmt.Sprintf("%d results found", len(s.results)),
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

	_, _ = fmt.Fprintln(s.logger.Out, root.Print())
}

func (s *ResultSummarizer) log(result scan.Result) {
	statusCode := result.StatusCode

	l := s.logger.WithFields(logrus.Fields{
		"status-code": statusCode,
		"method":      result.Target.Method,
		"url":         result.URL.String(),
	})

	if statusCode >= http.StatusInternalServerError {
		l.Warn(breakingText)
	} else {
		l.Info(foundText)
	}
}

func keyForResult(result scan.Result) string {
	return fmt.Sprintf("%s~%s", result.URL.String(), result.Target.Method)
}
