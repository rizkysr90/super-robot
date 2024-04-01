-- migrate:up
ALTER TABLE employees
ADD COLUMN IF NOT EXISTS user_id UUID NULL;

ALTER TABLE employees
ADD CONSTRAINT fk_user_employees FOREIGN KEY (user_id) 
    REFERENCES users(id);

-- migrate:down

