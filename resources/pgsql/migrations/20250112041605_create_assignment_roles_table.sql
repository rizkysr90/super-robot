-- migrate:up
CREATE TABLE IF NOT EXISTS user_assignment_roles (
   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   tenant_id UUID NOT NULL,
   assignment_id UUID NOT NULL,
   tenant_role_id UUID NOT NULL,
   created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP,
   deleted_at TIMESTAMP,
   created_by UUID NOT NULL,
   updated_by UUID,
   deleted_by UUID,
   CONSTRAINT fk_tenant 
       FOREIGN KEY (tenant_id) 
       REFERENCES tenants(id),
   CONSTRAINT fk_assignment 
       FOREIGN KEY (assignment_id) 
       REFERENCES user_assignments(id),
   CONSTRAINT fk_tenant_role
       FOREIGN KEY (tenant_role_id) 
       REFERENCES tenant_roles(id)
);

-- Create indexes for foreign keys
CREATE INDEX IF NOT EXISTS idx_user_assignment_roles_tenant_id ON user_assignment_roles(tenant_id);
CREATE INDEX IF NOT EXISTS idx_user_assignment_roles_assignment_id ON user_assignment_roles(assignment_id);
CREATE INDEX IF NOT EXISTS idx_user_assignment_roles_tenant_role_id ON user_assignment_roles(tenant_role_id);

-- migrate:down

