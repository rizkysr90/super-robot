-- migrate:up
-- Create the products table with updated structure
CREATE TABLE IF NOT EXISTS products (
    product_id CHAR(13) PRIMARY KEY,
    product_name VARCHAR(100) NOT NULL,
    price DECIMAL(17, 2) NOT NULL,
    base_price DECIMAL(17, 2) NOT NULL,
    stock_quantity INT NOT NULL DEFAULT 0,
    category_id UUID,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories(id)
);

-- Create indexes
CREATE INDEX idx_product_name ON products(product_name);
CREATE INDEX idx_category_id ON products(category_id);
CREATE INDEX idx_deleted_at ON products(deleted_at);


-- migrate:down

