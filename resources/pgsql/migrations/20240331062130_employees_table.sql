-- migrate:up
CREATE TABLE IF NOT EXISTS employees (
    id UUID PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
    contact VARCHAR(30) NOT NULL,
    username VARCHAR(50) NOT NULL,
    password TEXT NOT NULL,
    store_id UUID NOT NULL,
    role INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ NULL,

    CONSTRAINT fk_store_employees FOREIGN KEY (store_id) 
    REFERENCES stores(id)
);

-- migrate:down

