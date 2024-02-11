CREATE TABLE access_tokens (
    value TEXT PRIMARY KEY NOT NULL,
    admin_id INT NOT NULL
        REFERENCES admins(id)
        ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL
);
