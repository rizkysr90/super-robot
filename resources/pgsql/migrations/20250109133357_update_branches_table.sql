-- migrate:up
ALTER TABLE branches
    ALTER COLUMN created_by DROP NOT NULL,
    ADD COLUMN created_by_owner VARCHAR(255);

-- migrate:down

