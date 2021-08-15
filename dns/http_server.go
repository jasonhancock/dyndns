package dns

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// SVC is the interface the HTTP service requires.
type SVC interface {
	DNS(ctx context.Context, req Request) error
}

// HTTPServer is the transport layer for the service.
type HTTPServer struct {
	svc    SVC
	router *mux.Router
}

// NewHTTPServer initializes an HTTPServer
func NewHTTPServer(r *mux.Router, svc SVC) *HTTPServer {
	s := &HTTPServer{
		svc:    svc,
		router: r,
	}

	s.router.Methods(http.MethodPost).Path("/dns").HandlerFunc(s.handler)

	return s
}

type Request struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ServeHTTP fulfills the http.Handler interface and allows to use the HTTPServer type in a handler call.
func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *HTTPServer) handler(w http.ResponseWriter, r *http.Request) {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := s.svc.DNS(r.Context(), req)
	if err != nil {
		if err == errNotAllowed {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
