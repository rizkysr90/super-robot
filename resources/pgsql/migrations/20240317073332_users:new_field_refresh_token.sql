-- migrate:up
ALTER TABLE users
ADD COLUMN IF NOT EXISTS refresh_token TEXT;


-- migrate:down

