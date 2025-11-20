package web

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mrlyc/inu/pkg/anonymizer"
	"github.com/mrlyc/inu/pkg/web/handlers"
	"github.com/mrlyc/inu/pkg/web/middleware"
)

const version = "v0.1.0"

// Server represents the HTTP API server
type Server struct {
	config     *Config
	anonymizer *anonymizer.Anonymizer
	engine     *gin.Engine
	httpServer *http.Server
}

// NewServer creates a new web server instance
func NewServer(anon *anonymizer.Anonymizer, config *Config) (*Server, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// Set Gin mode to release for production
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())

	s := &Server{
		config:     config,
		anonymizer: anon,
		engine:     engine,
	}

	s.setupRoutes()

	return s, nil
}

// setupRoutes configures all HTTP routes and middleware
func (s *Server) setupRoutes() {
	// Health check endpoint (no auth required)
	s.engine.GET("/health", handlers.HealthHandler(version))

	// API v1 endpoints (auth required)
	v1 := s.engine.Group("/api/v1")
	v1.Use(middleware.BasicAuth(s.config.AdminUser, s.config.AdminToken))
	{
		v1.POST("/anonymize", handlers.AnonymizeHandler(s.anonymizer))
		v1.POST("/restore", handlers.RestoreHandler(s.anonymizer))
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.httpServer = &http.Server{
		Addr:    s.config.Addr,
		Handler: s.engine,
	}

	log.Printf("Starting Inu Web API Server")
	log.Printf("  Version: %s", version)
	log.Printf("  Listening on: %s", s.config.Addr)
	log.Printf("  Admin user: %s", s.config.AdminUser)
	log.Printf("  Available endpoints:")
	log.Printf("    GET  /health")
	log.Printf("    POST /api/v1/anonymize")
	log.Printf("    POST /api/v1/restore")

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

// Stop gracefully shuts down the HTTP server
func (s *Server) Stop() error {
	if s.httpServer == nil {
		return nil
	}

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	log.Println("Server stopped gracefully")
	return nil
}
