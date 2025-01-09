-- migrate:up
CREATE TABLE IF NOT EXISTS branches (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    address TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by VARCHAR(255) NOT NULL,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    CONSTRAINT unique_tenant_branch UNIQUE (tenant_id, name)
);

-- migrate:down

