package web

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mrlyc/inu/pkg/anonymizer"
	"github.com/mrlyc/inu/pkg/web/handlers"
	"github.com/mrlyc/inu/pkg/web/middleware"
)

//go:embed static/*
var staticFS embed.FS

const version = "v0.1.0"

// Server represents the HTTP API server
type Server struct {
	config      *Config
	anonymizer  anonymizer.Anonymizer
	engine      *gin.Engine
	httpServer  *http.Server
	entityTypes []string
}

// NewServer creates a new web server instance
func NewServer(anon anonymizer.Anonymizer, config *Config) (*Server, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// Set Gin mode to release for production
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())

	s := &Server{
		config:      config,
		anonymizer:  anon,
		engine:      engine,
		entityTypes: anonymizer.DefaultEntityTypes, // 使用默认实体类型
	}

	s.setupRoutes()

	return s, nil
}

// SetEntityTypes sets the entity types for the server
func (s *Server) SetEntityTypes(types []string) {
	if len(types) > 0 {
		s.entityTypes = types
	}
}

// setupRoutes configures all HTTP routes and middleware
func (s *Server) setupRoutes() {
	// Create embedded static file system
	staticSubFS, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatalf("Failed to create static filesystem: %v", err)
	}
	httpFS := http.FS(staticSubFS)

	// Determine if auth is enabled
	authEnabled := s.config.IsAuthEnabled()

	// UI routes (auth required if enabled)
	ui := s.engine.Group("/")
	if authEnabled {
		ui.Use(middleware.BasicAuth(s.config.AdminUser, s.config.AdminToken))
	}
	{
		ui.GET("/", func(c *gin.Context) {
			c.FileFromFS("index.html", httpFS)
		})
		ui.StaticFS("/static", httpFS)
	}

	// Health check endpoint (no auth required)
	s.engine.GET("/health", handlers.HealthHandler(version))

	// API v1 endpoints (auth required if enabled)
	v1 := s.engine.Group("/api/v1")
	if authEnabled {
		v1.Use(middleware.BasicAuth(s.config.AdminUser, s.config.AdminToken))
	}
	{
		v1.GET("/config", handlers.ConfigHandler(s.entityTypes))
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

	if s.config.IsAuthEnabled() {
		log.Printf("  Authentication: ENABLED")
		log.Printf("  Admin user: %s", s.config.AdminUser)
	} else {
		log.Printf("  Authentication: DISABLED")
		log.Printf("  ⚠️  WARNING: Running without authentication!")
	}

	log.Printf("  Available endpoints:")
	log.Printf("    GET  /              (Web UI)")
	log.Printf("    GET  /health        (No auth)")
	log.Printf("    GET  /api/v1/config")
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
