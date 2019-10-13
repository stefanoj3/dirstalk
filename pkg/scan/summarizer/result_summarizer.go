package summarizer

import (
	"fmt"
	"net/http"
	"sort"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/stefanoj3/dirstalk/pkg/scan"
)

const (
	breakingText = "Found something breaking"
	foundText    = "Found"
)

func NewResultSummarizer(treePrinter ResultTree, logger *logrus.Logger) *ResultSummarizer {
	return &ResultSummarizer{
		treePrinter: treePrinter,
		logger:      logger,
		resultMap:   make(map[string]struct{}),
	}
}

type ResultSummarizer struct {
	treePrinter ResultTree
	logger      *logrus.Logger
	results     []scan.Result
	resultMap   map[string]struct{}
	mux         sync.RWMutex
}

func (s *ResultSummarizer) Add(result scan.Result) {
	s.mux.Lock()
	defer s.mux.Unlock()

	key := keyForResult(result)

	if _, found := s.resultMap[key]; found {
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
	_, _ = fmt.Fprintln(s.logger.Out, s.treePrinter.String(s.results))
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
