-- migrate:up
CREATE TABLE IF NOT EXIST states {
    state_id TEXT PRIMARY KEY
}

-- migrate:down

