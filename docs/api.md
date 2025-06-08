# Database Schema Builder API Documentation

## Base URL
```
http://localhost:8080/api/v1
```

## Setup & Configuration

### Environment Variables

Create a `.env` file in the project root with the following configuration:

```env
# Server Configuration
PORT=8080
ENVIRONMENT=development

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=vdt_dashboard

# Frontend Configuration
FRONTEND_URL=http://localhost:3000
```

### Database Setup

1. **Install PostgreSQL**: Make sure PostgreSQL is installed and running on your system.

2. **Create Database**: 
   ```bash
   make db-create
   ```

3. **Run Migrations**:
   ```bash
   make migrate
   ```

4. **Seed Sample Data** (optional):
   ```bash
   make migrate-seed
   ```

5. **Complete Setup** (all steps at once):
   ```bash
   make db-setup
   ```

### Database Schema

The application uses a main PostgreSQL database to store schema metadata and creates separate databases for each user-defined schema.

**Main Database Tables:**

- **schemas**: Stores schema definitions and metadata
  - `id` (UUID): Unique identifier
  - `name` (VARCHAR): Human-readable schema name
  - `description` (TEXT): Optional description
  - `database_name` (VARCHAR): Generated database name
  - `status` (VARCHAR): Schema status (created, updated, error)
  - `version` (VARCHAR): Schema version
  - `schema_definition` (JSONB): Complete schema structure
  - `created_at`, `updated_at`, `deleted_at`: Timestamps

**Indexes**: Optimized for fast lookups by name, status, and creation date.

## Response Format
All API responses follow this structure:
```json
{
  "success": true|false,
  "message": "Human readable message",
  "data": {} | [] | null,
  "error": {
    "code": "ERROR_CODE",
    "details": "Detailed error information"
  } | null
}
```

## HTTP Status Codes
- `200` - Success
- `201` - Created
- `400` - Bad Request (validation errors)
- `404` - Not Found
- `409` - Conflict (duplicate names, etc.)
- `500` - Internal Server Error

---

## Schema Management Endpoints

### 1. Create Schema
Create a new database schema and generate the actual database.

**Endpoint:** `POST /schemas`

**Request Body:**
```json
{
  "name": "my_blog_schema",
  "description": "Blog database schema",
  "tables": [
    {
      "id": "users_table",
      "name": "users",
      "columns": [
        {
          "id": "user_id",
          "name": "id",
          "dataType": "INT",
          "nullable": false,
          "primaryKey": true,
          "autoIncrement": true
        },
        {
          "id": "user_email",
          "name": "email",
          "dataType": "VARCHAR",
          "length": 255,
          "nullable": false,
          "unique": true
        }
      ],
      "position": {"x": 100, "y": 100}
    }
  ],
  "foreignKeys": []
}
```

**Response (201):**
```json
{
  "success": true,
  "message": "Schema created successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "my_blog_schema",
    "description": "Blog database schema",
    "databaseName": "schema_550e8400_e29b_41d4_a716_446655440000",
    "status": "created",
    "createdAt": "2024-01-01T10:00:00Z",
    "updatedAt": "2024-01-01T10:00:00Z",
    "version": "1.0"
  }
}
```

---

### 2. Get All Schemas
Retrieve all schemas with basic metadata.

**Endpoint:** `GET /schemas`

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10, max: 100)
- `search` (optional): Search by name or description

**Response (200):**
```json
{
  "success": true,
  "message": "Schemas retrieved successfully",
  "data": [
    {
      "id": "b3e24e22-f1a3-4503-a3b1-cbae1d6a76ea",
      "name": "Blog Schema",
      "description": "A simple blog database schema with users, posts, and comments",
      "databaseName": "schema_blog_example",
      "status": "created",
      "tableCount": 3,
      "createdAt": "2025-06-09T10:22:04.057181+07:00",
      "updatedAt": "2025-06-09T10:22:04.057181+07:00",
      "version": "1.0"
    },
    {
      "id": "b144e70e-6705-47b4-8316-45d00ccec9a6",
      "name": "E-commerce Schema", 
      "description": "Basic e-commerce database with products and orders",
      "databaseName": "schema_ecommerce_example",
      "status": "created",
      "tableCount": 2,
      "createdAt": "2025-06-09T10:22:04.06736+07:00",
      "updatedAt": "2025-06-09T10:22:04.06736+07:00",
      "version": "1.0"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 3,
    "totalPages": 1
  }
}
```

