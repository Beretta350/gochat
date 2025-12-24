# 💬 GoChat

A full-stack real-time chat application built with **Go** and **Next.js**.

> ⚠️ **Security Note:** This repository contains example/development credentials. Always use strong, unique secrets in production environments.

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                              GoChat                                      │
├─────────────────────────────────┬───────────────────────────────────────┤
│           Frontend              │              Backend                   │
│  ┌───────────────────────────┐  │  ┌─────────────────────────────────┐  │
│  │        Next.js            │  │  │           Go + Fiber            │  │
│  │   TypeScript + Tailwind   │◄─┼──►   WebSocket + REST API         │  │
│  │        React              │  │  │     JWT Authentication         │  │
│  └───────────────────────────┘  │  └───────────────┬─────────────────┘  │
│                                 │                  │                    │
│                                 │    ┌─────────────┼─────────────┐      │
│                                 │    ▼             ▼             ▼      │
│                                 │ ┌──────┐   ┌─────────┐   ┌──────────┐ │
│                                 │ │Redis │   │ Redis   │   │PostgreSQL│ │
│                                 │ │Pub/Sub│   │ Stream  │   │          │ │
│                                 │ └──────┘   └─────────┘   └──────────┘ │
└─────────────────────────────────┴───────────────────────────────────────┘
```

## 📁 Project Structure

```
gochat/
├── backend/                # Go API (Fiber + WebSocket)
│   ├── cmd/                # Application entrypoint
│   ├── internal/           # Private application code
│   ├── pkg/                # Reusable packages
│   ├── database/           # SQL migrations & schema
│   ├── docs/               # API documentation
│   └── README.md           # Backend-specific docs
│
├── frontend/               # Next.js Web App
│   ├── src/
│   │   ├── app/            # Next.js App Router
│   │   ├── components/     # React components
│   │   ├── hooks/          # Custom hooks
│   │   ├── stores/         # State management
│   │   ├── lib/            # Utilities
│   │   └── types/          # TypeScript types
│   └── README.md           # Frontend-specific docs
│
├── docker-compose.yml      # Full stack orchestration
├── Makefile                # Project commands
└── README.md               # You are here
```

## 🚀 Tech Stack

| Layer | Technology |
|-------|------------|
| **Frontend** | Next.js, React, TypeScript, Tailwind CSS |
| **Backend** | Go 1.23, Fiber v2, Uber Fx |
| **Database** | PostgreSQL 16 |
| **Cache/Realtime** | Redis 7 (Pub/Sub + Streams) |
| **Auth** | JWT (Access + Refresh tokens) |
| **Realtime** | WebSocket |
| **Infrastructure** | Docker, Docker Compose |

## 🛠️ Getting Started

### Prerequisites

- Docker & Docker Compose
- Make
- Go 1.23+ (for backend development)
- Node.js 20+ (for frontend development)

### Quick Start (Docker)

```bash
# Clone the repository
git clone https://github.com/Beretta350/gochat.git
cd gochat

# Start all services
make up

# Services will be available at:
# - Frontend: http://localhost:3000
# - Backend API: http://localhost:8080
# - Redis Commander: http://localhost:8081
```

### Development Mode

```bash
# Start infrastructure (PostgreSQL + Redis)
make infra

# In one terminal - run backend with hot reload
make dev-api

# In another terminal - run frontend with hot reload
make dev-web
```

## 📋 Available Commands

```bash
make help           # Show all commands

# Full Stack
make up             # Start all services (Docker)
make down           # Stop all services
make logs           # View all logs

# Infrastructure
make infra          # Start only PostgreSQL + Redis
make infra-down     # Stop infrastructure

# Development
make dev-api        # Run backend with hot reload
make dev-web        # Run frontend with hot reload

# Backend specific
make api-build      # Build backend
make api-test       # Run backend tests
make api-lint       # Lint backend code

# Frontend specific
make web-build      # Build frontend
make web-lint       # Lint frontend code
```

## 🔗 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/register` | Create account |
| POST | `/api/v1/auth/login` | Login |
| POST | `/api/v1/auth/refresh` | Refresh token |
| GET | `/api/v1/auth/me` | Get current user |
| POST | `/api/v1/conversations` | Create conversation |
| GET | `/api/v1/conversations` | List conversations |
| GET | `/api/v1/conversations/:id/messages` | Get messages |
| WS | `/ws?token=<jwt>` | WebSocket connection |

> 📖 See [backend/README.md](backend/README.md) for detailed API documentation.

## ✨ Features

- [x] User authentication (register, login, JWT)
- [x] Real-time messaging via WebSocket
- [x] Direct messages (1:1)
- [x] Group conversations
- [x] Message history with pagination
- [x] Multi-device support (Redis Pub/Sub)
- [x] Offline message queue
- [ ] Typing indicators
- [ ] Read receipts
- [ ] File sharing
- [ ] Push notifications

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

