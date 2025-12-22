# Database Documentation

## Overview

GoChat uses PostgreSQL for persistent storage of users, conversations, and messages.

## Entity Relationship Diagram

```
┌─────────────────────────────────────────────────────────────────────────┐
│                              PostgreSQL                                  │
│                                                                         │
│  ┌─────────────────┐          ┌─────────────────────┐                  │
│  │     users       │          │   conversations     │                  │
│  ├─────────────────┤          ├─────────────────────┤                  │
│  │ id          PK  │◄─────┐   │ id              PK  │                  │
│  │ email           │      │   │ type   (direct/grp) │                  │
│  │ username        │      │   │ name   (groups only)│                  │
│  │ password_hash   │      └───│ created_by      FK  │                  │
│  │ is_active       │          │ created_at          │                  │
│  │ created_at      │          │ updated_at          │                  │
│  │ updated_at      │          └──────────┬──────────┘                  │
│  └────────┬────────┘                     │                             │
│           │                              │                             │
│           │         ┌────────────────────┴────────────────────┐        │
│           │         │                                         │        │
│           │         ▼                                         ▼        │
│           │  ┌─────────────────────────┐          ┌─────────────────┐  │
│           │  │ conversation_participants│          │    messages     │  │
│           │  ├─────────────────────────┤          ├─────────────────┤  │
│           └──│ user_id           PK,FK │          │ id          PK  │  │
│              │ conversation_id   PK,FK │◄─────────│ conversation_id │  │
│              │ role (groups only)      │          │ sender_id   FK  │──┘
│              │ joined_at               │          │ content         │
│              │ left_at                 │          │ type            │
│              └─────────────────────────┘          │ sent_at         │
│                                                   │ created_at      │
│                                                   └─────────────────┘
└─────────────────────────────────────────────────────────────────────────┘
```

## Tables

### users

Stores user account information.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PK, DEFAULT | Unique identifier |
| `email` | VARCHAR(255) | UNIQUE, NOT NULL | User email for login |
| `username` | VARCHAR(100) | UNIQUE, NOT NULL | Display name |
| `password_hash` | VARCHAR(255) | NOT NULL | Bcrypt hashed password |
| `is_active` | BOOLEAN | DEFAULT true | Soft delete flag |
| `created_at` | TIMESTAMPTZ | DEFAULT NOW() | Creation timestamp |
| `updated_at` | TIMESTAMPTZ | DEFAULT NOW() | Last update timestamp |

**Indexes:**
- `idx_users_email` - Fast lookup by email (login)
- `idx_users_username` - Fast lookup by username

---

### conversations

Stores chat conversations (1:1 or groups).

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PK, DEFAULT | Unique identifier |
| `type` | VARCHAR(20) | NOT NULL, CHECK | `'direct'` or `'group'` |
| `name` | VARCHAR(255) | | Group name (NULL for direct) |
| `created_by` | UUID | FK → users | User who created the conversation |
| `created_at` | TIMESTAMPTZ | DEFAULT NOW() | Creation timestamp |
| `updated_at` | TIMESTAMPTZ | DEFAULT NOW() | Last update timestamp |

**Constraints:**
- `chk_conversation_type`: type must be 'direct' or 'group'

---

### conversation_participants

Junction table linking users to conversations.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `conversation_id` | UUID | PK, FK | Reference to conversation |
| `user_id` | UUID | PK, FK | Reference to user |
| `role` | VARCHAR(20) | CHECK | NULL for direct, `'admin'`/`'member'` for groups |
| `joined_at` | TIMESTAMPTZ | DEFAULT NOW() | When user joined |
| `left_at` | TIMESTAMPTZ | | When user left (NULL = active) |

**Indexes:**
- `idx_participants_user_active` - Find user's active conversations
- `idx_participants_conversation` - Find conversation members

**Constraints:**
- `chk_participant_role`: role must be NULL, 'admin', or 'member'

---

### messages

Stores all chat messages.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | UUID | PK, DEFAULT | Unique identifier |
| `conversation_id` | UUID | FK, NOT NULL | Reference to conversation |
| `sender_id` | UUID | FK, NOT NULL | Reference to sender |
| `content` | TEXT | NOT NULL | Message content |
| `type` | VARCHAR(20) | CHECK, DEFAULT | `'text'`, `'image'`, `'file'`, `'audio'` |
| `sent_at` | TIMESTAMPTZ | NOT NULL | When message was sent (client time) |
| `created_at` | TIMESTAMPTZ | DEFAULT NOW() | Database insertion time |

