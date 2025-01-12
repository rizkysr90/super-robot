-- migrate:up
CREATE TABLE IF NOT EXISTS user_assignments (
   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   tenant_id UUID NOT NULL,
   user_id UUID NOT NULL,
   branch_id UUID,
   created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP,
   deleted_at TIMESTAMP,
   created_by UUID NOT NULL,
   updated_by UUID,
   deleted_by UUID,
   CONSTRAINT fk_tenant 
       FOREIGN KEY (tenant_id) 
       REFERENCES tenants(id),
   CONSTRAINT fk_user 
       FOREIGN KEY (user_id) 
       REFERENCES users(id),
   CONSTRAINT fk_branch 
       FOREIGN KEY (branch_id) 
       REFERENCES branches(id),
   CONSTRAINT fk_created_by
       FOREIGN KEY (created_by)
       REFERENCES users(id),
   CONSTRAINT fk_updated_by
       FOREIGN KEY (updated_by)
       REFERENCES users(id),
   CONSTRAINT fk_deleted_by
       FOREIGN KEY (deleted_by)
       REFERENCES users(id)
);

-- Create indexes for foreign keys
CREATE INDEX IF NOT EXISTS idx_user_assignments_tenant_id ON user_assignments(tenant_id);
CREATE INDEX IF NOT EXISTS idx_user_assignments_user_id ON user_assignments(user_id);
CREATE INDEX IF NOT EXISTS idx_user_assignments_branch_id ON user_assignments(branch_id);
-- Add table comments
COMMENT ON TABLE user_assignments IS 'Tracks user assignments to branches. NULL branch_id means access to all branches';
COMMENT ON COLUMN user_assignments.id IS 'Primary key';
COMMENT ON COLUMN user_assignments.tenant_id IS 'Reference to the tenant';
COMMENT ON COLUMN user_assignments.user_id IS 'Reference to the user';
COMMENT ON COLUMN user_assignments.branch_id IS 'Reference to the branch. NULL means access to all branches';
COMMENT ON COLUMN user_assignments.created_by IS 'User who created the assignment';
COMMENT ON COLUMN user_assignments.updated_by IS 'User who last updated the assignment';
COMMENT ON COLUMN user_assignments.deleted_by IS 'User who deleted the assignment';
-- migrate:down

