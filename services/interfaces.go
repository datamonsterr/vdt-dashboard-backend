package services

import (
	"fmt"
	"strings"
	"time"

	"vdt-dashboard-backend/config"
	"vdt-dashboard-backend/models"
	"vdt-dashboard-backend/repositories"

	"github.com/google/uuid"
)

// SchemaService defines the interface for schema business logic
type SchemaService interface {
	CreateSchema(request models.CreateSchemaRequest) (*models.Schema, error)
	GetSchema(id uuid.UUID) (*models.Schema, error)
	UpdateSchema(id uuid.UUID, request models.UpdateSchemaRequest) (*models.Schema, error)
	DeleteSchema(id uuid.UUID) error
	ListSchemas(pagination models.PaginationRequest) ([]models.SchemaListResponse, *models.PaginationResponse, error)
	ExportSQL(id uuid.UUID) (*models.SQLExportResponse, error)
}

// ValidatorService defines the interface for schema validation
type ValidatorService interface {
	ValidateSchema(request models.SchemaValidationRequest) (*models.ValidationResult, error)
}

// SQLGeneratorService defines the interface for SQL generation
type SQLGeneratorService interface {
	GenerateCreateDatabase(databaseName string) (string, error)
	GenerateCreateTables(schemaData models.SchemaData) ([]string, error)
	GenerateForeignKeys(schemaData models.SchemaData) ([]string, error)
}

// DatabaseManagerService defines the interface for database management
type DatabaseManagerService interface {
	CreateDatabase(databaseName string) error
	DropDatabase(databaseName string) error
	GetDatabaseStatus(databaseName string) (*models.DatabaseStatus, error)
	RegenerateDatabase(schemaData models.SchemaData, databaseName string) error
}

// NewSchemaService creates a new schema service
func NewSchemaService(repo repositories.SchemaRepository, cfg *config.Config) SchemaService {
	return &schemaService{
		repo:   repo,
		config: cfg,
	}
}

// NewValidatorService creates a new validator service
func NewValidatorService() ValidatorService {
	return &validatorService{}
}

// NewSQLGeneratorService creates a new SQL generator service
func NewSQLGeneratorService() SQLGeneratorService {
	return &sqlGeneratorService{}
}

// NewDatabaseManagerService creates a new database manager service
func NewDatabaseManagerService(cfg *config.Config) DatabaseManagerService {
	return &databaseManagerService{
		config: cfg,
	}
}

// Service implementations
type schemaService struct {
	repo   repositories.SchemaRepository
	config *config.Config
}

type validatorService struct{}

type sqlGeneratorService struct{}

type databaseManagerService struct {
	config *config.Config
}

