-- migrate:up
ALTER TABLE users 
ADD CONSTRAINT users_tenant_id_fkey 
FOREIGN KEY (tenant_id) REFERENCES tenants(id);

ALTER TABLE tenants 
ADD CONSTRAINT tenants_owner_id_fkey 
FOREIGN KEY (owner_id) REFERENCES users(id);

-- migrate:down

