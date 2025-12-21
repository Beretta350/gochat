# 💬 GoChat

Real-time chat application built with Go, Fiber, Redis Pub/Sub, and WebSocket.

## 🚀 Tech Stack

| Technology | Purpose |
|------------|---------|
| **Go 1.23** | Backend language |
| **Fiber v2** | Web framework |
| **Uber Fx** | Dependency injection |
| **PostgreSQL** | Persistent storage (users, messages) |
| **Redis** | Pub/Sub & Streams for real-time |
| **JWT** | Stateless authentication |
| **WebSocket** | Real-time communication |
| **Docker** | Containerization |

## 📁 Project Structure

```
gochat-backend/
├── cmd/
│   └── main.go                  # Application entrypoint
├── internal/
│   ├── app/
│   │   ├── app.go               # Fiber app with Fx lifecycle
│   │   ├── fx/
│   │   │   └── module.go        # Fx dependency module
│   │   ├── auth/
│   │   │   ├── jwt.go           # JWT token service
│   │   │   └── service.go       # Auth service (register/login)
│   │   ├── chat/
│   │   │   └── service.go       # Chat service with Redis Pub/Sub
│   │   ├── handler/
│   │   │   ├── auth.go          # Auth endpoints
│   │   │   ├── health.go        # Health check handler
│   │   │   └── websocket.go     # WebSocket handler
│   │   ├── middleware/
│   │   │   ├── auth.go          # JWT auth middleware
│   │   │   ├── error_handler.go # Custom error handler
│   │   │   └── middlewares.go   # Fiber middlewares setup
│   │   ├── model/
│   │   │   ├── user.go          # User model
│   │   │   ├── conversation.go  # Conversation model
│   │   │   └── message.go       # Message model
│   │   ├── repository/
│   │   │   ├── user_repository.go         # User persistence
│   │   │   ├── conversation_repository.go # Conversation persistence
│   │   │   └── message_repository.go      # Message persistence
│   │   └── worker/
│   │       └── message_worker.go          # Redis Stream consumer
│   └── config/
│       └── config.go            # Configuration (Fx provider)
├── pkg/
│   ├── envutil/                 # Environment utilities
│   ├── logger/                  # Zap logger wrapper
│   ├── postgres/                # PostgreSQL client (Fx provider)
│   └── redisclient/             # Redis client (Fx provider)
├── database/
│   ├── schema.sql               # Complete database schema
│   └── migrations/              # Versioned SQL migrations
├── docs/
│   ├── AUTH.md                  # Authentication documentation
│   └── DATABASE.md              # Database documentation
├── configs/
│   └── local.env                # Local environment variables
├── docker-compose.yml           # Redis + PostgreSQL containers
├── Makefile                     # Build and dev commands
├── .air.toml                    # Hot reload config
└── .golangci.yml                # Linter config
```

> 📖 **Documentation:**
> - [docs/AUTH.md](docs/AUTH.md) - Authentication & JWT
> - [docs/DATABASE.md](docs/DATABASE.md) - Database schema

## 🛠️ Getting Started

### Prerequisites

- Go 1.23+
- Docker & Docker Compose
- Make

### Setup

```bash
# Clone the repository
git clone https://github.com/Beretta350/gochat.git
cd gochat

# Install development tools (golangci-lint, air, goimports)
make setup

# Start Redis
make docker-up

# Run the server
make run

# Or with hot reload
make dev
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `ENV` | `dev` | Environment (local, dev, prod) |
| `SERVER_PORT` | `8080` | Server port |
| `DATABASE_URL` | | PostgreSQL connection string |
| `REDIS_ADDR` | `localhost:6379` | Redis address |
| `REDIS_PASSWORD` | `` | Redis password |
| `REDIS_DB` | `0` | Redis database |
| `JWT_SECRET` | | Secret key for JWT signing |
| `JWT_ACCESS_EXPIRY` | `15m` | Access token expiration |
| `JWT_REFRESH_EXPIRY` | `168h` | Refresh token expiration (7 days) |

## 📡 API

### Authentication

```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@test.com","username":"alice","password":"12345678"}'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@test.com","password":"12345678"}'

# Refresh Token
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"eyJhbG..."}'

# Get Current User (protected)
curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer eyJhbG..."
```

### Health Check

```bash
curl http://localhost:8080/api/v1/health
```

### Metrics Dashboard

```bash
# Open in browser
http://localhost:8080/metrics
```

### WebSocket Connection

```bash
# Connect with JWT token
wscat -c "ws://localhost:8080/ws?token=<access_token>"
```

### Message Format

```json
{
  "recipient": "user-uuid",
  "content": "Hello!"
}
```

> 📖 See [docs/AUTH.md](docs/AUTH.md) for complete authentication documentation.

## 🧪 Testing Chat

### 1. Create two users

```bash
# Register Alice
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@test.com","username":"alice","password":"12345678"}'

# Register Bob
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"bob@test.com","username":"bob","password":"12345678"}'
```

### 2. Connect via WebSocket

```bash
# Terminal 1 - Alice (use access_token from register response)
wscat -c "ws://localhost:8080/ws?token=<alice_access_token>"

