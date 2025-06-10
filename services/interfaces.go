package services

import (
	"fmt"
	"log"
	"strings"
	"time"

	"vdt-dashboard-backend/config"
	"vdt-dashboard-backend/models"
	"vdt-dashboard-backend/repositories"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SchemaService defines the interface for schema business logic
type SchemaService interface {
	CreateSchema(request models.CreateSchemaRequest, userID uuid.UUID) (*models.Schema, error)
	GetSchema(id, userID uuid.UUID) (*models.Schema, error)
	UpdateSchema(id, userID uuid.UUID, request models.UpdateSchemaRequest) (*models.Schema, error)
	DeleteSchema(id, userID uuid.UUID) error
	ListSchemas(pagination models.PaginationRequest, userID uuid.UUID) ([]models.SchemaListResponse, *models.PaginationResponse, error)
	ExportSQL(id, userID uuid.UUID) (*models.SQLExportResponse, error)
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
func NewSchemaService(repo repositories.SchemaRepository, databaseManager DatabaseManagerService, cfg *config.Config) SchemaService {
	return &schemaService{
		repo:            repo,
		databaseManager: databaseManager,
		config:          cfg,
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
	repo            repositories.SchemaRepository
	databaseManager DatabaseManagerService
	config          *config.Config
}

type validatorService struct{}

type sqlGeneratorService struct{}

type databaseManagerService struct {
	config *config.Config
}

// SchemaService implementation
func (s *schemaService) CreateSchema(request models.CreateSchemaRequest, userID uuid.UUID) (*models.Schema, error) {
	// Check if schema name already exists for this user
	if _, err := s.repo.GetByNameAndUserID(request.Name, userID); err == nil {
		return nil, fmt.Errorf("schema with name '%s' already exists", request.Name)
	}

	// Generate unique database name
	databaseName := fmt.Sprintf("schema_%s", strings.ReplaceAll(uuid.New().String(), "-", "_"))

	schema := &models.Schema{
		ID:           uuid.New(),
		Name:         request.Name,
		Description:  request.Description,
		DatabaseName: databaseName,
		Status:       "creating",
		Version:      "1.0",
		UserID:       userID,
		SchemaDefinition: models.SchemaData{
			Tables:      request.Tables,
			ForeignKeys: request.ForeignKeys,
			Version:     "1.0",
			ExportedAt:  time.Now().Format(time.RFC3339),
		},
	}

	// Create schema metadata first
	if err := s.repo.Create(schema); err != nil {
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	// Generate the actual database
	if err := s.databaseManager.RegenerateDatabase(schema.SchemaDefinition, schema.DatabaseName); err != nil {
		// Update status to error
		schema.Status = "error"
		s.repo.Update(schema)
		return nil, fmt.Errorf("failed to generate database: %w", err)
	}

	// Update status to created
	schema.Status = "created"
	if err := s.repo.Update(schema); err != nil {
		log.Printf("Warning: failed to update schema status: %v", err)
	}

	return schema, nil
}

func (s *schemaService) GetSchema(id, userID uuid.UUID) (*models.Schema, error) {
	return s.repo.GetByIDAndUserID(id, userID)
}

func (s *schemaService) UpdateSchema(id, userID uuid.UUID, request models.UpdateSchemaRequest) (*models.Schema, error) {
	schema, err := s.repo.GetByIDAndUserID(id, userID)
	if err != nil {
		return nil, err
	}

	// Check if new name conflicts with existing schema for this user (excluding current schema)
	if schema.Name != request.Name {
		if existing, err := s.repo.GetByNameAndUserID(request.Name, userID); err == nil && existing.ID != id {
			return nil, fmt.Errorf("schema with name '%s' already exists", request.Name)
		}
	}

	// Update schema definition
	schema.Name = request.Name
	schema.Description = request.Description
	schema.Status = "updating"
	schema.SchemaDefinition = models.SchemaData{
		Tables:      request.Tables,
		ForeignKeys: request.ForeignKeys,
		Version:     "1.1",
		ExportedAt:  time.Now().Format(time.RFC3339),
	}

	// Save schema metadata first
	if err := s.repo.Update(schema); err != nil {
		return nil, fmt.Errorf("failed to update schema: %w", err)
	}

	// Regenerate the database with new definition
	if err := s.databaseManager.RegenerateDatabase(schema.SchemaDefinition, schema.DatabaseName); err != nil {
		// Update status to error
		schema.Status = "error"
		s.repo.Update(schema)
		return nil, fmt.Errorf("failed to regenerate database: %w", err)
	}

	// Update status to updated
	schema.Status = "updated"
	if err := s.repo.Update(schema); err != nil {
		log.Printf("Warning: failed to update schema status: %v", err)
	}

	return schema, nil
}

func (s *schemaService) DeleteSchema(id, userID uuid.UUID) error {
	return s.repo.DeleteByIDAndUserID(id, userID)
}

func (s *schemaService) ListSchemas(pagination models.PaginationRequest, userID uuid.UUID) ([]models.SchemaListResponse, *models.PaginationResponse, error) {
	schemas, total, err := s.repo.ListByUserID(pagination, userID)
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

func (s *schemaService) ExportSQL(id, userID uuid.UUID) (*models.SQLExportResponse, error) {
	schema, err := s.repo.GetByIDAndUserID(id, userID)
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
		var columns []string
		var primaryKeys []string
		var uniqueConstraints []string

		// Generate column definitions
		for _, column := range table.Columns {
			columnDef := g.generateColumnDefinition(column)
			columns = append(columns, columnDef)

			if column.PrimaryKey {
				primaryKeys = append(primaryKeys, column.Name)
			}

			if column.Unique && !column.PrimaryKey {
				uniqueConstraints = append(uniqueConstraints, fmt.Sprintf("UNIQUE (%s)", column.Name))
			}
		}

		// Build CREATE TABLE statement
		statement := fmt.Sprintf("CREATE TABLE %s (\n", table.Name)
		statement += "    " + strings.Join(columns, ",\n    ")

		// Add primary key constraint
		if len(primaryKeys) > 0 {
			statement += fmt.Sprintf(",\n    PRIMARY KEY (%s)", strings.Join(primaryKeys, ", "))
		}

		// Add unique constraints
		for _, constraint := range uniqueConstraints {
			statement += fmt.Sprintf(",\n    %s", constraint)
		}

		statement += "\n);"
		statements = append(statements, statement)
	}

	return statements, nil
}

func (g *sqlGeneratorService) GenerateForeignKeys(schemaData models.SchemaData) ([]string, error) {
	var statements []string

	// First, create a map of table IDs to table names for lookup
	tableMap := make(map[string]string)
	columnMap := make(map[string]string)

	for _, table := range schemaData.Tables {
		tableMap[table.ID] = table.Name
		for _, column := range table.Columns {
			columnMap[column.ID] = column.Name
		}
	}

	for _, fk := range schemaData.ForeignKeys {
		sourceTable, sourceTableExists := tableMap[fk.SourceTableId]
		targetTable, targetTableExists := tableMap[fk.TargetTableId]
		sourceColumn, sourceColumnExists := columnMap[fk.SourceColumnId]
		targetColumn, targetColumnExists := columnMap[fk.TargetColumnId]

		if !sourceTableExists || !targetTableExists || !sourceColumnExists || !targetColumnExists {
			continue // Skip invalid foreign keys
		}

		constraintName := fk.Name
		if constraintName == "" {
			constraintName = fmt.Sprintf("fk_%s_%s", sourceTable, sourceColumn)
		}

		onDelete := "RESTRICT"
		if fk.OnDelete != "" && models.ValidForeignKeyActions[fk.OnDelete] {
			onDelete = fk.OnDelete
		}

		onUpdate := "RESTRICT"
		if fk.OnUpdate != "" && models.ValidForeignKeyActions[fk.OnUpdate] {
			onUpdate = fk.OnUpdate
		}

		statement := fmt.Sprintf(
			"ALTER TABLE %s ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s (%s) ON DELETE %s ON UPDATE %s;",
			sourceTable,
			constraintName,
			sourceColumn,
			targetTable,
			targetColumn,
			onDelete,
			onUpdate,
		)
		statements = append(statements, statement)
	}

	return statements, nil
}

// generateColumnDefinition creates SQL column definition from column model
func (g *sqlGeneratorService) generateColumnDefinition(column models.Column) string {
	var def strings.Builder

	def.WriteString(column.Name)
	def.WriteString(" ")

	// Data type mapping
	switch column.DataType {
	case "INT":
		if column.AutoIncrement {
			def.WriteString("SERIAL")
		} else {
			def.WriteString("INTEGER")
		}
	case "BIGINT":
		if column.AutoIncrement {
			def.WriteString("BIGSERIAL")
		} else {
			def.WriteString("BIGINT")
		}
	case "VARCHAR":
		length := 255
		if column.Length != nil && *column.Length > 0 {
			length = *column.Length
		}
		def.WriteString(fmt.Sprintf("VARCHAR(%d)", length))
	case "TEXT":
		def.WriteString("TEXT")
	case "BOOLEAN":
		def.WriteString("BOOLEAN")
	case "TIMESTAMP":
		def.WriteString("TIMESTAMP WITH TIME ZONE")
	case "DATE":
		def.WriteString("DATE")
	case "TIME":
		def.WriteString("TIME")
	case "DECIMAL":
		precision := 10
		scale := 2
		if column.Precision != nil {
			precision = *column.Precision
		}
		if column.Scale != nil {
			scale = *column.Scale
		}
		def.WriteString(fmt.Sprintf("DECIMAL(%d,%d)", precision, scale))
	case "FLOAT":
		def.WriteString("REAL")
	case "DOUBLE":
		def.WriteString("DOUBLE PRECISION")
	case "JSON":
		def.WriteString("JSONB")
	case "UUID":
		def.WriteString("UUID")
	default:
		def.WriteString("TEXT") // Fallback
	}

	// Nullable constraint
	if !column.Nullable {
		def.WriteString(" NOT NULL")
	}

	// Default value
	if column.DefaultValue != nil {
		switch v := column.DefaultValue.(type) {
		case string:
			if v != "" {
				def.WriteString(fmt.Sprintf(" DEFAULT '%s'", v))
			}
		case bool:
			def.WriteString(fmt.Sprintf(" DEFAULT %t", v))
		case float64:
			def.WriteString(fmt.Sprintf(" DEFAULT %v", v))
		}
	}

	// UUID default for UUID columns
	if column.DataType == "UUID" && column.DefaultValue == nil {
		def.WriteString(" DEFAULT gen_random_uuid()")
	}

	// Timestamp defaults
	if column.DataType == "TIMESTAMP" && column.DefaultValue == nil {
		def.WriteString(" DEFAULT CURRENT_TIMESTAMP")
	}

	return def.String()
}

// DatabaseManagerService implementation
func (d *databaseManagerService) CreateDatabase(databaseName string) error {
	return config.CreateDynamicDatabase(d.config, databaseName)
}

func (d *databaseManagerService) DropDatabase(databaseName string) error {
	return config.DropDynamicDatabase(d.config, databaseName)
}

func (d *databaseManagerService) GetDatabaseStatus(databaseName string) (*models.DatabaseStatus, error) {
	// Connect to the user's database to check status
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		d.config.DatabaseHost,
		d.config.DatabasePort,
		d.config.DatabaseUser,
		d.config.DatabasePass,
		databaseName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return &models.DatabaseStatus{
			DatabaseName: databaseName,
			Status:       "error",
			TableCount:   0,
			LastChecked:  time.Now(),
		}, nil
	}

	// Count tables
	var tableCount int64
	err = db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE'").Scan(&tableCount).Error
	if err != nil {
		tableCount = 0
	}

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		d.config.DatabaseUser,
		"***", // Hide password
		d.config.DatabaseHost,
		d.config.DatabasePort,
		databaseName,
	)

	return &models.DatabaseStatus{
		DatabaseName:     databaseName,
		Status:           "healthy",
		TableCount:       int(tableCount),
		LastChecked:      time.Now(),
		ConnectionString: connectionString,
	}, nil
}

func (d *databaseManagerService) RegenerateDatabase(schemaData models.SchemaData, databaseName string) error {
	// Create SQL generator
	sqlGen := &sqlGeneratorService{}

	// Drop existing database
	if err := d.DropDatabase(databaseName); err != nil {
		// Ignore error if database doesn't exist
		log.Printf("Warning: Failed to drop database %s: %v", databaseName, err)
	}

	// Create new database
	if err := d.CreateDatabase(databaseName); err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	// Connect to the new database
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		d.config.DatabaseHost,
		d.config.DatabasePort,
		d.config.DatabaseUser,
		d.config.DatabasePass,
		databaseName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to new database: %w", err)
	}

	// Generate and execute table creation statements
	tableStatements, err := sqlGen.GenerateCreateTables(schemaData)
	if err != nil {
		return fmt.Errorf("failed to generate table statements: %w", err)
	}

	for _, statement := range tableStatements {
		if err := db.Exec(statement).Error; err != nil {
			return fmt.Errorf("failed to execute table statement: %w\nStatement: %s", err, statement)
		}
	}

	// Generate and execute foreign key statements
	fkStatements, err := sqlGen.GenerateForeignKeys(schemaData)
	if err != nil {
		return fmt.Errorf("failed to generate foreign key statements: %w", err)
	}

	for _, statement := range fkStatements {
		if err := db.Exec(statement).Error; err != nil {
			return fmt.Errorf("failed to execute foreign key statement: %w\nStatement: %s", err, statement)
		}
	}

	log.Printf("Successfully regenerated database %s with %d tables", databaseName, len(schemaData.Tables))
	return nil
}
