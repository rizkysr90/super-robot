-- migrate:up
CREATE TABLE IF NOT EXISTS sessions (
    session_id UUID PRIMARY KEY,
    access_token VARCHAR(255) NOT NULL,
    refresh_token VARCHAR(255) NOT NULL,
    user_email VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    user_fullname VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);

-- Create index on user_id for faster lookups
CREATE INDEX idx_sessions_user_id ON sessions(user_id);

-- Create index on user_email for faster lookups
CREATE INDEX idx_sessions_user_email ON sessions(user_email);

CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);
-- migrate:down

