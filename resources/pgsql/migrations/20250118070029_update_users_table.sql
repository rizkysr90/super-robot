-- migrate:up
ALTER TABLE users 
ADD COLUMN created_by TEXT;

-- migrate:down

