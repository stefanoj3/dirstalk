package test

import (
	"net/http"
	"net/http/httptest"
	"sync"
)

func NewServerWithAssertion(handler http.HandlerFunc) (*httptest.Server, *ServerAssertion) {
	serverAssertion := &ServerAssertion{}

	server := httptest.NewServer(serverAssertion.wrap(handler))

	return server, serverAssertion
}

func NewTSLServerWithAssertion(handler http.HandlerFunc) (*httptest.Server, *ServerAssertion) {
	serverAssertion := &ServerAssertion{}

	server := httptest.NewUnstartedServer(serverAssertion.wrap(handler))
	server.StartTLS()

	return server, serverAssertion
}

type ServerAssertion struct {
	requests   []http.Request
	requestsMx sync.RWMutex
}

func (s *ServerAssertion) wrap(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)

		s.requestsMx.Lock()
		defer s.requestsMx.Unlock()

		s.requests = append(s.requests, *r)
	})
}

func (s *ServerAssertion) Range(fn func(index int, r http.Request)) {
	s.requestsMx.RLock()
	defer s.requestsMx.RUnlock()

	for i, r := range s.requests {
		fn(i, r)
	}
}

func (s *ServerAssertion) At(index int, fn func(r http.Request)) {
	s.requestsMx.RLock()
	defer s.requestsMx.RUnlock()

	fn(s.requests[index])
}

func (s *ServerAssertion) Len() int {
	s.requestsMx.RLock()
	defer s.requestsMx.RUnlock()

	return len(s.requests)
}
