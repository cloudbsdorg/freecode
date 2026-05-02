package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	Addr    string
	Router  *chi.Mux
	Server  *http.Server
}

func New(addr string) *Server {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	return &Server{
		Addr:   addr,
		Router: r,
		Server: &http.Server{
			Addr:         addr,
			Handler:      r,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

func (s *Server) Handle(pattern string, handler http.Handler) {
	s.Router.Handle(pattern, handler)
}

func (s *Server) HandleFunc(pattern string, handler http.HandlerFunc) {
	s.Router.HandleFunc(pattern, handler)
}

func (s *Server) Mount(pattern string, handler http.Handler) {
	s.Router.Mount(pattern, handler)
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		s.Server.Close()
	}()

	return s.Server.ListenAndServe()
}

func (s *Server) Stop() error {
	return s.Server.Close()
}

func (s *Server) URL() string {
	return fmt.Sprintf("http://%s", s.Addr)
}
