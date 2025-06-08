# Database Schema Builder - Backend API
## Product Requirements Document (PRD)

### Project Overview
A backend API service for a no-code database schema builder tool that allows users to visually design database schemas through drag-and-drop interface and automatically generate SQL databases.

### Product Vision
Enable users to design complex database schemas visually without SQL knowledge, with automatic database creation and management capabilities.

### Target Users
- **Primary**: Developers and database designers who want to quickly prototype database schemas
- **Secondary**: Non-technical users who need to design simple database structures
- **Tertiary**: Teams collaborating on database design

### Core Features

#### 1. Schema Management
- **Create Schema**: Accept JSON schema definitions from frontend and store them
- **Update Schema**: Modify existing schemas with versioning support
- **Delete Schema**: Remove schemas and optionally drop associated databases
- **List Schemas**: Retrieve all schemas for a user with metadata
- **Get Schema Details**: Fetch complete schema definition by ID

#### 2. Database Generation
- **Auto-Create Database**: Generate PostgreSQL databases from schema definitions
- **Schema Validation**: Validate schema structure before database creation
- **SQL Generation**: Convert JSON schema to PostgreSQL DDL statements
- **Database Naming**: Auto-generate unique database names with user prefixes

#### 3. Schema Validation & Constraints
- **Data Type Validation**: Support common SQL data types (INT, VARCHAR, TEXT, TIMESTAMP, BOOLEAN, etc.)
- **Primary Key Validation**: Ensure each table has appropriate primary keys
- **Foreign Key Validation**: Validate referential integrity constraints
- **Column Constraints**: Support nullable, auto-increment, default values
- **Table Relationships**: Handle one-to-many, many-to-many relationships

#### 4. Metadata Management
- **Schema Versioning**: Track schema versions and changes
- **Creation Timestamps**: Record when schemas are created/modified
- **User Association**: Link schemas to user accounts
- **Database Status**: Track database creation status and health

### Technical Requirements

#### Performance Requirements
- **Response Time**: API responses < 500ms for CRUD operations
- **Database Creation**: Schema-to-database conversion < 5 seconds
- **Concurrent Users**: Support 100+ concurrent schema operations
- **Storage**: Efficient JSON storage and retrieval

#### Security Requirements
- **Input Validation**: Sanitize all JSON inputs to prevent injection attacks
- **SQL Injection Protection**: Use parameterized queries and ORM
- **Database Isolation**: Each generated database should be isolated
- **Error Handling**: Prevent information leakage in error messages

#### Scalability Requirements
- **Horizontal Scaling**: Design for multiple server instances
- **Database Connections**: Efficient connection pooling
- **Memory Management**: Optimize for large schema processing
- **Caching**: Cache frequently accessed schemas

### Data Structure Specification

#### Schema JSON Format
```json
{
  "id": "schema_uuid",
  "name": "my_database_schema",
  "description": "Description of the schema",
  "tables": [
    {
      "id": "table_uuid",
      "name": "table_name",
      "columns": [
        {
          "id": "column_uuid",
          "name": "column_name",
          "dataType": "VARCHAR|INT|TEXT|TIMESTAMP|BOOLEAN|DECIMAL|DATE",
          "length": 255, // optional, for VARCHAR
          "precision": 10, // optional, for DECIMAL
          "scale": 2, // optional, for DECIMAL
          "nullable": true|false,
          "primaryKey": true|false,
          "autoIncrement": true|false,
          "defaultValue": "any", // optional
          "unique": true|false // optional
        }
      ],
      "position": {"x": 100, "y": 200}, // for UI positioning
      "indexes": [ // optional
        {
          "name": "index_name",
          "columns": ["column1", "column2"],
          "unique": true|false
        }
      ]
    }
  ],
  "foreignKeys": [
    {
      "id": "fk_uuid",
      "name": "fk_name", // optional
      "sourceTableId": "source_table_uuid",
      "sourceColumnId": "source_column_uuid",
      "targetTableId": "target_table_uuid",
      "targetColumnId": "target_column_uuid",
      "onDelete": "CASCADE|RESTRICT|SET NULL|NO ACTION",
      "onUpdate": "CASCADE|RESTRICT|SET NULL|NO ACTION"
    }
  ],
  "version": "1.0",
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

### Success Metrics
- **Schema Creation Success Rate**: > 99%
- **Database Generation Success Rate**: > 95%
- **API Uptime**: > 99.9%
- **Average Response Time**: < 300ms
- **User Schema Retention**: > 80% after 30 days

### Constraints & Assumptions
- **Database Engine**: PostgreSQL only (initially)
- **Single Server Deployment**: Same server for app and generated databases
- **Authentication**: Assumed to be handled by frontend/proxy
- **File Storage**: No file/blob storage requirements
- **Real-time Features**: No WebSocket/real-time requirements initially

### Future Enhancements (Out of Scope for V1)
- Multi-database engine support (MySQL, SQLite)
- Schema migration tools
- Database backup/restore
- Advanced indexing strategies
- Performance monitoring
- Schema sharing/collaboration
- Import from existing databases
- GraphQL API generation
