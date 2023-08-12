CREATE TABLE threads (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    board_code TEXT NOT NULL
        REFERENCES boards(code)
        ON DELETE CASCADE,
    updated_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE INDEX idx_thread_updated_at ON threads(updated_at);
