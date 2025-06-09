# VDT Dashboard Backend

A powerful backend API for a no-code database schema builder that enables visual database design through drag-and-drop interface with automatic SQL database generation.

## Features

- 🎨 **Visual Schema Design**: Accept JSON schema definitions from frontend
- 🗄️ **Automatic Database Creation**: Generate PostgreSQL databases from schema definitions
- ✅ **Schema Validation**: Comprehensive validation of schema structure and constraints
- 🔄 **CRUD Operations**: Complete schema management with versioning
- 📊 **SQL Export**: Export schemas as SQL DDL statements
- 🔍 **Health Monitoring**: Database status checking and health endpoints

## Tech Stack

- **Language**: Go 1.24
- **Web Framework**: Gin-Gonic
- **Database**: PostgreSQL (main app) + Dynamic PostgreSQL databases (generated)
- **ORM**: GORM
- **Development**: Air (hot reloading)

## 🚀 Quick Start

### Prerequisites
- Go 1.24+
- PostgreSQL 12+
- Make (for build commands)

### Environment Setup

1. **Clone the repository**
```bash
git clone <repository-url>
cd vdt-dashboard-backend
```

2. **Install dependencies**
```bash
make deps
```

3. **Create environment file**
```bash
cp .env.example .env
# Edit .env with your database credentials
```

Required environment variables:
```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=vdt_dashboard

# Server Configuration  
PORT=8080
ENVIRONMENT=development

# Frontend Configuration
FRONTEND_URL=http://localhost:3000

# Clerk Authentication (Required)
CLERK_SECRET_KEY=sk_test_your_clerk_secret_key_here
```

### Authentication Setup

This application uses **Clerk** for authentication. You need to:

1. **Create a Clerk account** at [clerk.com](https://clerk.com)
2. **Create a new application** in your Clerk dashboard
3. **Get your Secret Key** from the API Keys section
4. **Add the secret key** to your `.env` file as `CLERK_SECRET_KEY`

**Frontend Integration:**
- Install Clerk's React SDK in your frontend
- Wrap your app with `<ClerkProvider>`
- Use `useAuth()` to get session tokens
- Send tokens in `Authorization: Bearer <token>` header

**Testing:**
- Use Clerk's development tools to generate test tokens
- Or create test users through Clerk's dashboard

### Database Setup

4. **Create and setup database**
```bash
# Complete database setup (recommended)
make db-setup

# Or step by step:
make db-create      # Create database
make migrate        # Run migrations
make migrate-seed   # Add sample data
```

5. **Start development server**
```bash
make dev
```

The API will be available at `http://localhost:8080/api/v1`

## 📚 Documentation

- [API Documentation](docs/api.md) - Complete API reference with examples
- [Engineering Documentation](docs/engineer.md) - System architecture and design
- [Product Requirements](docs/prd.md) - Features and specifications
- [Development Tasks](docs/todos.md) - Current progress and roadmap

## 🛠️ Available Commands

### Development
```bash
make dev            # Start with hot reloading
make build          # Build application
make run            # Build and run
make clean          # Clean build artifacts
```

### Database Operations
```bash
make db-create      # Create database
make db-drop        # Drop database (destructive)
make db-reset       # Drop and recreate database
make migrate        # Run migrations
make migrate-reset  # Reset with fresh migrations  
make migrate-seed   # Seed sample data
make db-setup       # Complete setup (create + migrate + seed)
```

### Code Quality
```bash
make test           # Run tests
make test-coverage  # Run tests with coverage
make fmt            # Format code
make lint           # Run linter
```

### Utilities
```bash
make help           # Show all available commands
make check-env      # Verify environment setup
```

## 🏗️ Architecture

The application follows Clean Architecture principles with clear separation of concerns:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   API Handlers  │───▶│    Services     │───▶│  Repositories   │
│                 │    │                 │    │                 │
│ • Schema CRUD   │    │ • Business Logic│    │ • Data Access   │
│ • Validation    │    │ • SQL Generation│    │ • DB Management │
│ • Export/Import │    │ • DB Operations │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 📋 API Endpoints

### Schema Management
- `POST /api/v1/schemas` - Create new schema
- `GET /api/v1/schemas` - List all schemas  
- `GET /api/v1/schemas/{id}` - Get schema details
- `PUT /api/v1/schemas/{id}` - Update schema
- `DELETE /api/v1/schemas/{id}` - Delete schema

### Database Operations
- `GET /api/v1/schemas/{id}/database/status` - Check database status
- `POST /api/v1/schemas/{id}/database/regenerate` - Regenerate database

### Utilities
- `POST /api/v1/schemas/validate` - Validate schema
- `GET /api/v1/schemas/{id}/export/sql` - Export as SQL
- `GET /api/v1/health` - Health check

## 🗄️ Database Schema

The application uses PostgreSQL with the following main table:

**schemas** table:
- `id` (UUID) - Unique identifier
- `name` (VARCHAR) - Schema name
- `description` (TEXT) - Optional description  
- `database_name` (VARCHAR) - Generated database name
- `status` (VARCHAR) - Status (created, updated, error)
- `schema_definition` (JSONB) - Complete schema structure
- `created_at`, `updated_at`, `deleted_at` - Timestamps

## 🧪 Sample Data

After running `make migrate-seed`, you'll have sample schemas available:
- **Blog Schema** - Users, posts, and comments with relationships
- **E-commerce Schema** - Products and orders structure

## 🔧 Development Tools

- **Air** - Hot reloading for development
- **GORM** - Go ORM for database operations
- **Gin** - High-performance HTTP web framework
- **PostgreSQL** - Primary database
- **UUID** - Unique identifiers
- **JSONB** - Flexible schema storage

## 🚀 Production Deployment

```bash
# Build for production
make build-prod

# Or build for all platforms
make build-all
```

## 📝 Contributing

1. Read the [engineering documentation](docs/engineer.md)
2. Check [current tasks](docs/todos.md)  
3. Follow the established architecture patterns
4. Add tests for new features
5. Update documentation as needed

## 📄 License

[Add your license here]

## Support

For questions and support, please refer to the documentation in the `docs/` directory or create an issue in the repository. 