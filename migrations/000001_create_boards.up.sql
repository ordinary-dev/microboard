CREATE TABLE boards (
    code TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    page_limit SMALLINT NOT NULL DEFAULT 0,
    bump_limit SMALLINT NOT NULL DEFAULT 0,
    unlisted BOOLEAN DEFAULT FALSE
);
