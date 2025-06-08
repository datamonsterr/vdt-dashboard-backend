package handlers

import (
	"net/http"
	"time"

	"vdt-dashboard-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	db *gorm.DB
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

// HealthCheck handles GET /health
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	// Check database connection
	sqlDB, err := h.db.DB()
	var dbStatus string
	if err != nil {
		dbStatus = "disconnected"
	} else if err := sqlDB.Ping(); err != nil {
		dbStatus = "unhealthy"
	} else {
		dbStatus = "connected"
	}

	health := gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"database":  dbStatus,
		"version":   "1.0.0",
	}

	statusCode := http.StatusOK
	if dbStatus != "connected" {
		health["status"] = "unhealthy"
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, models.SuccessResponse("Service health check", health))
}
