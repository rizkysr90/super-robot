-- migrate:up
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
    price NUMERIC NOT NULL,
    stock_quantity INTEGER DEFAULT 0, -- Default to 0 for new products
    description TEXT,
    branch_id UUID,
    user_id UUID,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    FOREIGN KEY (branch_id) REFERENCES branches(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);


-- migrate:down

