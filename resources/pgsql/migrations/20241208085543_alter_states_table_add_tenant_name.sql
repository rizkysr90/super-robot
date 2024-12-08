-- migrate:up
ALTER TABLE states 
ADD COLUMN tenant_name VARCHAR(255) NULL;

-- migrate:down

