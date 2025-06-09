-- Migration: Add users table and update schemas table for user association
-- Up migration

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clerk_user_id VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) NOT NULL,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    profile_image_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create index on clerk_user_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_clerk_user_id ON users(clerk_user_id);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- Add user_id column to schemas table
ALTER TABLE schemas ADD COLUMN IF NOT EXISTS user_id UUID;

-- Create index on user_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_schemas_user_id ON schemas(user_id);

-- Add foreign key constraint
ALTER TABLE schemas ADD CONSTRAINT fk_schemas_user_id 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- Drop the old unique constraint on name and database_name
ALTER TABLE schemas DROP CONSTRAINT IF EXISTS schemas_name_key;
ALTER TABLE schemas DROP CONSTRAINT IF EXISTS schemas_database_name_key;

-- Add new unique constraint for name per user
ALTER TABLE schemas ADD CONSTRAINT unique_schema_name_per_user 
    UNIQUE (name, user_id);

-- Add unique constraint for database_name per user
ALTER TABLE schemas ADD CONSTRAINT unique_database_name_per_user 
    UNIQUE (database_name, user_id);

-- Update trigger for updated_at on users table
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); 