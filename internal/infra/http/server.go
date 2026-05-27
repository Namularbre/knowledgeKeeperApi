package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	httpSwagger "github.com/swaggo/http-swagger/v2"

	"github.com/Namularbre/knowledgeKeeperApi/internal/version"
)

type Server struct {
	httpServer *http.Server
	router     *http.ServeMux
}

type RouteRegistrar func(mux *http.ServeMux)

// VersionResponse is the payload returned by GET /version.
type VersionResponse struct {
	Version string `json:"version" example:"1.0.0"`
	API     string `json:"api" example:"knowledgeKeeperApi"`
}

// HealthResponse is the payload returned by GET /health.
type HealthResponse struct {
	Status string `json:"status" example:"ok"`
}

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
	s.router.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
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

// handleVersion godoc
// @Summary      API version
// @Description  Returns the API name and the build-stamped semantic version.
// @Tags         meta
// @Produce      json
// @Success      200  {object}  VersionResponse
// @Router       /version [get]
func (s *Server) handleVersion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"version":%q,"api":%q}`, version.Version, version.APIName)
}

// handleHealth godoc
// @Summary      Liveness probe
// @Description  Returns ok when the server is up. Does not check downstream dependencies.
// @Tags         meta
// @Produce      json
// @Success      200  {object}  HealthResponse
// @Router       /health [get]
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
