ALTER TABLE boards
    ADD CONSTRAINT code_length
    CHECK (TRIM(code) <> ''),
    ADD CONSTRAINT name_length
    CHECK (TRIM(name) <> ''),
    ADD CONSTRAINT bump_limit_not_negative
    CHECK (bump_limit >= 0),
    ADD CONSTRAINT page_limit_not_negative
    CHECK (page_limit >= 0);

ALTER TABLE admins
    ADD CONSTRAINT username_length
    CHECK (TRIM(username) <> ''),
    ADD CONSTRAINT salt_size
    CHECK (length(salt) > 0),
    ADD CONSTRAINT hash_size
    CHECK (length(hash) > 0);

ALTER TABLE access_tokens
    ADD CONSTRAINT value_length
    CHECK (value <> '');

ALTER TABLE files
    ADD CONSTRAINT filepath_length
    CHECK (filepath <> ''),
    ADD CONSTRAINT size_is_positive
    CHECK (size > 0),
    ADD CONSTRAINT mimetype_length
    CHECK (mimetype <> '');
