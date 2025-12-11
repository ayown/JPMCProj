-- Create verifications table
CREATE TABLE IF NOT EXISTS verifications (
    id UUID PRIMARY KEY,
    message_id UUID NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    is_fraud BOOLEAN NOT NULL,
    fraud_score DECIMAL(5,4) NOT NULL,
    fraud_type VARCHAR(50),
    confidence DECIMAL(5,4) NOT NULL,
    model_version VARCHAR(50) NOT NULL,
    ml_predictions JSONB,
    header_verified BOOLEAN DEFAULT FALSE,
    header_score DECIMAL(5,4) DEFAULT 0,
    rbi_compliant BOOLEAN DEFAULT TRUE,
    rbi_verification_result JSONB,
    explanation TEXT,
    recommendations JSONB,
    processing_time_ms INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_verifications_message_id ON verifications(message_id);
CREATE INDEX idx_verifications_user_id ON verifications(user_id);
CREATE INDEX idx_verifications_is_fraud ON verifications(is_fraud);
CREATE INDEX idx_verifications_fraud_score ON verifications(fraud_score);
CREATE INDEX idx_verifications_created_at ON verifications(created_at);

