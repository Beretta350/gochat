-- ============================================================================
-- GoChat Database Schema
-- PostgreSQL 15+
-- ============================================================================

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ============================================================================
-- USERS
-- ============================================================================
CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email           VARCHAR(255) UNIQUE NOT NULL,
    username        VARCHAR(100) UNIQUE NOT NULL,
    password_hash   VARCHAR(255) NOT NULL,
    is_active       BOOLEAN DEFAULT true,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);

-- ============================================================================
-- CONVERSATIONS
-- ============================================================================
-- type: 'direct' (1:1) or 'group'
-- name: only for groups, NULL for direct
CREATE TABLE conversations (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type        VARCHAR(20) NOT NULL DEFAULT 'direct',
    name        VARCHAR(255),
    created_by  UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    updated_at  TIMESTAMPTZ DEFAULT NOW(),
    
    CONSTRAINT chk_conversation_type CHECK (type IN ('direct', 'group'))
);

-- ============================================================================
-- CONVERSATION PARTICIPANTS
-- ============================================================================
-- role: NULL for direct chats, 'admin'/'member' for groups
-- left_at: NULL means still in conversation, timestamp means left
CREATE TABLE conversation_participants (
    conversation_id  UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    user_id          UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role             VARCHAR(20),
    joined_at        TIMESTAMPTZ DEFAULT NOW(),
    left_at          TIMESTAMPTZ,
    
    PRIMARY KEY (conversation_id, user_id),
    CONSTRAINT chk_participant_role CHECK (role IS NULL OR role IN ('admin', 'member'))
);

-- Index for finding user's conversations quickly
CREATE INDEX idx_participants_user_active ON conversation_participants(user_id) 
    WHERE left_at IS NULL;

-- Index for finding conversation members
CREATE INDEX idx_participants_conversation ON conversation_participants(conversation_id)
    WHERE left_at IS NULL;

-- ============================================================================
-- MESSAGES
-- ============================================================================
-- type: 'text', 'image', 'file', 'audio' (for future use)
CREATE TABLE messages (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    sender_id       UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content         TEXT NOT NULL,
    type            VARCHAR(20) DEFAULT 'text',
    sent_at         TIMESTAMPTZ NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    
    CONSTRAINT chk_message_type CHECK (type IN ('text', 'image', 'file', 'audio'))
);

-- Index for fetching conversation history (most recent first)
CREATE INDEX idx_messages_conversation_time ON messages(conversation_id, sent_at DESC);

-- Index for cursor-based pagination
CREATE INDEX idx_messages_conversation_cursor ON messages(conversation_id, sent_at, id);

-- ============================================================================
-- HELPER FUNCTIONS
-- ============================================================================

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers for updated_at
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_conversations_updated_at
    BEFORE UPDATE ON conversations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
