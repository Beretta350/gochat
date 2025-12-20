# 💬 GoChat

Real-time chat application built with Go, Fiber, Redis Pub/Sub, and WebSocket.

## 🚀 Tech Stack

| Technology | Purpose |
|------------|---------|
| **Go 1.23** | Backend language |
| **Fiber v2** | Web framework |
| **Uber Fx** | Dependency injection |
| **Redis** | Pub/Sub & Streams for messaging |
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
├── configs/
│   └── local.env                # Local environment variables
├── docker-compose.yml           # Redis container
├── Makefile                     # Build and dev commands
├── .air.toml                    # Hot reload config
└── .golangci.yml                # Linter config
```

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
| `REDIS_ADDR` | `localhost:6379` | Redis address |
| `REDIS_PASSWORD` | `` | Redis password |
| `REDIS_DB` | `0` | Redis database |

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
┌─────────┐     WebSocket      ┌─────────────┐
│ Client  │◄──────────────────►│   Fiber     │
└─────────┘                    │   Server    │
                               └──────┬──────┘
                                      │
                         ┌────────────┼────────────┐
                         ▼            ▼            ▼
                   ┌──────────┐ ┌──────────┐ ┌──────────┐
                   │  Pub/Sub │ │  Stream  │ │  Lists   │
                   │(realtime)│ │(persist) │ │(pending) │
                   └──────────┘ └──────────┘ └──────────┘
                               Redis
```

### Message Flow

1. User A connects via WebSocket with `?token=alice`
2. Server subscribes to Redis channel `user:alice`
3. User A sends message to User B
4. Message is added to Redis Stream (for persistence worker)
5. If User B is **online**: publish to Redis channel `user:bob`
6. If User B is **offline**: add to pending queue `pending:bob`
7. When User B connects, pending messages are delivered first

## 📝 TODO

- [ ] PostgreSQL for message persistence
- [ ] JWT Authentication
- [ ] Group chats
- [ ] Message history
- [ ] Read receipts
- [ ] Typing indicators
- [ ] File sharing

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.
