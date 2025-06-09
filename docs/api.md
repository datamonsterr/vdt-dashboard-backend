# Database Schema Builder API Documentation

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
This API uses **Clerk** for authentication. All protected endpoints require a valid JWT token in the Authorization header.

### Authentication Header
```
Authorization: Bearer <clerk-session-token>
```

### Getting the Token
- **Frontend**: Use Clerk's client-side SDK (React, Vue, etc.) via `useAuth().getToken()`
- **Testing**: Use Clerk's Dashboard or development tools to generate test tokens
- **Development**: See `docs/authentication-example.md` for complete integration examples

### Token Validation
- Tokens are verified against Clerk's servers on each request
- Invalid or expired tokens return `401 Unauthorized`
- User information is automatically synced from Clerk on each authenticated request

### Protected Endpoints
All schema management endpoints require authentication. Users can only access their own schemas.

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
- `401` - Unauthorized (missing or invalid token)
- `403` - Forbidden (insufficient permissions)
- `404` - Not Found
- `409` - Conflict (duplicate names, etc.)
- `500` - Internal Server Error

---

## User Endpoints

### Get Current User
Get information about the currently authenticated user.

**Endpoint:** `GET /user/me`  
**Authentication:** Required

**Response (200):**
```json
{
  "success": true,
  "message": "User retrieved successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "clerkUserId": "user_2abc123def456",
    "email": "user@example.com",
    "firstName": "John",
    "lastName": "Doe",
    "profileImageUrl": "https://images.clerk.dev/...",
    "fullName": "John Doe",
    "createdAt": "2024-01-01T10:00:00Z",
    "updatedAt": "2024-01-01T10:00:00Z"
  }
}
```

---

## Schema Management Endpoints

### 1. Create Schema
Create a new database schema and generate the actual database.

**Endpoint:** `POST /schemas`  
**Authentication:** Required

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
Retrieve all schemas for the authenticated user with basic metadata.

**Endpoint:** `GET /schemas`  
**Authentication:** Required

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
Retrieve complete schema definition including all tables, columns, and relationships. Only returns schemas owned by the authenticated user.

**Endpoint:** `GET /schemas/{id}`  
**Authentication:** Required

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
Update an existing schema owned by the authenticated user. This will modify the schema definition and regenerate the database.

**Endpoint:** `PUT /schemas/{id}`  
**Authentication:** Required

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
Delete a schema owned by the authenticated user and optionally drop the associated database.

**Endpoint:** `DELETE /schemas/{id}`  
**Authentication:** Required

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
Check the status of the generated database for a schema owned by the authenticated user.

**Endpoint:** `GET /schemas/{id}/database/status`  
**Authentication:** Required

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
Force regeneration of the database from the schema definition for a schema owned by the authenticated user.

**Endpoint:** `POST /schemas/{id}/database/regenerate`  
**Authentication:** Required

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
Export the schema definition as SQL DDL statements for a schema owned by the authenticated user.

**Endpoint:** `GET /schemas/{id}/export/sql`  
**Authentication:** Required

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

## Security Notes
- All authentication is handled directly by the API using Clerk JWT verification
- No upstream proxy or gateway authentication is required
- User identity is extracted from the verified JWT token
- Each request requiring authentication must include a valid Clerk session token

## CORS
The API supports CORS for browser-based applications. Preflight requests are handled automatically.
