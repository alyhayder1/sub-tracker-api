CREATE TABLE users (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clerk_id         TEXT NOT NULL UNIQUE,
    email            TEXT NOT NULL UNIQUE,
    full_name        TEXT,
    default_currency CHAR(3) NOT NULL DEFAULT 'EUR',
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);