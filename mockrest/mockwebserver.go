// mockrest is scriptable server which wraps httptest.
// Allows for client testing by responding back with appropriate schemas
package mockrest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"
)

// Server holds state relaced to the handlers and tests.
type Server struct {
	testServer *httptest.Server
	handlers   chan http.HandlerFunc
	requests   chan *http.Request
	URL        string
}

// Create a new Server but don't start it
func New() *Server {
	return &Server{
		handlers: make(chan http.HandlerFunc),
		requests: make(chan *http.Request),
	}
}

// Create a new Server and start it with the specified body as the response
func StartNewWithBody(body string) *Server {
	s := New()
	s.URL = s.Start()
	s.Enqueue(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, body)
	})
	return s
}

// Create a new Server and start it with the specified status as the response
func StartNewWithStatusCode(status int) *Server {
	s := New()
	s.URL = s.Start()
	s.Enqueue(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
	})
	return s
}

// Start a new server with the specified file as the response
func StartNewWithFile(file string) *Server {
	s := New()
	s.URL = s.Start()

	var output string
	b, err := ioutil.ReadFile(file)
	if err != nil {
		output = err.Error()
	} else {
		output = string(b)
	}

	s.Enqueue(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, output)
	})

	return s
}

// Start a Server
func (s *Server) Start() string {
	s.testServer = httptest.NewServer(s)
	s.URL = s.testServer.URL
	return s.testServer.URL
}

// Stop a running Server
func (s *Server) Stop() {
	s.testServer.Close()
}


func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	go func() {
		s.requests <- r
	}()

	select {
	case h := <-s.handlers:
		h.ServeHTTP(w, r)
	default:
		w.WriteHeader(200)
	}
}

func (s *Server) Enqueue(h http.HandlerFunc) {
	go func() {
		s.handlers <- h
	}()
}

func (s *Server) TakeRequest() *http.Request {
	return <-s.requests
}

func (s *Server) TakeRequestWithTimeout(duration time.Duration) *http.Request {
	select {
	case r := <-s.requests:
		return r
	case <-time.After(duration):
		return nil
	}
}
