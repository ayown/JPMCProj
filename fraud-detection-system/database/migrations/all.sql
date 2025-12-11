-- Apply all schema migrations

-- 001_create_users.sql
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20),
    is_active BOOLEAN DEFAULT TRUE,
    is_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_login_at TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);

-- 002_create_messages.sql
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
CREATE INDEX IF NOT EXISTS idx_messages_user_id ON messages(user_id);
CREATE INDEX IF NOT EXISTS idx_messages_sender_header ON messages(sender_header);
CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at);
CREATE INDEX IF NOT EXISTS idx_messages_has_links ON messages(has_links);

-- 003_create_verifications.sql
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
CREATE INDEX IF NOT EXISTS idx_verifications_message_id ON verifications(message_id);
CREATE INDEX IF NOT EXISTS idx_verifications_user_id ON verifications(user_id);
CREATE INDEX IF NOT EXISTS idx_verifications_is_fraud ON verifications(is_fraud);
CREATE INDEX IF NOT EXISTS idx_verifications_fraud_score ON verifications(fraud_score);
CREATE INDEX IF NOT EXISTS idx_verifications_created_at ON verifications(created_at);

-- 004_create_reports.sql
CREATE TABLE IF NOT EXISTS reports (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    message_id UUID REFERENCES messages(id) ON DELETE SET NULL,
    verification_id UUID REFERENCES verifications(id) ON DELETE SET NULL,
    report_type VARCHAR(20) NOT NULL CHECK (report_type IN ('FRAUD', 'FALSE_POSITIVE', 'FEEDBACK')),
    content TEXT NOT NULL,
    sender_header VARCHAR(50) NOT NULL,
    description TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'REVIEWED', 'RESOLVED', 'DISMISSED')),
    priority VARCHAR(20) NOT NULL DEFAULT 'MEDIUM' CHECK (priority IN ('LOW', 'MEDIUM', 'HIGH', 'CRITICAL')),
    reviewed_by UUID REFERENCES users(id) ON DELETE SET NULL,
    reviewed_at TIMESTAMP,
    review_notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_reports_user_id ON reports(user_id);
CREATE INDEX IF NOT EXISTS idx_reports_message_id ON reports(message_id);
CREATE INDEX IF NOT EXISTS idx_reports_verification_id ON reports(verification_id);
CREATE INDEX IF NOT EXISTS idx_reports_status ON reports(status);
CREATE INDEX IF NOT EXISTS idx_reports_priority ON reports(priority);
CREATE INDEX IF NOT EXISTS idx_reports_created_at ON reports(created_at);

-- 005_create_rbi_circulars.sql
CREATE TABLE IF NOT EXISTS rbi_circulars (
    id UUID PRIMARY KEY,
    circular_number VARCHAR(50) UNIQUE NOT NULL,
    title VARCHAR(500) NOT NULL,
    content TEXT NOT NULL,
    issued_date DATE NOT NULL,
    effective_date DATE,
    expiry_date DATE,
    category VARCHAR(50) NOT NULL,
    keywords TEXT[],
    is_active BOOLEAN DEFAULT TRUE,
    source_url VARCHAR(500),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_rbi_circulars_circular_number ON rbi_circulars(circular_number);
CREATE INDEX IF NOT EXISTS idx_rbi_circulars_category ON rbi_circulars(category);
CREATE INDEX IF NOT EXISTS idx_rbi_circulars_is_active ON rbi_circulars(is_active);
CREATE INDEX IF NOT EXISTS idx_rbi_circulars_issued_date ON rbi_circulars(issued_date);
CREATE INDEX IF NOT EXISTS idx_rbi_circulars_keywords ON rbi_circulars USING GIN(keywords);

-- 006_create_sender_registry.sql
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
CREATE INDEX IF NOT EXISTS idx_sender_registry_sender_id ON sender_registry(sender_id);
CREATE INDEX IF NOT EXISTS idx_sender_registry_bank_name ON sender_registry(bank_name);
CREATE INDEX IF NOT EXISTS idx_sender_registry_is_verified ON sender_registry(is_verified);
CREATE INDEX IF NOT EXISTS idx_sender_registry_is_active ON sender_registry(is_active);
CREATE INDEX IF NOT EXISTS idx_sender_registry_reputation_score ON sender_registry(reputation_score);
