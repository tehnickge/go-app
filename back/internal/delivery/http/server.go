package httpdelivery

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	Router chi.Router
}

func NewServer(router chi.Router) *Server {
	return &Server{Router: router}
}

func (s *Server) Handler() http.Handler {
	return s.Router
}

