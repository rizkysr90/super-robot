-- migrate:up
CREATE TABLE IF NOT EXISTS stores (
    id UUID PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
    address TEXT NOT NULL,
    contact VARCHAR(50) NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ NULL,

    CONSTRAINT fk_user_store FOREIGN KEY (user_id) 
    REFERENCES users(id)
);

-- migrate:down

