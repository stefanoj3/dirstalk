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
	ignoredText  = "Ignored"
)

func NewResultSummarizer(httpStatusesToIgnore []int, logger *logrus.Logger) *ResultSummarizer {
	httpStatusesToIgnoreMap := make(map[int]struct{}, len(httpStatusesToIgnore))

	for _, statusToIgnore := range httpStatusesToIgnore {
		httpStatusesToIgnoreMap[statusToIgnore] = struct{}{}
	}

	return &ResultSummarizer{
		httpStatusesToIgnoreMap: httpStatusesToIgnoreMap,
		logger:                  logger,
		resultMap:               make(map[string]struct{}),
	}
}

type ResultSummarizer struct {
	httpStatusesToIgnoreMap map[int]struct{}
	logger                  *logrus.Logger
	results                 []scan.Result
	resultMap               map[string]struct{}
	resultsReceived         int
	mux                     sync.RWMutex
}

func (s *ResultSummarizer) Add(result scan.Result) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.resultsReceived++

	key := keyForResult(result)
	_, found := s.resultMap[key]
	if found {
		return
	}

	s.log(result)
	if _, found := s.httpStatusesToIgnoreMap[result.StatusCode]; found {
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
	fmt.Fprintln(
		s.logger.Out,
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

	fmt.Fprintln(s.logger.Out, root.Print())
}

func (s *ResultSummarizer) log(result scan.Result) {
	statusCode := result.StatusCode

	l := s.logger.WithFields(logrus.Fields{
		"status-code": statusCode,
		"method":      result.Target.Method,
		"url":         result.URL.String(),
	})

	if _, found := s.httpStatusesToIgnoreMap[result.StatusCode]; found {
		l.Debug(ignoredText)
	} else if statusCode >= http.StatusInternalServerError {
		l.Warn(breakingText)
	} else {
		l.Info(foundText)
	}
}

func keyForResult(result scan.Result) string {
	return fmt.Sprintf("%s~%s", result.URL.String(), result.Target.Method)
}
