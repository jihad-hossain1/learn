# HTTP Deep Dive and APIs for Full Stack Development

## Table of Contents
1. [HTTP Protocol Deep Dive](#http-protocol-deep-dive)
2. [HTTP Headers](#http-headers)
3. [RESTful API Design](#restful-api-design)
4. [API Authentication](#api-authentication)
5. [CORS (Cross-Origin Resource Sharing)](#cors-cross-origin-resource-sharing)
6. [WebSockets](#websockets)
7. [GraphQL vs REST](#graphql-vs-rest)
8. [API Versioning](#api-versioning)

## HTTP Protocol Deep Dive

### HTTP Request Structure
```
GET /api/users/123 HTTP/1.1
Host: api.example.com
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64)
Accept: application/json
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{"name": "John Doe"}
```

**Components:**
1. **Request Line**: Method + URL + HTTP Version
2. **Headers**: Metadata about the request
3. **Body**: Data payload (optional)

### HTTP Response Structure
```
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: 85
Cache-Control: max-age=3600
Set-Cookie: sessionId=abc123; HttpOnly; Secure

{"id": 123, "name": "John Doe", "email": "john@example.com"}
```

**Components:**
1. **Status Line**: HTTP Version + Status Code + Reason Phrase
2. **Headers**: Metadata about the response
3. **Body**: Response data

### HTTP/1.1 vs HTTP/2 vs HTTP/3

| Feature | HTTP/1.1 | HTTP/2 | HTTP/3 |
|---------|----------|--------|---------|
| **Multiplexing** | No | Yes | Yes |
| **Header Compression** | No | HPACK | QPACK |
| **Server Push** | No | Yes | Yes |
| **Protocol** | TCP | TCP | UDP (QUIC) |
| **Binary Protocol** | No | Yes | Yes |
| **Connection Reuse** | Limited | Full | Full |

## HTTP Headers

### Request Headers

#### Authentication Headers
```http
Authorization: Bearer <token>
Authorization: Basic <base64-encoded-credentials>
Authorization: Digest <digest-credentials>
```

#### Content Headers
```http
Content-Type: application/json
Content-Length: 1234
Content-Encoding: gzip
Content-Language: en-US
```

#### Cache Headers
```http
Cache-Control: no-cache
If-None-Match: "etag-value"
If-Modified-Since: Wed, 21 Oct 2015 07:28:00 GMT
```

#### CORS Headers
```http
Origin: https://example.com
Access-Control-Request-Method: POST
Access-Control-Request-Headers: Content-Type
```

### Response Headers

#### Security Headers
```http
Strict-Transport-Security: max-age=31536000; includeSubDomains
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Content-Security-Policy: default-src 'self'
```

#### Cache Control Headers
```http
Cache-Control: public, max-age=3600
ETag: "33a64df551425fcc55e4d42a148795d9f25f89d4"
Expires: Wed, 21 Oct 2015 07:28:00 GMT
Last-Modified: Wed, 21 Oct 2015 07:28:00 GMT
```

#### CORS Response Headers
```http
Access-Control-Allow-Origin: https://example.com
Access-Control-Allow-Methods: GET, POST, PUT, DELETE
Access-Control-Allow-Headers: Content-Type, Authorization
Access-Control-Max-Age: 86400
```

## RESTful API Design

### REST Principles

1. **Stateless**: Each request contains all necessary information
2. **Client-Server**: Separation of concerns
3. **Cacheable**: Responses should be cacheable when appropriate
4. **Uniform Interface**: Consistent API design
5. **Layered System**: Architecture can have multiple layers
6. **Code on Demand** (optional): Server can send executable code

### RESTful URL Design

#### Good Examples
```
GET    /api/users              # Get all users
GET    /api/users/123          # Get user with ID 123
POST   /api/users              # Create a new user
PUT    /api/users/123          # Update user 123 (full update)
PATCH  /api/users/123          # Update user 123 (partial update)
DELETE /api/users/123          # Delete user 123

# Nested resources
GET    /api/users/123/posts    # Get posts by user 123
POST   /api/users/123/posts    # Create post for user 123
GET    /api/users/123/posts/456 # Get specific post
```

#### Bad Examples
```
GET    /api/getUsers           # Don't use verbs in URLs
POST   /api/user/delete/123    # Use DELETE method instead
GET    /api/users/123/delete   # Wrong method for deletion
```

### HTTP Status Code Usage

#### Success Responses
```javascript
// GET request - return data
res.status(200).json({ users: [...] });

// POST request - resource created
res.status(201).json({ id: 123, message: 'User created' });

// PUT/PATCH request - resource updated
res.status(200).json({ message: 'User updated' });

// DELETE request - resource deleted
res.status(204).send(); // No content
```

#### Error Responses
```javascript
// Bad request - validation errors
res.status(400).json({
  error: 'Validation failed',
  details: {
    email: 'Invalid email format',
    password: 'Password too short'
  }
});

// Unauthorized - authentication required
res.status(401).json({
  error: 'Authentication required',
  message: 'Please provide valid credentials'
});

// Forbidden - insufficient permissions
res.status(403).json({
  error: 'Insufficient permissions',
  message: 'Admin access required'
});

// Not found - resource doesn't exist
res.status(404).json({
  error: 'Resource not found',
  message: 'User with ID 123 not found'
});

// Server error - internal issues
res.status(500).json({
  error: 'Internal server error',
  message: 'Something went wrong'
});
```

### API Response Formats

#### Consistent Response Structure
```javascript
// Success response
{
  "success": true,
  "data": {
    "id": 123,
    "name": "John Doe",
    "email": "john@example.com"
  },
  "message": "User retrieved successfully",
  "timestamp": "2023-12-07T10:30:00Z"
}

// Error response
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input data",
    "details": {
      "email": "Email is required",
      "password": "Password must be at least 8 characters"
    }
  },
  "timestamp": "2023-12-07T10:30:00Z"
}
```

#### Pagination
```javascript
{
  "success": true,
  "data": {
    "users": [...],
    "pagination": {
      "page": 2,
      "limit": 10,
      "total": 150,
      "totalPages": 15,
      "hasNext": true,
      "hasPrev": true
    }
  }
}
```

## API Authentication

### 1. API Keys
```javascript
// Header-based
Authorization: Bearer your-api-key-here

// Query parameter (less secure)
GET /api/users?api_key=your-api-key-here
```

### 2. JWT (JSON Web Tokens)
```javascript
// Token structure: header.payload.signature
const token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c'

// Usage
Authorization: Bearer ${token}
```

**JWT Payload Example:**
```javascript
{
  "sub": "1234567890",    // Subject (user ID)
  "name": "John Doe",     // User name
  "iat": 1516239022,      // Issued at
  "exp": 1516242622,      // Expiration
  "role": "admin"         // User role
}
```

### 3. OAuth 2.0
```javascript
// Authorization Code Flow
// Step 1: Redirect to authorization server
https://auth.example.com/oauth/authorize?
  response_type=code&
  client_id=your-client-id&
  redirect_uri=https://yourapp.com/callback&
  scope=read write&
  state=random-string

// Step 2: Exchange code for token
POST /oauth/token
Content-Type: application/x-www-form-urlencoded

grant_type=authorization_code&
code=received-auth-code&
client_id=your-client-id&
client_secret=your-client-secret&
redirect_uri=https://yourapp.com/callback
```

### 4. Session-Based Authentication
```javascript
// Login endpoint
POST /api/auth/login
{
  "email": "user@example.com",
  "password": "password123"
}

// Response sets session cookie
Set-Cookie: sessionId=abc123; HttpOnly; Secure; SameSite=Strict

// Subsequent requests include cookie automatically
Cookie: sessionId=abc123
```

## CORS (Cross-Origin Resource Sharing)

### What is CORS?
CORS is a security feature that allows or restricts web pages to access resources from other domains.

### CORS Headers

#### Simple Requests
```javascript
// Frontend request
fetch('https://api.example.com/users', {
  method: 'GET',
  headers: {
    'Content-Type': 'application/json'
  }
});

// Server response
Access-Control-Allow-Origin: https://myapp.com
Access-Control-Allow-Credentials: true
```

#### Preflight Requests
```javascript
// Frontend request (triggers preflight)
fetch('https://api.example.com/users', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer token'
  },
  body: JSON.stringify({ name: 'John' })
});

// Preflight request (automatic)
OPTIONS /users HTTP/1.1
Origin: https://myapp.com
Access-Control-Request-Method: POST
Access-Control-Request-Headers: Content-Type, Authorization

// Preflight response
Access-Control-Allow-Origin: https://myapp.com
Access-Control-Allow-Methods: GET, POST, PUT, DELETE
Access-Control-Allow-Headers: Content-Type, Authorization
Access-Control-Max-Age: 86400
```

### CORS Configuration Examples

#### Express.js
```javascript
const cors = require('cors');

// Allow all origins (development only)
app.use(cors());

// Production configuration
app.use(cors({
  origin: ['https://myapp.com', 'https://admin.myapp.com'],
  credentials: true,
  methods: ['GET', 'POST', 'PUT', 'DELETE'],
  allowedHeaders: ['Content-Type', 'Authorization']
}));
```

## WebSockets

### When to Use WebSockets
- Real-time chat applications
- Live data feeds (stock prices, sports scores)
- Collaborative editing
- Gaming applications
- Live notifications

### WebSocket vs HTTP

| Feature | HTTP | WebSocket |
|---------|------|----------|
| **Connection** | Request-response | Persistent |
| **Communication** | Half-duplex | Full-duplex |
| **Overhead** | High (headers) | Low |
| **Real-time** | No | Yes |
| **Caching** | Yes | No |

### WebSocket Implementation

#### Client-side (JavaScript)
```javascript
// Establish connection
const ws = new WebSocket('wss://api.example.com/chat');

// Connection opened
ws.onopen = function(event) {
  console.log('Connected to WebSocket');
  ws.send(JSON.stringify({ type: 'join', room: 'general' }));
};

// Receive messages
ws.onmessage = function(event) {
  const data = JSON.parse(event.data);
  console.log('Received:', data);
};

// Handle errors
ws.onerror = function(error) {
  console.error('WebSocket error:', error);
};

// Connection closed
ws.onclose = function(event) {
  console.log('WebSocket connection closed');
};

// Send message
function sendMessage(message) {
  if (ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify({
      type: 'message',
      content: message,
      timestamp: new Date().toISOString()
    }));
  }
}
```

#### Server-side (Node.js with ws library)
```javascript
const WebSocket = require('ws');
const wss = new WebSocket.Server({ port: 8080 });

wss.on('connection', function connection(ws, req) {
  console.log('New client connected');

  // Send welcome message
  ws.send(JSON.stringify({
    type: 'welcome',
    message: 'Connected to chat server'
  }));

  // Handle incoming messages
  ws.on('message', function incoming(data) {
    try {
      const message = JSON.parse(data);
      
      // Broadcast to all clients
      wss.clients.forEach(function each(client) {
        if (client !== ws && client.readyState === WebSocket.OPEN) {
          client.send(JSON.stringify(message));
        }
      });
    } catch (error) {
      console.error('Invalid message format:', error);
    }
  });

  // Handle disconnection
  ws.on('close', function close() {
    console.log('Client disconnected');
  });
});
```

## GraphQL vs REST

### REST API Example
```javascript
// Multiple requests needed
GET /api/users/123           // Get user
GET /api/users/123/posts     // Get user's posts
GET /api/users/123/followers // Get user's followers
```

### GraphQL Example
```javascript
// Single request
POST /graphql
{
  "query": """
    query GetUser($id: ID!) {
      user(id: $id) {
        id
        name
        email
        posts {
          id
          title
          content
        }
        followers {
          id
          name
        }
      }
    }
  """,
  "variables": {
    "id": "123"
  }
}
```

### Comparison

| Aspect | REST | GraphQL |
|--------|------|----------|
| **Requests** | Multiple endpoints | Single endpoint |
| **Over-fetching** | Common | Eliminated |
| **Under-fetching** | Common | Eliminated |
| **Caching** | Easy (HTTP) | Complex |
| **Learning Curve** | Low | High |
| **Tooling** | Mature | Growing |

## API Versioning

### 1. URL Versioning
```
GET /api/v1/users
GET /api/v2/users
```

### 2. Header Versioning
```
GET /api/users
API-Version: v2
```

### 3. Query Parameter Versioning
```
GET /api/users?version=2
```

### 4. Content Negotiation
```
GET /api/users
Accept: application/vnd.api.v2+json
```

### Best Practices
1. **Semantic Versioning**: Use major.minor.patch format
2. **Backward Compatibility**: Support old versions for reasonable time
3. **Deprecation Warnings**: Inform clients about upcoming changes
4. **Documentation**: Clearly document version differences

## Key Takeaways

1. **HTTP is foundational**: Master headers, methods, and status codes
2. **REST principles**: Follow consistent patterns for API design
3. **Authentication matters**: Choose appropriate auth strategy
4. **CORS is crucial**: Understand cross-origin policies
5. **WebSockets for real-time**: Use when HTTP isn't sufficient
6. **Version your APIs**: Plan for evolution and changes

## Next Steps

In the next part, we'll cover:
- Network security fundamentals
- SSL/TLS and HTTPS implementation
- API security best practices
- Common vulnerabilities and protection

---

*Continue to [Part 3: Network Security](03-network-security.md)*