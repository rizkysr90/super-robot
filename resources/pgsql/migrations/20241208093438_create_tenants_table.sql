-- migrate:up
CREATE TABLE IF NOT EXISTS tenants (
    id uuid PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    owner_id uuid NULL,  -- Allow NULL initially
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

-- migrate:down

