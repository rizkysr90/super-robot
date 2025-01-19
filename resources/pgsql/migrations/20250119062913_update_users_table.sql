-- migrate:up
ALTER TABLE users 
ADD COLUMN is_verified SMALLINT DEFAULT 1;

-- migrate:down