# Terminal 2 - Bob (use access_token from register response)
wscat -c "ws://localhost:8080/ws?token=<bob_access_token>"
```

### 3. Send messages

```bash
# In Alice's terminal, send to Bob's user ID:
{"recipient": "<bob_user_id>", "content": "Hey Bob!"}

# Bob receives the message! ✅
```

## 🔧 Development

```bash
make help           # Show all commands
make run            # Run server
make dev            # Run with hot reload
make build          # Build binary
make test           # Run tests
make lint           # Run linter
make lint-fix       # Run linter with auto-fix
make fmt            # Format code
make docker-up      # Start Redis
make docker-down    # Stop Redis
make docker-logs    # View Redis logs
make all            # fmt + lint + test + build
```

## 🏗️ Architecture

### Dependency Injection (Uber Fx)

The application uses **Uber Fx** for dependency injection, providing:
- Automatic dependency resolution
- Clean lifecycle management (start/stop hooks)
- Testability through constructor injection

```
Config → RedisClient → ChatService → WebSocketHandler
                    ↘               ↗
              MessageRepository → MessageWorker
```

### System Overview

```
┌──────────────┐                      ┌──────────────────────────────────┐
│   Clients    │                      │            Server                │
│              │     WebSocket        │                                  │
│  ┌────────┐  │◄────────────────────►│  ┌────────────────────────────┐ │
│  │ Web/App│  │                      │  │      Fiber + Handlers      │ │
│  └────────┘  │     REST API         │  └─────────────┬──────────────┘ │
│  ┌────────┐  │◄────────────────────►│                │                │
│  │ Mobile │  │                      │                ▼                │
│  └────────┘  │                      │  ┌────────────────────────────┐ │
└──────────────┘                      │  │      Chat Service          │ │
                                      │  └─────────────┬──────────────┘ │
                                      │                │                │
                                      │    ┌───────────┼───────────┐    │
                                      │    ▼           ▼           ▼    │
                                      │ ┌──────┐  ┌────────┐ ┌───────┐  │
                                      │ │Pub/Sub│  │ Stream │ │Worker │  │
                                      │ └──┬───┘  └───┬────┘ └───┬───┘  │
                                      └────┼──────────┼──────────┼──────┘
                                           │          │          │
                                           ▼          │          ▼
                                      ┌─────────┐     │    ┌───────────┐
                                      │  Redis  │◄────┘    │ PostgreSQL│
                                      │(realtime)│         │ (persist) │
                                      └─────────┘          └───────────┘
```

### Message Flow

```
Alice sends message to Bob
           │
           ▼
┌─────────────────────────────┐
│  1. Save to PostgreSQL      │ ──► Persistent storage
└──────────────┬──────────────┘
               │
               ▼
┌─────────────────────────────┐
│  2. Add to Redis Stream     │ ──► For async processing (optional)
└──────────────┬──────────────┘
               │
               ▼
┌─────────────────────────────┐
│  3. Publish to Pub/Sub      │ ──► channel: user:{bob_id}
└──────────────┬──────────────┘
               │
       ┌───────┴───────┐
       ▼               ▼
   Bob's PC       Bob's Phone
   (online)        (online)
       │               │
    RECEIVES        RECEIVES
    via WS          via WS

If Bob is offline → He fetches history from PostgreSQL when reconnects
```

### Authentication Flow (JWT)

```
┌──────────┐                              ┌──────────┐
│  Client  │                              │  Server  │
└────┬─────┘                              └────┬─────┘
     │                                         │
     │  POST /api/v1/auth/register             │
     │  { email, username, password }          │
     │────────────────────────────────────────►│
     │                                         │ bcrypt hash
     │                                         │ save to PostgreSQL
     │  { user, tokens }                       │
     │◄────────────────────────────────────────│
     │                                         │
     │  POST /api/v1/auth/login                │
     │  { email, password }                    │
     │────────────────────────────────────────►│
     │                                         │ verify password
     │  { user, tokens }                       │ generate JWT
     │◄────────────────────────────────────────│
     │                                         │
     │  access_token expires...                │
     │                                         │
     │  POST /api/v1/auth/refresh              │
     │  { refresh_token }                      │
     │────────────────────────────────────────►│
     │                                         │
     │  { new tokens }                         │
     │◄────────────────────────────────────────│
     │                                         │
     │  WS /ws?token={access_token}            │
     │────────────────────────────────────────►│
     │                                         │ validate JWT
     │  Connection established                 │ extract user_id
     │◄═══════════════════════════════════════►│
     │                                         │
```

> 📖 See [docs/AUTH.md](docs/AUTH.md) for complete authentication documentation.

## 📝 TODO

- [x] WebSocket real-time messaging
- [x] Redis Pub/Sub for multi-device support
- [x] Redis Streams for async processing
- [x] Uber Fx dependency injection
- [x] Database schema design
- [x] PostgreSQL integration
- [x] JWT Authentication (register, login, refresh)
- [x] User management (CRUD)
- [ ] Conversation management (create, list)
- [ ] Message history with cursor pagination
- [ ] Group chats
- [ ] Typing indicators
- [ ] Read receipts
- [ ] File sharing

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.
