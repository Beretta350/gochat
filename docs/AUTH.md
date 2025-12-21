# Authentication Documentation

GoChat uses JWT (JSON Web Tokens) for stateless authentication.

## Overview

```
┌──────────────────────────────────────────────────────────────────┐
│                     Authentication Flow                          │
│                                                                  │
│   Register ──► Login ──► Access Token ──► API/WebSocket          │
│                              │                                   │
│                              ▼                                   │
│                        (expires 15min)                           │
│                              │                                   │
│                              ▼                                   │
│                    Refresh Token ──► New Access Token            │
│                        (expires 7d)                              │
└──────────────────────────────────────────────────────────────────┘
```

## Token Types

| Token | Expiration | Purpose |
|-------|------------|---------|
| **Access Token** | 15 minutes | API requests, WebSocket auth |
| **Refresh Token** | 7 days | Get new access token without login |

## Endpoints

### Register

Create a new user account.

```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "alice@example.com",
  "username": "alice",
  "password": "securepassword123"
}
```

**Response (201 Created):**

```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "alice@example.com",
    "username": "alice",
    "is_active": true,
    "created_at": "2025-01-20T10:30:00Z"
  },
  "tokens": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 900
  }
}
```

**Errors:**

| Status | Message |
|--------|---------|
| 400 | Invalid request body |
| 400 | Email, username and password are required |
| 400 | Password must be at least 8 characters |
| 409 | Email already exists |

---

### Login

Authenticate with email and password.

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "alice@example.com",
  "password": "securepassword123"
}
```

**Response (200 OK):**

```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "alice@example.com",
    "username": "alice",
    "is_active": true,
    "created_at": "2025-01-20T10:30:00Z"
  },
  "tokens": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 900
  }
}
```

**Errors:**

| Status | Message |
|--------|---------|
| 400 | Invalid request body |
| 400 | Email and password are required |
| 401 | Invalid email or password |

---

### Refresh Token

Get a new access token using refresh token.

```http
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Response (200 OK):**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 900
}
```

**Errors:**

| Status | Message |
|--------|---------|
| 400 | Invalid request body |
| 400 | Refresh token is required |
| 401 | Invalid or expired refresh token |
| 401 | User not found |

---

### Get Current User

Get authenticated user information.

```http
GET /api/v1/auth/me
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

**Response (200 OK):**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "alice@example.com",
  "username": "alice"
}
```

**Errors:**

| Status | Message |
|--------|---------|
| 401 | Missing authorization header |
| 401 | Invalid authorization format |
| 401 | Token expired |
| 401 | Invalid token |

---

## WebSocket Authentication

WebSocket connections require a valid JWT access token in the query string.

```
ws://localhost:8080/ws?token=eyJhbGciOiJIUzI1NiIs...
```

**Connection Flow:**

```
1. Client connects with token in query string
2. Server validates JWT
3. Server extracts user_id from token
4. Connection established with user context
5. User subscribes to their Redis Pub/Sub channel
```

**Errors:**

| Status | Message |
|--------|---------|
| 401 | Missing token |
| 401 | Token expired |
| 401 | Invalid token |
| 426 | Upgrade Required (not a WebSocket request) |

---

## JWT Token Structure

### Access Token Claims

```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "alice@example.com",
  "username": "alice",
  "type": "access",
  "exp": 1705750200,
  "iat": 1705749300,
  "nbf": 1705749300
}
```

### Refresh Token Claims

```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "alice@example.com",
  "username": "alice",
  "type": "refresh",
  "exp": 1706354100,
  "iat": 1705749300,
  "nbf": 1705749300
}
```

---

## Using with Postman

### 1. Register/Login

Send POST request to get tokens.

### 2. Set Access Token

In Postman, go to **Authorization** tab:
- Type: `Bearer Token`
- Token: `<paste access_token here>`

### 3. WebSocket Connection

In Postman WebSocket tab:
- URL: `ws://localhost:8080/ws?token=<access_token>`

---

## Security Considerations

### Password Storage

- Passwords are hashed using **bcrypt** with default cost (10 rounds)
- Plain text passwords are never stored

### Token Security

- Tokens are signed with **HS256** (HMAC-SHA256)
- Secret key should be at least 32 characters in production
- Access tokens are short-lived (15 min) to minimize exposure
- Refresh tokens allow re-authentication without password

### Best Practices

1. **Store tokens securely** - Use httpOnly cookies or secure storage
2. **Always use HTTPS** - Tokens can be intercepted over HTTP
3. **Rotate secrets** - Change JWT_SECRET periodically in production
4. **Validate on every request** - Middleware checks token on protected routes

---

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `JWT_SECRET` | (required) | Secret key for signing tokens |
| `JWT_ACCESS_EXPIRY` | `15m` | Access token lifetime |
| `JWT_REFRESH_EXPIRY` | `168h` | Refresh token lifetime (7 days) |

**Example `.env`:**

```env
JWT_SECRET=your-super-secret-key-at-least-32-chars
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h
```
