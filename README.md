# 💬 GoChat Backend

Real-time chat application built with Go, Fiber, Redis Pub/Sub, and WebSocket.

## 🚀 Tech Stack

| Technology | Purpose |
|------------|---------|
| **Go 1.23** | Backend language |
| **Fiber v2** | Web framework |
| **Uber Fx** | Dependency injection |
| **PostgreSQL** | Persistent storage (users, conversations, messages) |
| **Redis** | Pub/Sub (real-time) & Streams (async processing) |
| **JWT** | Stateless authentication |
| **WebSocket** | Real-time bidirectional communication |
| **Docker** | Containerization |
| **Air** | Hot reload for development |

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
│   │   │   ├── conversation.go  # Conversation endpoints
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
│   │       └── message_worker.go          # Redis Stream → PostgreSQL
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
├── scripts/
│   └── dev/                     # Development scripts (gitignored)
├── configs/
│   └── local.env                # Local environment variables
├── docker-compose.yml           # Redis + PostgreSQL containers
├── Dockerfile                   # Production image
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

### Quick Start

```bash
# Clone the repository
git clone https://github.com/Beretta350/gochat-backend.git
cd gochat-backend

# Start Redis + PostgreSQL
make docker-up

# Run the server (with hot reload)
make dev

# Or without hot reload
make run
```

### Install Development Tools

```bash
# Install Air (hot reload)
go install github.com/air-verse/air@latest

# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
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

## 📡 API Endpoints

### Authentication

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/api/v1/auth/register` | ❌ | Create new user |
| POST | `/api/v1/auth/login` | ❌ | Login and get tokens |
| POST | `/api/v1/auth/refresh` | ❌ | Refresh access token |
| GET | `/api/v1/auth/me` | ✅ | Get current user |

### Conversations

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/api/v1/conversations` | ✅ | Create conversation (direct or group) |
| GET | `/api/v1/conversations` | ✅ | List user's conversations |
| GET | `/api/v1/conversations/:id` | ✅ | Get conversation details |
| GET | `/api/v1/conversations/:id/messages` | ✅ | Get messages (with pagination) |

### WebSocket

| Endpoint | Auth | Description |
|----------|------|-------------|
| `ws://localhost:8080/ws?token=<jwt>` | ✅ | Real-time messaging |

### Other

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/health` | Health check |
| GET | `/metrics` | Metrics dashboard |

## 🔐 Authentication

```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@test.com","username":"alice","password":"12345678"}'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@test.com","password":"12345678"}'

# Response structure
{
  "user": { "id": "...", "email": "...", "username": "..." },
  "tokens": {
    "access_token": "eyJhbG...",
    "refresh_token": "eyJhbG...",
    "expires_in": 900
  }
}
```

## 💬 Conversations

### Create Direct Conversation (1:1)

```bash
curl -X POST http://localhost:8080/api/v1/conversations \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"participant_id": "<other_user_id>"}'
```

### Create Group Conversation

```bash
curl -X POST http://localhost:8080/api/v1/conversations \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Group",
    "participant_ids": ["user_id_1", "user_id_2"]
  }'
```

### List Conversations

```bash
curl http://localhost:8080/api/v1/conversations \
  -H "Authorization: Bearer <token>"
```

### Get Messages (with cursor pagination)

```bash
curl "http://localhost:8080/api/v1/conversations/<id>/messages?limit=50" \
  -H "Authorization: Bearer <token>"

# For next page, use next_cursor from response
curl "http://localhost:8080/api/v1/conversations/<id>/messages?cursor=<next_cursor>&limit=50" \
  -H "Authorization: Bearer <token>"
```

## 🔌 WebSocket

### Connect

```bash
# Using wscat
wscat -c "ws://localhost:8080/ws?token=<access_token>"

# Using websocat
websocat "ws://localhost:8080/ws?token=<access_token>"
```

### Send Message

```json
{
  "conversation_id": "<conversation_uuid>",
  "content": "Hello!"
}
```

### Receive Message

```json
{
  "id": "msg-uuid",
  "conversation_id": "conv-uuid",
  "sender_id": "user-uuid",
  "sender_username": "alice",
  "content": "Hello!",
  "type": "text",
  "sent_at": 1705834567890
}
```

## 🧪 Testing Chat

### Quick Test Flow

```bash
# 1. Register two users
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@test.com","username":"alice","password":"12345678"}'

curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"bob@test.com","username":"bob","password":"12345678"}'

# 2. Create conversation (Alice creates with Bob)
curl -X POST http://localhost:8080/api/v1/conversations \
  -H "Authorization: Bearer <alice_token>" \
  -H "Content-Type: application/json" \
  -d '{"participant_id": "<bob_id>"}'

# 3. Connect both via WebSocket
wscat -c "ws://localhost:8080/ws?token=<alice_token>"
wscat -c "ws://localhost:8080/ws?token=<bob_token>"

# 4. Alice sends message
{"conversation_id": "<conv_id>", "content": "Hey Bob!"}

# 5. Bob receives it in real-time! ✅
```

## 🔧 Development

```bash
make help           # Show all commands
make run            # Run server
make dev            # Run with hot reload (Air)
make build          # Build binary
make test           # Run tests
make lint           # Run linter
make lint-fix       # Run linter with auto-fix
make fmt            # Format code
make docker-up      # Start Redis + PostgreSQL
make docker-down    # Stop containers
make docker-logs    # View container logs
make docker-build   # Build Docker image
make all            # fmt + lint + test + build
```

## 🏗️ Architecture

### Dependency Injection (Uber Fx)

```
Config ─────────────────────────────────────────────────────────┐
   │                                                            │
   ├──► PostgresClient ──► UserRepository ──► AuthService       │
   │                    ├──► ConversationRepository ──┐         │
   │                    └──► MessageRepository ───────┼──► Handlers
   │                                                  │         │
   └──► RedisClient ────► ChatService ────────────────┘         │
                      └──► MessageWorker ───────────────────────┘
```

### Message Flow

```
Alice sends message
        │
        ▼
┌───────────────────┐
│   WebSocket       │
│   Handler         │
└─────────┬─────────┘
          │
          ▼
┌───────────────────┐     ┌─────────────────┐
│   Chat Service    │────►│  Redis Stream   │
└─────────┬─────────┘     └────────┬────────┘
          │                        │
          │                        ▼
          │               ┌─────────────────┐
          │               │  Message Worker │
          │               └────────┬────────┘
          │                        │
          │                        ▼
          │               ┌─────────────────┐
          │               │   PostgreSQL    │
          │               │  (persistence)  │
          │               └─────────────────┘
          │
          ▼
┌───────────────────┐
│  Redis Pub/Sub    │
│  channel:user:X   │
└─────────┬─────────┘
          │
    ┌─────┴─────┐
    ▼           ▼
 Bob's PC   Bob's Phone
 (online)    (online)
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

## 📝 Features

- [x] JWT Authentication (register, login, refresh)
- [x] WebSocket real-time messaging
- [x] Redis Pub/Sub for multi-device support
- [x] Redis Streams for async message processing
- [x] PostgreSQL persistence
- [x] Conversation management (direct & groups)
- [x] Message history with cursor pagination
- [x] Uber Fx dependency injection
- [x] Hot reload development (Air)
- [x] Docker support
- [ ] Typing indicators
- [ ] Read receipts
- [ ] File sharing
- [ ] Push notifications

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.
