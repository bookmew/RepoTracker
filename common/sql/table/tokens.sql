CREATE TABLE tokens (
    id SERIAL PRIMARY KEY,
    symbol TEXT NOT NULL,
    name TEXT NOT NULL,
    address TEXT NOT NULL,
    metadata JSONB,
    repo_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_token_symbol ON tokens(symbol);
CREATE INDEX idx_token_address ON tokens(address);

ALTER TABLE tokens ADD CONSTRAINT unique_address UNIQUE (address);

-- Insert sample data
-- INSERT INTO tokens (symbol, name, address, metadata, repo_url)
-- VALUES (
--     'AGNT', 
--     'AI Agent Token', 
--     '1234567890abcdef1234567890abcdef12345678', 
--     '{"description": "An AI agent meme token", "image": "https://example.com/image.png"]}',
--     'https://github.com/example/ai-agent-repo'
-- );