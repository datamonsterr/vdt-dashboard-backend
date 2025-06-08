package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"vdt-dashboard-backend/config"
	"vdt-dashboard-backend/models"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Get command line argument
	command := "up"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	// Initialize configuration
	cfg := config.Load()

	// Initialize database
	db, err := config.InitDatabase(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	switch command {
	case "up":
		if err := runMigrations(db); err != nil {
			log.Fatal("Migration failed:", err)
		}
		log.Println("✅ Migrations completed successfully")
	case "create-models":
		if err := createModels(db); err != nil {
			log.Fatal("Failed to create models:", err)
		}
		log.Println("✅ Models created successfully")
	case "seed":
		if err := seedData(db); err != nil {
			log.Fatal("Seeding failed:", err)
		}
		log.Println("✅ Data seeded successfully")
	case "reset":
		if err := resetDatabase(db); err != nil {
			log.Fatal("Reset failed:", err)
		}
		log.Println("✅ Database reset successfully")
	default:
		log.Printf("Unknown command: %s", command)
		log.Println("Available commands: up, create-models, seed, reset")
		os.Exit(1)
	}
}

// runMigrations runs all SQL migration files
func runMigrations(db *gorm.DB) error {
	log.Println("🔄 Running SQL migrations...")

	// Get current directory
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Read migration files
	files, err := filepath.Glob(filepath.Join(dir, "*.sql"))
	if err != nil {
		return fmt.Errorf("failed to read migration files: %w", err)
	}

	// Sort files to ensure proper order
	sort.Strings(files)

	// Execute each migration file
	for _, file := range files {
		log.Printf("📄 Executing migration: %s", filepath.Base(file))
		
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		// Split file content by statements (separated by semicolons followed by newlines)
		statements := strings.Split(string(content), ";\n")
		
		for i, statement := range statements {
			statement = strings.TrimSpace(statement)
			if statement == "" || strings.HasPrefix(statement, "--") {
				continue
			}

			if err := db.Exec(statement).Error; err != nil {
				return fmt.Errorf("failed to execute statement %d in file %s: %w\nStatement: %s", i+1, file, err, statement)
			}
		}
	}

	return nil
}

// createModels creates database tables using GORM AutoMigrate
func createModels(db *gorm.DB) error {
	log.Println("🔄 Creating models with GORM AutoMigrate...")

	// AutoMigrate will create tables, missing columns, missing indexes
	// It will NOT delete unused columns to protect data
	if err := db.AutoMigrate(&models.Schema{}); err != nil {
		return fmt.Errorf("failed to migrate Schema model: %w", err)
	}

	log.Println("✅ Models created/updated successfully")
	return nil
}

// seedData runs seed migration files (files starting with seed_)
func seedData(db *gorm.DB) error {
	log.Println("🔄 Seeding database with sample data...")

	// Check if we already have data
	var count int64
	if err := db.Model(&models.Schema{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check existing data: %w", err)
	}

	if count > 0 {
		log.Printf("📊 Found %d existing schemas, skipping seed data", count)
		return nil
	}

	// Run seed migrations
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Look for seed files
	files, err := filepath.Glob(filepath.Join(dir, "*seed*.sql"))
	if err != nil {
		return fmt.Errorf("failed to read seed files: %w", err)
	}

	// Also run 002_seed_data.sql if it exists
	seedFile := filepath.Join(dir, "002_seed_data.sql")
	if _, err := os.Stat(seedFile); err == nil {
		files = append(files, seedFile)
	}

	if len(files) == 0 {
		log.Println("📝 No seed files found")
		return nil
	}

	sort.Strings(files)

	for _, file := range files {
		log.Printf("📄 Executing seed file: %s", filepath.Base(file))
		
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read seed file %s: %w", file, err)
		}

		if err := db.Exec(string(content)).Error; err != nil {
			return fmt.Errorf("failed to execute seed file %s: %w", file, err)
		}
	}

	return nil
}

// resetDatabase drops all tables and recreates them
func resetDatabase(db *gorm.DB) error {
	log.Println("⚠️  Resetting database (this will delete all data)...")

	// Drop tables
	if err := db.Migrator().DropTable(&models.Schema{}); err != nil {
		log.Printf("Warning: failed to drop schemas table: %v", err)
	}

	// Recreate tables
	if err := createModels(db); err != nil {
		return fmt.Errorf("failed to recreate models: %w", err)
	}

	// Run migrations
	if err := runMigrations(db); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Seed data
	if err := seedData(db); err != nil {
		return fmt.Errorf("failed to seed data: %w", err)
	}

	return nil
} 