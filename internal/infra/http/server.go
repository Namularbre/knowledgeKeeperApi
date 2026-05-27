package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Namularbre/knowledgeKeeperApi/internal/version"
)

type Server struct {
	httpServer *http.Server
	router     *http.ServeMux
}

type RouteRegistrar func(mux *http.ServeMux)

func NewServer(port string) *Server {
	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		httpServer: srv,
		router:     mux,
	}
}

func (s *Server) RegisterRoutes(extra ...RouteRegistrar) {
	s.router.HandleFunc("/version", s.handleVersion)
	s.router.HandleFunc("/health", s.handleHealth)
	for _, reg := range extra {
		reg(s.router)
	}
}

func (s *Server) Start() error {
	log.Printf("Server starting on %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Server shutting down...")
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) handleVersion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"version":%q,"api":%q}`, version.Version, version.APIName)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
