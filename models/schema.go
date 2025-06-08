package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Schema represents a database schema definition
type Schema struct {
	ID               uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name             string         `json:"name" gorm:"uniqueIndex;not null"`
	Description      string         `json:"description"`
	DatabaseName     string         `json:"databaseName" gorm:"uniqueIndex;not null"`
	Status           string         `json:"status" gorm:"not null;default:'created'"`
	Version          string         `json:"version" gorm:"not null;default:'1.0'"`
	SchemaDefinition SchemaData     `json:"schemaDefinition" gorm:"type:jsonb"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`
}

// SchemaData represents the complete schema definition structure
type SchemaData struct {
	Tables      []Table      `json:"tables"`
	ForeignKeys []ForeignKey `json:"foreignKeys"`
	Version     string       `json:"version"`
	ExportedAt  string       `json:"exportedAt,omitempty"`
}

// Value implements the driver.Valuer interface for database storage
func (s SchemaData) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Scan implements the sql.Scanner interface for database retrieval
func (s *SchemaData) Scan(value interface{}) error {
	if value == nil {
		*s = SchemaData{Tables: []Table{}, ForeignKeys: []ForeignKey{}}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	case nil:
		*s = SchemaData{Tables: []Table{}, ForeignKeys: []ForeignKey{}}
		return nil
	default:
		return errors.New("cannot scan SchemaData from non-byte value")
	}

	if len(bytes) == 0 {
		*s = SchemaData{Tables: []Table{}, ForeignKeys: []ForeignKey{}}
		return nil
	}

	err := json.Unmarshal(bytes, s)
	if err != nil {
		// If unmarshal fails, initialize with empty values
		*s = SchemaData{Tables: []Table{}, ForeignKeys: []ForeignKey{}}
	}
	return nil
}

// Table represents a database table definition
type Table struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Columns  []Column `json:"columns"`
	Position Position `json:"position"`
	Indexes  []Index  `json:"indexes,omitempty"`
}

// Column represents a database column definition
type Column struct {
	ID            string      `json:"id"`
	Name          string      `json:"name"`
	DataType      string      `json:"dataType"`
	Length        *int        `json:"length,omitempty"`
	Precision     *int        `json:"precision,omitempty"`
	Scale         *int        `json:"scale,omitempty"`
	Nullable      bool        `json:"nullable"`
	PrimaryKey    bool        `json:"primaryKey"`
	AutoIncrement bool        `json:"autoIncrement"`
	Unique        bool        `json:"unique,omitempty"`
	DefaultValue  interface{} `json:"defaultValue,omitempty"`
}

// ForeignKey represents a foreign key relationship
type ForeignKey struct {
	ID             string `json:"id"`
	Name           string `json:"name,omitempty"`
	SourceTableId  string `json:"sourceTableId"`
	SourceColumnId string `json:"sourceColumnId"`
	TargetTableId  string `json:"targetTableId"`
	TargetColumnId string `json:"targetColumnId"`
	OnDelete       string `json:"onDelete"`
	OnUpdate       string `json:"onUpdate"`
}

// Position represents UI positioning for tables
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Index represents a database index
type Index struct {
	Name    string   `json:"name"`
	Columns []string `json:"columns"`
	Unique  bool     `json:"unique"`
}

// CreateSchemaRequest represents the request structure for creating a schema
type CreateSchemaRequest struct {
	Name        string       `json:"name" binding:"required,min=1,max=100"`
	Description string       `json:"description" binding:"max=500"`
	Tables      []Table      `json:"tables" binding:"required,min=1"`
	ForeignKeys []ForeignKey `json:"foreignKeys"`
}

// UpdateSchemaRequest represents the request structure for updating a schema
type UpdateSchemaRequest struct {
	Name        string       `json:"name" binding:"required,min=1,max=100"`
	Description string       `json:"description" binding:"max=500"`
	Tables      []Table      `json:"tables" binding:"required,min=1"`
	ForeignKeys []ForeignKey `json:"foreignKeys"`
}

// SchemaListResponse represents a simplified schema for listing
type SchemaListResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	DatabaseName string    `json:"databaseName"`
	Status       string    `json:"status"`
	TableCount   int       `json:"tableCount"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	Version      string    `json:"version"`
}

// SchemaValidationRequest represents the request for schema validation
type SchemaValidationRequest struct {
	Name        string       `json:"name" binding:"required"`
	Tables      []Table      `json:"tables" binding:"required,min=1"`
	ForeignKeys []ForeignKey `json:"foreignKeys"`
}

// ValidationResult represents the result of schema validation
type ValidationResult struct {
	Valid        bool              `json:"valid"`
	Errors       []ValidationError `json:"errors,omitempty"`
	Warnings     []string          `json:"warnings,omitempty"`
	GeneratedSQL []string          `json:"generatedSQL,omitempty"`
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

// DatabaseStatus represents the status of a generated database
type DatabaseStatus struct {
	SchemaID         uuid.UUID `json:"schemaId"`
	DatabaseName     string    `json:"databaseName"`
	Status           string    `json:"status"`
	TableCount       int       `json:"tableCount"`
	LastChecked      time.Time `json:"lastChecked"`
	ConnectionString string    `json:"connectionString,omitempty"`
}

// SQLExportResponse represents the response for SQL export
type SQLExportResponse struct {
	SchemaID    uuid.UUID `json:"schemaId"`
	SQL         string    `json:"sql"`
	GeneratedAt time.Time `json:"generatedAt"`
}

// PaginationRequest represents pagination parameters
type PaginationRequest struct {
	Page   int    `form:"page,default=1" binding:"min=1"`
	Limit  int    `form:"limit,default=10" binding:"min=1,max=100"`
	Search string `form:"search"`
}

// Supported data types
var SupportedDataTypes = map[string]bool{
	"INT":       true,
	"BIGINT":    true,
	"VARCHAR":   true,
	"TEXT":      true,
	"BOOLEAN":   true,
	"TIMESTAMP": true,
	"DATE":      true,
	"TIME":      true,
	"DECIMAL":   true,
	"FLOAT":     true,
	"DOUBLE":    true,
	"JSON":      true,
	"UUID":      true,
}

// Valid foreign key actions
var ValidForeignKeyActions = map[string]bool{
	"CASCADE":   true,
	"RESTRICT":  true,
	"SET NULL":  true,
	"NO ACTION": true,
}

// BeforeCreate sets up UUID before creating the schema
func (s *Schema) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}