---

### 3. Get Schema by ID
Retrieve complete schema definition including all tables, columns, and relationships.

**Endpoint:** `GET /schemas/{id}`

**Response (200):**
```json
{
  "success": true,
  "message": "Schema retrieved successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "my_blog_schema",
    "description": "Blog database schema",
    "databaseName": "schema_550e8400_e29b_41d4_a716_446655440000",
    "status": "created",
    "tables": [
      {
        "id": "users_table",
        "name": "users",
        "columns": [
          {
            "id": "user_id",
            "name": "id",
            "dataType": "INT",
            "nullable": false,
            "primaryKey": true,
            "autoIncrement": true
          }
        ],
        "position": {"x": 100, "y": 100}
      }
    ],
    "foreignKeys": [],
    "createdAt": "2024-01-01T10:00:00Z",
    "updatedAt": "2024-01-01T10:00:00Z",
    "version": "1.0"
  }
}
```

---

### 4. Update Schema
Update an existing schema. This will modify the schema definition and regenerate the database.

**Endpoint:** `PUT /schemas/{id}`

**Request Body:** Same format as Create Schema

**Response (200):**
```json
{
  "success": true,
  "message": "Schema updated successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "my_blog_schema",
    "description": "Updated blog database schema",
    "databaseName": "schema_550e8400_e29b_41d4_a716_446655440000",
    "status": "updated",
    "createdAt": "2024-01-01T10:00:00Z",
    "updatedAt": "2024-01-01T11:00:00Z",
    "version": "1.1"
  }
}
```

---

### 5. Delete Schema
Delete a schema and optionally drop the associated database.

**Endpoint:** `DELETE /schemas/{id}`

**Query Parameters:**
- `dropDatabase` (optional): Boolean, whether to drop the actual database (default: true)

**Response (200):**
```json
{
  "success": true,
  "message": "Schema deleted successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "databaseDropped": true
  }
}
```

---

## Database Management Endpoints

### 6. Get Database Status
Check the status of the generated database.

**Endpoint:** `GET /schemas/{id}/database/status`

**Response (200):**
```json
{
  "success": true,
  "message": "Database status retrieved",
  "data": {
    "schemaId": "550e8400-e29b-41d4-a716-446655440000",
    "databaseName": "schema_550e8400_e29b_41d4_a716_446655440000",
    "status": "healthy",
    "tableCount": 3,
    "lastChecked": "2024-01-01T12:00:00Z",
    "connectionString": "postgres://localhost:5432/schema_550e8400_e29b_41d4_a716_446655440000"
  }
}
```

---

### 7. Regenerate Database
Force regeneration of the database from the schema definition.

**Endpoint:** `POST /schemas/{id}/database/regenerate`

**Response (200):**
```json
{
  "success": true,
  "message": "Database regenerated successfully",
  "data": {
    "schemaId": "550e8400-e29b-41d4-a716-446655440000",
    "databaseName": "schema_550e8400_e29b_41d4_a716_446655440000",
    "status": "regenerated",
    "regeneratedAt": "2024-01-01T12:30:00Z"
  }
}
```

---

## Validation & Utility Endpoints

### 8. Validate Schema
Validate a schema definition without creating it.

**Endpoint:** `POST /schemas/validate`

**Request Body:** Same as Create Schema

