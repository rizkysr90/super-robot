-- migrate:up
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY NOT NULL,
    category_name VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ
);

-- migrate:down

