package main

import (
	"log"

	"vdt-dashboard-backend/api"
	"vdt-dashboard-backend/config"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize configuration
	cfg := config.Load()

	// Initialize database
	db, err := config.InitDatabase(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize server
	server := api.NewServer(db, cfg)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := server.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
