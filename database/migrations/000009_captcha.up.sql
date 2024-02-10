CREATE TABLE captchas (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    answer TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
)
