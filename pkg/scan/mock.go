package scan

import "net/http"

type DoerMock struct {
	DoFunc func(*http.Request) (*http.Response, error)
}

func (m *DoerMock) Do(r *http.Request) (*http.Response, error) {
	return m.DoFunc(r)
}
