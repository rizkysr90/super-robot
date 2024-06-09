-- migrate:up
CREATE TABLE IF NOT EXISTS branches (
    id UUID PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
    address TEXT,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    user_id UUID REFERENCES users(id) -- Foreign key constraint referencing the users table
);

-- migrate:down

