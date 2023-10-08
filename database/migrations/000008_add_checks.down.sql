ALTER TABLE boards
    DROP CONSTRAINT code_length,
    DROP CONSTRAINT name_length,
    DROP CONSTRAINT bump_limit_not_negative,
    DROP CONSTRAINT page_limit_not_negative;

ALTER TABLE admins
    DROP CONSTRAINT username_length,
    DROP CONSTRAINT salt_size,
    DROP CONSTRAINT hash_size;

ALTER TABLE access_tokens
    DROP CONSTRAINT value_length;

ALTER TABLE files
    DROP CONSTRAINT filepath_length,
    DROP CONSTRAINT size_is_positive,
    DROP CONSTRAINT mimetype_length;
