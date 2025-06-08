package api

import (
	"vdt-dashboard-backend/api/handlers"
	"vdt-dashboard-backend/config"
	"vdt-dashboard-backend/repositories"
	"vdt-dashboard-backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes configures all API routes
func SetupRoutes(router *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	// Initialize repositories
	schemaRepo := repositories.NewSchemaRepository(db)

	// Initialize services
	schemaService := services.NewSchemaService(schemaRepo, cfg)
	validatorService := services.NewValidatorService()
	sqlGeneratorService := services.NewSQLGeneratorService()
	databaseManagerService := services.NewDatabaseManagerService(cfg)

	// Initialize handlers
	schemaHandler := handlers.NewSchemaHandler(schemaService)
	healthHandler := handlers.NewHealthHandler(db)
	validatorHandler := handlers.NewValidatorHandler(validatorService, sqlGeneratorService)
	databaseHandler := handlers.NewDatabaseHandler(databaseManagerService, schemaService)

	// Health check
	router.GET("/health", healthHandler.HealthCheck)

	// Schema management routes
	schemaRoutes := router.Group("/schemas")
	{
		schemaRoutes.POST("", schemaHandler.CreateSchema)
		schemaRoutes.GET("", schemaHandler.ListSchemas)
		schemaRoutes.GET("/:id", schemaHandler.GetSchema)
		schemaRoutes.PUT("/:id", schemaHandler.UpdateSchema)
		schemaRoutes.DELETE("/:id", schemaHandler.DeleteSchema)

		// Schema export
		schemaRoutes.GET("/:id/export/sql", schemaHandler.ExportSQL)

		// Database management
		schemaRoutes.GET("/:id/database/status", databaseHandler.GetDatabaseStatus)
		schemaRoutes.POST("/:id/database/regenerate", databaseHandler.RegenerateDatabase)
	}

	// Validation routes
	router.POST("/schemas/validate", validatorHandler.ValidateSchema)
}
