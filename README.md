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
│   │   ├── chat/
│   │   │   └── service.go       # Chat service with Redis Pub/Sub
│   │   ├── handler/
│   │   │   ├── health.go        # Health check handler
│   │   │   └── websocket.go     # WebSocket handler
│   │   ├── middleware/
│   │   │   ├── error_handler.go # Custom error handler
│   │   │   └── middlewares.go   # Fiber middlewares setup
│   │   ├── model/
│   │   │   └── chat_message_model.go
│   │   ├── repository/
│   │   │   └── message_repository.go  # Message persistence
│   │   └── worker/
│   │       └── message_worker.go      # Redis Stream consumer
│   └── config/
│       └── config.go            # Configuration (Fx provider)
├── pkg/
│   ├── envutil/                 # Environment utilities
│   ├── logger/                  # Zap logger wrapper
│   └── redisclient/             # Redis client (Fx provider)
├── database/
│   ├── schema.sql               # Complete database schema
│   └── migrations/              # Versioned SQL migrations
├── docs/
│   └── DATABASE.md              # Database documentation
├── configs/
│   └── local.env                # Local environment variables
├── docker-compose.yml           # Redis + PostgreSQL containers
├── Makefile                     # Build and dev commands
├── .air.toml                    # Hot reload config
└── .golangci.yml                # Linter config
```

> 📖 See [docs/DATABASE.md](docs/DATABASE.md) for complete database documentation.

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
# Connect with wscat
wscat -c "ws://localhost:8080/ws?token=alice"
```

### Message Format

```json
{
  "recipient": "bob",
  "content": "Hello Bob!"
}
```

## 🧪 Testing Chat

Open two terminals:

```bash
# Terminal 1 - Alice
wscat -c "ws://localhost:8080/ws?token=alice"

# Terminal 2 - Bob
wscat -c "ws://localhost:8080/ws?token=bob"

# In Alice's terminal, send:
{"recipient": "bob", "content": "Hey Bob!"}

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
     │  POST /auth/register                    │
     │  { email, username, password }          │
     │────────────────────────────────────────►│
     │                                         │
     │  { user }                               │
     │◄────────────────────────────────────────│
     │                                         │
     │  POST /auth/login                       │
     │  { email, password }                    │
     │────────────────────────────────────────►│
     │                                         │
     │  { access_token (15min),                │
     │    refresh_token (7d) }                 │
     │◄────────────────────────────────────────│
     │                                         │
     │  WS /ws?token={access_token}            │
     │────────────────────────────────────────►│
     │                                         │
     │  Connection established                 │
     │◄═══════════════════════════════════════►│
     │                                         │
```

## 📝 TODO

- [x] WebSocket real-time messaging
- [x] Redis Pub/Sub for multi-device support
- [x] Redis Streams for async processing
- [x] Uber Fx dependency injection
- [x] Database schema design
- [ ] PostgreSQL integration
- [ ] JWT Authentication (register, login, refresh)
- [ ] User management (CRUD)
- [ ] Conversation management (create, list)
- [ ] Message history with cursor pagination
- [ ] Group chats
- [ ] Typing indicators
- [ ] Read receipts
- [ ] File sharing

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.
