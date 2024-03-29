CREATE TABLE files (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    post_id BIGINT NOT NULL
        REFERENCES posts(id)
        ON DELETE CASCADE,
    filepath TEXT NOT NULL,
    name TEXT NOT NULL,
    size INTEGER NOT NULL,
    mimetype TEXT NOT NULL
);
