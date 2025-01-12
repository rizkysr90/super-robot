-- migrate:up
CREATE TABLE IF NOT EXISTS permissions (
    code VARCHAR(100) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    category VARCHAR(50) NOT NULL
);

-- Create an index on category for faster lookups
CREATE INDEX idx_permissions_category ON permissions(category);

COMMENT ON TABLE permissions IS 'Stores system-wide permissions and their categories';
COMMENT ON COLUMN permissions.code IS 'Unique identifier for the permission';
COMMENT ON COLUMN permissions.name IS 'Human-readable name of the permission';
COMMENT ON COLUMN permissions.category IS 'Category grouping for the permission';

-- migrate:down

