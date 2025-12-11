-- Create sender registry table
CREATE TABLE IF NOT EXISTS sender_registry (
    id UUID PRIMARY KEY,
    sender_id VARCHAR(50) UNIQUE NOT NULL,
    bank_name VARCHAR(255) NOT NULL,
    bank_code VARCHAR(20) NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    verified_by VARCHAR(50) NOT NULL CHECK (verified_by IN ('TELECOM', 'BANK', 'MANUAL')),
    telecom_operator VARCHAR(50),
    registration_date DATE NOT NULL,
    last_verified_at TIMESTAMP,
    reputation_score DECIMAL(3,2) DEFAULT 1.0 CHECK (reputation_score >= 0 AND reputation_score <= 1),
    message_count INTEGER DEFAULT 0,
    fraud_report_count INTEGER DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_sender_registry_sender_id ON sender_registry(sender_id);
CREATE INDEX idx_sender_registry_bank_name ON sender_registry(bank_name);
CREATE INDEX idx_sender_registry_is_verified ON sender_registry(is_verified);
CREATE INDEX idx_sender_registry_is_active ON sender_registry(is_active);
CREATE INDEX idx_sender_registry_reputation_score ON sender_registry(reputation_score);

