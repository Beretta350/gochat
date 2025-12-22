# Conversations API Documentation

GoChat supports both direct (1:1) and group conversations.

## Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                    Conversation Types                            │
│                                                                 │
│   Direct (1:1)           │          Group                       │
│   ────────────           │          ─────                       │
│   • 2 participants       │   • 2+ participants                  │
│   • No name              │   • Has name                         │
│   • Idempotent creation  │   • Creator becomes admin            │
│   • No roles             │   • Roles: admin, member             │
└─────────────────────────────────────────────────────────────────┘
```

## Endpoints

### Create Conversation

Creates a new conversation or returns existing one (for direct chats).

```http
POST /api/v1/conversations
Authorization: Bearer <access_token>
Content-Type: application/json
```

#### Direct Conversation (1:1)

```json
{
  "participant_id": "other-user-uuid"
}
```

**Response (201 Created / 200 OK):**

```json
{
  "conversation": {
    "id": "8b3d468f-d93d-431e-ba9c-9ca14b4ece77",
    "type": "direct",
    "created_at": "2025-12-22T21:59:55.025Z",
    "updated_at": "2025-12-22T21:59:55.025Z"
  },
  "participants": [
    {
      "id": "ff97a765-7471-4740-a28e-6866dbee6706",
      "email": "alice@test.com",
      "username": "alice",
      "is_active": true,
      "joined_at": "2025-12-22T21:59:55.025Z",
      "created_at": "2025-12-22T21:17:29.380Z"
    },
    {
      "id": "fd14141e-4576-4ab9-9fa1-f832ef1afc7a",
      "email": "bob@test.com",
      "username": "bob",
      "is_active": true,
      "joined_at": "2025-12-22T21:59:55.025Z",
      "created_at": "2025-12-21T20:08:32.667Z"
    }
  ],
  "is_new": true
}
```

> **Note:** If a direct conversation already exists between the two users, it returns the existing one with `"is_new": false`.

#### Group Conversation

```json
{
  "name": "Project Team",
  "participant_ids": ["user-uuid-1", "user-uuid-2", "user-uuid-3"]
}
```

**Response (201 Created):**

```json
{
  "conversation": {
    "id": "9c4e579g-e04e-542f-cb0d-0db25c5fdf88",
    "type": "group",
    "name": "Project Team",
    "created_by": "your-user-uuid",
    "created_at": "2025-12-22T22:00:00.000Z",
    "updated_at": "2025-12-22T22:00:00.000Z"
  },
  "participants": [
    {
      "id": "your-user-uuid",
      "username": "you",
      "role": "admin",
      "joined_at": "2025-12-22T22:00:00.000Z"
    },
    {
      "id": "user-uuid-1",
      "username": "member1",
      "role": "member",
      "joined_at": "2025-12-22T22:00:00.000Z"
    }
  ]
}
```

**Errors:**

| Status | Message |
|--------|---------|
| 400 | Invalid request body |
| 400 | Group name is required |
| 401 | Invalid/missing token |
| 404 | Participant not found |

---

### List Conversations

Get all conversations for the authenticated user.

```http
GET /api/v1/conversations
Authorization: Bearer <access_token>
```

**Response (200 OK):**

```json
{
  "conversations": [
    {
      "conversation": {
        "id": "8b3d468f-d93d-431e-ba9c-9ca14b4ece77",
        "type": "direct",
        "created_at": "2025-12-22T21:59:55.025Z",
        "updated_at": "2025-12-22T21:59:55.025Z"
      },
      "participants": [...]
    },
    {
      "conversation": {
        "id": "9c4e579g-e04e-542f-cb0d-0db25c5fdf88",
        "type": "group",
        "name": "Project Team",
        "created_at": "2025-12-22T22:00:00.000Z",
        "updated_at": "2025-12-22T22:00:00.000Z"
      },
      "participants": [...]
    }
  ],
  "count": 2
}
```

---

### Get Conversation

Get details of a specific conversation.

```http
GET /api/v1/conversations/:id
Authorization: Bearer <access_token>
```

**Response (200 OK):**

```json
{
  "conversation": {
    "id": "8b3d468f-d93d-431e-ba9c-9ca14b4ece77",
    "type": "direct",
    "created_at": "2025-12-22T21:59:55.025Z",
    "updated_at": "2025-12-22T21:59:55.025Z"
  },
  "participants": [...]
}
```

**Errors:**

| Status | Message |
|--------|---------|
| 401 | Invalid/missing token |
| 403 | You are not a participant of this conversation |
| 404 | Conversation not found |

---

### Get Messages

Get messages for a conversation with cursor-based pagination.

```http
GET /api/v1/conversations/:id/messages
Authorization: Bearer <access_token>
```

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `limit` | int | 50 | Max messages to return (max: 100) |
| `cursor` | string | | Timestamp cursor for pagination |

**Response (200 OK):**

```json
{
  "messages": [
    {
      "id": "e1f063f6-4912-41b5-9aeb-9e93528e7a18",
      "conversation_id": "8b3d468f-d93d-431e-ba9c-9ca14b4ece77",
      "sender_id": "fd14141e-4576-4ab9-9fa1-f832ef1afc7a",
      "sender_username": "gabriel",
      "content": "Hey! How are you?",
      "type": "text",
      "sent_at": "2025-12-22T22:16:21.203Z"
    },
    {
      "id": "f2g174g7-5023-52c6-0bfc-0f04639f8b29",
      "conversation_id": "8b3d468f-d93d-431e-ba9c-9ca14b4ece77",
      "sender_id": "ff97a765-7471-4740-a28e-6866dbee6706",
      "sender_username": "stefany",
      "content": "I'm good! You?",
      "type": "text",
      "sent_at": "2025-12-22T22:16:25.456Z"
    }
  ],
  "has_more": true,
  "next_cursor": "2025-12-22T22:16:21.203000000Z"
}
```

**Pagination:**

To get the next page, use `next_cursor` from the response:

```http
GET /api/v1/conversations/:id/messages?cursor=2025-12-22T22:16:21.203000000Z&limit=50
```

**Errors:**

| Status | Message |
|--------|---------|
| 401 | Invalid/missing token |
| 403 | You are not a participant of this conversation |
| 404 | Conversation not found |

---

## WebSocket Messaging

### Send Message

Once connected to WebSocket, send messages with:

```json
{
  "conversation_id": "8b3d468f-d93d-431e-ba9c-9ca14b4ece77",
  "content": "Hello!"
}
```

### Receive Message

Messages are delivered to all online participants:

```json
{
  "id": "msg-uuid",
  "conversation_id": "8b3d468f-d93d-431e-ba9c-9ca14b4ece77",
  "sender_id": "sender-uuid",
  "sender_username": "alice",
  "content": "Hello!",
  "type": "text",
  "sent_at": 1705834567890
}
```

### Error Messages

If something goes wrong, you'll receive:

```json
{
  "error": true,
  "message": "conversation_id is required"
}
```

**Possible Errors:**

| Message |
|---------|
| `Invalid message format` |
| `conversation_id is required` |
| `content is required` |
| `Conversation not found` |
| `You are not a participant of this conversation` |

---

## Message Types

| Type | Description |
|------|-------------|
| `text` | Plain text message (default) |
| `image` | Image attachment (future) |
| `file` | File attachment (future) |
| `audio` | Audio message (future) |

---

## Message Flow

```
Alice sends message
        │
        ▼
┌───────────────────┐
│   Validate        │ ◄── Is Alice a participant?
│   conversation_id │
└─────────┬─────────┘
          │
          ▼
┌───────────────────┐
│  Add to Redis     │ ◄── For async persistence
│  Stream           │
└─────────┬─────────┘
          │
          ▼
┌───────────────────┐
│  Get participants │
│  from DB          │
└─────────┬─────────┘
          │
          ▼
    For each participant:
          │
    ┌─────┴─────┐
    ▼           ▼
 Online?     Offline?
    │           │
    ▼           ▼
 Pub/Sub     Pending
 channel     queue
    │           │
    ▼           ▼
 Receives    Receives
 instantly   on reconnect
```

---

## Best Practices

1. **Always create conversation first** before sending messages
2. **Cache conversation_id** to avoid repeated API calls
3. **Use cursor pagination** for large message histories
4. **Handle WebSocket errors** gracefully in your client
5. **Refresh tokens** before they expire to maintain connection
