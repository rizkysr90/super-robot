-- migrate:up
CREATE TABLE IF NOT EXISTS tenant_role_permissions (
   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   permission_code VARCHAR(100) NOT NULL,
   tenant_role_id UUID NOT NULL,
   tenant_id UUID NOT NULL,
   created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP,
   deleted_at TIMESTAMP,
   created_by TEXT NOT NULL,
   deleted_by TEXT,
   updated_by TEXT,
   CONSTRAINT fk_permission_code 
       FOREIGN KEY (permission_code) 
       REFERENCES permissions(code),
   CONSTRAINT fk_tenant_role 
       FOREIGN KEY (tenant_role_id) 
       REFERENCES tenant_roles(id),
   CONSTRAINT fk_tenant 
       FOREIGN KEY (tenant_id) 
       REFERENCES tenants(id)
);

-- Create indexes for foreign keys
CREATE INDEX IF NOT EXISTS idx_trp_permission_code ON tenant_role_permissions(permission_code);
CREATE INDEX IF NOT EXISTS idx_trp_tenant_role_id ON tenant_role_permissions(tenant_role_id); 
CREATE INDEX IF NOT EXISTS idx_trp_tenant_id ON tenant_role_permissions(tenant_id);

-- migrate:down

