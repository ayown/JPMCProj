-- Create RBI circulars table
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

-- Create indexes
CREATE INDEX idx_rbi_circulars_circular_number ON rbi_circulars(circular_number);
CREATE INDEX idx_rbi_circulars_category ON rbi_circulars(category);
CREATE INDEX idx_rbi_circulars_is_active ON rbi_circulars(is_active);
CREATE INDEX idx_rbi_circulars_issued_date ON rbi_circulars(issued_date);
CREATE INDEX idx_rbi_circulars_keywords ON rbi_circulars USING GIN(keywords);

