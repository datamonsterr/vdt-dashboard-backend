# Authentication Example

This document shows how to integrate with the VDT Dashboard Backend API using Clerk authentication.

## Frontend Integration (React)

### 1. Install Clerk React SDK

```bash
npm install @clerk/clerk-react
```

### 2. Setup Clerk Provider

```jsx
// App.jsx
import { ClerkProvider } from '@clerk/clerk-react'

const clerkPubKey = process.env.REACT_APP_CLERK_PUBLISHABLE_KEY

function App() {
  return (
    <ClerkProvider publishableKey={clerkPubKey}>
      <YourAppComponents />
    </ClerkProvider>
  )
}
```

### 3. Get Session Token and Make API Calls

```jsx
// components/SchemaManager.jsx
import { useAuth } from '@clerk/clerk-react'
import { useState, useEffect } from 'react'

function SchemaManager() {
  const { getToken, isSignedIn } = useAuth()
  const [schemas, setSchemas] = useState([])

  const apiCall = async (endpoint, options = {}) => {
    if (!isSignedIn) return null
    
    const token = await getToken()
    
    return fetch(`http://localhost:8080/api/v1${endpoint}`, {
      ...options,
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
        ...options.headers,
      },
    })
  }

  // Get current user info
  const getCurrentUser = async () => {
    const response = await apiCall('/user/me')
    return response.json()
  }

  // List user's schemas
  const listSchemas = async () => {
    const response = await apiCall('/schemas')
    const data = await response.json()
    setSchemas(data.data)
  }

  // Create a new schema
  const createSchema = async (schemaData) => {
    const response = await apiCall('/schemas', {
      method: 'POST',
      body: JSON.stringify(schemaData),
    })
    return response.json()
  }

  useEffect(() => {
    if (isSignedIn) {
      listSchemas()
    }
  }, [isSignedIn])

  return (
    <div>
      {/* Your schema management UI */}
    </div>
  )
}
```

## Backend Environment Setup

### 1. Get Clerk Secret Key

1. Go to [Clerk Dashboard](https://dashboard.clerk.com)
2. Select your application
3. Go to "API Keys" section
4. Copy the "Secret Key" (starts with `sk_test_` or `sk_live_`)

### 2. Environment Configuration

```env
# .env file
CLERK_SECRET_KEY=sk_test_your_secret_key_here
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=vdt_dashboard
PORT=8080
ENVIRONMENT=development
FRONTEND_URL=http://localhost:3000
```

## API Usage Examples

### 1. Get Current User

```bash
curl -X GET http://localhost:8080/api/v1/user/me \
  -H "Authorization: Bearer YOUR_CLERK_TOKEN"
```

Response:
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
    "fullName": "John Doe"
  }
}
```

### 2. Create a Schema

```bash
curl -X POST http://localhost:8080/api/v1/schemas \
  -H "Authorization: Bearer YOUR_CLERK_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Blog",
    "description": "A simple blog schema",
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
  }'
```

### 3. List User's Schemas

```bash
curl -X GET "http://localhost:8080/api/v1/schemas?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_CLERK_TOKEN"
```

## Testing with Postman

1. **Set up Environment Variables**:
   - `base_url`: `http://localhost:8080/api/v1`
   - `clerk_token`: Your Clerk session token

2. **Get Token from Clerk**:
   - Use Clerk's development tools
   - Or create a simple frontend to get tokens

3. **Add Authorization Header**:
   ```
   Authorization: Bearer {{clerk_token}}
   ```

## Error Handling

### Common Authentication Errors

```json
// Missing token
{
  "success": false,
  "message": "Authorization header is required",
  "error": {
    "code": "UNAUTHORIZED",
    "details": "Missing Authorization header"
  }
}

// Invalid token
{
  "success": false,
  "message": "Invalid token",
  "error": {
    "code": "UNAUTHORIZED",
    "details": "Token verification failed"
  }
}

// Accessing another user's schema
{
  "success": false,
  "message": "Schema not found",
  "error": {
    "code": "SCHEMA_NOT_FOUND",
    "details": "No schema found with the given ID"
  }
}
```

## Security Features

- **User Isolation**: Users can only access their own schemas
- **Token Verification**: All tokens are verified against Clerk's servers
- **Automatic User Creation**: Users are automatically created in our database when first authenticated
- **User Data Sync**: User information is kept in sync with Clerk on each request 