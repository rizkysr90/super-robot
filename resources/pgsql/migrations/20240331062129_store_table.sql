-- migrate:up
CREATE TABLE IF NOT EXISTS stores (
    id UUID PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
    address TEXT NOT NULL,
    contact VARCHAR(50) NOT NULL,
    employee_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ NULL

    -- CONSTRAINT fk_employee_store FOREIGN KEY (employee_id) 
    -- REFERENCES employees(id)
);

-- migrate:down

