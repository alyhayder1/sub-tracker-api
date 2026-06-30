CREATE TABLE subscriptions (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id           UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider_key      TEXT REFERENCES providers(key),
    name              TEXT NOT NULL,
    category          TEXT NOT NULL DEFAULT 'other',
    amount            NUMERIC(12,2) NOT NULL,
    currency          CHAR(3) NOT NULL DEFAULT 'EUR',
    billing_cycle     TEXT NOT NULL,
    next_billing_date DATE,
    status            TEXT NOT NULL DEFAULT 'active',
    source            TEXT NOT NULL DEFAULT 'manual',
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);