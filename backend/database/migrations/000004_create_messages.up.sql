-- Create messages table
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
