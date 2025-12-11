-- Create messages table
CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    content TEXT NOT NULL,
    sender_header VARCHAR(50) NOT NULL,
    received_at TIMESTAMP,
    message_type VARCHAR(20) DEFAULT 'SMS',
    phone_number VARCHAR(20),
    has_links BOOLEAN DEFAULT FALSE,
    link_count INTEGER DEFAULT 0,
    extracted_urls TEXT[],
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_messages_user_id ON messages(user_id);
CREATE INDEX idx_messages_sender_header ON messages(sender_header);
CREATE INDEX idx_messages_created_at ON messages(created_at);
CREATE INDEX idx_messages_has_links ON messages(has_links);