**Response (200):**
```json
{
  "success": true,
  "message": "Schema is valid",
  "data": {
    "valid": true,
    "warnings": [
      "Table 'posts' has no primary key defined"
    ],
    "generatedSQL": [
      "CREATE TABLE users (id SERIAL PRIMARY KEY, email VARCHAR(255) UNIQUE NOT NULL);",
      "CREATE TABLE posts (id SERIAL PRIMARY KEY, title VARCHAR(255) NOT NULL);"
    ]
  }
}
```

**Response (400) - Invalid Schema:**
```json
{
  "success": false,
  "message": "Schema validation failed",
  "error": {
    "code": "VALIDATION_ERROR",
    "details": "Primary key is required for all tables"
  },
  "data": {
    "valid": false,
    "errors": [
      {
        "field": "tables[0].columns",
        "message": "Primary key is required",
        "code": "MISSING_PRIMARY_KEY"
      }
    ]
  }
}
```

---

### 9. Export Schema as SQL
Export the schema definition as SQL DDL statements.

**Endpoint:** `GET /schemas/{id}/export/sql`

**Response (200):**
```json
{
  "success": true,
  "message": "SQL export generated",
  "data": {
    "schemaId": "550e8400-e29b-41d4-a716-446655440000",
    "sql": "-- Generated SQL for schema: my_blog_schema\nCREATE TABLE users (\n  id SERIAL PRIMARY KEY,\n  email VARCHAR(255) UNIQUE NOT NULL\n);\n\nCREATE TABLE posts (\n  id SERIAL PRIMARY KEY,\n  user_id INTEGER NOT NULL,\n  title VARCHAR(255) NOT NULL,\n  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE\n);",
    "generatedAt": "2024-01-01T13:00:00Z"
  }
}
```

---

## Health Check

### 10. Health Check
Check API health status.

**Endpoint:** `GET /health`

**Response (200):**
```json
{
  "success": true,
  "message": "Service is healthy",
  "data": {
    "status": "healthy",
    "timestamp": "2024-01-01T13:00:00Z",
    "database": "connected",
    "version": "1.0.0"
  }
}
```

---

## Error Codes

| Error Code | Description |
|------------|-------------|
| `VALIDATION_ERROR` | Schema validation failed |
| `SCHEMA_NOT_FOUND` | Schema with given ID not found |
| `DATABASE_ERROR` | Database operation failed |
| `DUPLICATE_NAME` | Schema name already exists |
| `INVALID_JSON` | Malformed JSON in request body |
| `MISSING_REQUIRED_FIELD` | Required field is missing |
| `UNSUPPORTED_DATA_TYPE` | Data type not supported |
| `FOREIGN_KEY_ERROR` | Foreign key constraint error |
| `DATABASE_CREATION_FAILED` | Failed to create database |
| `INTERNAL_ERROR` | Unexpected server error |

---

## Data Types Supported

| Type | PostgreSQL Mapping | Optional Parameters |
|------|-------------------|-------------------|
| `INT` | INTEGER | - |
| `BIGINT` | BIGINT | - |
| `VARCHAR` | VARCHAR(n) | length (required) |
| `TEXT` | TEXT | - |
| `BOOLEAN` | BOOLEAN | - |
| `TIMESTAMP` | TIMESTAMP | - |
| `DATE` | DATE | - |
| `TIME` | TIME | - |
| `DECIMAL` | DECIMAL(p,s) | precision, scale |
| `FLOAT` | REAL | - |
| `DOUBLE` | DOUBLE PRECISION | - |
| `JSON` | JSONB | - |
| `UUID` | UUID | - |

---

## Rate Limiting
- **Schema Creation/Update**: 10 requests per minute
- **Schema Retrieval**: 100 requests per minute
- **Validation**: 50 requests per minute

## Authentication
Authentication is handled upstream (proxy/gateway). The API expects user information to be available in request headers:
- `X-User-ID`: User identifier
- `X-User-Email`: User email (optional)

## CORS
The API supports CORS for browser-based applications. Preflight requests are handled automatically.
