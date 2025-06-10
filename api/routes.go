package api

import (
	"vdt-dashboard-backend/api/handlers"
	"vdt-dashboard-backend/api/middleware"
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
	userRepo := repositories.NewUserRepository(db)

	// Initialize services
	databaseManagerService := services.NewDatabaseManagerService(cfg)
	schemaService := services.NewSchemaService(schemaRepo, databaseManagerService, cfg)
	validatorService := services.NewValidatorService()
	sqlGeneratorService := services.NewSQLGeneratorService()

	// Initialize handlers
	schemaHandler := handlers.NewSchemaHandler(schemaService)
	healthHandler := handlers.NewHealthHandler(db)
	validatorHandler := handlers.NewValidatorHandler(validatorService, sqlGeneratorService)
	databaseHandler := handlers.NewDatabaseHandler(databaseManagerService, schemaService)
	userHandler := handlers.NewUserHandler()

	// Health check
	router.GET("/health", healthHandler.HealthCheck)

	// User routes (protected)
	userRoutes := router.Group("/user")
	userRoutes.Use(middleware.AuthMiddleware(userRepo, cfg.ClerkSecretKey)) // Apply authentication middleware
	{
		userRoutes.GET("/me", userHandler.GetCurrentUser)
	}

	// Schema management routes (protected)
	schemaRoutes := router.Group("/schemas")
	schemaRoutes.Use(middleware.AuthMiddleware(userRepo, cfg.ClerkSecretKey)) // Apply authentication middleware
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
