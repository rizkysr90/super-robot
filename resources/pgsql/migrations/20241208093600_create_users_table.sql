-- migrate:up
CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    google_id VARCHAR(255),
    password_hash VARCHAR(255),
    auth_type VARCHAR(20) NOT NULL,
    user_type VARCHAR(20) NOT NULL,
    tenant_id uuid NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    last_login_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- CONSTRAINT users_auth_type_check CHECK (auth_type IN ('google', 'password')),
-- CONSTRAINT users_user_type_check CHECK (user_type IN ('owner', 'admin'))
-- migrate:down