**Indexes:**
- `idx_messages_conversation_time` - Fetch history (most recent first)
- `idx_messages_conversation_cursor` - Cursor-based pagination

**Constraints:**
- `chk_message_type`: type must be 'text', 'image', 'file', or 'audio'

---

## Common Queries

### Get user's conversations

```sql
SELECT c.*, 
       (SELECT content FROM messages m 
        WHERE m.conversation_id = c.id 
        ORDER BY sent_at DESC LIMIT 1) as last_message
FROM conversations c
JOIN conversation_participants cp ON c.id = cp.conversation_id
WHERE cp.user_id = $1 AND cp.left_at IS NULL
ORDER BY c.updated_at DESC;
```

### Get conversation messages (cursor pagination)

```sql
-- Most recent first (default)
SELECT m.id, m.conversation_id, m.sender_id, u.username, 
       m.content, m.type, m.sent_at
FROM messages m
JOIN users u ON m.sender_id = u.id
WHERE m.conversation_id = $1
  AND m.sent_at < $2  -- cursor (optional)
ORDER BY m.sent_at DESC
LIMIT 50;
```

### Find existing direct conversation

```sql
SELECT c.id FROM conversations c
WHERE c.type = 'direct'
  AND EXISTS (
    SELECT 1 FROM conversation_participants cp1
    WHERE cp1.conversation_id = c.id 
    AND cp1.user_id = $1 
    AND cp1.left_at IS NULL
  )
  AND EXISTS (
    SELECT 1 FROM conversation_participants cp2
    WHERE cp2.conversation_id = c.id 
    AND cp2.user_id = $2 
    AND cp2.left_at IS NULL
  );
```

### Create direct conversation

```sql
-- 1. Create conversation
INSERT INTO conversations (type) 
VALUES ('direct') 
RETURNING id;

-- 2. Add participants
INSERT INTO conversation_participants (conversation_id, user_id) 
VALUES ($conv_id, $user1), ($conv_id, $user2);
```

### Batch insert messages (from Redis Stream worker)

```sql
INSERT INTO messages (conversation_id, sender_id, content, type, sent_at)
VALUES 
  ($1, $2, $3, $4, $5),
  ($6, $7, $8, $9, $10),
  ...;
```

---

## Message Response Format

When fetching messages via API, the response is flattened:

```json
{
  "id": "msg-uuid",
  "conversation_id": "conv-uuid",
  "sender_id": "user-uuid",
  "sender_username": "alice",
  "content": "Hello!",
  "type": "text",
  "sent_at": "2025-12-22T22:16:21.203Z"
}
```

> **Note:** `sender_username` is joined from the `users` table for convenience. The full sender object is not included to keep responses lightweight.

---

## Migrations

Migrations are located in `database/migrations/` and follow the naming convention:

```
{version}_{description}.up.sql   - Apply migration
{version}_{description}.down.sql - Rollback migration
```

To run migrations, use [golang-migrate](https://github.com/golang-migrate/migrate):

```bash
# Install
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
migrate -path database/migrations \
  -database "postgres://user:pass@localhost:5432/gochat?sslmode=disable" up

# Rollback last migration
migrate -path database/migrations \
  -database "postgres://user:pass@localhost:5432/gochat?sslmode=disable" down 1

# Check migration version
migrate -path database/migrations \
  -database "postgres://user:pass@localhost:5432/gochat?sslmode=disable" version
```

---

## Redis Integration

While PostgreSQL handles persistence, Redis provides real-time capabilities:

```
┌─────────────────────────────────────────────────────────────┐
│                      Data Flow                               │
│                                                             │
│  WebSocket ──► Redis Stream ──► Worker ──► PostgreSQL       │
│      │                                                      │
│      └──► Redis Pub/Sub ──► Online users (real-time)        │
│                                                             │
│  API GET ◄────────────────────────────── PostgreSQL         │
└─────────────────────────────────────────────────────────────┘
```

| Component | Purpose |
|-----------|---------|
| **Redis Stream** | Buffer messages for batch insertion |
| **Redis Pub/Sub** | Real-time delivery to online users |
| **Redis Lists** | Pending messages for offline users |
| **PostgreSQL** | Permanent storage, history queries |
