-- migrate:up
ALTER TABLE employees
ADD COLUMN IF NOT EXISTS refresh_token TEXT;


-- migrate:down

