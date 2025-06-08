-- Migration: 001_create_schemas.sql
-- Description: Create schemas table for storing user-defined database schemas
-- Created: 2024-06-08

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create schemas table
CREATE TABLE IF NOT EXISTS schemas (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    database_name VARCHAR(255) NOT NULL UNIQUE,
    status VARCHAR(50) NOT NULL DEFAULT 'created',
    version VARCHAR(20) NOT NULL DEFAULT '1.0',
    schema_definition JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_schemas_name ON schemas(name);
CREATE INDEX IF NOT EXISTS idx_schemas_created_at ON schemas(created_at);
CREATE INDEX IF NOT EXISTS idx_schemas_status ON schemas(status);
CREATE INDEX IF NOT EXISTS idx_schemas_database_name ON schemas(database_name);
CREATE INDEX IF NOT EXISTS idx_schemas_deleted_at ON schemas(deleted_at) WHERE deleted_at IS NULL;

-- Create function to automatically update updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_schemas_updated_at 
    BEFORE UPDATE ON schemas 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Add comments for documentation
COMMENT ON TABLE schemas IS 'Stores user-defined database schema definitions';
COMMENT ON COLUMN schemas.id IS 'Unique identifier for the schema';
COMMENT ON COLUMN schemas.name IS 'Human-readable name for the schema';
COMMENT ON COLUMN schemas.description IS 'Optional description of the schema';
COMMENT ON COLUMN schemas.database_name IS 'Generated database name for the schema';
COMMENT ON COLUMN schemas.status IS 'Current status: created, updated, error, etc.';
COMMENT ON COLUMN schemas.version IS 'Schema version for tracking changes';
COMMENT ON COLUMN schemas.schema_definition IS 'Complete JSON schema definition including tables, columns, foreign keys';
COMMENT ON COLUMN schemas.created_at IS 'Timestamp when schema was created';
COMMENT ON COLUMN schemas.updated_at IS 'Timestamp when schema was last updated';
COMMENT ON COLUMN schemas.deleted_at IS 'Soft delete timestamp'; 