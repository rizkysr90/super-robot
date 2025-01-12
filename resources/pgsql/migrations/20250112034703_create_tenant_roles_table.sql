-- migrate:up
CREATE TABLE IF NOT EXISTS tenant_roles (
   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   tenant_id UUID NOT NULL,
   name VARCHAR(255) NOT NULL,
   is_head_office_role BOOLEAN NOT NULL,
   description TEXT,
   created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP,
   deleted_at TIMESTAMP,
   created_by TEXT NOT NULL,
   deleted_by TEXT,
   updated_by TEXT,
   CONSTRAINT fk_tenant 
       FOREIGN KEY (tenant_id) 
       REFERENCES tenants(id)
);

-- Create index for foreign key
CREATE INDEX IF NOT EXISTS idx_tenant_roles_tenant_id ON tenant_roles(tenant_id);


-- migrate:down

