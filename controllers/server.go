package controllers

import (
	"log"
	"net/http"

	"api_monitoring_stats/controllers/oas"
)

type Server struct {
	httpServer *http.Server
	mux        *http.ServeMux
}

func NewServer(handler *Handler, address string) (*Server, error) {
	server, err := oas.NewServer(handler)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.Handle("/", server)

	return &Server{
		mux: mux,
		httpServer: &http.Server{
			Addr:    address,
			Handler: mux,
		},
	}, nil
}

func (s *Server) Run() {
	err := s.httpServer.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