// SchemaService implementation
func (s *schemaService) CreateSchema(request models.CreateSchemaRequest) (*models.Schema, error) {
	// Generate unique database name
	databaseName := fmt.Sprintf("schema_%s", strings.ReplaceAll(uuid.New().String(), "-", "_"))

	schema := &models.Schema{
		ID:           uuid.New(),
		Name:         request.Name,
		Description:  request.Description,
		DatabaseName: databaseName,
		Status:       "created",
		Version:      "1.0",
		SchemaDefinition: models.SchemaData{
			Tables:      request.Tables,
			ForeignKeys: request.ForeignKeys,
			Version:     "1.0",
			ExportedAt:  time.Now().Format(time.RFC3339),
		},
	}

	if err := s.repo.Create(schema); err != nil {
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	return schema, nil
}

func (s *schemaService) GetSchema(id uuid.UUID) (*models.Schema, error) {
	return s.repo.GetByID(id)
}

func (s *schemaService) UpdateSchema(id uuid.UUID, request models.UpdateSchemaRequest) (*models.Schema, error) {
	schema, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	schema.Name = request.Name
	schema.Description = request.Description
	schema.SchemaDefinition = models.SchemaData{
		Tables:      request.Tables,
		ForeignKeys: request.ForeignKeys,
		Version:     "1.1",
		ExportedAt:  time.Now().Format(time.RFC3339),
	}

	if err := s.repo.Update(schema); err != nil {
		return nil, fmt.Errorf("failed to update schema: %w", err)
	}

	return schema, nil
}

func (s *schemaService) DeleteSchema(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *schemaService) ListSchemas(pagination models.PaginationRequest) ([]models.SchemaListResponse, *models.PaginationResponse, error) {
	schemas, total, err := s.repo.List(pagination)
	if err != nil {
		return nil, nil, err
	}

	totalPages := (total + pagination.Limit - 1) / pagination.Limit
	paginationResp := &models.PaginationResponse{
		Page:       pagination.Page,
		Limit:      pagination.Limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return schemas, paginationResp, nil
}

func (s *schemaService) ExportSQL(id uuid.UUID) (*models.SQLExportResponse, error) {
	schema, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Basic SQL generation placeholder
	sql := fmt.Sprintf("-- Generated SQL for schema: %s\n-- TODO: Implement SQL generation", schema.Name)

	return &models.SQLExportResponse{
		SchemaID:    schema.ID,
		SQL:         sql,
		GeneratedAt: time.Now(),
	}, nil
}

// ValidatorService implementation
func (v *validatorService) ValidateSchema(request models.SchemaValidationRequest) (*models.ValidationResult, error) {
	var errors []models.ValidationError
	var warnings []string

	// Basic validation
	if len(request.Tables) == 0 {
		errors = append(errors, models.ValidationError{
			Field:   "tables",
			Message: "At least one table is required",
			Code:    "MISSING_TABLES",
		})
	}

	// Validate each table has at least one primary key
	for i, table := range request.Tables {
		hasPrimaryKey := false
		for _, column := range table.Columns {
			if column.PrimaryKey {
				hasPrimaryKey = true
				break
			}
		}
		if !hasPrimaryKey {
			warnings = append(warnings, fmt.Sprintf("Table '%s' has no primary key defined", table.Name))
		}

		// Validate data types
		for j, column := range table.Columns {
			if !models.SupportedDataTypes[column.DataType] {
				errors = append(errors, models.ValidationError{
					Field:   fmt.Sprintf("tables[%d].columns[%d].dataType", i, j),
					Message: fmt.Sprintf("Unsupported data type: %s", column.DataType),
					Code:    "UNSUPPORTED_DATA_TYPE",
				})
			}
		}
	}

	return &models.ValidationResult{
		Valid:    len(errors) == 0,
		Errors:   errors,
		Warnings: warnings,
	}, nil
}

// SQLGeneratorService implementation
func (g *sqlGeneratorService) GenerateCreateDatabase(databaseName string) (string, error) {
	return fmt.Sprintf("CREATE DATABASE %s;", databaseName), nil
}

func (g *sqlGeneratorService) GenerateCreateTables(schemaData models.SchemaData) ([]string, error) {
	var statements []string

	for _, table := range schemaData.Tables {
		statement := fmt.Sprintf("CREATE TABLE %s (", table.Name)
		// TODO: Implement full column definition generation
		statement += "\n  -- TODO: Generate column definitions"
		statement += "\n);"
		statements = append(statements, statement)
	}

	return statements, nil
}

func (g *sqlGeneratorService) GenerateForeignKeys(schemaData models.SchemaData) ([]string, error) {
	var statements []string

	for _, fk := range schemaData.ForeignKeys {
		statement := fmt.Sprintf("-- TODO: Generate foreign key for %s", fk.ID)
		statements = append(statements, statement)
	}

	return statements, nil
}

// DatabaseManagerService implementation
func (d *databaseManagerService) CreateDatabase(databaseName string) error {
	return config.CreateDynamicDatabase(d.config, databaseName)
}

func (d *databaseManagerService) DropDatabase(databaseName string) error {
	return config.DropDynamicDatabase(d.config, databaseName)
}

func (d *databaseManagerService) GetDatabaseStatus(databaseName string) (*models.DatabaseStatus, error) {
	return &models.DatabaseStatus{
		DatabaseName: databaseName,
		Status:       "healthy",
		TableCount:   0,
		LastChecked:  time.Now(),
	}, nil
}

func (d *databaseManagerService) RegenerateDatabase(schemaData models.SchemaData, databaseName string) error {
	// TODO: Implement database regeneration logic
	return nil
}
