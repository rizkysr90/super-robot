-- migrate:up
CREATE TABLE IF NOT EXISTS tasks (
    task_id INT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    name TEXT NOT NULL,
    created_at TIMESTAMP
);

-- migrate:down

