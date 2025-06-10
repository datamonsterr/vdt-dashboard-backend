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

## How to Get Clerk Token for Testing

### Method 1: Using Clerk Dashboard (Quick Test)

1. **Go to Clerk Dashboard**:
   - Visit [https://dashboard.clerk.com](https://dashboard.clerk.com)
   - Select your application

2. **Navigate to Users**:
   - Go to "Users" section in the sidebar
   - Find or create a test user

3. **Generate Session Token**:
   - Click on a user
   - Go to "Sessions" tab
   - Click "Generate JWT" or "Copy token"
   - Use this token in your API requests

⚠️ **Note**: Dashboard tokens expire quickly (usually 1 hour)

### Method 2: Using Browser Console (Development)

1. **Set up a minimal frontend** (if you don't have one):
   ```html
   <!DOCTYPE html>
   <html>
   <head>
       <script src="https://unpkg.com/@clerk/clerk-js@latest/dist/clerk.browser.js"></script>
   </head>
   <body>
       <div id="sign-in"></div>
       <button onclick="getToken()" id="getTokenBtn" style="display:none;">Get Token</button>
       <div id="token"></div>
       
       <script>
           window.addEventListener('load', async () => {
               const clerk = window.Clerk
               await clerk.load({
                   publishableKey: 'YOUR_CLERK_PUBLISHABLE_KEY'
               })
               
               if (clerk.user) {
                   document.getElementById('getTokenBtn').style.display = 'block'
               } else {
                   clerk.mountSignIn(document.getElementById('sign-in'))
               }
           })
           
           async function getToken() {
               const token = await window.Clerk.session?.getToken()
               document.getElementById('token').innerHTML = `
                   <h3>Your Token:</h3>
                   <textarea rows="10" cols="80">${token}</textarea>
                   <p><strong>Copy this token and use it in your API requests!</strong></p>
               `
           }
       </script>
   </body>
   </html>
   ```

2. **Open in browser and sign in**
3. **Click "Get Token" button**
4. **Copy the displayed token**

### Method 3: Using React App (Production-like)

Create a simple token extractor component:

```jsx
import { useAuth, useUser } from '@clerk/clerk-react'
import { useState } from 'react'

function TokenExtractor() {
  const { getToken } = useAuth()
  const { user } = useUser()
  const [token, setToken] = useState('')

  const handleGetToken = async () => {
    const token = await getToken()
    setToken(token)
    console.log('Token:', token)
  }

  if (!user) {
    return <div>Please sign in first</div>
  }

  return (
    <div style={{ padding: '20px' }}>
      <h2>Get Clerk Token</h2>
      <p>User: {user.emailAddresses[0]?.emailAddress}</p>
      <button onClick={handleGetToken}>Get Current Token</button>
      {token && (
        <div style={{ marginTop: '20px' }}>
          <h3>Your Token:</h3>
          <textarea
            value={token}
            readOnly
            rows={10}
            cols={80}
            style={{ fontFamily: 'monospace', fontSize: '12px' }}
          />
          <p><strong>Copy this token for API testing!</strong></p>
        </div>
      )}
    </div>
  )
}
```

### Method 4: Using cURL to Get Token (Advanced)

You can also get tokens programmatically using Clerk's API:

```bash
# Get user sessions (requires secret key)
curl -X GET "https://api.clerk.com/v1/users/USER_ID/sessions" \
  -H "Authorization: Bearer sk_test_YOUR_SECRET_KEY"
  
# Extract session_id and create token
curl -X POST "https://api.clerk.com/v1/sessions/SESSION_ID/tokens" \
  -H "Authorization: Bearer sk_test_YOUR_SECRET_KEY"
```

## Testing with Postman

### Setup Environment Variables

1. **Create New Environment** in Postman
2. **Add these variables**:
   ```
   base_url: http://localhost:8080/api/v1
   clerk_token: YOUR_ACTUAL_TOKEN_HERE
   ```

### Get Your Token

1. **Use any method above** to get a valid Clerk token
2. **Copy the token** (starts with `eyJ...`)
3. **Update the `clerk_token` variable** in Postman

### Make API Requests

1. **Set Authorization Header**:
   ```
   Authorization: Bearer {{clerk_token}}
   ```

2. **Test the API**:
   - GET `{{base_url}}/user/me`
   - GET `{{base_url}}/schemas`
   - POST `{{base_url}}/schemas`

### Example Postman Request

```
Method: GET
URL: {{base_url}}/user/me
Headers:
  Authorization: Bearer {{clerk_token}}
  Content-Type: application/json
```

## Token Troubleshooting

### Common Issues

1. **Token Expired**: Get a new token (they typically last 1 hour)
2. **Invalid Token**: Make sure you copied the full token
3. **Wrong Environment**: Ensure you're using test tokens with test secret key
4. **CORS Issues**: Make sure your frontend URL is configured in Clerk dashboard

### Validating Your Token

You can verify your token is working:

```bash
# Test with a simple API call
curl -X GET http://localhost:8080/api/v1/user/me \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -v
```

**Expected Success Response (200)**:
```json
{
  "success": true,
  "message": "User retrieved successfully",
  "data": {
    "id": "...",
    "email": "user@example.com"
  }
}
```

**Expected Error Response (401)**:
```json
{
  "success": false,
  "message": "Invalid token",
  "error": {
    "code": "UNAUTHORIZED"
  }
}
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