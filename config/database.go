package config

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDatabase initializes the database connection
func InitDatabase(config *Config) (*gorm.DB, error) {
	var dsn string

	// Use DATABASE_URL if provided, otherwise construct from individual components
	if config.DatabaseURL != "" {
		dsn = config.DatabaseURL
	} else {
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
			config.DatabaseHost,
			config.DatabasePort,
			config.DatabaseUser,
			config.DatabasePass,
			config.DatabaseName,
		)
	}

	// Configure GORM logger
	var gormLogger logger.Interface
	if config.Environment == "development" {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Silent)
	}

	// Open database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connected successfully")
	return db, nil
}

// CreateDynamicDatabase creates a new database for user schemas
func CreateDynamicDatabase(config *Config, databaseName string) error {
	// Connect to postgres database to create new database
	var dsn string

	if config.DatabaseURL != "" {
		// For DATABASE_URL, we need to connect to the default postgres database
		dsn = config.DatabaseURL + "_postgres"
	} else {
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
			config.DatabaseHost,
			config.DatabasePort,
			config.DatabaseUser,
			config.DatabasePass,
		)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to postgres database: %w", err)
	}

	// Create the new database
	createSQL := fmt.Sprintf("CREATE DATABASE %s", databaseName)
	if err := db.Exec(createSQL).Error; err != nil {
		return fmt.Errorf("failed to create database %s: %w", databaseName, err)
	}

	log.Printf("Database %s created successfully", databaseName)
	return nil
}

// DropDynamicDatabase drops a user schema database
func DropDynamicDatabase(config *Config, databaseName string) error {
	// Connect to postgres database to drop database
	var dsn string

	if config.DatabaseURL != "" {
		dsn = config.DatabaseURL + "_postgres"
	} else {
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
			config.DatabaseHost,
			config.DatabasePort,
			config.DatabaseUser,
			config.DatabasePass,
		)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to postgres database: %w", err)
	}

	// Drop the database
	dropSQL := fmt.Sprintf("DROP DATABASE IF EXISTS %s", databaseName)
	if err := db.Exec(dropSQL).Error; err != nil {
		return fmt.Errorf("failed to drop database %s: %w", databaseName, err)
	}

	log.Printf("Database %s dropped successfully", databaseName)
	return nil
}
