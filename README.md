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
| **Frontend** | Next.js 14, React 18, TypeScript, Tailwind CSS |
| **State Management** | Redux Toolkit, RTK Query |
| **UI Components** | Radix UI, Framer Motion |
| **Forms** | React Hook Form, Zod |
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
make help                # Show all commands

# Docker - Full Stack
make docker-up           # Start all services (web + api + infra)
make docker-down         # Stop all services
make docker-logs         # View all logs
make docker-build        # Build all images
make docker-restart      # Rebuild and restart all

# Docker - Infrastructure
make docker-infra        # Start only PostgreSQL + Redis
make docker-infra-down   # Stop infrastructure

# Docker - Backend Only
make docker-api-up       # Start API + infra (no frontend)
make docker-api-build    # Build API image
make docker-api-logs     # View API logs
make docker-api-restart  # Rebuild and restart API

# Docker - Frontend Only  
make docker-web-up       # Start Web + API + infra
make docker-web-build    # Build Web image
make docker-web-logs     # View Web logs
make docker-web-restart  # Rebuild and restart Web

# Development (Local)
make dev-api             # Run backend with hot reload
make dev-web             # Run frontend dev server

# Backend (Go)
make api-build           # Build backend binary
make api-test            # Run backend tests
make api-lint            # Lint backend code
make api-fmt             # Format backend code

# Frontend (Next.js)
make web-install         # Install dependencies
make web-build           # Build frontend
make web-lint            # Lint frontend code
make web-test            # Run frontend tests
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

### Backend
- [x] User authentication (register, login, JWT)
- [x] Real-time messaging via WebSocket
- [x] Direct messages (1:1)
- [x] Group conversations
- [x] Message history with pagination
- [x] Multi-device support (Redis Pub/Sub)
- [x] Offline message queue

### Frontend
- [x] Modern responsive UI (mobile-first)
- [x] Dark theme with custom color palette
- [x] Smooth animations (Framer Motion)
- [x] Form validation (React Hook Form + Zod)
- [x] State management (Redux Toolkit)
- [x] Data caching (RTK Query)

### Coming Soon
- [ ] Typing indicators
- [ ] Read receipts
- [ ] File sharing
- [ ] Push notifications
- [ ] PWA support

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

