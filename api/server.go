package api

import (
	"vdt-dashboard-backend/api/middleware"
	"vdt-dashboard-backend/config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Server represents the HTTP server
type Server struct {
	router *gin.Engine
	db     *gorm.DB
	config *config.Config
}

// NewServer creates a new HTTP server
func NewServer(db *gorm.DB, cfg *config.Config) *Server {
	server := &Server{
		db:     db,
		config: cfg,
	}

	server.setupRouter()
	return server
}

// setupRouter configures the Gin router with middleware and routes
func (s *Server) setupRouter() {
	// Create router
	s.router = gin.New()

	// Add middleware
	s.router.Use(middleware.Logger())
	s.router.Use(middleware.Recovery())
	s.router.Use(middleware.CORS(s.config.AllowOrigins))
	s.router.Use(middleware.ErrorHandler())

	// Setup routes
	s.setupRoutes()
}

// setupRoutes configures all API routes
func (s *Server) setupRoutes() {
	// API v1 group
	v1 := s.router.Group("/api/v1")

	// Initialize routes
	SetupRoutes(v1, s.db, s.config)
}

// Run starts the HTTP server
func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}

// GetRouter returns the Gin router instance
func (s *Server) GetRouter() *gin.Engine {
	return s.router
}
