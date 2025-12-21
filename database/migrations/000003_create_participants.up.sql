-- Create conversation_participants table
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
